#
# This example covers all use cases for the jumpcloud provider
#
terraform {
  required_providers {
    jumpcloud = {
      source = "Spotnana-Tech/jumpcloud"
    }
  }
}
variable "api_key" {
  type      = string
  sensitive = true
}
provider "jumpcloud" {
  api_key = var.api_key
}

# Pulls all usergroups from the JumpCloud API
data "jumpcloud_usergroups" "all_usergroups" {}

# Pulls all apps from the JumpCloud API
data "jumpcloud_apps" "all_apps" {}

# Create a new usergroup
resource "jumpcloud_usergroup" "new_usergroup" {
  name        = "tf-provider-test-new_usergroup"
  description = "This is a new usergroup from the Terraform provider"
  members     = []
}

# Importing the app association via applicationID
import {
  to = jumpcloud_app.test_app
  id = "65bc1fdaf6fc2af5f541a4c3"
}

# Associate the user groups with the app
resource "jumpcloud_app" "test_app" {
  associated_groups = [
    jumpcloud_usergroup.new_usergroup.id
  ]
}

output "num_usergroups" {
  value = length(data.jumpcloud_usergroups.all_usergroups.usergroups)
}

output "num_apps" {
  value = length(data.jumpcloud_apps.all_apps.apps)
}