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
  name        = "example-terraform-group"
  description = "This group was created by Spotnana Terraform Provider!"
}

output "group_details" {
  value       = snjumpcloud_usergroup.example_group
  description = "The ID of the group"
}