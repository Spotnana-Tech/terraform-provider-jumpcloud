package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUserGroups_GetAllGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Ensure that data source snjumpcloud_usergroups returns at least one user group
				Config: providerConfig + `data "snjumpcloud_usergroups" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.snjumpcloud_usergroups.test",
						"usergroups.#",
						regexp.MustCompile(`^0*[1-9]\d*$`)), // regex for a positive integer
				),
			},
		},
	})
}
