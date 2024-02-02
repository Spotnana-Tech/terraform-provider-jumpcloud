terraform {
  required_providers {
    snjumpcloud = {
      source = "test.com/Spotnana-Tech/snjumpcloud"
    }
  }
}

variable "api_key" {
  type      = string
  sensitive = true
}
provider "snjumpcloud" {
  apikey = var.api_key
}

resource "snjumpcloud_usergroup" "my_group" {
  name        = "example-terraform-group"
  description = "A test group created by Terraform"
}

output "usergroup_id" {
  value = snjumpcloud_usergroup.my_group.id
  description = "The ID of the group"
}
