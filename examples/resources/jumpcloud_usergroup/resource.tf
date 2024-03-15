resource "jumpcloud_usergroup" "example" {
  name        = "example-name"
  description = "example description"
  members = [
    "example-member-1",
    "example-member-2",
  ]
}