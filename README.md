# Jumpcloud Terraform Provider
A Terraform provider for managing Jumpcloud resources.

### Requirements
- [Spotnana Jumpcloud Go Client](https://github.com/Spotnana-Tech/sec-jumpcloud-client-go) >= 1.0.0
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

---

# Getting Started
## Installation
Clone the repository locally
```shell
git clone https://github.com/Spotnana-Tech/terraform-provider-jumpcloud.git
cd terraform-provider-jumpcloud
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
     "Spotnana-Tech/jumpcloud" = "$SN_GOPATH"  
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF
```

## Testing
See [examples](examples/jumpcloud) for usage and consult [Spotnana Security & Trust](https://spotnana.slack.com/archives/C03SV2FGLN7) team for help

While using local build of the provider, compact warnings to avoid long warnings in the output
```shell
export TF_CLI_ARGS_plan="-compact-warnings"
export TF_CLI_ARGS_apply="-compact-warnings"
```
Set the required environment variables
```shell
export TF_VAR_api_key=<<YOUR_JUMPCLOUD_API_KEY>>
```
Navigate to the example directory and check the plan. If the plan is successful, the provider is installed correctly.
```shell
cd ./examples/jumpcloud/confirm_install && terraform plan
```
---
# Usage
### All Features Example
See the [core example](examples/jumpcloud/main.tf) to see all features executed in a single plan. Alternatively browse the [documentation](/docs) for a detailed list.

### Importing Existing Resources

Add an import block to the Terraform configuration file for the resource you want to import.
```terraform
import {
  to = jumpcloud_app.example_app
  id = "6abcd1230987654321" # The `app_id` of the application in Jumpcloud
}
```
### Generate Configuration
Generate a `.tf` file for the resource you want to import.
```shell
terraform plan -generate-config-out="generated.tf"
```

## Get support
If you need help with this provider, please reach out to the [Spotnana Security & Trust](https://spotnana.slack.com/archives/C03SV2FGLN7) team.
