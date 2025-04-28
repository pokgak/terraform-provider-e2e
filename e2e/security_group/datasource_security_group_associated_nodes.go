package security_group

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceSecurityGroupAssociatedNodes() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of nodes (VMs) associated with a security group",

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "ID of the security group to get associated nodes for",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID of the project",
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Location of the resources",
				ValidateFunc: validation.StringInSlice([]string{"Delhi", "Mumbai"}, false),
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of nodes associated with the security group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vm_id": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "ID of the virtual machine",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the node",
						},
						"ip_address_private": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP address of the node",
						},
						"status_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the node (Running, Stopped, etc.)",
						},
						"ip_address_public": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP address of the node",
						},
						"security_group_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of security group association",
						},
					},
				},
			},
		},

		ReadContext: dataSourceSecurityGroupAssociatedNodesRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceSecurityGroupAssociatedNodesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	apiClient := m.(*client.Client)
	securityGroupId := d.Get("security_group_id").(float64)
	projectId := d.Get("project_id").(int)
	location := d.Get("location").(string)

	log.Printf("[INFO] Getting associated nodes for security group ID: %.0f", securityGroupId)

	response, err := apiClient.GetSecurityGroupAssociatedNodes(securityGroupId, projectId, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting associated nodes for security group %.0f", securityGroupId))
	}

	d.Set("nodes", flattenAssociatedNodes(&response.Data))
	d.SetId(strconv.FormatFloat(securityGroupId, 'f', -1, 64) + "-" + strconv.Itoa(projectId) + "-" + location)

	return diags
}

func flattenAssociatedNodes(nodes *[]models.AssociatedNode) []interface{} {
	if nodes != nil {
		ois := make([]interface{}, len(*nodes))

		for i, node := range *nodes {
			oi := make(map[string]interface{})
			oi["vm_id"] = node.VmId
			oi["name"] = node.Name
			oi["ip_address_private"] = node.IpAddressPrivate
			oi["status_name"] = node.StatusName
			oi["ip_address_public"] = node.IpAddressPublic
			oi["security_group_status"] = node.SecurityGroupStatus

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
