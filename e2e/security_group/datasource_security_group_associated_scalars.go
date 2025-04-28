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

func DataSourceSecurityGroupAssociatedScalars() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of scalars associated with a security group",

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "ID of the security group to get associated scalars for",
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
			"scalars": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of scalars associated with the security group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the scalar",
						},
						"scaler_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the scalar (Running, Stopped, etc.)",
						},
					},
				},
			},
		},

		ReadContext: dataSourceSecurityGroupAssociatedScalarsRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceSecurityGroupAssociatedScalarsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	apiClient := m.(*client.Client)
	securityGroupId := d.Get("security_group_id").(float64)
	projectId := d.Get("project_id").(int)
	location := d.Get("location").(string)

	log.Printf("[INFO] Getting associated scalars for security group ID: %.0f", securityGroupId)

	response, err := apiClient.GetSecurityGroupAssociatedScalars(securityGroupId, projectId, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting associated scalars for security group %.0f", securityGroupId))
	}

	d.Set("scalars", flattenAssociatedScalars(&response.Data))
	d.SetId(strconv.FormatFloat(securityGroupId, 'f', -1, 64) + "-" + strconv.Itoa(projectId) + "-" + location)

	return diags
}

func flattenAssociatedScalars(scalars *[]models.AssociatedScalar) []interface{} {
	if scalars != nil {
		ois := make([]interface{}, len(*scalars))

		for i, scalar := range *scalars {
			oi := make(map[string]interface{})
			oi["name"] = scalar.Name
			oi["scaler_status"] = scalar.ScalerStatus

			ois[i] = oi
		}

		return ois
	}

	return make([]interface{}, 0)
}
