package hedvig

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHedvigMount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHedvigMountDestroy("hedvig_mount.test-mount"),
		Steps: []resource.TestStep{
			{
				Config: testAccHedvigMountConfig,
				Check:  testAccCheckHedvigMountExists("hedvig_mount.test-mount"),
			},
		},
	})
}

var testAccHedvigMountConfig = fmt.Sprintf(`
provider "hedvig" {
  node = "%s"
  username = "%s"
  password = "%s"
}

resource "hedvig_vdisk" "test-mount-vdisk" {
  name = "%s"
  size = 11
  type = "NFS"
}

resource "hedvig_mount" "test-mount" {
  vdisk = "${hedvig_vdisk.test-mount-vdisk.name}"
  controller = "%s"
}
`, os.Getenv("HV_TESTNODE"), os.Getenv("HV_TESTUSER"), os.Getenv("HV_TESTPASS"),
	genRandomVdiskName(), os.Getenv("HV_TESTCONT"))

func testAccCheckHedvigMountExists(n string) resource.TestCheckFunc {
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

func testAccCheckHedvigMountDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hedvig_mount" {
				continue
			}
			name := rs.Primary.ID
			if name == n {
				return fmt.Errorf("Found resource: %s", name)
			}
		}
		return nil
	}
}
