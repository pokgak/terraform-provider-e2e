package scaler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Required:    true,
				// Default:     false,
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
				Required:    true,
				// Default:     true,
				ForceNew:    true,
				Description: "Whether to assign a public IP to nodes. Defaults to true.",
			},
			"provision_status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped"}, false),
				Description: "Set to 'Stopped' to stop the Scaler Group, or 'Running' to start it.",
			},



			"min_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(int)
					if v < 1 {
						errs = append(errs, fmt.Errorf("%q must be at least 1, got: %d", key, v))
					}
					return
				},
			},

			"max_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				// ForceNew: true,
			},
			"desired": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
			},
			"vpc": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"network_id": {
								Type:     schema.TypeInt,
								Computed: true, // if fetched from name
							},
							"ipv4_cidr": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"state": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"subnets": {
								Type:     schema.TypeList,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"id": {
											Type:     schema.TypeInt,
											Computed: true,
										},
										"subnet_name": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"cidr": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"used_ips": {
											Type:     schema.TypeInt,
											Computed: true,
										},

										"total_ips": {
											Type:     schema.TypeInt,
											Computed: true,
										},
									},
								},
							},
						},
					},
				},


			


			"policy": {
				Type:     schema.TypeList,
				Required: true,
				// ForceNew: true,
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
				// ForceNew: true,
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
		UpdateContext: resourceUpdateScalerGroup,
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
			min := diff.Get("min_nodes").(int)
			desired := diff.Get("desired").(int)
			max := diff.Get("max_nodes").(int)
		
			if min > desired {
				return fmt.Errorf("min_nodes (%d) cannot be greater than desired (%d)", min, desired)
			}
		
			if desired > max {
				return fmt.Errorf("desired (%d) cannot be greater than max_nodes (%d)", desired, max)
			}
		
			return nil
		},
		
		
		
		
		
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
	// Normalize provision_status to avoid diff when API returns transitional state
	normalizedStatus := group.ProvisionStatus
	if group.ProvisionStatus == "Stopping" {
		log.Printf("[INFO] Normalizing provision_status 'Stopping' → 'Stopped'")
		normalizedStatus = "Stopped"
	} else if group.ProvisionStatus == "Starting" {
		log.Printf("[INFO] Normalizing provision_status 'Starting' → 'Running'")
		normalizedStatus = "Running"
	}

	if err := d.Set("provision_status", normalizedStatus); err != nil {
		return diag.FromErr(fmt.Errorf("failed to set provision_status: %v", err))
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

	// Expand elastic policies
	var elasticPolicies []models.ElasticPolicy
	if v, ok := d.GetOk("policy"); ok {
		for _, p := range v.([]interface{}) {
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
	}

	// Expand scheduled policies
	var schedPolicies []models.ScheduledPolicy
	if v, ok := d.GetOk("scheduled_policy"); ok {
		for _, s := range v.([]interface{}) {
			sMap := s.(map[string]interface{})
			schedPolicies = append(schedPolicies, models.ScheduledPolicy{
				Type:       sMap["type"].(string),
				Adjust:     sMap["adjust"].(string),
				Recurrence: sMap["recurrence"].(string),
			})
		}
	}

	// Expand VPC block
	var vpcDetails []models.VPCDetail
	if v, ok := d.GetOk("vpc"); ok {
		for _, vRaw := range v.([]interface{}) {
			vMap := vRaw.(map[string]interface{})
			vpcName := vMap["name"].(string)

			vpcMeta, err := client.GetVpcDetailsByName(projectID, location, vpcName)
			if err != nil {
				return nil, fmt.Errorf("failed to get VPC details for %s: %v", vpcName, err)
			}

			// Convert subnet_ids into lookup map
			subnetIDSet := make(map[int]struct{})
			if rawList, ok := vMap["subnet_ids"].([]interface{}); ok {
				for _, sid := range rawList {
					subnetIDSet[sid.(int)] = struct{}{}
				}
			}

			// Filter and collect subnets
			var selectedSubnets []models.SubnetDetail
			for _, subnet := range vpcMeta.Subnets {
				if _, ok := subnetIDSet[subnet.ID]; ok {
					selectedSubnets = append(selectedSubnets, models.SubnetDetail{
						ID:         subnet.ID,
						SubnetName: subnet.SubnetName,
						CIDR:       subnet.CIDR,
						UsedIPs:    subnet.UsedIPs,
						TotalIPs:   subnet.TotalIPs,
					})
				}
			}

			vpcDetails = append(vpcDetails, models.VPCDetail{
				Name:      vpcMeta.Name,
				NetworkID: vpcMeta.NetworkID,
				IPv4CIDR:  vpcMeta.IPv4CIDR,
				State:     vpcMeta.State,
				Subnets:   selectedSubnets,
			})
		}
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
		VPC:                  vpcDetails,
	}, nil
}


func resourceUpdateScalerGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	id := d.Id()

	// Handle stopping the scaler group if provision_status is set to "Stopped"
	// Handle provision_status update separately
	if d.HasChange("provision_status") {
		oldStatus, newStatus := d.GetChange("provision_status")
		log.Printf("[INFO] Changing provision_status from %s → %s", oldStatus, newStatus)

		intID, err := strconv.Atoi(id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid scaler group ID: %v", err))
		}

		if err := apiClient.UpdateScalerGroupStatus(intID, newStatus.(string), projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update provision_status to %s: %v", newStatus, err))
		}

		return resourceReadScalerGroup(ctx, d, m)
	}


	// Validate desired within min-max range
	minNodes := d.Get("min_nodes").(int)
	maxNodes := d.Get("max_nodes").(int)
	desired := d.Get("desired").(int)

	if desired < minNodes || desired > maxNodes {
		return diag.Errorf("desired node count (%d) must be between min_nodes (%d) and max_nodes (%d)", desired, minNodes, maxNodes)
	}

	// If only desired changed, call separate API
	if d.HasChange("desired") &&
		!(d.HasChange("min_nodes") || d.HasChange("max_nodes") || d.HasChange("policy_type") || d.HasChange("policy") || d.HasChange("scheduled_policy")) {
		log.Printf("[INFO] Only desired node count changed; using separate API.")
		intID, err := strconv.Atoi(id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid scaler group ID: %v", err))
		}
		if err := apiClient.UpdateDesiredNodeCount(intID, desired, projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update desired node count: %v", err))
		}
		return resourceReadScalerGroup(ctx, d, m)
	}

	// Skip update if nothing relevant changed
	if !(d.HasChange("min_nodes") || d.HasChange("max_nodes") || d.HasChange("policy_type") || d.HasChange("policy") || d.HasChange("scheduled_policy")) {
		log.Println("[INFO] No relevant changes detected, skipping update.")
		return nil
	}

	// Expand policy fields
	policies := []models.ElasticPolicy{}
	for _, p := range d.Get("policy").([]interface{}) {
		pMap := p.(map[string]interface{})
		policies = append(policies, models.ElasticPolicy{
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

	schedPolicies := []models.ScheduledPolicy{}
	for _, s := range d.Get("scheduled_policy").([]interface{}) {
		sMap := s.(map[string]interface{})
		schedPolicies = append(schedPolicies, models.ScheduledPolicy{
			Type:       sMap["type"].(string),
			Adjust:     sMap["adjust"].(string),
			Recurrence: sMap["recurrence"].(string),
		})
	}

	// Create request object
	req := &models.UpdateScalerGroupRequest{
		Name:            d.Get("name").(string),
		PlanID:          d.Get("plan_id").(string),
		MinNodes:        minNodes,
		MaxNodes:        maxNodes,
		PolicyType:      d.Get("policy_type").(string),
		Policy:          policies,
		ScheduledPolicy: schedPolicies,
	}

	log.Printf("[INFO] Updating ScalerGroup %s with new configuration...", id)
	if err := apiClient.UpdateScalerGroup(id, req, projectID, location); err != nil {
		return diag.FromErr(fmt.Errorf("failed to update scaler group: %v", err))
	}

	return resourceReadScalerGroup(ctx, d, m)
}



// func expandVPCList(d *schema.ResourceData) []VPC {
// 	rawList := d.Get("vpcs").([]interface{})
// 	var vpcs []VPC
  
// 	for _, raw := range rawList {
// 	  data := raw.(map[string]interface{})
// 	  vpc := VPC{
// 		Name:      data["name"].(string),
// 		NetworkID: data["network_id"].(int),
// 		IPv4CIDR:  data["ipv4_cidr"].(string),
// 		State:     data["state"].(string),
// 	  }
  
// 	  if subnetsRaw, ok := data["subnets"]; ok && subnetsRaw != nil {
// 		for _, subnetRaw := range subnetsRaw.([]interface{}) {
// 		  subnet := subnetRaw.(map[string]interface{})
// 		  vpc.Subnets = append(vpc.Subnets, subnet) // You may want to map to struct
// 		}
// 	  }
  
// 	  vpcs = append(vpcs, vpc)
// 	}
  
// 	return vpcs
//   }
  

// func expandCreateScalerGroupRequestWithVPCNames(d *schema.ResourceData, client *client.Client, projectID, location string, sgID int) (*models.CreateScalerGroupRequest, error) {
// 	req, err := expandCreateScalerGroupRequest(d, client, projectID, location, sgID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Resolve vpc_names to vpc_ids
// 	vpcNamesIface, ok := d.GetOk("vpc_names")
// 	if !ok {
// 		return req, nil // No vpc_names provided
// 	}

// 	vpcNames := vpcNamesIface.([]interface{})
// 	var vpcIDs []string
// 	for _, name := range vpcNames {
// 		vpcID, err := client.GetVpcIDByName(name.(string), projectID, location)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to resolve VPC name %q: %v", name, err)
// 		}
// 		vpcIDs = append(vpcIDs, vpcID)
// 	}

// 	req.VPC = vpcIDs
// 	return req, nil
// }











