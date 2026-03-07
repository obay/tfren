package azure

import (
	"fmt"
	"strings"
	"unicode"
)

// CharClass represents the set of valid characters for a resource name.
type CharClass int

const (
	// AlphanumHyphen allows alphanumerics and hyphens.
	AlphanumHyphen CharClass = iota
	// AlphanumHyphenUnderscore allows alphanumerics, hyphens, and underscores.
	AlphanumHyphenUnderscore
	// AlphanumHyphenUnderscorePeriod allows alphanumerics, hyphens, underscores, and periods.
	AlphanumHyphenUnderscorePeriod
	// LowercaseNumber allows only lowercase letters and numbers.
	LowercaseNumber
	// LowercaseNumberHyphen allows lowercase letters, numbers, and hyphens.
	LowercaseNumberHyphen
	// Alphanumeric allows only alphanumerics (no hyphens).
	Alphanumeric
	// Broad allows most characters (used for resources with few restrictions).
	Broad
)

// StartEndRule describes what the first or last character must be.
type StartEndRule int

const (
	Any StartEndRule = iota
	Letter
	Alphanum
	LowercaseLetter
	AlphanumUnderscore
)

// CAFRule defines the Azure CAF naming convention rule for a resource type.
type CAFRule struct {
	Abbreviation         string
	AltAbbreviations     []string
	AllowsHyphens        bool
	Pattern              string
	MinSegments          int
	MinLength            int
	MaxLength            int
	ValidChars           CharClass
	MustStartWith        StartEndRule
	MustEndWith          StartEndRule
	CantEndWithHyphen    bool
	CantEndWithPeriod    bool
	NoConsecutiveHyphens bool
}

// CAFRules maps azurerm Terraform resource types to their CAF naming rules.
// Sources:
//   - https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations
//   - https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules
var CAFRules = map[string]CAFRule{
	// Management and governance
	"azurerm_resource_group": {
		Abbreviation: "rg", AllowsHyphens: true, Pattern: "rg-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 90,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
		CantEndWithPeriod: true,
	},
	"azurerm_management_group": {
		Abbreviation: "mg", AllowsHyphens: true, Pattern: "mg-<purpose>",
		MinSegments: 2, MinLength: 1, MaxLength: 90,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
		CantEndWithPeriod: true,
	},
	"azurerm_policy_definition": {
		Abbreviation: "", AllowsHyphens: true, Pattern: "<descriptive>",
		MinSegments: 1, MinLength: 1, MaxLength: 64,
		ValidChars: Broad, MustStartWith: Any,
		CantEndWithPeriod: true,
	},
	"azurerm_automation_account": {
		Abbreviation: "aa", AllowsHyphens: true, Pattern: "aa-<workload>-<env>",
		MinSegments: 2, MinLength: 6, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Letter, MustEndWith: Alphanum,
	},

	// Networking
	"azurerm_virtual_network": {
		Abbreviation: "vnet", AllowsHyphens: true, Pattern: "vnet-<purpose>-<region>-<###>",
		MinSegments: 2, MinLength: 2, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_subnet": {
		Abbreviation: "snet", AllowsHyphens: true, Pattern: "snet-<purpose>-<region>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_network_interface": {
		Abbreviation: "nic", AllowsHyphens: true, Pattern: "nic-<##>-<vm>-<purpose>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_public_ip": {
		Abbreviation: "pip", AllowsHyphens: true, Pattern: "pip-<name>-<env>-<region>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_public_ip_prefix": {
		Abbreviation: "ippre", AllowsHyphens: true, Pattern: "ippre-<name>-<env>-<region>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_network_security_group": {
		Abbreviation: "nsg", AllowsHyphens: true, Pattern: "nsg-<policy>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_network_security_rule": {
		Abbreviation: "nsgsr", AllowsHyphens: true, Pattern: "nsgsr-<policy>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_application_gateway": {
		Abbreviation: "agw", AllowsHyphens: true, Pattern: "agw-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_application_security_group": {
		Abbreviation: "asg", AllowsHyphens: true, Pattern: "asg-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_load_balancer": {
		Abbreviation: "lbi", AltAbbreviations: []string{"lbe"}, AllowsHyphens: true,
		Pattern: "lbi-<name>-<env>-<###> or lbe-<name>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_lb_rule": {
		Abbreviation: "rule", AllowsHyphens: true, Pattern: "rule-<name>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_local_network_gateway": {
		Abbreviation: "lgw", AllowsHyphens: true, Pattern: "lgw-<purpose>-<region>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_virtual_network_gateway": {
		Abbreviation: "vgw", AllowsHyphens: true, Pattern: "vgw-<purpose>-<region>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_vpn_gateway": {
		Abbreviation: "vpng", AllowsHyphens: true, Pattern: "vpng-<purpose>-<region>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_vpn_gateway_connection": {
		Abbreviation: "vcn", AllowsHyphens: true, Pattern: "vcn-<from>-to-<to>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_vpn_site": {
		Abbreviation: "vst", AllowsHyphens: true, Pattern: "vst-<purpose>-<region>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_express_route_circuit": {
		Abbreviation: "erc", AllowsHyphens: true, Pattern: "erc-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_firewall": {
		Abbreviation: "afw", AllowsHyphens: true, Pattern: "afw-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_firewall_policy": {
		Abbreviation: "afwp", AllowsHyphens: true, Pattern: "afwp-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_route_table": {
		Abbreviation: "rt", AllowsHyphens: true, Pattern: "rt-<name>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_route": {
		Abbreviation: "udr", AllowsHyphens: true, Pattern: "udr-<name>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_nat_gateway": {
		Abbreviation: "ng", AllowsHyphens: true, Pattern: "ng-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_traffic_manager_profile": {
		Abbreviation: "traf", AllowsHyphens: true, Pattern: "traf-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 63,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_private_endpoint": {
		Abbreviation: "pep", AllowsHyphens: true, Pattern: "pep-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_private_link_service": {
		Abbreviation: "pl", AllowsHyphens: true, Pattern: "pl-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_virtual_wan": {
		Abbreviation: "vwan", AllowsHyphens: true, Pattern: "vwan-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_virtual_hub": {
		Abbreviation: "vhub", AllowsHyphens: true, Pattern: "vhub-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_cdn_profile": {
		Abbreviation: "cdnp", AllowsHyphens: true, Pattern: "cdnp-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 260,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_cdn_endpoint": {
		Abbreviation: "cdne", AllowsHyphens: true, Pattern: "cdne-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_cdn_frontdoor_profile": {
		Abbreviation: "afd", AllowsHyphens: true, Pattern: "afd-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 260,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_cdn_frontdoor_endpoint": {
		Abbreviation: "fde", AllowsHyphens: true, Pattern: "fde-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_network_watcher": {
		Abbreviation: "nw", AllowsHyphens: true, Pattern: "nw-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_virtual_network_peering": {
		Abbreviation: "peer", AllowsHyphens: true, Pattern: "peer-<vnet1>-to-<vnet2>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_ip_group": {
		Abbreviation: "ipg", AllowsHyphens: true, Pattern: "ipg-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_network_manager": {
		Abbreviation: "vnm", AllowsHyphens: true, Pattern: "vnm-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_private_dns_resolver": {
		Abbreviation: "dnspr", AllowsHyphens: true, Pattern: "dnspr-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscore, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_route_filter": {
		Abbreviation: "rf", AllowsHyphens: true, Pattern: "rf-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},

	// Compute and web
	"azurerm_virtual_machine": {
		Abbreviation: "vm", AllowsHyphens: true, Pattern: "vm-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
		CantEndWithHyphen: true, CantEndWithPeriod: true,
	},
	"azurerm_linux_virtual_machine": {
		Abbreviation: "vm", AllowsHyphens: true, Pattern: "vm-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
		CantEndWithHyphen: true, CantEndWithPeriod: true,
	},
	"azurerm_windows_virtual_machine": {
		Abbreviation: "vm", AllowsHyphens: true, Pattern: "vm-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 15,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum,
		CantEndWithHyphen: true,
	},
	"azurerm_virtual_machine_scale_set": {
		Abbreviation: "vmss", AllowsHyphens: true, Pattern: "vmss-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
		CantEndWithHyphen: true, CantEndWithPeriod: true,
	},
	"azurerm_linux_virtual_machine_scale_set": {
		Abbreviation: "vmss", AllowsHyphens: true, Pattern: "vmss-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
		CantEndWithHyphen: true, CantEndWithPeriod: true,
	},
	"azurerm_windows_virtual_machine_scale_set": {
		Abbreviation: "vmss", AllowsHyphens: true, Pattern: "vmss-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 15,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
		CantEndWithHyphen: true, CantEndWithPeriod: true,
	},
	"azurerm_availability_set": {
		Abbreviation: "avail", AllowsHyphens: true, Pattern: "avail-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_managed_disk": {
		Abbreviation: "disk", AllowsHyphens: true, Pattern: "disk-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscore,
	},
	"azurerm_snapshot": {
		Abbreviation: "snap", AllowsHyphens: true, Pattern: "snap-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_service_plan": {
		Abbreviation: "asp", AllowsHyphens: true, Pattern: "asp-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 60,
		ValidChars: AlphanumHyphen,
	},
	"azurerm_linux_web_app": {
		Abbreviation: "app", AllowsHyphens: true, Pattern: "app-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 2, MaxLength: 60,
		ValidChars: AlphanumHyphen, CantEndWithHyphen: true,
	},
	"azurerm_windows_web_app": {
		Abbreviation: "app", AllowsHyphens: true, Pattern: "app-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 2, MaxLength: 60,
		ValidChars: AlphanumHyphen, CantEndWithHyphen: true,
	},
	"azurerm_linux_function_app": {
		Abbreviation: "func", AllowsHyphens: true, Pattern: "func-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 2, MaxLength: 60,
		ValidChars: AlphanumHyphen, CantEndWithHyphen: true,
	},
	"azurerm_windows_function_app": {
		Abbreviation: "func", AllowsHyphens: true, Pattern: "func-<workload>-<env>-<###>",
		MinSegments: 2, MinLength: 2, MaxLength: 60,
		ValidChars: AlphanumHyphen, CantEndWithHyphen: true,
	},
	"azurerm_static_web_app": {
		Abbreviation: "stapp", AllowsHyphens: true, Pattern: "stapp-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 60,
		ValidChars: AlphanumHyphen,
	},
	"azurerm_batch_account": {
		Abbreviation: "ba", AllowsHyphens: false, Pattern: "ba<workload><env>",
		MinLength: 3, MaxLength: 24, ValidChars: LowercaseNumber,
	},
	"azurerm_shared_image_gallery": {
		Abbreviation: "gal", AllowsHyphens: false, Pattern: "gal<workload><env>",
		MinLength: 1, MaxLength: 80, ValidChars: Alphanumeric,
		MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_disk_encryption_set": {
		Abbreviation: "des", AllowsHyphens: true, Pattern: "des-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscore,
	},
	"azurerm_proximity_placement_group": {
		Abbreviation: "ppg", AllowsHyphens: true, Pattern: "ppg-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},

	// Containers
	"azurerm_kubernetes_cluster": {
		Abbreviation: "aks", AllowsHyphens: true, Pattern: "aks-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 63,
		ValidChars: AlphanumHyphenUnderscore, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_container_registry": {
		Abbreviation: "cr", AllowsHyphens: false, Pattern: "cr<workload><env><###>",
		MinLength: 5, MaxLength: 50, ValidChars: Alphanumeric,
	},
	"azurerm_container_group": {
		Abbreviation: "ci", AllowsHyphens: true, Pattern: "ci-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
		NoConsecutiveHyphens: true,
	},
	"azurerm_container_app": {
		Abbreviation: "ca", AllowsHyphens: true, Pattern: "ca-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 32,
		ValidChars: LowercaseNumberHyphen, MustStartWith: LowercaseLetter, MustEndWith: Alphanum,
	},
	"azurerm_container_app_environment": {
		Abbreviation: "cae", AllowsHyphens: true, Pattern: "cae-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 32,
		ValidChars: LowercaseNumberHyphen, MustStartWith: LowercaseLetter, MustEndWith: Alphanum,
	},
	"azurerm_container_app_job": {
		Abbreviation: "caj", AllowsHyphens: true, Pattern: "caj-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 32,
		ValidChars: LowercaseNumberHyphen, MustStartWith: LowercaseLetter, MustEndWith: Alphanum,
	},
	"azurerm_service_fabric_cluster": {
		Abbreviation: "sf", AllowsHyphens: true, Pattern: "sf-<workload>-<env>",
		MinSegments: 2, MinLength: 4, MaxLength: 23,
		ValidChars: LowercaseNumberHyphen, MustStartWith: LowercaseLetter, MustEndWith: Alphanum,
	},
	"azurerm_service_fabric_managed_cluster": {
		Abbreviation: "sfmc", AllowsHyphens: true, Pattern: "sfmc-<workload>-<env>",
		MinSegments: 2, MinLength: 4, MaxLength: 23,
		ValidChars: LowercaseNumberHyphen, MustStartWith: LowercaseLetter, MustEndWith: Alphanum,
	},

	// Storage
	"azurerm_storage_account": {
		Abbreviation: "st", AllowsHyphens: false, Pattern: "st<workload><env>",
		MinLength: 3, MaxLength: 24, ValidChars: LowercaseNumber,
	},
	"azurerm_storage_share": {
		Abbreviation: "share", AllowsHyphens: true, Pattern: "share-<name>",
		MinSegments: 1, MinLength: 3, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
		NoConsecutiveHyphens: true,
	},
	"azurerm_data_protection_backup_vault": {
		Abbreviation: "bvault", AllowsHyphens: true, Pattern: "bvault-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Letter,
	},

	// Databases
	"azurerm_mssql_server": {
		Abbreviation: "sql", AllowsHyphens: true, Pattern: "sql-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
	},
	"azurerm_mssql_database": {
		Abbreviation: "sqldb", AllowsHyphens: true, Pattern: "sqldb-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 128,
		ValidChars: Broad, CantEndWithPeriod: true,
	},
	"azurerm_mssql_elasticpool": {
		Abbreviation: "sqlep", AllowsHyphens: true, Pattern: "sqlep-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 128,
		ValidChars: Broad, CantEndWithPeriod: true,
	},
	"azurerm_mssql_managed_instance": {
		Abbreviation: "sqlmi", AllowsHyphens: true, Pattern: "sqlmi-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
	},
	"azurerm_mysql_server": {
		Abbreviation: "mysql", AllowsHyphens: true, Pattern: "mysql-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
	},
	"azurerm_mysql_flexible_server": {
		Abbreviation: "mysql", AllowsHyphens: true, Pattern: "mysql-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
	},
	"azurerm_postgresql_server": {
		Abbreviation: "psql", AllowsHyphens: true, Pattern: "psql-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
	},
	"azurerm_postgresql_flexible_server": {
		Abbreviation: "psql", AllowsHyphens: true, Pattern: "psql-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: LowercaseNumberHyphen, CantEndWithHyphen: true,
	},
	"azurerm_cosmosdb_account": {
		Abbreviation: "cosmos", AllowsHyphens: true, Pattern: "cosmos-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 44,
		ValidChars: LowercaseNumberHyphen, MustStartWith: Alphanum,
	},
	"azurerm_cosmosdb_sql_database": {
		Abbreviation: "cosmos", AllowsHyphens: true, Pattern: "cosmos-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 44,
		ValidChars: LowercaseNumberHyphen, MustStartWith: Alphanum,
	},
	"azurerm_cosmosdb_mongo_database": {
		Abbreviation: "cosmon", AllowsHyphens: true, Pattern: "cosmon-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 44,
		ValidChars: LowercaseNumberHyphen, MustStartWith: Alphanum,
	},
	"azurerm_cosmosdb_cassandra_keyspace": {
		Abbreviation: "coscas", AllowsHyphens: true, Pattern: "coscas-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 44,
		ValidChars: LowercaseNumberHyphen, MustStartWith: Alphanum,
	},
	"azurerm_cosmosdb_gremlin_database": {
		Abbreviation: "cosgrm", AllowsHyphens: true, Pattern: "cosgrm-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 44,
		ValidChars: LowercaseNumberHyphen, MustStartWith: Alphanum,
	},
	"azurerm_cosmosdb_table": {
		Abbreviation: "costab", AllowsHyphens: true, Pattern: "costab-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 44,
		ValidChars: LowercaseNumberHyphen, MustStartWith: Alphanum,
	},
	"azurerm_redis_cache": {
		Abbreviation: "redis", AllowsHyphens: true, Pattern: "redis-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 63,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
		NoConsecutiveHyphens: true,
	},

	// Security
	"azurerm_key_vault": {
		Abbreviation: "kv", AllowsHyphens: true, Pattern: "kv-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 24,
		ValidChars: AlphanumHyphen, MustStartWith: Letter, MustEndWith: Alphanum,
		NoConsecutiveHyphens: true,
	},
	"azurerm_user_assigned_identity": {
		Abbreviation: "id", AllowsHyphens: true, Pattern: "id-<workload>-<env>-<region>-<###>",
		MinSegments: 2, MinLength: 3, MaxLength: 128,
		ValidChars: AlphanumHyphenUnderscore, MustStartWith: Alphanum,
	},
	"azurerm_bastion_host": {
		Abbreviation: "bas", AllowsHyphens: true, Pattern: "bas-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_ssh_public_key": {
		Abbreviation: "sshkey", AllowsHyphens: true, Pattern: "sshkey-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
	},

	// Management and monitoring
	"azurerm_log_analytics_workspace": {
		Abbreviation: "log", AllowsHyphens: true, Pattern: "log-<workload>-<env>",
		MinSegments: 2, MinLength: 4, MaxLength: 63,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_application_insights": {
		Abbreviation: "appi", AllowsHyphens: true, Pattern: "appi-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 260,
		ValidChars: Broad, CantEndWithPeriod: true,
	},
	"azurerm_monitor_action_group": {
		Abbreviation: "ag", AllowsHyphens: true, Pattern: "ag-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 260,
		ValidChars: Broad, CantEndWithPeriod: true,
	},
	"azurerm_monitor_data_collection_rule": {
		Abbreviation: "dcr", AllowsHyphens: true, Pattern: "dcr-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 260,
		ValidChars: Broad, CantEndWithPeriod: true,
	},
	"azurerm_purview_account": {
		Abbreviation: "pview", AllowsHyphens: true, Pattern: "pview-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 80,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
	},
	"azurerm_recovery_services_vault": {
		Abbreviation: "rsv", AllowsHyphens: true, Pattern: "rsv-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Letter,
	},

	// AI + Machine Learning
	"azurerm_search_service": {
		Abbreviation: "srch", AllowsHyphens: true, Pattern: "srch-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 60,
		ValidChars: LowercaseNumberHyphen, MustStartWith: LowercaseLetter,
		CantEndWithHyphen: true,
	},
	"azurerm_cognitive_account": {
		Abbreviation: "oai", AllowsHyphens: true, Pattern: "oai-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 64,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_machine_learning_workspace": {
		Abbreviation: "mlw", AllowsHyphens: true, Pattern: "mlw-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 33,
		ValidChars: AlphanumHyphenUnderscore,
	},
	"azurerm_bot_service_azure_bot": {
		Abbreviation: "bot", AllowsHyphens: true, Pattern: "bot-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
	},

	// Analytics and IoT
	"azurerm_data_factory": {
		Abbreviation: "adf", AllowsHyphens: true, Pattern: "adf-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_databricks_workspace": {
		Abbreviation: "dbw", AllowsHyphens: true, Pattern: "dbw-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscore,
	},
	"azurerm_data_lake_store": {
		Abbreviation: "dls", AllowsHyphens: false, Pattern: "dls<workload><env>",
		MinLength: 3, MaxLength: 24, ValidChars: LowercaseNumber,
	},
	"azurerm_kusto_cluster": {
		Abbreviation: "dec", AllowsHyphens: false, Pattern: "dec<workload><env>",
		MinLength: 4, MaxLength: 22, ValidChars: LowercaseNumber, MustStartWith: LowercaseLetter,
	},
	"azurerm_stream_analytics_job": {
		Abbreviation: "asa", AllowsHyphens: true, Pattern: "asa-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: AlphanumHyphenUnderscore,
	},
	"azurerm_synapse_workspace": {
		Abbreviation: "synw", AllowsHyphens: true, Pattern: "synw-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 50,
		ValidChars: LowercaseNumberHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_synapse_sql_pool": {
		Abbreviation: "syndp", AllowsHyphens: true, Pattern: "syndp-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 60,
		ValidChars: Broad, CantEndWithPeriod: true,
	},
	"azurerm_synapse_spark_pool": {
		Abbreviation: "synsp", AllowsHyphens: false, Pattern: "synsp<workload><env>",
		MinLength: 1, MaxLength: 15, ValidChars: Alphanumeric, MustStartWith: Letter, MustEndWith: Alphanum,
	},
	"azurerm_eventhub_namespace": {
		Abbreviation: "evhns", AllowsHyphens: true, Pattern: "evhns-<workload>-<env>",
		MinSegments: 2, MinLength: 6, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Letter, MustEndWith: Alphanum,
	},
	"azurerm_eventhub": {
		Abbreviation: "evh", AllowsHyphens: true, Pattern: "evh-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 256,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_eventgrid_domain": {
		Abbreviation: "evgd", AllowsHyphens: true, Pattern: "evgd-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 50,
		ValidChars: AlphanumHyphen,
	},
	"azurerm_eventgrid_topic": {
		Abbreviation: "evgt", AllowsHyphens: true, Pattern: "evgt-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 50,
		ValidChars: AlphanumHyphen,
	},
	"azurerm_eventgrid_system_topic": {
		Abbreviation: "egst", AllowsHyphens: true, Pattern: "egst-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 50,
		ValidChars: AlphanumHyphen,
	},
	"azurerm_iothub": {
		Abbreviation: "iot", AllowsHyphens: true, Pattern: "iot-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 50,
		ValidChars: AlphanumHyphen, CantEndWithHyphen: true,
	},
	"azurerm_hdinsight_hadoop_cluster": {
		Abbreviation: "hadoop", AllowsHyphens: true, Pattern: "hadoop-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 59,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_hdinsight_hbase_cluster": {
		Abbreviation: "hbase", AllowsHyphens: true, Pattern: "hbase-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 59,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_hdinsight_kafka_cluster": {
		Abbreviation: "kafka", AllowsHyphens: true, Pattern: "kafka-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 59,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_hdinsight_spark_cluster": {
		Abbreviation: "spark", AllowsHyphens: true, Pattern: "spark-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 59,
		ValidChars: AlphanumHyphen, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_analysis_services_server": {
		Abbreviation: "as", AllowsHyphens: false, Pattern: "as<workload><env>",
		MinLength: 3, MaxLength: 63, ValidChars: LowercaseNumber, MustStartWith: LowercaseLetter,
	},
	"azurerm_powerbi_embedded": {
		Abbreviation: "pbi", AllowsHyphens: false, Pattern: "pbi<workload><env>",
		MinLength: 3, MaxLength: 63, ValidChars: LowercaseNumber, MustStartWith: LowercaseLetter,
	},

	// Integration
	"azurerm_api_management": {
		Abbreviation: "apim", AllowsHyphens: true, Pattern: "apim-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Letter, MustEndWith: Alphanum,
	},
	"azurerm_logic_app_workflow": {
		Abbreviation: "logic", AllowsHyphens: true, Pattern: "logic-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 43,
		ValidChars: AlphanumHyphenUnderscorePeriod,
	},
	"azurerm_servicebus_namespace": {
		Abbreviation: "sbns", AllowsHyphens: true, Pattern: "sbns-<workload>-<env>",
		MinSegments: 2, MinLength: 6, MaxLength: 50,
		ValidChars: AlphanumHyphen, MustStartWith: Letter, MustEndWith: Alphanum,
	},
	"azurerm_servicebus_queue": {
		Abbreviation: "sbq", AllowsHyphens: true, Pattern: "sbq-<workload>",
		MinSegments: 2, MinLength: 1, MaxLength: 260,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_servicebus_topic": {
		Abbreviation: "sbt", AllowsHyphens: true, Pattern: "sbt-<workload>",
		MinSegments: 2, MinLength: 1, MaxLength: 260,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},
	"azurerm_servicebus_subscription": {
		Abbreviation: "sbts", AllowsHyphens: true, Pattern: "sbts-<workload>",
		MinSegments: 2, MinLength: 1, MaxLength: 50,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: Alphanum,
	},

	// Developer tools
	"azurerm_app_configuration": {
		Abbreviation: "appcs", AllowsHyphens: true, Pattern: "appcs-<workload>-<env>",
		MinSegments: 2, MinLength: 5, MaxLength: 50,
		ValidChars: AlphanumHyphen, CantEndWithHyphen: true,
		NoConsecutiveHyphens: true,
	},
	"azurerm_signalr_service": {
		Abbreviation: "sigr", AllowsHyphens: true, Pattern: "sigr-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: AlphanumHyphen, MustStartWith: Letter, MustEndWith: Alphanum,
	},
	"azurerm_web_pubsub": {
		Abbreviation: "wps", AllowsHyphens: true, Pattern: "wps-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 63,
		ValidChars: AlphanumHyphen, MustStartWith: Letter, MustEndWith: Alphanum,
	},
	"azurerm_maps_account": {
		Abbreviation: "map", AllowsHyphens: true, Pattern: "map-<workload>-<env>",
		MinSegments: 2, MinLength: 1, MaxLength: 98,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
	},

	// Virtual Desktop Infrastructure
	"azurerm_virtual_desktop_host_pool": {
		Abbreviation: "vdpool", AllowsHyphens: true, Pattern: "vdpool-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_virtual_desktop_application_group": {
		Abbreviation: "vdag", AllowsHyphens: true, Pattern: "vdag-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_virtual_desktop_workspace": {
		Abbreviation: "vdws", AllowsHyphens: true, Pattern: "vdws-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},
	"azurerm_virtual_desktop_scaling_plan": {
		Abbreviation: "vdscaling", AllowsHyphens: true, Pattern: "vdscaling-<workload>-<env>",
		MinSegments: 2, MinLength: 3, MaxLength: 64,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum, MustEndWith: AlphanumUnderscore,
	},

	// Migration
	"azurerm_database_migration_service": {
		Abbreviation: "dms", AllowsHyphens: true, Pattern: "dms-<workload>-<env>",
		MinSegments: 2, MinLength: 2, MaxLength: 62,
		ValidChars: AlphanumHyphenUnderscorePeriod, MustStartWith: Alphanum,
	},
}

// LookupRule returns the CAF rule for the given azurerm resource type.
func LookupRule(resourceType string) (CAFRule, bool) {
	rule, ok := CAFRules[resourceType]
	return rule, ok
}

// ValidateResourceName checks a resource name against all CAF naming rules.
// Returns a list of issues found (empty means the name is valid).
func ValidateResourceName(resourceType, name string) []string {
	rule, ok := CAFRules[resourceType]
	if !ok {
		return nil
	}

	var issues []string

	// 1. Prefix check
	if rule.Abbreviation != "" {
		prefixOk := false
		allPrefixes := []string{rule.Abbreviation}
		allPrefixes = append(allPrefixes, rule.AltAbbreviations...)

		for _, prefix := range allPrefixes {
			if rule.AllowsHyphens {
				if strings.HasPrefix(name, prefix+"-") {
					prefixOk = true
					break
				}
			} else {
				if strings.HasPrefix(name, prefix) {
					prefixOk = true
					break
				}
			}
		}

		if !prefixOk {
			if rule.AllowsHyphens {
				issues = append(issues, fmt.Sprintf("name should start with prefix \"%s-\"", rule.Abbreviation))
			} else {
				issues = append(issues, fmt.Sprintf("name should start with prefix \"%s\"", rule.Abbreviation))
			}
		}
	}

	// 2. Structure check (hyphen segments)
	if rule.AllowsHyphens && rule.MinSegments > 0 {
		segments := strings.Split(name, "-")
		if len(segments) < rule.MinSegments {
			issues = append(issues, fmt.Sprintf("expected at least %d hyphen-separated segments, got %d", rule.MinSegments, len(segments)))
		}
	}
	if !rule.AllowsHyphens && strings.Contains(name, "-") {
		issues = append(issues, "contains hyphens (not allowed for this resource type)")
	}

	// 3. Length check
	if rule.MaxLength > 0 {
		if len(name) < rule.MinLength {
			issues = append(issues, fmt.Sprintf("length %d is below minimum %d", len(name), rule.MinLength))
		}
		if len(name) > rule.MaxLength {
			issues = append(issues, fmt.Sprintf("length %d exceeds maximum %d", len(name), rule.MaxLength))
		}
	}

	// 4. Character check
	charIssue := validateChars(name, rule.ValidChars)
	if charIssue != "" {
		issues = append(issues, charIssue)
	}

	// 5. Start/end constraints
	if len(name) > 0 {
		if issue := checkStartChar(name, rule.MustStartWith); issue != "" {
			issues = append(issues, issue)
		}
		if issue := checkEndChar(name, rule.MustEndWith); issue != "" {
			issues = append(issues, issue)
		}
		if rule.CantEndWithHyphen && strings.HasSuffix(name, "-") {
			issues = append(issues, "name can't end with a hyphen")
		}
		if rule.CantEndWithPeriod && strings.HasSuffix(name, ".") {
			issues = append(issues, "name can't end with a period")
		}
	}

	// 6. Consecutive hyphens
	if rule.NoConsecutiveHyphens && strings.Contains(name, "--") {
		issues = append(issues, "contains consecutive hyphens (not allowed)")
	}

	return issues
}

// GetExpectedPrefix returns the abbreviation for display purposes.
func GetExpectedPrefix(resourceType string) string {
	rule, ok := CAFRules[resourceType]
	if !ok {
		return ""
	}
	return rule.Abbreviation
}

// GetExpectedPattern returns the pattern for display purposes.
func GetExpectedPattern(resourceType string) string {
	rule, ok := CAFRules[resourceType]
	if !ok {
		return ""
	}
	return rule.Pattern
}

func validateChars(name string, class CharClass) string {
	for _, r := range name {
		if !isValidChar(r, class) {
			return fmt.Sprintf("contains invalid character '%c' (%s)", r, charClassDescription(class))
		}
	}
	return ""
}

func isValidChar(r rune, class CharClass) bool {
	switch class {
	case AlphanumHyphen:
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-'
	case AlphanumHyphenUnderscore:
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
	case AlphanumHyphenUnderscorePeriod:
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.'
	case LowercaseNumber:
		return unicode.IsLower(r) || unicode.IsDigit(r)
	case LowercaseNumberHyphen:
		return unicode.IsLower(r) || unicode.IsDigit(r) || r == '-'
	case Alphanumeric:
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	case Broad:
		return true
	}
	return true
}

func charClassDescription(class CharClass) string {
	switch class {
	case AlphanumHyphen:
		return "only alphanumerics and hyphens allowed"
	case AlphanumHyphenUnderscore:
		return "only alphanumerics, hyphens, and underscores allowed"
	case AlphanumHyphenUnderscorePeriod:
		return "only alphanumerics, hyphens, underscores, and periods allowed"
	case LowercaseNumber:
		return "only lowercase letters and numbers allowed"
	case LowercaseNumberHyphen:
		return "only lowercase letters, numbers, and hyphens allowed"
	case Alphanumeric:
		return "only alphanumerics allowed"
	case Broad:
		return "most characters allowed"
	}
	return ""
}

func checkStartChar(name string, rule StartEndRule) string {
	if len(name) == 0 {
		return ""
	}
	r := rune(name[0])
	switch rule {
	case Letter:
		if !unicode.IsLetter(r) {
			return "name must start with a letter"
		}
	case Alphanum:
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "name must start with an alphanumeric character"
		}
	case LowercaseLetter:
		if !unicode.IsLower(r) {
			return "name must start with a lowercase letter"
		}
	}
	return ""
}

func checkEndChar(name string, rule StartEndRule) string {
	if len(name) == 0 {
		return ""
	}
	r := rune(name[len(name)-1])
	switch rule {
	case Alphanum:
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return "name must end with an alphanumeric character"
		}
	case AlphanumUnderscore:
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return "name must end with an alphanumeric character or underscore"
		}
	case Letter:
		if !unicode.IsLetter(r) {
			return "name must end with a letter"
		}
	case LowercaseLetter:
		if !unicode.IsLower(r) {
			return "name must end with a lowercase letter"
		}
	}
	return ""
}
