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

## Using the provider

See [examples](examples) for usage and consult [Spotnana Security & Trust](https://spotnana.slack.com/archives/C03SV2FGLN7) team for help
```terraform
terraform {
  required_providers {
    snjumpcloud = {
      # This dummy URL is used to force Terraform to use the local build
      source = "test.com/Spotnana-Tech/snjumpcloud" 
    }
  }
}
# TF_VAR_api_key=$JUMPCLOUD_API_KEY
variable "api_key" {
  type      = string
  sensitive = true
}
provider "snjumpcloud" {
  apikey = var.api_key
}

resource "snjumpcloud_usergroup" "example_group" {
  name        = "example-terraform-group"
  description = "This group was created by Spotnana Terraform Provider!"
}

output "group_id" {
  value = snjumpcloud_usergroup.example_group.id
  description = "The ID of the group"
}
```

## Import existing resources
To simply manage the state of a resource, import it via the CLI
```shell
terraform import snjumpcloud_usergroup.example_groupname <<EXAMPLE_GROUP_ID>>
```

To prepare for a whole organization import, see
[this documentation](https://developer.hashicorp.com/terraform/language/import)
