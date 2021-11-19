package rustack_terraform

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
)

func (args *Arguments) injectCreateFirewallRule() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the firewall rule",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the firewall rule",
		},
		"destination_ip": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "destination ip address",
		},
		"port_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "max range of port",
		},
		"protocol": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "protocol tcp/upd/icmp",
		},
	})
}

func ruleToMap(firewallRules rustack.FirewallRule) (rule map[string]interface{}) {

	rule = map[string]interface{}{
		"id":             firewallRules.ID,
		"name":           firewallRules.Name,
		"destination_ip": firewallRules.DestinationIp,
		"protocol":       firewallRules.Protocol,
	}
	rule["port_range"] = nil
	if firewallRules.DstPortRangeMin != nil {
		rule["port_range"] = fmt.Sprintf("%d", *firewallRules.DstPortRangeMin)
	}
	if firewallRules.DstPortRangeMax != nil {
		rule["port_range"] = fmt.Sprintf("%s:%d", rule["port_range"], *firewallRules.DstPortRangeMax)
	}
	return
}

func rulesToMap(firewallRules []*rustack.FirewallRule) (rules map[string][]map[string]interface{}) {
	rules = make(map[string][]map[string]interface{}, 2)

	for i := 0; i < len(firewallRules); i++ {
		rule := ruleToMap(*firewallRules[i])
		rules[firewallRules[i].Direction] = append(rules[firewallRules[i].Direction], rule)
	}

	for _, ruleType := range []string{"ingress", "egress"} {
		if len(rules[ruleType]) == 0 {
			rules[ruleType] = nil
		}
	}

	return
}
