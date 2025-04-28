package security_group

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Manages an E2E Networks security group and its rules",

		CreateContext: resourceSecurityGroupCreate,
		ReadContext:   resourceSecurityGroupRead,
		UpdateContext: resourceSecurityGroupUpdate,
		DeleteContext: resourceSecurityGroupDelete,

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The ID of the security group assigned by the API",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the security group",
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Description of the security group",
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether this is the default security group",
			},
			"project_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				Description:  "ID of the project",
				ValidateFunc: validation.IntAtLeast(1),
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Location of the security group",
				ValidateFunc: validation.StringInSlice([]string{"Delhi", "Mumbai"}, false),
			},
			"rules": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Set of security group rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "ID of the rule",
						},
						"network": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Network type (any, manual)",
							ValidateFunc: validation.StringInSlice([]string{"any", "manual"}, false),
						},
						"rule_type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Direction of traffic (Inbound, Outbound)",
							ValidateFunc: validation.StringInSlice([]string{"Inbound", "Outbound"}, false),
						},
						"protocol_name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Protocol name (All, ICMP, Custom_TCP, etc.)",
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
						"port_range": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Port range (All, single port, or comma-separated list)",
							ValidateFunc: validation.StringLenBetween(0, 255),
						},
						"network_cidr": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "CIDR block when network is 'manual'",
							ValidateFunc: validation.IsCIDR,
						},
						"network_size": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Description:  "Network size when network is 'manual'",
							ValidateFunc: validation.FloatAtLeast(1),
						},
					},
				},
				Set: resourceSecurityGroupRuleHash,
			},
		},
	}
}

// resourceSecurityGroupRuleHash generates a unique hash for rules
func resourceSecurityGroupRuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["network"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["rule_type"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["protocol_name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["port_range"].(string)))
	if cidr, ok := m["network_cidr"].(string); ok && cidr != "" {
		buf.WriteString(fmt.Sprintf("%s-", cidr))
	}
	if size, ok := m["network_size"].(float64); ok && size > 0 {
		buf.WriteString(fmt.Sprintf("%f-", size))
	}
	return schema.HashString(buf.String())
}

func resourceSecurityGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	projectID := d.Get("project_id").(int)
	location := d.Get("location").(string)

	request := &models.SecurityGroupCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       expandSecurityGroupRules(d.Get("rules").(*schema.Set).List()),
	}

	log.Printf("[DEBUG] Creating security group: %#v", request)

	// Create the security group
	_, err := client.CreateSecurityGroup(request, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating security group: %w", err))
	}

	// Find the security group by name to get its ID
	securityGroup, err := findSecurityGroupByName(client, request.Name, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error finding created security group: %w", err))
	}

	d.Set("security_group_id", securityGroup.Id)
	d.SetId(fmt.Sprintf("%.0f-%d-%s", securityGroup.Id, projectID, location))

	// Mark as default if specified
	if d.Get("default").(bool) {
		if err := client.MarkSecurityGroupAsDefault(securityGroup.Id, projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("error marking security group as default: %w", err))
		}
	}

	log.Printf("[INFO] Security group created with ID %.0f", securityGroup.Id)
	return resourceSecurityGroupRead(ctx, d, m)
}

func resourceSecurityGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	securityGroupID := d.Get("security_group_id").(float64)
	projectID := d.Get("project_id").(int)
	location := d.Get("location").(string)

	securityGroup, err := client.GetSecurityGroup(securityGroupID, projectID, location)
	if err != nil {
		if _, ok := err.(interface{ Error() string }); ok && err.Error() == fmt.Sprintf("security group with ID %.0f not found", securityGroupID) {
			log.Printf("[INFO] Security group %.0f not found, removing from state", securityGroupID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading security group: %w", err))
	}

	d.Set("name", securityGroup.Name)
	d.Set("description", securityGroup.Description)
	d.Set("default", securityGroup.Is_default)
	d.Set("security_group_id", securityGroup.Id)
	d.Set("rules", flattenSecurityGroupRules(securityGroup.Rules))

	return nil
}

func resourceSecurityGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	securityGroupID := d.Get("security_group_id").(float64)
	projectID := d.Get("project_id").(int)
	location := d.Get("location").(string)

	// Handle name change
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		log.Printf("[DEBUG] Renaming security group to: %s", newName)
		if err := client.RenameSecurityGroup(securityGroupID, newName, projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("error renaming security group: %w", err))
		}
	}

	// Handle description change
	if d.HasChange("description") {
		newDescription := d.Get("description").(string)
		log.Printf("[DEBUG] Updating description to: %s", newDescription)
		if err := client.UpdateSecurityGroupDescription(securityGroupID, newDescription, projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("error updating description: %w", err))
		}
	}

	// Handle rules change
	if d.HasChange("rules") {
		request := &models.SecurityGroupCreateRequest{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Rules:       expandSecurityGroupRules(d.Get("rules").(*schema.Set).List()),
		}
		log.Printf("[DEBUG] Updating security group rules: %#v", request)
		if err := client.UpdateSecurityGroup(securityGroupID, request, projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("error updating rules: %w", err))
		}
	}

	// Handle default change
	if d.HasChange("default") && d.Get("default").(bool) {
		log.Printf("[DEBUG] Marking security group as default")
		if err := client.MarkSecurityGroupAsDefault(securityGroupID, projectID, location); err != nil {
			return diag.FromErr(fmt.Errorf("error marking as default: %w", err))
		}
	}

	return resourceSecurityGroupRead(ctx, d, m)
}

func resourceSecurityGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	securityGroupID := d.Get("security_group_id").(float64)
	projectID := d.Get("project_id").(int)
	location := d.Get("location").(string)

	log.Printf("[DEBUG] Deleting security group ID %.0f", securityGroupID)

	_, err := client.DeleteSecurityGroup(securityGroupID, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting security group: %w", err))
	}

	log.Printf("[INFO] Security group deleted successfully")
	d.SetId("")
	return nil
}

// expandSecurityGroupRules converts Terraform rule set to API format
func expandSecurityGroupRules(rules []interface{}) []models.RuleCreate {
	result := make([]models.RuleCreate, len(rules))
	for i, rule := range rules {
		r := rule.(map[string]interface{})
		rc := models.RuleCreate{
			Network:       r["network"].(string),
			Rule_type:     r["rule_type"].(string),
			Protocol_name: r["protocol_name"].(string),
			Port_range:    r["port_range"].(string),
		}
		if v, ok := r["network_cidr"].(string); ok && v != "" {
			rc.Network_cidr = &v
		}
		if v, ok := r["network_size"].(float64); ok && v > 0 {
			rc.Network_size = &v
		}
		result[i] = rc
	}
	return result
}

// flattenSecurityGroupRules converts API rules to Terraform state
func flattenSecurityGroupRules(rules []models.Rule) []interface{} {
	result := make([]interface{}, len(rules))
	for i, rule := range rules {
		r := make(map[string]interface{})
		r["id"] = rule.Id
		r["network"] = rule.Network
		r["rule_type"] = rule.Rule_type
		r["protocol_name"] = rule.Protocol_name
		r["port_range"] = rule.Port_range
		r["network_cidr"] = rule.Network_cidr
		r["network_size"] = rule.Network_size
		result[i] = r
	}
	return result
}

// findSecurityGroupByName locates a security group by name
func findSecurityGroupByName(c *client.Client, name string, projectID int, location string) (*models.SecurityGroup, error) {
	groups, err := c.GetSecurityGroups(projectID, location)
	if err != nil {
		return nil, fmt.Errorf("error listing security groups: %w", err)
	}

	for _, group := range groups.Data {
		if group.Name == name {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("security group with name '%s' not found", name)
}
