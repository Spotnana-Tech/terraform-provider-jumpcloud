# Jumpcloud Terraform Provider
A Terraform provider for managing Jumpcloud resources.

_"Jumpcloud" is a trademark of Jumpcloud, Inc.
"Terraform" is a trademark of HashiCorp, Inc.
"Go" is a trademark of Google LLC or its affiliate ("Google") for its programming language (see https://go.dev/brand).
These marks are used nominatively to indicate the nature and function of Spotnana's
terraform provider, which is neither sponsored or endorsed by Jumpcloud, Inc., Hashicorp, Inc. or Google._

### Requirements
- [Spotnana Jumpcloud Go Client](https://github.com/Spotnana-Tech/sec-jumpcloud-client-go) >= 1.0.0
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

---

# Getting Started
```terraform
terraform {
  required_providers {
    jumpcloud = {
      source = "Spotnana-Tech/jumpcloud"
    }
  }
}

provider "jumpcloud" {
  api_key = var.api_key
}

resource "jumpcloud_usergroup" "example_group" {
  name = "example"
}
```
# Usage
See the [core example](examples/jumpcloud/main.tf) to see all features executed in a single plan.

### Importing Existing Resources
Add an import block to the Terraform configuration file for the resource you want to import.
```terraform
import {
  to = jumpcloud_app.example_app
  id = "6abcd1230987654321" # The `app_id` of the application in Jumpcloud
}
```
Generate a `.tf` file for the resource you want to import.
```shell
terraform plan -generate-config-out="generated.tf"
```
---
## Installation for Local Development
Clone the repository locally
```shell
git clone https://github.com/Spotnana-Tech/terraform-provider-jumpcloud.git
cd terraform-provider-jumpcloud
```
Build the provider using the Go `install` command:

```shell
go install .
```
Use this command to force terraform to reference the local build within `~/.terraformrc`:
```shell
export SN_GOPATH=~/go/bin 

cat > ~/.terraformrc <<EOF
provider_installation {

  dev_overrides {
     "Spotnana-Tech/jumpcloud" = "$SN_GOPATH"  
  }
  direct {}
}
EOF
```

## Testing
Run `make testacc` to run the full suite of Acceptance tests.

Set the required environment variables, navigate to the example directory, and run `terraform plan`.
```shell
export TF_VAR_api_key=<<YOUR_JUMPCLOUD_API_KEY>>

cd ./examples/jumpcloud/confirm_install && terraform plan
```
---


## Get Support
Open a new issue on this repository.
