package scaler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
    "strings"	
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceScalerGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The project ID associated with the scaler group.",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The location where the scaler group will be created.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"plan_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The internal ID of the plan derived from plan_name.",
			},
			"sku_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SKU ID (same as plan_id) used for the scaler group.",
				},
			"slug_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The slug representation of the plan.",
				},

			"vm_image_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					trim := func(name string) string {
						if idx := strings.Index(name, "_"); idx != -1 {
							return name[:idx]
						}
						return name
					}
					trimmedOld := trim(old)
					trimmedNew := trim(new)
				
					log.Printf("[DEBUG] DiffSuppressFunc: old=%s → %s, new=%s → %s", old, trimmedOld, new, trimmedNew)
				
					return trimmedOld == trimmedNew
				},
				
			},
			"vm_image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"my_account_sg_id": {
				Type:     schema.TypeInt,
				Optional: true,         
				Computed: true,         
				Description: "The Security Group ID to attach to the scaler group. If not provided, a default will be fetched from the API.",
			},
			"is_encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Enable encryption for the scaler group. Defaults to false.",
			},
			"encryption_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				ForceNew:    true,
				Description: "Passphrase for encryption (if enabled). Defaults to empty string.",
			},
			"is_public_ip_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "Whether to assign a public IP to nodes. Defaults to true.",
			},

			"min_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"max_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"desired": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				Description: "List of VPC IDs. Defaults to an empty list.",
			},
			"policy": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":           {Type: schema.TypeString, Required: true},
						"adjust":         {Type: schema.TypeInt, Required: true},
						"parameter":      {Type: schema.TypeString, Required: true},
						"operator":       {Type: schema.TypeString, Required: true},
						"value":          {Type: schema.TypeString, Required: true},
						"period_number":  {Type: schema.TypeString, Required: true},
						"period_seconds": {Type: schema.TypeString, Required: true},
						"cooldown":       {Type: schema.TypeString, Required: true},
					},
				},
			},
			"scheduled_policy": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type":       {Type: schema.TypeString, Required: true},
						"adjust":     {Type: schema.TypeString, Required: true},
						"recurrence": {Type: schema.TypeString, Required: true},
					},
				},
			},
		},
		CreateContext: resourceCreateScalerGroup,
		ReadContext:   resourceReadScalerGroup,
		DeleteContext: resourceDeleteScalerGroup,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[INFO] Starting CreateScalerGroup operation")

	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	imageName := d.Get("vm_image_name").(string)

	// Fetch image metadata
	savedImage, err := apiClient.GetSavedImageByName(imageName, projectID, location)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch saved image details: %v", err)
		return diag.FromErr(fmt.Errorf("failed to fetch saved image details for '%s': %v", imageName, err))
	}

	if err := d.Set("vm_image_id", savedImage.ImageID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vm_image_id: %v", err))
	}
	if err := d.Set("vm_template_id", savedImage.TemplateID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vm_template_id: %v", err))
	}

	log.Printf("[DEBUG] Image Details → ID: %s, TemplateID: %d, Distro: %s", savedImage.ImageID, savedImage.TemplateID, savedImage.Distro)

	// Determine security group ID
	var sgID int
	if v, ok := d.GetOk("my_account_sg_id"); ok {
		sgID = v.(int)
		log.Printf("[INFO] Using user-provided Security Group ID: %d", sgID)
	} else {
		sgID, err = apiClient.GetDefaultSecurityGroupID(projectID, location)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch default Security Group ID: %v", err)
			return diag.FromErr(fmt.Errorf("failed to fetch default security group ID: %v", err))
		}
		log.Printf("[INFO] Using default Security Group ID from API: %d", sgID)
		if err := d.Set("my_account_sg_id", sgID); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set my_account_sg_id: %v", err))
		}
	}

	// Expand and create the request
	req, err := expandCreateScalerGroupRequest(d, m.(*client.Client), projectID, location, sgID)
if err != nil {
	return diag.FromErr(err)
}

	requestJSON, _ := json.MarshalIndent(req, "", "  ")
	log.Printf("[DEBUG] CreateScalerGroup Request JSON:\n%s", requestJSON)

	resp, err := apiClient.CreateScalerGroup(req, projectID, location)
	if err != nil {
		log.Printf("[ERROR] Failed to create ScalerGroup: %v", err)
		return diag.FromErr(fmt.Errorf("failed to create scaler group: %v", err))
	}

	log.Printf("[INFO] ScalerGroup created with ID: %s", resp.ID)
	d.SetId(resp.ID)

	return resourceReadScalerGroup(ctx, d, m)
}

func resourceReadScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Reading ScalerGroup ID: %s", d.Id())

	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	id := d.Id()

	group, err := apiClient.GetScalerGroup(id, projectID, location)
	if err != nil {
		log.Printf("[ERROR] Failed to read ScalerGroup: %v", err)
		return diag.FromErr(fmt.Errorf("failed to read scaler group: %v", err))
	}

	log.Printf("[DEBUG] Retrieved ScalerGroup: %+v", group)

	// Suppress refresh diff on vm_image_name
	stateVMImageName := d.Get("vm_image_name").(string)
	apiVMImageName := group.VMImageName
	if !strings.HasPrefix(stateVMImageName, apiVMImageName) {
		log.Printf("[INFO] Updating vm_image_name to: %s", apiVMImageName)
		if err := d.Set("vm_image_name", apiVMImageName); err != nil {
			return diag.FromErr(fmt.Errorf("failed to set vm_image_name: %v", err))
		}
	} else {
		log.Printf("[INFO] Keeping existing vm_image_name: %s", stateVMImageName)
	}

	// Set core fields
	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set name: %v", err))
	}
	if err := d.Set("desired", group.Desired); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set desired: %v", err))
	}
	if err := d.Set("min_nodes", group.MinNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set min_nodes: %v", err))
	}
	if err := d.Set("max_nodes", group.MaxNodes); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set max_nodes: %v", err))
	}
	if err := d.Set("plan_name", group.PlanName); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set plan_name: %v", err))
	}
	if err := d.Set("plan_id", strconv.Itoa(group.PlanID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set plan_id: %v", err))
	}
	if err := d.Set("sku_id", strconv.Itoa(group.PlanID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set sku_id: %v", err))
	}
	if err := d.Set("policy_type", group.PolicyType); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set policy_type: %v", err))
	}
	if err := d.Set("vm_image_id", strconv.Itoa(group.ImageID)); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set vm_image_id: %v", err))
	}

	// Optional: recompute slug_name using template_id and plan_name
	if templateID, ok := d.Get("vm_template_id").(int); ok && templateID > 0 {
		_, slugName, err := apiClient.GetPlanDetailsFromPlanName(templateID, group.PlanName, projectID, location)
		if err == nil {
			d.Set("slug_name", slugName)
		} else {
			log.Printf("[WARN] Failed to recompute slug_name: %v", err)
		}
	}

	log.Printf("[INFO] ScalerGroup ID %s state synced successfully", id)
	return nil
}

func resourceDeleteScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting ScalerGroup ID: %s", d.Id())

	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	id := d.Id()

	if err := apiClient.DeleteScalerGroup(id, projectID, location); err != nil {
		log.Printf("[ERROR] Failed to delete ScalerGroup: %v", err)
		return diag.FromErr(fmt.Errorf("failed to delete scaler group: %v", err))
	}

	d.SetId("")
	log.Println("[INFO] ScalerGroup deleted successfully")
	return nil
}

func expandCreateScalerGroupRequest(d *schema.ResourceData, client *client.Client, projectID, location string, sgID int) (*models.CreateScalerGroupRequest, error) {
	planName := d.Get("plan_name").(string)
	imageName := d.Get("vm_image_name").(string)

	image, err := client.GetSavedImageByName(imageName, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch saved image: %v", err)
	}

	planID, slugName, err := client.GetPlanDetailsFromPlanName(image.TemplateID, planName, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch plan details: %v", err)
	}

	var elasticPolicies []models.ElasticPolicy
	for _, p := range d.Get("policy").([]interface{}) {
		pMap := p.(map[string]interface{})
		elasticPolicies = append(elasticPolicies, models.ElasticPolicy{
			Type:          pMap["type"].(string),
			Adjust:        pMap["adjust"].(int),
			Parameter:     pMap["parameter"].(string),
			Operator:      pMap["operator"].(string),
			Value:         pMap["value"].(string),
			PeriodNumber:  pMap["period_number"].(string),
			PeriodSeconds: pMap["period_seconds"].(string),
			Cooldown:      pMap["cooldown"].(string),
		})
	}

	var schedPolicies []models.ScheduledPolicy
	for _, s := range d.Get("scheduled_policy").([]interface{}) {
		sMap := s.(map[string]interface{})
		schedPolicies = append(schedPolicies, models.ScheduledPolicy{
			Type:       sMap["type"].(string),
			Adjust:     sMap["adjust"].(string),
			Recurrence: sMap["recurrence"].(string),
		})
	}

	var vpcList []string
	for _, v := range d.Get("vpc").([]interface{}) {
		vpcList = append(vpcList, v.(string))
	}

	return &models.CreateScalerGroupRequest{
		Name:                 d.Get("name").(string),
		PlanName:             planName,
		PlanID:               planID,
		SKUID:                planID,
		SlugName:             slugName,
		VMImageName:          imageName,
		VMImageID:            image.ImageID,
		VMTemplateID:         image.TemplateID,
		MyAccountSGID:        sgID,
		IsEncryptionEnabled:  d.Get("is_encryption_enabled").(bool),
		EncryptionPassphrase: d.Get("encryption_passphrase").(string),
		IsPublicIPRequired:   d.Get("is_public_ip_required").(bool),
		MinNodes:             strconv.Itoa(d.Get("min_nodes").(int)),
		MaxNodes:             strconv.Itoa(d.Get("max_nodes").(int)),
		Desired:              strconv.Itoa(d.Get("desired").(int)),
		PolicyType:           d.Get("policy_type").(string),
		Policy:               elasticPolicies,
		ScheduledPolicy:      schedPolicies,
		VPC:                  vpcList,
	}, nil
}





// func trimVMImageName(imageName string) string {
// 	if idx := indexOfFirstUnderscore(imageName); idx != -1 {
// 		return imageName[:idx]
// 	}
// 	return imageName
// }

// func indexOfFirstUnderscore(s string) int {
// 	for i := range s {
// 		if s[i] == '_' {
// 			return i
// 		}
// 	}
// 	return -1
// }



