package hedvig

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

func TestAccHedvigVdisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccHedvigVdiskConfig,
				Check:  resource.ComposeTestCheckFunc(testAccCheckHedvigVdiskExists("hedvig_vdisk.test-vdisk1"), testAccCheckHedvigVdiskExists("hedvig_vdisk.test-vdisk2")),
			},
		},
	})
}

var testAccHedvigVdiskConfig = fmt.Sprintf(`
provider "hedvig" {
  node = "%s"
  username = "%s"
  password = "%s"
}

resource "hedvig_vdisk" "test-vdisk1" {
  cluster = "%s"
  name = "%s"
  size = 9
  type = "BLOCK"
}

resource "hedvig_vdisk" "test-vdisk2" {
  cluster = "%s"
  name = "%s"
  size = 11
  type = "NFS"
}
`, os.Getenv("HV_TESTNODE"), os.Getenv("HV_TESTUSER"), os.Getenv("HV_TESTPASS"),
	os.Getenv("HV_TESTCLUST"), genRandomVdiskName(),
	os.Getenv("HV_TESTCLUST"), genRandomVdiskName())

func testAccCheckHedvigVdiskExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No vdisk ID is set")
		}

		return nil
	}
}
