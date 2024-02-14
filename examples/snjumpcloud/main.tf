#
# This example covers all use cases for the snjumpcloud provider
#
terraform {
  required_providers {
    snjumpcloud = {
      source = "github.com/Spotnana-Tech/snjumpcloud"
    }
  }
}
variable "api_key" {
  type      = string
  sensitive = true
}
provider "snjumpcloud" {
  api_key = var.api_key
}

# Pulls all usergroups from the JumpCloud API
data "snjumpcloud_usergroups" "all_usergroups" {}

# Pulls all apps from the JumpCloud API
data "snjumpcloud_apps" "all_apps" {}

# Create a new usergroup
resource "snjumpcloud_usergroup" "new_usergroup" {
  name        = "new_usergroup"
  description = "This is a new usergroup"
}

# Importing the app association via applicationID
import {
  to = snjumpcloud_app.test_app
  id = "65bc1fdaf6fc2af5f541a4c3"
}

# Associate the user groups with the app
resource "snjumpcloud_app" "test_app" {
  associated_groups = [
    snjumpcloud_usergroup.new_usergroup.id
  ]
}

output "num_usergroups" {
  value = length(data.snjumpcloud_usergroups.all_usergroups.usergroups)
}

output "num_apps" {
  value = length(data.snjumpcloud_apps.all_apps.apps)
}