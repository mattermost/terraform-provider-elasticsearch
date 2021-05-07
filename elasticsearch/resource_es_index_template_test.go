package elasticsearch

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIndexTemplateValid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIndexTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIndexTemplateConfigBasic("logstash"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIndexTemplateExists("elasticsearch_template.new"),
				),
			},
		},
	})
}

func TestAccIndexTemplate_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIndexTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIndexTemplateConfigBasic("logstash"),
			},
			{
				ResourceName:      "elasticsearch_template.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIndexTemplate_Invalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIndexTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckIndexTemplateConfigInvalid("logstash"),
				ExpectError: regexp.MustCompile("\"template\" contains an invalid JSON"),
			},
		},
	})
}

func testAccCheckIndexTemplateDestroy(s *terraform.State) error {
	opts := testAccProvider.Meta().(*providerOpts)

	c, err := getClientES(opts)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticsearch_template" {
			continue
		}

		id := rs.Primary.ID
		req := esapi.IndicesDeleteIndexTemplateRequest{
			Name:   id,
			Pretty: true,
		}
		_, err := req.Do(context.TODO(), c)
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckIndexTemplateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set")
		}

		return nil
	}
}

func testAccCheckIndexTemplateConfigBasic(name string) string {
	return fmt.Sprintf(`
	resource "elasticsearch_template" "new" {
		name = "%s"

		template = <<EOF
		{
			"index_patterns": ["foo*", "bar*"],
			"settings": {
				"index": {
					"number_of_shards": 1
				},
				"mapping": {
					"total_fields": 2000
				}
			}
		}
		EOF
	}
	`, name)
}

func testAccCheckIndexTemplateConfigInvalid(name string) string {
	return fmt.Sprintf(`
	resource "elasticsearch_template" "new" {
		name = "%s"

		template = "asdfa"
	}
	`, name)
}
