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

resource "snjumpcloud_usergroup" "group1" {
  name        = "Test-Terraform-Group1"
  description = "This group was created by Spotnana Terraform Provider!"
}
resource snjumpcloud_usergroup "group2" {
    name        = "Test-Terraform-Group2"
    description = "This group was also created by Spotnana Terraform Provider!"
}

import {
  to = snjumpcloud_app_association.test_app_association
  id = "65bc1fdaf6fc2af5f541a4c3"
}

resource "snjumpcloud_app_association" "test_app_association" {
  associated_groups = [
    snjumpcloud_usergroup.group1.id,
    snjumpcloud_usergroup.group2.id
  ]
}

output "app_details2" {
  value = snjumpcloud_app_association.test_app_association
  description = "The name of the group"
}
output "group_name" {
  value = snjumpcloud_usergroup.group1.id
  description = "The ID of the group"
}