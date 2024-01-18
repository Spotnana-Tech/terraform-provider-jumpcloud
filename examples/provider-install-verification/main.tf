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

data "snjumpcloud_group_persona" "example_persona" {}

// terraform plan -var "api_key=$JC_API_KEY"