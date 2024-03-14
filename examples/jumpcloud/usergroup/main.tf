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

resource "jumpcloud_usergroup" "example_group" {
  name        = "example-terraform-group"
  description = "This group was created by Spotnana Terraform Provider!"
  members     = [
    "kgibson@spotnana.com",
    "bgodard@spotnana.com",]
}

output "group_details" {
  value = jumpcloud_usergroup.example_group
}