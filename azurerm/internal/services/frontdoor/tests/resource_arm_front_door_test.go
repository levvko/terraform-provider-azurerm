package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMFrontDoor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMFrontDoor_requiresImport),
		},
	})
}

func TestAccAzureRMFrontDoor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFrontDoor_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMFrontDoor_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMFrontDoor_waf(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFrontDoorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMFrontDoor_waf(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMFrontDoorExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMFrontDoorExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Front Door not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Front Door %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on FrontDoorsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMFrontDoorDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Frontdoor.FrontDoorsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_front_door" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on FrontDoorsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMFrontDoor_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

locals {
  backend_name        = "backend-bing"
  endpoint_name       = "frontend-endpoint"
  health_probe_name   = "health-probe"
  load_balancing_name = "load-balancing-setting"
}

resource "azurerm_frontdoor" "test" {
  name                                         = "acctestfd-%d"
  location                                     = azurerm_resource_group.test.location
  resource_group_name                          = azurerm_resource_group.test.name
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = local.backend_name
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name = local.health_probe_name
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name                              = local.endpoint_name
    host_name                         = "acctestfd-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMFrontDoor_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_frontdoor" "import" {
  name                                         = azurerm_frontdoor.test.name
  location                                     = azurerm_frontdoor.test.location
  resource_group_name                          = azurerm_frontdoor.test.resource_group_name
  enforce_backend_pools_certificate_name_check = azurerm_frontdoor.test.enforce_backend_pools_certificate_name_check

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = local.backend_name
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name = local.health_probe_name
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name                              = local.endpoint_name
    host_name                         = "acctestfd-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMFrontDoor_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

locals {
  backend_name        = "backend-bing"
  endpoint_name       = "frontend-endpoint"
  health_probe_name   = "health-probe"
  load_balancing_name = "load-balancing-setting"
}

resource "azurerm_frontdoor" "test" {
  name                                         = "acctestfd-%d"
  location                                     = azurerm_resource_group.test.location
  resource_group_name                          = azurerm_resource_group.test.name
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = local.backend_name
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name = local.health_probe_name
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name                              = local.endpoint_name
    host_name                         = "acctestfd-%d.azurefd.net"
    custom_https_provisioning_enabled = false
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMFrontDoor_waf(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

locals {
  backend_name        = "backend-bing"
  endpoint_name       = "frontend-endpoint"
  health_probe_name   = "health-probe"
  load_balancing_name = "load-balancing-setting"
}

resource "azurerm_frontdoor_firewall_policy" "test" {
  name                = "acctestwafp%d"
  resource_group_name = azurerm_resource_group.test.name
  mode                = "Prevention"
}

resource "azurerm_frontdoor" "test" {
  name                                         = "acctestfd-%d"
  location                                     = azurerm_resource_group.test.location
  resource_group_name                          = azurerm_resource_group.test.name
  enforce_backend_pools_certificate_name_check = false

  routing_rule {
    name               = "routing-rule"
    accepted_protocols = ["Http", "Https"]
    patterns_to_match  = ["/*"]
    frontend_endpoints = [local.endpoint_name]
    forwarding_configuration {
      forwarding_protocol = "MatchRequest"
      backend_pool_name   = local.backend_name
    }
  }

  backend_pool_load_balancing {
    name = local.load_balancing_name
  }

  backend_pool_health_probe {
    name = local.health_probe_name
  }

  backend_pool {
    name = local.backend_name
    backend {
      host_header = "www.bing.com"
      address     = "www.bing.com"
      http_port   = 80
      https_port  = 443
    }

    load_balancing_name = local.load_balancing_name
    health_probe_name   = local.health_probe_name
  }

  frontend_endpoint {
    name                                    = local.endpoint_name
    host_name                               = "acctestfd-%d.azurefd.net"
    custom_https_provisioning_enabled       = false
    web_application_firewall_policy_link_id = azurerm_frontdoor_firewall_policy.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
