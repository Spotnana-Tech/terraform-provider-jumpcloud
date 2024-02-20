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
data "jumpcloud_usergroups" "all_usergroups" {}

output "number_of_usergroups" {
  value       = length(data.jumpcloud_usergroups.all_usergroups.usergroups)
  description = "The number of usergroups available in the JumpCloud API"
}