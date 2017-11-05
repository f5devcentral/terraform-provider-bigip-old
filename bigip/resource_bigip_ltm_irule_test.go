package bigip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/scottdware/go-bigip"
)

var TEST_IRULE_NAME = "/" + TEST_PARTITION + "/test-rule"

var TEST_IRULE_RESOURCE = `
	resource "bigip_ltm_irule" "test-rule" {
		name = "` + TEST_IRULE_NAME + `"
		irule = <<EOF
when CLIENT_ACCEPTED {
     log local0. "test"
}
EOF
	}`

func TestBigipLtmIRule_create(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIRulesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_IRULE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckIRuleExists(TEST_IRULE_NAME),
				),
			},
		},
	})
}

func TestBigipLtmIRule_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAcctPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testCheckIRulesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: TEST_IRULE_RESOURCE,
				Check: resource.ComposeTestCheckFunc(
					testCheckIRuleExists(TEST_IRULE_NAME),
				),
				ResourceName:      TEST_IRULE_NAME,
				ImportState:       false,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckIRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*bigip.BigIP)

		irule, err := client.IRule(name)
		if err != nil {
			return fmt.Errorf("Error while fetching irule: %v", err)

		}
		body := s.RootModule().Resources["bigip_ltm_irule.test-rule"].Primary.Attributes["irule"]
		if irule.Rule != body {
			return fmt.Errorf("IRule body does not match. Expecting %s got %s.", body, irule.Rule)
		}

		irule_name := fmt.Sprintf("/%s/%s", irule.Partition, irule.Name)
		if irule_name != name {
			return fmt.Errorf("IRule name does not match. Expecting %s got %s.", name, irule_name)
		}
		return nil
	}
}

func testCheckIRulesDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*bigip.BigIP)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "bigip_ltm_irule" {
			continue
		}

		name := rs.Primary.ID
		irule, err := client.IRule(name)

		if err != nil {
			return nil
		}
		if irule != nil {
			return fmt.Errorf("IRule %s not destroyed.", name)
		}
	}
	return nil
}
