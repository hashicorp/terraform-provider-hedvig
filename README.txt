Hedvig Terraform Provider

To run acceptance tests:

Set these environment variables:
TF_ACC=1
HV_TESTUSER= *admin username on target cluster*
HV_TESTPASS= *admin password on target cluster*
HV_TESTCLUST= *name of target cluster*
HV_TESTNODE= *FQDN of node in cluster*
HV_TESTCONT= *FQDN of cluster controller*

Run "go build -o terraform-provider-hedvig"
Run "go test -v"
