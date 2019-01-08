package hedvig

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHedvigVdisk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHedvigVdiskDestroy("hedvig_vdisk.test-vdisk1"),
		Steps: []resource.TestStep{
			{
				Config: testAccHedvigVdiskConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHedvigVdiskExists("hedvig_vdisk.test-vdisk1"),
					testAccCheckHedvigVdiskExists("hedvig_vdisk.test-vdisk2"),
					testAccCheckHedvigVdiskSize("hedvig_vdisk.test-vdisk1"),
				),
			},
		},
	})
}

// TODO: Add update vdisk test

var testAccHedvigVdiskConfig = fmt.Sprintf(`
provider "hedvig" {
  node = "%s"
  username = "%s"
  password = "%s"
}

resource "hedvig_vdisk" "test-vdisk1" {
  name = "%s"
  size = 9
  type = "BLOCK"
}

resource "hedvig_vdisk" "test-vdisk2" {
  name = "%s"
  size = 11
  type = "NFS"
}
`, os.Getenv("HV_TESTNODE"), os.Getenv("HV_TESTUSER"), os.Getenv("HV_TESTPASS"),
	genRandomVdiskName(),
	genRandomVdiskName())

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

func testAccCheckHedvigVdiskDestroy(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hedvig_vdisk" {
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

func testAccCheckHedvigVdiskSize(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.Attributes["size"] == "9" {
			return nil
		}
		if rs.Primary.Attributes["size"] == "" {
			return errors.New("Size expected to not be nil")
		}
		return errors.New("Unknown problem with size of vdisk")
	}
}
