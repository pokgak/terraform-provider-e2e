package models

type MySqlCreate struct {
	Name                string         `json:"name"`
	Plan                string         `json:"plan"`
	Database            DatabaseDetail `json:"database"`
	Vpcs                []VpcDetail    `json:"vpcs"`
	SoftwareID          int            `json:"software_id"`
	TemplateID          int            `json:"template_id"`
	Location            string         `json:"location"`
	IsEncryptionEnabled bool           `json:"is_encryption_enabled"`
	ParameterGroupId    int            `json:"pg_id,omitempty"`
	Version             string         `json:"db_version"`
	AttachPublicIp      bool           `json:"attach_public_ip"`
}

type ResponseMySql struct {
	Code    int    `json:"code"`
	Data    MySql  `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type MySql struct {
	ID                  int        `json:"id"`
	Name                string     `json:"name"`
	Status              string     `json:"status"`
	PublicIPAddress     string     `json:"public_ip_address"`
	PrivateIPAddress    string     `json:"private_ip_address"`
	IsEncryptionEnabled bool       `json:"is_encryption_enabled"`
	MasterNode          MasterNode `json:"master_node"`
}

type MasterNode struct {
	NodeName           string            `json:"node_name"`
	InstanceID         int               `json:"instance_id"`
	ClusterID          int               `json:"cluster_id"`
	NodeID             int               `json:"node_id"`
	VMID               int               `json:"vm_id"`
	Port               string            `json:"port"`
	PublicIPAddress    string            `json:"public_ip_address"`
	PrivateIPAddress   string            `json:"private_ip_address"`
	AllowedIPAddresses AllowedIPAddress  `json:"allowed_ip_address"`
	ZabbixHostID       *int              `json:"zabbix_host_id"`
	Database           DatabaseInfo      `json:"database"`
	RAM                string            `json:"ram"`
	CPU                string            `json:"cpu"`
	Disk               string            `json:"disk"`
	Status             string            `json:"status"`
	DBStatus           string            `json:"db_status"`
	CreatedAt          string            `json:"created_at"`
	Plan               Plan              `json:"plan"`
	SSL                bool              `json:"ssl"`
	Domain             *string           `json:"domain"`
	PublicPort         *string           `json:"public_port"`
	CommittedInfo      interface{}       `json:"committed_info"`
	CommittedDetails   []CommittedDetail `json:"committed_details"`
}

type DatabaseDetail struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type PlanResponse struct {
	Code    int            `json:"code"`
	Data    PlanData       `json:"data"`
	Errors  map[string]any `json:"errors"`
	Message string         `json:"message"`
}

type PlanData struct {
	TemplatePlans   []PlanTemplate     `json:"template_plans"`
	DatabaseEngines []EngineDefinition `json:"database_engines"`
}

type PlanTemplate struct {
	PlanName             string             `json:"name"`
	PlanDisplayPrice     string             `json:"price"`
	PlanTemplateID       int                `json:"template_id"`
	PlanRAMGB            string             `json:"ram"`
	PlanCPUCores         string             `json:"cpu"`
	PlanDiskGB           string             `json:"disk"`
	PlanCurrency         string             `json:"currency"`
	PlanSoftware         TemplateSoftware   `json:"software"`
	IsInventoryAvailable bool               `json:"available_inventory_status"`
	PlanHourlyPrice      float64            `json:"price_per_hour"`
	PlanMonthlyPrice     float64            `json:"price_per_month"`
	CommittedSKUs        []PlanCommittedSKU `json:"committed_sku"`
}

type EngineDefinition struct {
	EngineID          int       `json:"id"`
	EngineName        string    `json:"name"`
	EngineVersion     string    `json:"version"`
	EngineType        string    `json:"engine"`
	EngineDescription ***string `json:"description"`
}

type TemplateSoftware struct {
	SoftwareName    string `json:"name"`
	SoftwareVersion string `json:"version"`
	SoftwareEngine  string `json:"engine"`
}

type PlanCommittedSKU struct {
	SKUID           int     `json:"committed_sku_id"`
	SKUName         string  `json:"committed_sku_name"`
	SKUNodeMessage  string  `json:"committed_node_message"`
	SKUPrice        float64 `json:"committed_sku_price"`
	SKUEndDate      string  `json:"committed_upto_date"`
	SKUDurationDays int     `json:"committed_days"`
}

type AllowedIPAddress struct {
	WhitelistedIPs         []string `json:"whitelisted_ips"`
	TempIPs                []string `json:"temp_ips"`
	WhitelistedIPsTags     []string `json:"whitelisted_ips_tags"`
	TempIPsTags            []string `json:"temp_ips_tags"`
	WhitelistingInProgress bool     `json:"whitelisting_in_progress"`
}

type DatabaseInfo struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Database string   `json:"database"`
	PGDetail PGDetail `json:"pg_detail"`
}

type PGDetail struct {
	ID int `json:"pg_id"`
}

type Plan struct {
	Name                   string          `json:"name"`
	Price                  string          `json:"price"`
	TemplateID             int             `json:"template_id"`
	RAM                    string          `json:"ram"`
	CPU                    string          `json:"cpu"`
	Disk                   string          `json:"disk"`
	Currency               string          `json:"currency"`
	Software               SoftwareDetails `json:"software"`
	AvailableInventoryStat bool            `json:"available_inventory_status"`
	PricePerHour           float64         `json:"price_per_hour"`
	PricePerMonth          float64         `json:"price_per_month"`
	CommittedSKU           []CommittedSKU  `json:"committed_sku"`
}

type SoftwareDetails struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Engine  string `json:"engine"`
}

type CommittedSKU struct {
	CommittedSKUID       int     `json:"committed_sku_id"`
	CommittedSKUName     string  `json:"committed_sku_name"`
	CommittedNodeMessage string  `json:"committed_node_message"`
	CommittedSKUPrice    float64 `json:"committed_sku_price"`
	CommittedUptoDate    string  `json:"committed_upto_date"`
	CommittedDays        int     `json:"committed_days"`
}

type CommittedDetail struct {
	CommittedSKUID       int     `json:"committed_sku_id"`
	CommittedSKUName     string  `json:"committed_sku_name"`
	CommittedNodeMessage string  `json:"committed_node_message"`
	CommittedSKUPrice    float64 `json:"committed_sku_price"`
	CommittedUptoDate    string  `json:"committed_upto_date"`
	CommittedDays        int     `json:"committed_days"`
}

type AttachDetachVPC struct {
	Action string      `json:"action"`
	Vpcs   []VpcDetail `json:"vpcs"`
}

type MySQlPlanUpgradeAction struct {
	TemplateID int `json:"template_id"`
}

type MYSQLExpandDisk struct {
	Size int `json:"size"`
}

type VPC struct {
	VpcName    string  `json:"vpc_name,omitempty"`
	Ipv4_cidr  string  `json:"ipv4_cidr,omitempty"`
	Network_id float64 `json:"network_id,omitempty"`
}
