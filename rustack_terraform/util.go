package rustack_terraform

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pilat/rustack-go/rustack"
	"github.com/pkg/errors"
)

type Arguments map[string]*schema.Schema

func Defaults() Arguments {
	return make(Arguments)
}

func (args Arguments) merge(extraArg Arguments) {
	for key, val := range extraArg {
		args[key] = val
	}
}

func GetFirewallTemplateByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.FirewallTemplate, error) {
	firewallTemplateName := d.Get("name").(string)
	firewallTemplates, err := vdc.GetFirewallTemplates()

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of firewall templates")
	}

	for _, firewallTemplate := range firewallTemplates {
		if strings.EqualFold(firewallTemplate.Name, firewallTemplateName) {
			return firewallTemplate, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Firewall template with name '%s' not found", firewallTemplateName)
}

func GetKubernetesTemplateByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.KubernetesTemplate, error) {
	templateName := d.Get("name").(string)
	templates, err := vdc.GetKubernetesTemplates()

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of kubernetes templates")
	}

	for _, template := range templates {
		if strings.EqualFold(template.Name, templateName) {
			return template, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Kubernetes template with name '%s' not found", templateName)

}

func GetKubernetesTemplateById(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.KubernetesTemplate, error) {
	templateId := d.Get("template_id").(string)
	template, err := manager.GetKubernetesTemplate(templateId)

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of kubernetes templates")
	}

	return template, nil

}

func GetKubernetesByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Kubernetes, error) {
	kubernetes_name := d.Get("name").(string)
	kubernetes, err := vdc.GetKubernetes()

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of kubernetes templates")
	}

	for _, k8s := range kubernetes {
		if strings.EqualFold(k8s.Name, kubernetes_name) {
			return k8s, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Kubernetes template with name '%s' not found", kubernetes_name)

}

func GetTemplateByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Template, error) {
	templateName := d.Get("name").(string)
	templates, err := vdc.GetTemplates()

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of templates")
	}

	for _, template := range templates {
		if strings.EqualFold(template.Name, templateName) {
			return template, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Template with name '%s' not found", templateName)
}

func GetPlatformByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Platform, error) {
	templateName := d.Get("name").(string)

	platforms, err := manager.GetPlatforms(vdc.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of platfroms")
	}

	for _, platform := range platforms {
		if strings.EqualFold(platform.Name, templateName) {
			return platform, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Template with name '%s' not found", templateName)
}

func GetPubKeyByName(d *schema.ResourceData, manager *rustack.Manager) (*rustack.PubKey, error) {
	key_name := d.Get("name").(string)
	account := d.Get("account_id").(string)
	pub_keys, err := manager.GetPublicKeys(account)

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of public keys")
	}

	for _, pub_key := range pub_keys {
		if strings.EqualFold(pub_key.Name, key_name) {
			return pub_key, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Public key with name '%s' not found", key_name)

}

func GetTemplateById(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Template, error) {
	templateId := d.Get("template_id")
	template, err := manager.GetTemplate(templateId.(string))
	if err != nil {
		return nil, errors.Wrapf(err, "Template with id '%s' not found", templateId)
	}

	return template, nil

}

func GetDiskByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Disk, error) {
	diskName := d.Get("name").(string)
	disks, err := vdc.GetDisks()

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of disks")
	}

	for _, disk := range disks {
		if strings.EqualFold(disk.Name, diskName) {
			return disk, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Disk with name '%s' not found", diskName)

}

func GetDiskById(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Disk, error) {
	diskId := d.Get("id")
	disk, err := manager.GetDisk(diskId.(string))
	if err != nil {
		return nil, errors.Wrapf(err, "Disk with id '%s' not found", diskId)
	}

	return disk, nil
}

func GetStorageProfileById(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc, prefix *string) (*rustack.StorageProfile, error) {
	storageProfiles, err := vdc.GetStorageProfiles()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting list of storage profiles")
	}

	var targetStorageProfile *rustack.StorageProfile

	storageProfileId := d.Get(MakePrefix(prefix, "storage_profile_id")).(string)
	for _, storageProfile := range storageProfiles {
		if storageProfile.ID == storageProfileId {
			targetStorageProfile = storageProfile
			break
		}
	}

	if targetStorageProfile == nil {
		return nil, fmt.Errorf("ERROR. Storage profile with id '%s' not found", storageProfileId)
	}

	return targetStorageProfile, nil

}

func GetNetworkById(d *schema.ResourceData, manager *rustack.Manager, prefix *string) (*rustack.Network, error) {
	networkPrefix := MakePrefix(prefix, "network_id")
	networkId := d.Get(networkPrefix).(string)
	targetNetwork, err := manager.GetNetwork(networkId)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting network")
	}

	return targetNetwork, nil
}

func GetNetworkByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Network, error) {
	networks, err := manager.GetNetworks()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting list of networks")
	}

	var targetNetwork *rustack.Network

	networkName := d.Get("name")
	for _, network := range networks {
		if network.Name == networkName.(string) && network.Vdc.Id == vdc.ID {
			targetNetwork = network
			break
		}
	}

	if targetNetwork == nil {
		return nil, fmt.Errorf("ERROR. Network with name '%s' not found", networkName)
	}

	return targetNetwork, nil
}

func GetStorageProfileByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.StorageProfile, error) {
	storageProfiles, err := vdc.GetStorageProfiles()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting list of storage profiles")
	}

	var targetStorageProfile *rustack.StorageProfile

	storageProfileName := d.Get("name")
	for _, storageProfile := range storageProfiles {
		if strings.EqualFold(storageProfile.Name, storageProfileName.(string)) {
			targetStorageProfile = storageProfile
			break
		}
	}

	if targetStorageProfile == nil {
		return nil, fmt.Errorf("ERROR: Storage profile with name '%s' not found in vdc '%s'", storageProfileName, vdc.Name)
	}

	return targetStorageProfile, nil
}

func GetHypervisorByName(d *schema.ResourceData, manager *rustack.Manager, project *rustack.Project) (*rustack.Hypervisor, error) {
	hypervisors, err := project.GetAvailableHypervisors()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting list of hypervisors")
	}

	var targetHypervisor *rustack.Hypervisor

	hypervisorName := strings.ToLower(d.Get("name").(string))
	for _, hypervisor := range hypervisors {
		if strings.ToLower(hypervisor.Name) == hypervisorName {
			targetHypervisor = hypervisor
			break
		}
	}

	if targetHypervisor == nil {
		return nil, fmt.Errorf("ERROR: Hypervisor with name '%s' not found", hypervisorName)
	}

	return targetHypervisor, nil
}

func GetHypervisorById(d *schema.ResourceData, manager *rustack.Manager, project *rustack.Project) (*rustack.Hypervisor, error) {
	hypervisors, err := project.GetAvailableHypervisors()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting list of hypervisors")
	}

	var targetHypervisor *rustack.Hypervisor

	hypervisorId := d.Get("hypervisor_id")
	for _, hypervisor := range hypervisors {
		if hypervisor.ID == hypervisorId.(string) {
			targetHypervisor = hypervisor
			break
		}
	}

	if targetHypervisor == nil {
		return nil, fmt.Errorf("ERROR: Hypervisor with id '%s' not found", hypervisorId)
	}

	return targetHypervisor, nil
}

func GetHypervisorByIdRead(d *schema.ResourceData, manager *rustack.Manager, project *rustack.Project) (*rustack.Hypervisor, error) {
	hypervisors, err := project.GetAvailableHypervisors()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting list of hypervisors")
	}

	var targetHypervisor *rustack.Hypervisor

	hypervisorId := d.Get("id")
	for _, hypervisor := range hypervisors {
		if hypervisor.ID == hypervisorId.(string) {
			targetHypervisor = hypervisor
			break
		}
	}

	if targetHypervisor == nil {
		return nil, fmt.Errorf("ERROR: Hypervisor with id '%s' not found", hypervisorId)
	}

	return targetHypervisor, nil
}

func GetHypervisorByIdK8s(d *schema.ResourceData, manager *rustack.Manager, project *rustack.Project) (*rustack.Hypervisor, error) {
	hypervisors, err := project.GetAvailableHypervisors()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting list of hypervisors")
	}

	var targetHypervisor *rustack.Hypervisor

	hypervisorId := d.Get("platform")
	for _, hypervisor := range hypervisors {
		if hypervisor.ID == hypervisorId.(string) {
			targetHypervisor = hypervisor
			break
		}
	}

	if targetHypervisor == nil {
		return nil, fmt.Errorf("ERROR: Hypervisor with id '%s' not found", hypervisorId)
	}

	return targetHypervisor, nil
}

func GetProjectByName(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Project, error) {
	projectName := d.Get("name")
	projects, err := manager.GetProjects(rustack.Arguments{"name": projectName.(string)})

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of projects")
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Project with name '%s' not found", projectName)
}

func GetProjectById(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Project, error) {
	projectId := d.Get("project_id")
	project, err := manager.GetProject(projectId.(string))
	if err != nil {
		return nil, errors.Wrapf(err, "Project with id '%s' not found", projectId)
	}

	return project, nil
}

func GetVdcByName(d *schema.ResourceData, manager *rustack.Manager, project *rustack.Project) (*rustack.Vdc, error) {
	vdcName := d.Get("name")
	vdcs, err := manager.GetVdcs(rustack.Arguments{"name": vdcName.(string)})

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of vdcs")
	}

	for _, vdc := range vdcs {
		if vdc.Name == vdcName && (project == nil || vdc.Project.ID == project.ID) {
			return vdc, nil
		}
	}

	return nil, fmt.Errorf("VDC with name '%s' not found in project '%s'", vdcName, project.Name)

}

func GetVdcById(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Vdc, error) {
	vdcId := d.Get("vdc_id")
	vdc, err := manager.GetVdc(vdcId.(string))
	if err != nil {
		return nil, errors.Wrapf(err, "VDC with id '%s' not found", vdcId)
	}

	return vdc, nil
}

func GetVmByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Vm, error) {
	vmName := d.Get("name")
	vms, err := vdc.GetVms(rustack.Arguments{"name": vmName.(string)})

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of vms")
	}

	for _, vm := range vms {
		if vm.Name == vmName {
			return vm, nil
		}
	}

	return nil, fmt.Errorf("VM with name '%s' not found in vdc '%s'", vmName, vdc.Name)
}

func GetRouterByName(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Router, error) {
	routerName := d.Get("name").(string)
	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return nil, err
	}
	routers, err := vdc.GetRouters(rustack.Arguments{"name": routerName})

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of routers")
	}

	for _, router := range routers {
		if router.Name == routerName {
			return router, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Router with name '%s' not found in vdc '%s'", routerName, vdc.Name)
}

func MakePrefix(prefix *string, name string) string {
	if prefix == nil {
		return name
	}
	if name == "" {
		return *prefix
	}

	return fmt.Sprintf("%s.%s", *prefix, name)
}

func setResourceDataFromMap(d *schema.ResourceData, m map[string]interface{}) error {
	for key, value := range m {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("ERROR: Unable to set `%s` attribute: %s", key, err)
		}
	}
	return nil
}

func getRouterByNetwork(manager rustack.Manager, network rustack.Network) (router *rustack.Router, err error) {
	routers, err := manager.GetRouters()
	if err != nil {
		return nil, err
	}
	for _, router := range routers {
		for _, port := range router.Ports {
			if port.Network.ID == network.ID {
				return router, nil
			}
		}
	}
	return
}

func repeatOnError(f func() error, targerInterface interface{ WaitLock() error }) (err error) {
	for j := 0; j < 15; j++ {
		targerInterface.WaitLock()
		err = f()
		if err == nil {
			return
		}
		fmt.Printf("err: %v\n", err)
		time.Sleep(time.Second)
	}
	return
}

func GetServiseNetworkByVdc(vdc *rustack.Vdc) (*rustack.Network, error) {
	allNetworks, err := vdc.GetNetworks()
	if err != nil {
		return nil, err
	}
	for _, network := range allNetworks {
		if network.IsDefault {
			return network, nil
		}
	}
	err = errors.New("Cant find service network")
	return nil, err
}

func GetPortById(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Port, error) {
	portId := d.Get("id")
	port, err := manager.GetPort(portId.(string))
	if err != nil {
		return nil, errors.Wrapf(err, "VDC with id '%s' not found", portId)
	}

	return port, nil
}

func GetPortByIp(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.Port, error) {
	portIp := d.Get("ip_address")
	allPorts, err := vdc.GetPorts()
	if err != nil {
		return nil, err
	}
	for _, port := range allPorts {
		if *port.IpAddress == portIp {
			return port, nil
		}
	}

	err = errors.New("Cant find port with this ip")
	return nil, err
}

func GetDnsByName(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Dns, error) {
	dnsName := d.Get("name").(string)
	project, err := manager.GetProject(d.Get("project_id").(string))
	if err != nil {
		return nil, errors.Wrap(err, "Error getting project by id")
	}

	dnss, err := project.GetDnss()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of dns")
	}

	for _, dns := range dnss {
		if strings.EqualFold(dns.Name, dnsName) {
			return dns, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Dns with name '%s' not found", dnsName)

}

func GetDnsById(d *schema.ResourceData, manager *rustack.Manager) (*rustack.Dns, error) {
	dnsid := d.Get("id").(string)
	dns, err := manager.GetDnss()

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of dns")
	}

	for _, dns := range dns {
		if strings.EqualFold(dns.ID, dnsid) {
			return dns, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Dns with id '%s' not found", dnsid)

}

func GetLbaasByName(d *schema.ResourceData, manager *rustack.Manager, vdc *rustack.Vdc) (*rustack.LoadBalancer, error) {
	lbaasName := d.Get("name")
	lbs, err := vdc.GetLoadBalancers(rustack.Arguments{"name": lbaasName.(string)})

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of lbs")
	}

	for _, lb := range lbs {
		if lb.Name == lbaasName {
			return lb, nil
		}
	}

	return nil, fmt.Errorf("Load Balancer with name '%s' not found in vdc '%s'", lbaasName, vdc.Name)
}

func GetS3ByName(d *schema.ResourceData, manager *rustack.Manager) (*rustack.S3Storage, error) {
	s3_name := d.Get("name").(string)
	project, err := manager.GetProject(d.Get("project_id").(string))
	if err != nil {
		return nil, errors.Wrap(err, "Error getting project by id")
	}

	storages, err := project.GetS3Storages()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of dns")
	}

	for _, s3 := range storages {
		if strings.EqualFold(s3.Name, s3_name) {
			return s3, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Dns with name '%s' not found", s3_name)

}

func GetS3ById(d *schema.ResourceData, manager *rustack.Manager) (*rustack.S3Storage, error) {
	s3_id := d.Get("id").(string)
	storages, err := manager.GetS3Storages()

	if err != nil {
		return nil, errors.Wrap(err, "Error getting list of dns")
	}

	for _, s3 := range storages {
		if strings.EqualFold(s3.ID, s3_id) {
			return s3, nil
		}
	}

	return nil, fmt.Errorf("ERROR: Dns with id '%s' not found", s3_id)

}

func checkDatasourceNameOrId(d *schema.ResourceData) (search string, err error) {
	id := d.Get("id").(string)
	name := d.Get("name").(string)
	if name == "" && id == "" {
		return "",fmt.Errorf("Must be specified 'name' or 'id'")
	}
	if id != "" {
		return "id", nil
	}
	return "name", nil
}
