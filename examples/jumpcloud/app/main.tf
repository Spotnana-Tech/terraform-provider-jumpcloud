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

# Create a few user groups
resource "jumpcloud_usergroup" "group1" {
  name        = "Test-Terraform-Group1"
  description = "This group was created by Spotnana Terraform Provider!"
  members     = []
}
resource "jumpcloud_usergroup" "group2" {
  name        = "Test-Terraform-Group2"
  description = "This group was also created by Spotnana Terraform Provider!"
  members     = []
}
resource "jumpcloud_usergroup" "group3" {
  name        = "Test-Terraform-Group3"
  description = "This group was the 3rd created by Spotnana Terraform Provider!"
  members     = []
}

# Importing the app association via applicationID
import {
  to = jumpcloud_app.test_app
  id = "65bc1fdaf6fc2af5f541a4c3"
}

# Associate the user groups with the app
resource "jumpcloud_app" "test_app" {
  associated_groups = [
    jumpcloud_usergroup.group1.id,
    jumpcloud_usergroup.group2.id,
    jumpcloud_usergroup.group3.id
  ]
}
