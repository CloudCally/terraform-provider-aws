package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
)

func TestAccAWSPinpointSMSChannel_basic(t *testing.T) {
	var channel pinpoint.SMSChannelResponse
	resourceName := "aws_pinpoint_sms_channel.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSPinpointApp(t) },
		ErrorCheck:   acctest.ErrorCheck(t, pinpoint.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSPinpointSMSChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSPinpointSMSChannelConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSPinpointSMSChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// There can be a delay before these Computed values are returned
				// e.g. 0 on Create -> Read, 20 on Import
				// These seem non-critical for other Terraform resource references,
				// so ignoring them for now, but we can likely adjust the Read function
				// to wait until they are available on creation with retry logic.
				ImportStateVerifyIgnore: []string{
					"promotional_messages_per_second",
					"transactional_messages_per_second",
				},
			},
			{
				Config: testAccAWSPinpointSMSChannelConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSPinpointSMSChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func TestAccAWSPinpointSMSChannel_full(t *testing.T) {
	var channel pinpoint.SMSChannelResponse
	resourceName := "aws_pinpoint_sms_channel.test"
	senderId := "1234"
	shortCode := "5678"
	newShortCode := "7890"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSPinpointApp(t) },
		ErrorCheck:   acctest.ErrorCheck(t, pinpoint.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSPinpointSMSChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSPinpointSMSChannelConfig_full(senderId, shortCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSPinpointSMSChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "sender_id", senderId),
					resource.TestCheckResourceAttr(resourceName, "short_code", shortCode),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "promotional_messages_per_second"),
					resource.TestCheckResourceAttrSet(resourceName, "transactional_messages_per_second"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// There can be a delay before these Computed values are returned
				// e.g. 0 on Create -> Read, 20 on Import
				// These seem non-critical for other Terraform resource references,
				// so ignoring them for now, but we can likely adjust the Read function
				// to wait until they are available on creation with retry logic.
				ImportStateVerifyIgnore: []string{
					"promotional_messages_per_second",
					"transactional_messages_per_second",
				},
			},
			{
				Config: testAccAWSPinpointSMSChannelConfig_full(senderId, newShortCode),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSPinpointSMSChannelExists(resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, "sender_id", senderId),
					resource.TestCheckResourceAttr(resourceName, "short_code", newShortCode),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "promotional_messages_per_second"),
					resource.TestCheckResourceAttrSet(resourceName, "transactional_messages_per_second"),
				),
			},
		},
	})
}

func TestAccAWSPinpointSMSChannel_disappears(t *testing.T) {
	var channel pinpoint.SMSChannelResponse
	resourceName := "aws_pinpoint_sms_channel.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheckAWSPinpointApp(t) },
		ErrorCheck:   acctest.ErrorCheck(t, pinpoint.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckAWSPinpointSMSChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSPinpointSMSChannelConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSPinpointSMSChannelExists(resourceName, &channel),
					acctest.CheckResourceDisappears(acctest.Provider, ResourceSMSChannel(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckAWSPinpointSMSChannelExists(n string, channel *pinpoint.SMSChannelResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pinpoint SMS Channel with that application ID exists")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).PinpointConn

		// Check if the app exists
		params := &pinpoint.GetSmsChannelInput{
			ApplicationId: aws.String(rs.Primary.ID),
		}
		output, err := conn.GetSmsChannel(params)

		if err != nil {
			return err
		}

		*channel = *output.SMSChannelResponse

		return nil
	}
}

const testAccAWSPinpointSMSChannelConfig_basic = `
resource "aws_pinpoint_app" "test_app" {}

resource "aws_pinpoint_sms_channel" "test" {
  application_id = aws_pinpoint_app.test_app.application_id
}
`

func testAccAWSPinpointSMSChannelConfig_full(senderId, shortCode string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test_app" {}

resource "aws_pinpoint_sms_channel" "test" {
  application_id = aws_pinpoint_app.test_app.application_id
  enabled        = "false"
  sender_id      = "%s"
  short_code     = "%s"
}
`, senderId, shortCode)
}

func testAccCheckAWSPinpointSMSChannelDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).PinpointConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_pinpoint_sms_channel" {
			continue
		}

		// Check if the event stream exists
		params := &pinpoint.GetSmsChannelInput{
			ApplicationId: aws.String(rs.Primary.ID),
		}
		_, err := conn.GetSmsChannel(params)
		if err != nil {
			if tfawserr.ErrMessageContains(err, pinpoint.ErrCodeNotFoundException, "") {
				continue
			}
			return err
		}
		return fmt.Errorf("SMS Channel exists when it should be destroyed!")
	}

	return nil
}
