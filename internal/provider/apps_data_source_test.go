package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceApps_GetAllApps(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Ensure that data source jumpcloud_usergroups returns at least one user group
				Config: providerConfig + `data "jumpcloud_apps" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.jumpcloud_apps.test",
						"apps.#",
						regexp.MustCompile(`^0*[1-9]\d*$`)), // regex for a positive integer
				),
			},
		},
	})
}
