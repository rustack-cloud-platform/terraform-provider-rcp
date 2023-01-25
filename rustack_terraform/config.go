package rustack_terraform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/pilat/rustack-go/rustack"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Token            string
	APIEndpoint      string
	TerraformVersion string
	ClientID         string
}

type CombinedConfig struct {
	manager *rustack.Manager
}

func (c *CombinedConfig) rustackManager() *rustack.Manager { return c.manager }

func (c *Config) Client() (*CombinedConfig, diag.Diagnostics) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	manager := rustack.NewManager(c.Token)
	manager.Logger = logger
	manager.BaseURL = strings.TrimSuffix(c.APIEndpoint, "/")
	manager.ClientID = c.ClientID
	manager.UserAgent = fmt.Sprintf("Terraform/%s", c.TerraformVersion)

	return &CombinedConfig{
		manager: manager,
	}, nil
}
