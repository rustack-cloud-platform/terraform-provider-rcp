package rustack_terraform

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pilat/rustack-go/rustack"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRustackFirewallTemplate() *schema.Resource {
	args := Defaults()
	args.injectCreateFirewallTemplate()
	args.injectContextVdcById()

	return &schema.Resource{
		CreateContext: resourceRustackFirewallTemplateCreate,
		ReadContext:   resourceRustackFirewallTemplateRead,
		UpdateContext: resourceRustackFirewallTemplateUpdate,
		DeleteContext: resourceRustackFirewallTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: args,
	}
}

func resourceRustackFirewallTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}

	newFirewallTemplate := rustack.NewFirewallTemplate(d.Get("name").(string))
	err = targetVdc.CreateFirewallTemplate(&newFirewallTemplate)
	if err != nil {
		return diag.Errorf("Error creating Firewall Template: %s", err)
	}

	d.SetId(newFirewallTemplate.ID)
	log.Printf("[INFO] FirewallTemplate created, ID: %s", d.Id())

	for _, ruleType := range []string{"ingress", "egress"} {
		rulesCount := d.Get(fmt.Sprintf("%s_rule.#", ruleType)).(int)
		rules := make([]map[string]interface{}, rulesCount)
		for i := 0; i < rulesCount; i++ {
			rulePrefix := fmt.Sprintf("%s_rule.%d", ruleType, i)

			var newFirewallRule rustack.FirewallRule
			newFirewallRule.Name = rulePrefix
			newFirewallRule.Direction = ruleType
			setUpRule(&newFirewallRule, d)

			if err = newFirewallTemplate.Update(&newFirewallRule); err != nil {
				return diag.Errorf("Error creating Firewall rule: %s", err)
			}

			rules[i] = ruleToMap(newFirewallRule)

			log.Printf("Update firewall rule '%s' name\n", newFirewallRule.Name)
		}

		log.Printf("F rules '%s' \n", rules)

		if err = d.Set(fmt.Sprintf("%s_rule", ruleType), rules); err != nil {
			return diag.Errorf("Error setting %s: %s", ruleType, err)
		}
	}

	d.SetId(newFirewallTemplate.ID)
	log.Printf("[INFO] Firewall Template Updated, ID: %s", d.Id())

	return resourceRustackFirewallTemplateRead(ctx, d, meta)
}

func resourceRustackFirewallTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	firewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		return diag.Errorf("Error getting Firewall Template: %s", err)
	}
	firewallRules, err := manager.GetFirewallRules(d.Id())
	if err != nil {
		return diag.Errorf("Error getting Firewall Rule: %s", err)
	}
	rules := rulesToMap(firewallRules)

	d.SetId(firewallTemplate.ID)
	d.Set("name", firewallTemplate.Name)
	d.Set("ingress_rule", rules["ingress"])
	d.Set("egress_rule", rules["egress"])

	return nil
}

func resourceRustackFirewallTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()

	firewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		return diag.Errorf("Error getting FirewallTemplate: %s", err)
	}

	if d.HasChange("name") {
		if err = firewallTemplate.Rename(d.Get("name").(string)); err != nil {
			return diag.Errorf("Error rename Firewall Template: %s", err)
		}
	}

	rules := make(map[string]interface{}, 2)
	for _, ruleType := range []string{"ingress", "egress"} {
		a, b := d.GetChange(fmt.Sprintf("%s_rule", ruleType))

		rules[ruleType] = d.Get(fmt.Sprintf("%s_rule", ruleType))

		toUpdate, toDelete := stateDifference(a.([]interface{}), b.([]interface{}))

		if err = deleteRules(*firewallTemplate, toDelete); err != nil {
			return diag.Errorf("Error delete rules: %s", err)
		}
		if err = updateRules(d, *firewallTemplate, toUpdate, ruleType, rules[ruleType]); err != nil {
			return diag.Errorf("Error update rules: %s", err)
		}
		if err = createRules(d, *firewallTemplate, ruleType, rules[ruleType]); err != nil {
			return diag.Errorf("Error create rules: %s", err)
		}
	}
	d.Set("ingress_rule", rules["ingress"])
	d.Set("egress_rule", rules["egress"])

	return resourceRustackFirewallTemplateRead(ctx, d, meta)
}

func deleteRules(firewallTemplate rustack.FirewallTemplate, rules []interface{}) (err error) {
	for i := 0; i < len(rules); i++ {
		rule, err := firewallTemplate.GetRuleById(rules[i].(map[string]interface{})["id"].(string))
		if err != nil {
			return err
		}
		rule.Delete()
	}
	return
}

func updateRules(d *schema.ResourceData, firewallTemplate rustack.FirewallTemplate,
	toUpdate []interface{}, ruleType string, res interface{}) (err error) {
	// Update
	for i := 0; i < len(toUpdate); i++ {
		ruleId := toUpdate[i].(map[string]interface{})["id"].(string)

		if ruleId == "" {
			continue
		}
		firewallRule, err := firewallTemplate.GetRuleById(ruleId)
		if err != nil {
			return err
		}
		firewallRule.Name = toUpdate[i].(map[string]interface{})["name"].(string)
		firewallRule.Protocol = toUpdate[i].(map[string]interface{})["protocol"].(string)
		firewallRule.DestinationIp = toUpdate[i].(map[string]interface{})["destination_ip"].(string)
		portRange := toUpdate[i].(map[string]interface{})["port_range"].(string)
		var portRangeMax, portRangeMin int
		if portRange == "" {
			firewallRule.DstPortRangeMin = nil
			firewallRule.DstPortRangeMax = nil
		} else if strings.Contains(portRange, ":") {
			portRangeList := strings.Split(portRange, ":")
			portRangeMin, err = strconv.Atoi(portRangeList[0])
			if err != nil {
				return err
			}
			portRangeMax, err = strconv.Atoi(portRangeList[1])
			if err != nil {
				return err
			}
			firewallRule.DstPortRangeMin = &portRangeMin
			firewallRule.DstPortRangeMax = &portRangeMax
		} else {
			portRangeMax, err = strconv.Atoi(portRange)
			if err != nil {
				return err
			}
			firewallRule.DstPortRangeMin = nil
			firewallRule.DstPortRangeMax = &portRangeMax
		}

		if err = firewallRule.Update(); err != nil {
			return err
		}
	}
	return
}

func createRules(d *schema.ResourceData, firewallTemplate rustack.FirewallTemplate,
	ruleType string, res interface{}) (err error) {

	// Create new rule
	for i, item := range res.([]interface{}) {
		ruleId := item.(map[string]interface{})["id"].(string)
		if ruleId != "" {
			continue
		}
		rulePrefix := fmt.Sprintf("%s_rule.%d", ruleType, i)
		var newFirewallRule rustack.FirewallRule
		newFirewallRule.Name = rulePrefix
		newFirewallRule.Direction = ruleType
		setUpRule(&newFirewallRule, d)

		if err = firewallTemplate.CreateFirewallRule(&newFirewallRule); err != nil {
			return err
		}

		item.(map[string]interface{})["id"] = newFirewallRule.ID
		item.(map[string]interface{})["name"] = rulePrefix
	}
	return
}

func resourceRustackFirewallTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).rustackManager()
	FirewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		return diag.Errorf("Error getting FirewallTemplate: %s", err)
	}

	err = FirewallTemplate.Delete()
	if err != nil {
		return diag.Errorf("Error deleting FirewallTemplate: %s", err)
	}

	return nil
}

func setUpRule(rule *rustack.FirewallRule, d *schema.ResourceData) (err error) {
	rule.Protocol = d.Get(fmt.Sprintf("%s.protocol", rule.Name)).(string)
	rule.DestinationIp = d.Get(fmt.Sprintf("%s.destination_ip", rule.Name)).(string)
	if rule.Protocol == "icmp" {
		rule.DstPortRangeMax = nil
		rule.DstPortRangeMin = nil
	} else {
		// Two ways to set up port range
		// 1:40 - port range from min to max
		// 50 - single port
		portRange := d.Get(fmt.Sprintf("%s.port_range", rule.Name)).(string)
		var a, b int
		_, err := fmt.Sscanf(portRange, "%d:%d", &a, &b)
		if err != nil {
			_, err = fmt.Sscanf(portRange, "%d", &a)
			rule.DstPortRangeMax = nil
		} else {
			rule.DstPortRangeMax = &b
		}
		rule.DstPortRangeMin = &a
	}

	return
}

func stateDifference(slice1 []interface{}, slice2 []interface{}) (toUpdate []interface{}, toDelete []interface{}) {
	for _, s1 := range slice1 {
		found := false
		if s1.(map[string]interface{})["id"].(string) == "" {
			continue
		}
		for _, s2 := range slice2 {
			if reflect.DeepEqual(
				s1.(map[string]interface{})["id"].(string),
				s2.(map[string]interface{})["id"].(string),
			) {
				found = true
				if !reflect.DeepEqual(s1, s2) {
					toUpdate = append(toUpdate, s2)
				}
				break
			}
		}
		if !found {
			toDelete = append(toDelete, s1)
		}
	}

	return
}
