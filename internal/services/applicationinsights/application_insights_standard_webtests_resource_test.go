package applicationinsights_test

import (
	"context"
	"fmt"
	"testing"

	webtests "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApplicationInsightsStandardWebTestResource struct{}

func TestAccApplicationInsightsStandardWebTest_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_standard_web_test", "test")
	testResource := ApplicationInsightsStandardWebTestResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
	})
}

func TestAccApplicationInsightsStandardWebTest_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_standard_web_test", "test")
	testResource := ApplicationInsightsStandardWebTestResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.RequiresImportErrorStep(testResource.requiresImportConfig),
	})
}

func TestAccApplicationInsightsStandardWebTest_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_application_insights_standard_web_test", "test")
	testResource := ApplicationInsightsStandardWebTestResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.completeConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
		data.ApplyStep(testResource.completeConfig, testResource),
		data.ImportStep(),
	})
}

func (ApplicationInsightsStandardWebTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webtests.ParseWebTestID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppInsights.StandardWebTestsClient.WebTestsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Properties != nil), nil
}

func (ApplicationInsightsStandardWebTestResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_standard_web_test" "test" {
  name                    = "acctestappinsightswebtests-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  geo_locations           = ["us-tx-sn1-azr"]

  request {
    follow_redirects_enabled         = false
    http_verb                        = "GET"
    parse_dependent_requests_enabled = false
    url                              = "http://microsoft.com"

    header {
      name  = "x-header"
      value = "testheader"
    }
    header {
      name  = "x-header-2"
      value = "testheader2"
    }
  }

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApplicationInsightsStandardWebTestResource) requiresImportConfig(data acceptance.TestData) string {
	template := r.basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_application_insights_standard_web_test" "import" {
  name                    = azurerm_application_insights_standard_web_test.test.name
  location                = azurerm_application_insights_standard_web_test.test.location
  resource_group_name     = azurerm_application_insights_standard_web_test.test.resource_group_name
  application_insights_id = azurerm_application_insights_standard_web_test.test.application_insights_id
  geo_locations           = azurerm_application_insights_standard_web_test.test.geo_locations
  request {
    follow_redirects_enabled         = azurerm_application_insights_standard_web_test.test.request.0.follow_redirects_enabled
    http_verb                        = azurerm_application_insights_standard_web_test.test.request.0.http_verb
    parse_dependent_requests_enabled = azurerm_application_insights_standard_web_test.test.request.0.parse_dependent_requests_enabled
    url                              = azurerm_application_insights_standard_web_test.test.request.0.url

    header {
      name  = azurerm_application_insights_standard_web_test.test.request.0.header.0.name
      value = azurerm_application_insights_standard_web_test.test.request.0.header.0.value
    }

    header {
      name  = azurerm_application_insights_standard_web_test.test.request.0.header.1.name
      value = azurerm_application_insights_standard_web_test.test.request.0.header.1.value
    }
  }
}
`, template)
}

func (ApplicationInsightsStandardWebTestResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appinsights-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_application_insights_standard_web_test" "test" {
  name                    = "acctestappinsightswebtests-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  application_insights_id = azurerm_application_insights.test.id
  frequency               = 900
  timeout                 = 120
  enabled                 = true
  description             = "web_test"
  retry_enabled           = true
  tags = {
    ENV = "web_test"
  }
  geo_locations = ["us-tx-sn1-azr", "us-il-ch1-azr"]

  request {
    follow_redirects_enabled         = true
    http_verb                        = "POST"
    parse_dependent_requests_enabled = true
    url                              = "https://microsoft.com"

    body = "{\"test\": \"value\"}"

    header {
      name  = "x-header"
      value = "testheader"
    }
    header {
      name  = "x-header-2"
      value = "testheaderupdated"
    }
  }
  validation_rules {
    expected_status_code = 200

    ssl_cert_remaining_lifetime = 20
    ssl_check_enabled           = true

    content {
      content_match      = "Unknown"
      ignore_case        = true
      pass_if_text_found = true
    }
  }

  lifecycle {
    ignore_changes = ["tags"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
