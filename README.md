# Jumpcloud Terraform Provider
A Terraform provider for managing Jumpcloud resources.

### Requirements
- [Spotnana Jumpcloud Go Client](https://github.com/Spotnana-Tech/sec-jumpcloud-client-go) >= 0.0.2
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
     "github.com/Spotnana-Tech/snjumpcloud" = "$SN_GOPATH"  
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF
```

## Testing the provider
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
Navigate to the example directory and check the plan. If the plan is successful, the provider is installed correctly.
```shell
cd ./examples/jumpcloud/confirm_install && terraform plan
```

See the [examples](examples/snjumpcloud) for more provider usage examples.


## Import existing resources

Add an import block to the Terraform configuration file for the resource you want to import.
```terraform
import {
  to = snjumpcloud_app.example_app
  id = "6abcd1230987654321" # The `app_id` of the application in Jumpcloud
}
```
Generate a `.tf` file for the resource you want to import.
```shell
terraform plan -generate-config-out="generated.tf"
```

## Get support
If you need help with this provider, please reach out to the [Spotnana Security & Trust](https://spotnana.slack.com/archives/C03SV2FGLN7) team.
