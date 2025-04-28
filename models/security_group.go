package models

// BaseResponse represents a generic API response
type BaseResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
	Message string      `json:"message"`
}

// SecurityGroup represents a security group
type SecurityGroup struct {
	Id                  float64 `json:"id"`
	Name                string  `json:"name"`
	Description         string  `json:"description"`
	Rules               []Rule  `json:"rules"`
	Is_default          bool    `json:"is_default"`
	Is_all_traffic_rule bool    `json:"is_all_traffic_rule"`
}

// SecurityGroupListResponse represents the response for listing security groups
type SecurityGroupListResponse struct {
	Code    int             `json:"code"`
	Data    []SecurityGroup `json:"data"`
	Errors  interface{}     `json:"errors"`
	Message string          `json:"message"`
}

// Rule represents a security group rule
type Rule struct {
	Id            float64  `json:"id"`
	Rule_type     string   `json:"rule_type"`
	Protocol_name string   `json:"protocol_name"`
	Port_range    string   `json:"port_range"`
	Network       string   `json:"network"`
	Network_cidr  string   `json:"network_cidr"`
	Network_size  float64  `json:"network_size"`
	Vpc_id        *float64 `json:"vpc_id"`
}

// SecurityGroupCreateRequest represents the payload for creating/updating a security group
type SecurityGroupCreateRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Rules       []RuleCreate `json:"rules"`
}

// RuleCreate represents a rule in create/update requests
type RuleCreate struct {
	Id            *float64 `json:"id,omitempty"`
	Network       string   `json:"network"`
	Rule_type     string   `json:"rule_type"`
	Protocol_name string   `json:"protocol_name"`
	Port_range    string   `json:"port_range"`
	Network_cidr  *string  `json:"network_cidr,omitempty"`
	Network_size  *float64 `json:"network_size,omitempty"`
	Vpc_id        *float64 `json:"vpc_id,omitempty"`
}

type AssociatedNode struct {
	VmId                float64 `json:"vm_id"`
	Name                string  `json:"name"`
	IpAddressPrivate    string  `json:"ip_address_private"`
	StatusName          string  `json:"status_name"`
	IpAddressPublic     string  `json:"ip_address_public"`
	SecurityGroupStatus string  `json:"security_group_status"`
}

type SecurityGroupAssociatedNodesResponse struct {
	Code    int              `json:"code"`
	Data    []AssociatedNode `json:"data"`
	Errors  interface{}      `json:"errors"`
	Message string           `json:"message"`
}

type AssociatedScalar struct {
	Name         string `json:"name"`
	ScalerStatus string `json:"scaler_status"`
}

type SecurityGroupAssociatedScalarsResponse struct {
	Code    int                `json:"code"`
	Data    []AssociatedScalar `json:"data"`
	Errors  interface{}        `json:"errors"`
	Message string             `json:"message"`
}
