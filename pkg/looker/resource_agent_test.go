package looker

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func TestAcc_Agent(t *testing.T) {
	name1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	name2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: agentConfig(name1, "test description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_agent.test", "name", name1),
					resource.TestCheckResourceAttr("looker_agent.test", "description", "test description"),
					resource.TestCheckResourceAttr("looker_agent.test", "code_interpreter", "false"),
				),
			},
			{
				Config: agentConfig(name2, "updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_agent.test", "name", name2),
					resource.TestCheckResourceAttr("looker_agent.test", "description", "updated description"),
				),
			},
			{
				ResourceName:      "looker_agent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
		CheckDestroy: testAccCheckAgentDestroy,
	})
}

func TestAcc_AgentWithSources(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: agentConfigWithSources(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("looker_agent.test_sources", "name", name),
					resource.TestCheckResourceAttr("looker_agent.test_sources", "sources.#", "1"),
					resource.TestCheckResourceAttr("looker_agent.test_sources", "sources.0.model", "test_model"),
					resource.TestCheckResourceAttr("looker_agent.test_sources", "sources.0.explore", "test_explore"),
					resource.TestCheckResourceAttr("looker_agent.test_sources", "instructions", "test instructions"),
					resource.TestCheckResourceAttr("looker_agent.test_sources", "code_interpreter", "true"),
				),
			},
		},
		CheckDestroy: testAccCheckAgentDestroy,
	})
}

func testAccCheckAgentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*apiclient.LookerSDK)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "looker_agent" {
			continue
		}

		agentID := rs.Primary.ID

		_, err := client.GetAgent(agentID, "", nil)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				continue
			}
			return err
		}

		return fmt.Errorf("agent still exists: %s", rs.Primary.ID)
	}

	return nil
}

func agentConfig(name, description string) string {
	return fmt.Sprintf(`
	resource "looker_agent" "test" {
		name        = "%s"
		description = "%s"
	}
	`, name, description)
}

func agentConfigWithSources(name string) string {
	return fmt.Sprintf(`
	resource "looker_agent" "test_sources" {
		name        = "%s"
		description = "agent with sources"

		sources {
			model   = "test_model"
			explore = "test_explore"
		}

		instructions     = "test instructions"
		code_interpreter = true
	}
	`, name)
}
