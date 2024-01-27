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

resource "snjumpcloud_usergroup" "example_group" {
    name = "example-terraform-group"
    description = "This group was created by Spotnana Terraform Provider!"
}

// export TF_VAR_api_key=$JC_API_KEY
// export TF_LOG=TRACE
// terraform plan

// api_key=$JC_API_KEY TF_LOG=TRACE terraform plan