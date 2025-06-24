package models

/*
Package models

This file contains all struct definitions used for interacting with the MariaDB API.

It includes:
- Response models for unmarshalling API responses
- Request models for creating a MariaDB service
- Nested types for database credentials, plans, IP settings, etc.
*/

// --------------------
// API Response Structs
// --------------------

// Top-level API response wrapper
type MariaDBResponse struct {
	Code    int     `json:"code"`
	Data    MariaDB `json:"data"`
	Errors  any     `json:"errors"`
	Message string  `json:"message"`
}

// -----------------------------
// MariaDB Cluster State (GET)
// -----------------------------

type MariaDB struct {
	ID                  int           `json:"id"`
	Name                string        `json:"name"`
	Status              string        `json:"status"`
	StatusTitle         string        `json:"status_title"`
	StatusActions       []string      `json:"status_actions"`
	NumInstances        int           `json:"num_instances"`
	Software            Software      `json:"software"`
	MasterNode          MariaDBNode   `json:"master_node"`
	ConnectivityDetail  string        `json:"connectivity_detail"`
	VectorDBStatus      string        `json:"vector_database_status"`
	ProjectName         string        `json:"project_name"`
	SnapshotExist       bool          `json:"snapshot_exist"`
	ZookeeperInstances  int           `json:"zookeeper_instances"`
	SlaveInstances      int           `json:"slave_instances"`
	IsEncryptionEnabled bool          `json:"isEncryptionEnabled"`
}

// Software version info (name, version, engine)
type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Engine  string `json:"engine"`
}

// --------------------------
// Master Node Details (GET)
// --------------------------

type MariaDBNode struct {
	NodeName         string         `json:"node_name"`
	InstanceID       int            `json:"instance_id"`
	ClusterID        int            `json:"cluster_id"`
	NodeID           int            `json:"node_id"`
	VMID             int            `json:"vm_id"`
	Port             string         `json:"port"`
	PublicIPAddress  string         `json:"public_ip_address"`
	PrivateIPAddress string         `json:"private_ip_address"`
	AllowedIPs       AllowedIPs     `json:"allowed_ip_address"`
	ZabbixHostID     *int           `json:"zabbix_host_id"`
	Database         MariaDBCreds   `json:"database"`
	RAM              string         `json:"ram"`
	CPU              string         `json:"cpu"`
	Disk             string         `json:"disk"`
	Status           string         `json:"status"`
	DBStatus         string         `json:"db_status"`
	CreatedAt        string         `json:"created_at"`
	Plan             Plan           `json:"plan"`
	SSL              bool           `json:"ssl"`
	Domain           *string        `json:"domain"`
	PublicPort       *int           `json:"public_port"`
	CommittedInfo    any            `json:"committed_info"`
	CommittedDetails []CommittedSKU `json:"committed_details"`
}

// DB user credentials returned from GET response
type MariaDBCreds struct {
	ID       int       `json:"id"`
	Username string    `json:"username"`
	Database string    `json:"database"`
	PGDetail *PGDetail `json:"pg_detail"` // pointer for null safety
}

// Parameter group detail (pg_detail)
type PGDetail struct {
	Name   string `json:"name"`
	Family string `json:"family"`
	PGID   int    `json:"pg_id"`
}

// Allowed IPs for firewall rules
type AllowedIPs struct {
	WhitelistedIPs      []string `json:"whitelisted_ips"`
	TempIPs             []string `json:"temp_ips"`
	WhitelistedIPsTags  []string `json:"whitelisted_ips_tags"`
	TempIPsTags         []string `json:"temp_ips_tags"`
	WhitelistingRunning bool     `json:"whitelisting_in_progress"`
}

// -------------------------------
// Pricing Plan (Attached to Node)
// -------------------------------

type Plan struct {
	Name                     string         `json:"name"`
	Price                    string         `json:"price"`
	TemplateID               int            `json:"template_id"`
	RAM                      string         `json:"ram"`
	CPU                      string         `json:"cpu"`
	Disk                     string         `json:"disk"`
	Currency                 string         `json:"currency"`
	Software                 Software       `json:"software"`
	AvailableInventoryStatus bool           `json:"available_inventory_status"`
	PricePerHour             float64        `json:"price_per_hour"`
	PricePerMonth            float64        `json:"price_per_month"`
	CommittedSKUs            []CommittedSKU `json:"committed_sku"`
}

// Reserved instance pricing options
type CommittedSKU struct {
	ID       int     `json:"committed_sku_id"`
	Name     string  `json:"committed_sku_name"`
	Message  string  `json:"committed_node_message"`
	Price    float64 `json:"committed_sku_price"`
	UptoDate string  `json:"committed_upto_date"`
	Days     int     `json:"committed_days"`
}

// ------------------------------
// MariaDB Create Request (POST)
// ------------------------------

type MariaDBCreateRequest struct {
	Name                 string        `json:"name"`
	SoftwareID           int           `json:"software_id"`
	TemplateID           int           `json:"template_id"`
	PublicIPRequired     bool          `json:"public_ip_required"`      // user input
	Group                string        `json:"group"`
	VPCs                 []VPCMetadata `json:"vpcs,omitempty"`
	Database             DBConfig      `json:"database"`
	PGID                 int           `json:"pg_id"`                   // Parameter Group ID
	IsEncryptionEnabled  bool          `json:"isEncryptionEnabled"`     // Encryption flag
	EncryptionPassphrase string        `json:"encryption_passphrase"`   // Empty if not provided
}

// DB config (nested inside POST payload)
type DBConfig struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	DBaaSNumber int    `json:"dbaas_number"`
}

// ------------------------------
// VPC Attach/Detach Payloads
// ------------------------------

type AttachDetachVPCRequest struct {
	Action string        `json:"action"` // "attach" or "detach"
	VPCs   []VPCMetadata `json:"vpcs"`
}

type VPCMetadata struct {
	NetworkID string `json:"network_id"`
	VPCName   string `json:"vpc_name"`
	IPv4CIDR  string `json:"ipv4_cidr"`
}

// Used for attaching or detaching a parameter group
type ParameterGroupRequest struct {
	Action string `json:"action"` // "add" or "remove"
}

//tto upgrade plan of a MariaDB cluster
type UpgradePlanRequest struct {
	TemplateID int `json:"template_id"`
}


// DiskUpgradeRequest is used to specify additional disk size to add to a MariaDB cluster
type DiskUpgradeRequest struct {
	Size int `json:"size"`
}
