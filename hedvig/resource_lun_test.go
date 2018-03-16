package hedvig

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
	"os"
)

func testHedvigLun() error {
	return nil
}

func TestAccHedvigLun(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHedvigLunConfig,
				Check:  testAccCheckHedvigLunExists("hedvig_lun.test-lun"),
			},
		},
	})
}

var testAccHedvigLunConfig = fmt.Sprintf(`
provider "hedvig" {
  node = "%s"
  username = "%s"
  password = "%s"
}

resource "hedvig_vdisk" "test-lun-vdisk" {
  cluster = "%s"
  name = "HedvigVdiskTest3"
  size = 9
  type = "BLOCK"
}

resource "hedvig_lun" "test-lun" {
  cluster = "%s"
  vdisk = "${hedvig_vdisk.test-lun-vdisk.name}"
  controller = "%s"
}
`, os.Getenv("HV_TESTNODE"), os.Getenv("HV_TESTUSER"), os.Getenv("HV_TESTPASS"), os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTCLUST"), os.Getenv("HV_TESTCONT"))

func testAccCheckHedvigLunExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No lun ID is set")
		}

		return nil
	}
}
