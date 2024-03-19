package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceGroupLookup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create a new user group and verify that aspects of it are correct
				Config: providerConfig + `data "jumpcloud_group_lookup" "g" {
										  name = "a"
										}`,
				// Compose multiple test checks to verify the resource
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.jumpcloud_group_lookup.g",
						"groups.#",
						regexp.MustCompile(`^0*[1-9]\d*$`)), // regex for a positive integer
				),
			},
		},
	})
}
