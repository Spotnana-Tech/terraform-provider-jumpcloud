terraform {
  required_providers {
    snjumpcloud = {
      source = "test.com/Spotnana-Tech/snjumpcloud"
    }
  }
}

provider "snjumpcloud" {
  apikey = "abc123"
}

data "snjumpcloud_usergroup" "example_usergroup" {}
