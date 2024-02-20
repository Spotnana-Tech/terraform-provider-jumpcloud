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