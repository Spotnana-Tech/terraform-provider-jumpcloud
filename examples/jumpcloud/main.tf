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

# Pulls all usergroups from the JumpCloud API
data "snjumpcloud_usergroups" "all_usergroups" {}

locals {
  # filter the usergroups to only include those that start with "SEC-"
  sec_groups = [
    for g in data.snjumpcloud_usergroups.all_usergroups.usergroups : g.id
    if startswith(g.name, "SEC-") # && g.name != "SEC-Admins"
  ]
  # Different types of list comprehensions
  # [groupid, groupid, ...]
  all_test_groups = [
    for g in data.snjumpcloud_usergroups.all_usergroups.usergroups : g.id
    if startswith(g.name, "test")
  ]
  # {groupname: groupid, groupname: groupid, ...}
  other_groups = {
  for g in data.snjumpcloud_usergroups.all_usergroups.usergroups : g.name => g.id
  if startswith(g.name, "test")
  }
}

output "all_usergroups" {
  value = data.snjumpcloud_usergroups.all_usergroups
  description = "The usergroup data available in the JumpCloud API"
}
output "number_of_usergroups" {
  value = length(data.snjumpcloud_usergroups.all_usergroups.usergroups)
  description = "The number of usergroups available in the JumpCloud API"
}
output "sec_groups" {
  value = local.sec_groups
}
output "z_ll_test_groups" {
  value = local.all_test_groups
}
output "z_other_groups" {
  value = local.other_groups
}