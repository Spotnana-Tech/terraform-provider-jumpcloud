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

# Create a few user groups
resource "snjumpcloud_usergroup" "group1" {
  name        = "Test-Terraform-Group1"
  description = "This group was created by Spotnana Terraform Provider!"
}
resource "snjumpcloud_usergroup" "group2" {
  name        = "Test-Terraform-Group2"
  description = "This group was also created by Spotnana Terraform Provider!"
}
resource "snjumpcloud_usergroup" "group3" {
  name        = "Test-Terraform-Group3"
  description = "This group was the 3rd created by Spotnana Terraform Provider!"
}

# Importing the app association via applicationID
import {
  to = snjumpcloud_app_association.test_app
  id = "65bc1fdaf6fc2af5f541a4c3"
}

# Associate the user groups with the app
resource "snjumpcloud_app_association" "test_app" {
  associated_groups = [
    snjumpcloud_usergroup.group1.id,
    snjumpcloud_usergroup.group2.id,
    snjumpcloud_usergroup.group3.id
  ]
}

output "app_details" {
  value       = snjumpcloud_app_association.test_app
  description = "The details of the test app"
}