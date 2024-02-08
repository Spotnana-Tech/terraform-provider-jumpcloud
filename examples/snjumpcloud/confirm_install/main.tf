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
  api_key = var.api_key
}
data "snjumpcloud_usergroups" "all_usergroups" {}

output "number_of_usergroups" {
  value       = length(data.snjumpcloud_usergroups.all_usergroups.usergroups)
  description = "The number of usergroups available in the JumpCloud API"
}