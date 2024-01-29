terraform {
  required_providers {
    snjumpcloud = {
      source = "test.com/Spotnana-Tech/snjumpcloud"
    }
  }
}

variable "api_key" {
  type = string
  sensitive = true
}
provider "snjumpcloud" {
  apikey = var.api_key
}

data "snjumpcloud_usergroup" "allusergroups" {}

output "allusergroups" {
  value = data.snjumpcloud_usergroup.allusergroups
}
