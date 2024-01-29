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

While this provider is in alpha, we will be using a local build of the provider. Using this command to reference the local build using the dummy URL.

```shell
p=~/go/bin # or wherever your $GOPATH is

cat > ~/.terraformrc <<EOF
provider_installation {

  dev_overrides {
     "test.com/Spotnana-Tech/snjumpcloud" = "$p"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
EOF
```


## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/Spotnana-Tech/sec-jumpcloud-client-go
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

See [examples](examples) for usage and consult [Spotnana Security & Trust](https://spotnana.slack.com/archives/C03SV2FGLN7) team for help
