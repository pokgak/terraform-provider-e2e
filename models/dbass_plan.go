package models

/*
Package models

This file defines struct models to parse the response from the /rds/plans/ API.

Purpose:
- Used to look up software IDs and template IDs based on user-friendly values (like plan name/version).
- Required during MariaDB creation to map Terraform input to backend identifiers.
*/

// Response structure from /rds/plans/ endpoint
type PlanResponse struct {
	Code    int                    `json:"code"`
	Data    PlanData               `json:"data"`
	Errors  map[string]interface{} `json:"errors"`
	Message string                 `json:"message"`
}

// Holds available plans and supported database engines
type PlanData struct {
	TemplatePlans   []PlanTemplate     `json:"template_plans"`
	DatabaseEngines []EngineDefinition `json:"database_engines"`
}

// Pricing and resource configuration for a DB plan
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

// Software details attached to a plan
type TemplateSoftware struct {
	SoftwareName    string `json:"name"`
	SoftwareVersion string `json:"version"`
	SoftwareEngine  string `json:"engine"`
}

// Reserved/discounted pricing options (like 6-month commitment)
type PlanCommittedSKU struct {
	SKUID           int     `json:"committed_sku_id"`
	SKUName         string  `json:"committed_sku_name"`
	SKUNodeMessage  string  `json:"committed_node_message"`
	SKUPrice        float64 `json:"committed_sku_price"`
	SKUEndDate      string  `json:"committed_upto_date"`
	SKUDurationDays int     `json:"committed_days"`
}

// Supported database engine (e.g., MariaDB, PostgreSQL)
type EngineDefinition struct {
	EngineID          int     `json:"id"`
	EngineName        string  `json:"name"`
	EngineVersion     string  `json:"version"`
	EngineType        string  `json:"engine"`
	EngineDescription *string `json:"description"`
}
