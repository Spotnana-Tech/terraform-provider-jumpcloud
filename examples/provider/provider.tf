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