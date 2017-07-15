package buildkite

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPipeline_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPipeline_basic,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						ms := s.RootModule()
						rs := ms.Resources["buildkite_pipeline.test"]
						if rs == nil {
							return fmt.Errorf("no state for buildkite_pipeline.test")
						}
						if got, want := rs.Primary.Attributes["slug"], "tf-acc-basic"; got != want {
							return fmt.Errorf("incorrect slug %s; want %s", got, want)
						}
						return nil
					},
				),
			},
		},
	})
}

const testAccPipeline_basic = `
resource "buildkite_pipeline" "test" {
  name = "tf-acc-basic"
  repository = "git@github.com:saymedia/terraform-provider-buildkite.git"
  default_branch = "master"

  step {
    type = "script"
    name = "test"
    command = "echo 'Hello World'"
  }
}
`
