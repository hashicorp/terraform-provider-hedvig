package hedvig

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"hedvig": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	requiredVars := []string{"HV_TESTNODE", "HV_TESTCONT", "HV_TESTUSER",
		"HV_TESTPASS", "HV_TESTADDR", "HV_TESTADDR2"}
	missingVars := []string{}

	for _, v := range requiredVars {
		if _, ok := os.LookupEnv(v); !ok {
			missingVars = append(missingVars, v)
		}
	}

	if len(missingVars) > 0 {
		t.Fatalf("The following env vars must be set for acceptance tests: %s", missingVars)
	}
	err := testAccProvider.Configure(terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}
