package mariadb

import (
	"context"
	"fmt"
	"log"
	
	"strconv"
	"strings"

	"github.com/e2eterraformprovider/terraform-provider-e2e/client"
	
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/e2eterraformprovider/terraform-provider-e2e/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This resource allows users to:
//   - Provision a MariaDB instance with a specific version and plan
//   - Configure networking (VPCs, public IP), encryption, and parameter groups
//   - Control lifecycle (create, read, update, delete)
//   - Perform in-place operations such as start, stop, restart, upgrade, and disk expansion
func ResourceMariaDB() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{

			// Unique name for the MariaDB instance (must be unique per user/project)
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the MariaDB service instance.",
			},

			// Computed field: software ID (resolved internally from software_name + software_version)
			"software_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Computed field: template ID (resolved internally from plan + software_id)
			"template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Software name, e.g., "MariaDB"
			"software_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The software name (e.g., MariaDB).",
			},

			// Software version, e.g., "10.6"
			"software_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The software version (e.g., 10.6).",
			},

			// Desired plan name, e.g., DBS.16GB
			"plan_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plan name specifying CPU/memory (e.g. DBS.16GB).",
			},

			// Whether to attach a public IP (default true)
			"public_ip_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether a public IP should be attached during creation or update.",
			},

			// Computed: reflects whether a public IP is actually attached
			"public_ip_attached": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether a public IP is currently attached (backend state).",
			},

			// Optional: parameter group ID to attach; 0 means no parameter group
			"parameter_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "ID of the parameter group to attach. Use 0 to skip.",
			},

			// Group to which this database belongs (e.g. "Default")
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The group to which this database belongs (e.g. 'Default').",
			},

			// List of VPC network IDs to attach to the cluster
			"vpcs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of VPC IDs to associate (optional).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			// Nested database block: configuration for DB name, user, and password
			"database": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				ForceNew:    true,
				Description: "Database configuration (user, password, db name).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database username.",
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Password for the database user.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the database to create.",
						},
						"dbaas_number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "DBaaS number (typically 1 for single instance).",
						},
					},
				},
			},

			// Project where this MariaDB service will be created
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID under which the MariaDB cluster is provisioned.",
			},

			// Region or location of the database (e.g., "Delhi")
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Region where the MariaDB instance will be created.",
			},

			// Optional encryption toggle (default: false)
			"is_encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable encryption at rest for the MariaDB cluster.",
			},

			// Optional encryption passphrase (used only if encryption is enabled)
			"encryption_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "Passphrase for encryption. Leave empty if encryption is not enabled.",
			},

			// Current or desired operational state of the DB
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"STOPPED", "RUNNING", "RESTARTING"},
					false,
				),
				Description: "Operational status: STOPPED, RUNNING, or RESTARTING.",
			},

			// Computed: assigned public IP (if any)
			"public_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP assigned to the master node (if enabled).",
			},

			// Computed: internal private IP assigned to the master node
			"private_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private IP assigned to the master node.",
			},

			// Optional disk expansion (additional GBs to add on update)
			"disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Additional disk size (in GB) to expand during update.",
			},

			// Computed: service port (typically 3306)
			"port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Port number on which the MariaDB service is accessible.",
			},
		},

		CreateContext: resourceCreateMariaDB,
		ReadContext:   resourceReadMariaDB,
		UpdateContext: resourceUpdateMariaDB,
		DeleteContext: resourceDeleteMariaDB,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// This function performs the following steps:
//   1. Resolves required IDs like software_id and template_id
//   2. Prepares and validates the creation request payload
//   3. Converts Terraform inputs (like VPC IDs) into required API formats
//   4. Calls the backend API to create the MariaDB service
//   5. Updates Terraform state with the resulting service details
func resourceCreateMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	//  Extract required fields from Terraform schema 
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)
	softwareName := d.Get("software_name").(string)
	softwareVersion := d.Get("software_version").(string)
	planName := d.Get("plan_name").(string)

	//  Resolve software ID from software name + version 
	softwareID, err := apiClient.GetSoftwareId(projectID, location, softwareName, softwareVersion)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get software ID: %v", err))
	}

	//  Resolve template ID based on plan and software 
	templateID, err := apiClient.GetTemplateId(projectID, location, planName, strconv.Itoa(softwareID))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get template ID: %v", err))
	}

	//  Extract nested database config block (1 item expected) 
	dbConfigList := d.Get("database").([]interface{})
	dbConfigMap := dbConfigList[0].(map[string]interface{})

	//  Extract and convert list of VPC IDs (if any) 
	var vpcIDs []string
	for _, v := range d.Get("vpcs").([]interface{}) {
		vpcIDs = append(vpcIDs, v.(string))
	}

	// Expand VPC metadata required by backend (name, network_id, cidr)
	var vpcList []models.VPCMetadata
	if len(vpcIDs) > 0 {
		vpcList, err = apiClient.ExpandMariaDBVpcList(vpcIDs, projectID, location)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to expand VPC list during create: %v", err))
		}
	}

	//  Read optional Terraform parameters 

	publicIPEnabled := d.Get("public_ip_enabled").(bool)

	parameterGroupID := 0
	if v, ok := d.GetOk("parameter_group_id"); ok {
		parameterGroupID = v.(int)
	}

	isEncryptionEnabled := false
	if v, ok := d.GetOk("is_encryption_enabled"); ok {
		isEncryptionEnabled = v.(bool)
	}

	encryptionPassphrase := ""
	if v, ok := d.GetOk("encryption_passphrase"); ok {
		encryptionPassphrase = v.(string)
	}

	//  Construct MariaDB creation request 
	req := models.MariaDBCreateRequest{
		Name:                 d.Get("name").(string),
		SoftwareID:           softwareID,
		TemplateID:           templateID,
		PublicIPRequired:     publicIPEnabled,
		Group:                d.Get("group").(string),
		VPCs:                 vpcList,
		PGID:                 parameterGroupID,
		IsEncryptionEnabled:  isEncryptionEnabled,
		EncryptionPassphrase: encryptionPassphrase,
		Database: models.DBConfig{
			User:        dbConfigMap["user"].(string),
			Password:    dbConfigMap["password"].(string),
			Name:        dbConfigMap["name"].(string),
			DBaaSNumber: dbConfigMap["dbaas_number"].(int),
		},
	}

	//  Call backend API to create the MariaDB service 
	mariaDB, err := apiClient.CreateMariaDB(&req, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create MariaDB instance: %v", err))
	}

	//  Populate Terraform state with returned data 
	d.SetId(fmt.Sprintf("%d", mariaDB.ID))
	d.Set("name", mariaDB.Name)
	d.Set("status", mariaDB.Status)
	d.Set("public_ip_address", mariaDB.MasterNode.PublicIPAddress)
	d.Set("private_ip_address", mariaDB.MasterNode.PrivateIPAddress)
	d.Set("port", mariaDB.MasterNode.Port)
	d.Set("software_id", softwareID)
	d.Set("template_id", templateID)

	//  Update public_ip_attached (computed) 
	d.Set("public_ip_attached", mariaDB.MasterNode.PublicIPAddress != "")

	return diags
}

func resourceReadMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	// Extract identifiers from Terraform state
	id := d.Id()
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	// Call the E2E Cloud API to get the live MariaDB instance
	mariaDB, err := apiClient.ReadMariaDB(id, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read MariaDB instance: %v", err))
	}

	// Update Terraform state with backend values
	_ = d.Set("name", mariaDB.Name)

	// Normalize "SUSPENDED" (backend) to "STOPPED" (Terraform-expected value)
	status := mariaDB.Status
	if status == "SUSPENDED" {
		status = "STOPPED"
	}
	_ = d.Set("status", status)

	// Software metadata (used for recomputing template_id during plan upgrades)
	_ = d.Set("software_name", mariaDB.Software.Name)
	_ = d.Set("software_version", mariaDB.Software.Version)

	// Plan name (required for diff suppression and correct plan logic)
	_ = d.Set("plan_name", mariaDB.MasterNode.Plan.Name)

	// Network metadata
	_ = d.Set("public_ip_address", mariaDB.MasterNode.PublicIPAddress)
	_ = d.Set("private_ip_address", mariaDB.MasterNode.PrivateIPAddress)
	_ = d.Set("port", mariaDB.MasterNode.Port)

	// Computed field: whether public IP is currently attached
	_ = d.Set("public_ip_attached", mariaDB.MasterNode.PublicIPAddress != "")

	// Encryption state (read-only: backend-enforced)
	_ = d.Set("is_encryption_enabled", mariaDB.IsEncryptionEnabled)

	return diags
}

func resourceDeleteMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics

	id := d.Id()
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	// Trigger API call to delete the MariaDB cluster
	err := apiClient.DeleteMariaDB(id, projectID, location)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete MariaDB instance: %v", err))
	}

	// Clear resource ID from state to mark it as deleted
	d.SetId("")
	return diags
}

func resourceUpdateMariaDB(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	id := d.Id()
	projectID := d.Get("project_id").(string)
	location := d.Get("location").(string)

	//   Cluster Status Update 
	// Change operational status: STOPPED, RUNNING, or RESTARTING.
	if d.HasChange("status") {
		newStatus := d.Get("status").(string)
		switch strings.ToUpper(newStatus) {
		case "STOPPED":
			if err := apiClient.ShutdownMariaDB(id, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to shutdown MariaDB instance: %v", err))
			}
		case "RUNNING":
			if err := apiClient.ResumeMariaDB(id, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to resume MariaDB instance: %v", err))
			}
		case "RESTARTING":
			if err := apiClient.RestartMariaDB(id, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to restart MariaDB instance: %v", err))
			}
		default:
			return diag.FromErr(fmt.Errorf("unsupported status value: %s", newStatus))
		}
	}

	//  VPC Attach/Detach Handling 
	// Detect differences in VPC list and call the appropriate attach/detach APIs.
	if d.HasChange("vpcs") {
		oldRaw, newRaw := d.GetChange("vpcs")
		oldVPCSet := expandStringSet(oldRaw.([]interface{}))
		newVPCSet := expandStringSet(newRaw.([]interface{}))

		var toDetach, toAttach []string

		for vpc := range oldVPCSet {
			if _, exists := newVPCSet[vpc]; !exists {
				toDetach = append(toDetach, vpc)
			}
		}
		for vpc := range newVPCSet {
			if _, exists := oldVPCSet[vpc]; !exists {
				toAttach = append(toAttach, vpc)
			}
		}

		if len(toDetach) > 0 {
			if err := apiClient.DetachVPCFromMariaDB(id, projectID, location, toDetach); err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach VPC(s): %v", err))
			}
		}
		if len(toAttach) > 0 {
			if err := apiClient.AttachVPCToMariaDB(id, projectID, location, toAttach); err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach VPC(s): %v", err))
			}
		}
	}

	//   Public IP Attach/Detach 
	// If public_ip_enabled changes, call appropriate API.
	if d.HasChange("public_ip_enabled") {
		newVal := d.Get("public_ip_enabled").(bool)
		if newVal {
			if err := apiClient.AttachPublicIPToMariaDB(id, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach public IP: %v", err))
			}
		} else {
			if err := apiClient.DetachPublicIPFromMariaDB(id, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach public IP: %v", err))
			}
		}
	}

	//  Parameter Group Attachment/Detachment 
	// Handle attach, detach, or no-op based on changes in parameter_group_id.
	if d.HasChange("parameter_group_id") {
		oldRaw, newRaw := d.GetChange("parameter_group_id")
		oldPGID := oldRaw.(int)
		newPGID := newRaw.(int)

		switch {
		case oldPGID != 0 && newPGID == 0:
			if err := apiClient.DetachParameterGroupFromMariaDB(id, oldPGID, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to detach parameter group: %v", err))
			}
		case newPGID != 0 && newPGID != oldPGID:
			if err := apiClient.AttachParameterGroupToMariaDB(id, newPGID, projectID, location); err != nil {
				return diag.FromErr(fmt.Errorf("failed to attach parameter group: %v", err))
			}
		}
	}

	//  Plan Upgrade (Requires DB to be STOPPED) 
	if d.HasChange("plan_name") {
	oldPlan, newPlan := d.GetChange("plan_name")
	log.Printf("[INFO] Plan change detected: %s -> %s", oldPlan.(string), newPlan.(string))

	// Ensure cluster is STOPPED before upgrading
	status := d.Get("status").(string)
	if strings.ToUpper(status) != "STOPPED" {
		return diag.FromErr(fmt.Errorf("cannot upgrade plan: MariaDB must be STOPPED, current status is '%s'", status))
	}

	// Get software ID based on name/version
	softwareName := d.Get("software_name").(string)
	softwareVersion := d.Get("software_version").(string)
	softwareID, err := apiClient.GetSoftwareId(projectID, location, softwareName, softwareVersion)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get software ID for %s %s: %v", softwareName, softwareVersion, err))
	}

	// Get new template ID using the new plan
	templateID, err := apiClient.GetTemplateId(projectID, location, newPlan.(string), fmt.Sprintf("%d", softwareID))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get template ID for plan %s: %v", newPlan.(string), err))
	}

	// Trigger plan upgrade via API
	if err := apiClient.UpgradeMariaDBPlan(id, projectID, location, templateID); err != nil {
		return diag.FromErr(fmt.Errorf("failed to upgrade MariaDB plan: %v", err))
	}

	log.Printf("[INFO] Successfully upgraded %s %s to plan %s (template_id=%d)", softwareName, softwareVersion, newPlan, templateID)

	//  Update template_id in Terraform state to prevent drift
	_ = d.Set("template_id", templateID)
}



	if d.HasChange("disk_size") {
	additionalSize := d.Get("disk_size").(int)

	if additionalSize > 0 {
		status := d.Get("status").(string)
		if strings.ToUpper(status) != "STOPPED" {
			return diag.FromErr(fmt.Errorf("cannot expand disk: MariaDB must be STOPPED, current status is '%s'", status))
		}

		err := apiClient.ExpandMariaDBDisk(id, projectID, location, additionalSize)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to expand MariaDB disk: %v", err))
		}

		log.Printf("[INFO] Disk expanded by %d GB for cluster %s", additionalSize, id)

		// Reset so that apply doesn't re-trigger
		_ = d.Set("disk_size", 0)
	} else {
		log.Printf("[INFO] disk_size is 0, skipping expansion.")
	}
}


	
	return resourceReadMariaDB(ctx, d, m)
}

func expandStringSet(list []interface{}) map[string]struct{} {
	result := make(map[string]struct{})
	for _, v := range list {
		result[v.(string)] = struct{}{}
	}
	return result
}
