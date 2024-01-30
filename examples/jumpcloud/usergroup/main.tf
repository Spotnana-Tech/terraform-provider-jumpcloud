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

resource "snjumpcloud_usergroup" "example_group" {
  name        = "example-terraform-group-changed-TEST"
  description = "This group was created by Spotnana Terraform Provider!"
}

output "group_id" {
  value = snjumpcloud_usergroup.example_group.id
  description = "The ID of the group"
}