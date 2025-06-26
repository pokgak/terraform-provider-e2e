package models

type DBResponse struct {
	Code    int    `json:"code"`
	Data    DB     `json:"data"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type DB struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Status              string   `json:"status"`
	StatusTitle         string   `json:"status_title"`
	StatusActions       []string `json:"status_actions"`
	NumInstances        int      `json:"num_instances"`
	Software            Software `json:"software"`
	MasterNode          DBNode   `json:"master_node"`
	ConnectivityDetail  string   `json:"connectivity_detail"`
	VectorDBStatus      string   `json:"vector_database_status"`
	ProjectName         string   `json:"project_name"`
	SnapshotExist       bool     `json:"snapshot_exist"`
	ZookeeperInstances  int      `json:"zookeeper_instances"`
	SlaveInstances      int      `json:"slave_instances"`
	IsEncryptionEnabled bool     `json:"isEncryptionEnabled"`
}

type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Engine  string `json:"engine"`
}

type DBNode struct {
	NodeName           string         `json:"node_name"`
	InstanceID         int            `json:"instance_id"`
	ClusterID          int            `json:"cluster_id"`
	NodeID             int            `json:"node_id"`
	VMID               int            `json:"vm_id"`
	Port               string         `json:"port"`
	PublicIPAddress    string         `json:"public_ip_address"`
	PrivateIPAddress   string         `json:"private_ip_address"`
	AllowedIPAddresses AllowedIPs     `json:"allowed_ip_address"`
	ZabbixHostID       *int           `json:"zabbix_host_id"`
	Database           DBCreds        `json:"database"`
	RAM                string         `json:"ram"`
	CPU                string         `json:"cpu"`
	Disk               string         `json:"disk"`
	Status             string         `json:"status"`
	DBStatus           string         `json:"db_status"`
	CreatedAt          string         `json:"created_at"`
	Plan               Plan           `json:"plan"`
	SSL                bool           `json:"ssl"`
	Domain             *string        `json:"domain"`
	PublicPort         *string        `json:"public_port"`
	CommittedInfo      interface{}    `json:"committed_info"`
	CommittedDetails   []CommittedSKU `json:"committed_details"`
}

type DBCreds struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Database string   `json:"database"`
	PGDetail PGDetail `json:"pg_detail"`
}

type AllowedIPs struct {
	WhitelistedIPs         []string `json:"whitelisted_ips"`
	TempIPs                []string `json:"temp_ips"`
	WhitelistedIPsTags     []string `json:"whitelisted_ips_tags"`
	TempIPsTags            []string `json:"temp_ips_tags"`
	WhitelistingInProgress bool     `json:"whitelisting_in_progress"`
}

type Plan struct {
	Name                   string         `json:"name"`
	Price                  string         `json:"price"`
	TemplateID             int            `json:"template_id"`
	RAM                    string         `json:"ram"`
	CPU                    string         `json:"cpu"`
	Disk                   string         `json:"disk"`
	Currency               string         `json:"currency"`
	Software               Software       `json:"software"`
	AvailableInventoryStat bool           `json:"available_inventory_status"`
	PricePerHour           float64        `json:"price_per_hour"`
	PricePerMonth          float64        `json:"price_per_month"`
	CommittedSKU           []CommittedSKU `json:"committed_sku"`
}

type CommittedSKU struct {
	CommittedSKUID       int     `json:"committed_sku_id"`
	CommittedSKUName     string  `json:"committed_sku_name"`
	CommittedNodeMessage string  `json:"committed_node_message"`
	CommittedSKUPrice    float64 `json:"committed_sku_price"`
	CommittedUptoDate    string  `json:"committed_upto_date"`
	CommittedDays        int     `json:"committed_days"`
}

type MySqlCreate struct {
	Name             string      `json:"name"`
	Database         DBConfig    `json:"database"`
	Vpcs             []VpcDetail `json:"vpcs"`
	SoftwareID       int         `json:"software_id"`
	TemplateID       int         `json:"template_id"`
	ParameterGroupId int         `json:"pg_id,omitempty"`
	PublicIPRequired bool        `json:"public_ip_required"`
	Group            string      `json:"group"`
}

type DBConfig struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	DBaaSNumber int    `json:"dbaas_number"`
}

type PGDetail struct {
	ID int `json:"pg_id"`
}

type AttachVPCPayloadRequest struct {
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
