package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceUserGroups_CreateGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create a new user group and verify that aspects of it are correct
				Config: providerConfig + `resource "jumpcloud_usergroup" "new_usergroup" {
											name        = "new_usergroup_terraform_test"
											description = "This group made via terraform test"
										}`,
				// Compose multiple test checks to verify the resource
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("jumpcloud_usergroup.new_usergroup",
						"name",
						"new_usergroup_terraform_test"),
					resource.TestCheckResourceAttr("jumpcloud_usergroup.new_usergroup",
						"type",
						"user_group"),
					resource.TestCheckResourceAttr("jumpcloud_usergroup.new_usergroup",
						"membership_method",
						"STATIC"),
				),
			},
		},
	})
}
