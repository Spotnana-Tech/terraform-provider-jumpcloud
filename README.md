# Jumpcloud Terraform Provider
A Terraform provider for managing Jumpcloud resources.

## Requirements
- [Spotnana Jumpcloud Go Client](https://github.com/Spotnana-Tech/sec-jumpcloud-client-go) >= 0.0.2
- [Jumpcloud](https://console.jumpcloud.com/) API Key in environment variable `JC_API_KEY`
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider
Clone the repository locally
```shell
git clone https://github.com/Spotnana-Tech/sec-terraform-provider-snjumpcloud.git
cd sec-terraform-provider-snjumpcloud
```
Build the provider using the Go `install` command:

```shell
go install .
```

While this provider is in alpha, we will be using a local build of the provider. Using this command to force terraform to reference the local build using the dummy URL.

```shell
export SN_GOPATH=~/go/bin 

cat > ~/.terraformrc <<EOF
provider_installation {

  dev_overrides {
     "test.com/Spotnana-Tech/snjumpcloud" = "$SN_GOPATH"  
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF
```
### Test your installation
Navigate to the test directory and check the plan. If the plan is successful, the provider is installed correctly.
```shell
cd ./examples/jumpcloud/confirm_install && terraform plan
```
## Using the provider
See [examples](examples) for usage and consult [Spotnana Security & Trust](https://spotnana.slack.com/archives/C03SV2FGLN7) team for help

While using local build of the provider, compact warnings to avoid long warnings in the output
```shell
export TF_CLI_ARGS_plan="-compact-warnings"
export TF_CLI_ARGS_apply="-compact-warnings"
```
Set the required environment variables
```shell
export TF_VAR_api_key=<<YOUR_JUMPCLOUD_API_KEY>>
```


```terraform
terraform {
  required_providers {
    snjumpcloud = {
      # This dummy URL is used to force Terraform to use the local build
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
  # filter the usergroups to only include those that start with "test"
  test_groups = [
    for g in data.snjumpcloud_usergroups.all_usergroups.usergroups : g.id
    if startswith(g.name, "test")
  ]
}

resource "snjumpcloud_usergroup" "example_group" {
  name        = "example-terraform-group"
  description = "This group was created by Spotnana Terraform Provider!"
}

output "group_id" {
  value = snjumpcloud_usergroup.example_group.id
  description = "The ID of the created group"
}
```
See the [examples](examples/jumpcloud) for more provider usage examples.


## Import existing resources
To simply manage the state of a resource, import it via the CLI
```shell
terraform import snjumpcloud_app_association.test_app <<EXAMPLE_APP_ID>>
```
Or import resources in your Terraform configuration file
```terraform
# Importing the app association via applicationID
import {
  to = snjumpcloud_app_association.test_app
  id = "65bc1fdaf6fc2af5f541a4c3" # The ID of the app association
}

# Associate the user groups with the app
resource "snjumpcloud_app_association" "test_app" {
  associated_groups = [
    snjumpcloud_usergroup.group1.id,
    snjumpcloud_usergroup.group2.id,
    snjumpcloud_usergroup.group3.id
  ]
}
```
See the [examples](examples/jumpcloud) for more provider import examples.

## Get support
If you need help with this provider, please reach out to the [Spotnana Security & Trust](https://spotnana.slack.com/archives/C03SV2FGLN7) team.
