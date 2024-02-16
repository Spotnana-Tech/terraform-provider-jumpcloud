default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	@TF_ACC=1 TF_VAR_api_key=$(JC_API_KEY) go test ./... -v $(TESTARGS) -timeout 120m
