package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMIotHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHub_ipFilterRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_ipFilterRules(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHub_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHub_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub"),
			},
		},
	})
}

func TestAccAzureRMIotHub_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_standard(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHub_customRoutes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_customRoutes(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint.0.type", "AzureIotHub.StorageContainer"),
					resource.TestCheckResourceAttr(data.ResourceName, "endpoint.1.type", "AzureIotHub.EventHub"),
					resource.TestCheckResourceAttr(data.ResourceName, "route.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHub_fileUpload(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_fileUpload(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "file_upload.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "file_upload.0.lock_duration", "PT5M"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHub_fallbackRoute(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHub_fallbackRoute(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "fallback_route.0.source", "DeviceMessages"),
					resource.TestCheckResourceAttr(data.ResourceName, "fallback_route.0.endpoint_names.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "fallback_route.0.enabled", "true"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMIotHubDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("IotHub %s still exists in resource group %s", name, resourceGroup)
		}
	}
	return nil
}

func testCheckAzureRMIotHubExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		iothubName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IotHub: %s", iothubName)
		}

		resp, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: IotHub %q (resource group: %q) does not exist", iothubName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMIotHub_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "B1"
    tier     = "Basic"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotHub_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotHub_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub" "import" {
  name                = "${azurerm_iothub.test.name}"
  resource_group_name = "${azurerm_iothub.test.name}"
  location            = "${azurerm_iothub.test.location}"

  sku {
    name     = "B1"
    tier     = "Basic"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, template)
}

func testAccAzureRMIotHub_standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotHub_ipFilterRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  ip_filter_rule {
    name    = "test"
    ip_mask = "10.0.0.0/31"
    action  = "Accept"
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotHub_customRoutes(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_eventhub_namespace" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  name                = "acctest-%d"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest"
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  eventhub_name       = "${azurerm_eventhub.test.name}"
  name                = "acctest"
  send                = true
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  endpoint {
    type                       = "AzureIotHub.StorageContainer"
    connection_string          = "${azurerm_storage_account.test.primary_blob_connection_string}"
    name                       = "export"
    batch_frequency_in_seconds = 60
    max_chunk_size_in_bytes    = 10485760
    container_name             = "test"
    encoding                   = "Avro"
    file_name_format           = "{iothub}/{partition}_{YYYY}_{MM}_{DD}_{HH}_{mm}"
  }

  endpoint {
    type              = "AzureIotHub.EventHub"
    connection_string = "${azurerm_eventhub_authorization_rule.test.primary_connection_string}"
    name              = "export2"
  }

  route {
    name           = "export"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export"]
    enabled        = true
  }

  route {
    name           = "export2"
    source         = "DeviceMessages"
    condition      = "true"
    endpoint_names = ["export2"]
    enabled        = true
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMIotHub_fallbackRoute(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  fallback_route {
    source         = "DeviceMessages"
    endpoint_names = ["events"]
    enabled        = true
  }

  tags = {
    purpose = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotHub_fileUpload(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "test"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  file_upload {
    connection_string  = "${azurerm_storage_account.test.primary_blob_connection_string}"
    container_name     = "${azurerm_storage_container.test.name}"
    notifications      = true
    max_delivery_count = 12
    sas_ttl            = "PT2H"
    default_ttl        = "PT3H"
    lock_duration      = "PT5M"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}
