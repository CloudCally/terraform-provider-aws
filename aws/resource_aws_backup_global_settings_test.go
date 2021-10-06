package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

func TestAccAwsBackupGlobalSettings_basic(t *testing.T) {
	var settings backup.DescribeGlobalSettingsOutput

	resourceName := "aws_backup_global_settings.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			testAccPreCheckAWSBackup(t)
		},
		ErrorCheck:   acctest.ErrorCheck(t, backup.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBackupGlobalSettingsConfig("true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupGlobalSettingsExists(&settings),
					resource.TestCheckResourceAttr(resourceName, "global_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "global_settings.isCrossAccountBackupEnabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBackupGlobalSettingsConfig("false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupGlobalSettingsExists(&settings),
					resource.TestCheckResourceAttr(resourceName, "global_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "global_settings.isCrossAccountBackupEnabled", "false"),
				),
			},
			{
				Config: testAccBackupGlobalSettingsConfig("true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsBackupGlobalSettingsExists(&settings),
					resource.TestCheckResourceAttr(resourceName, "global_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "global_settings.isCrossAccountBackupEnabled", "true"),
				),
			},
		},
	})
}

func testAccCheckAwsBackupGlobalSettingsExists(settings *backup.DescribeGlobalSettingsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		conn := acctest.Provider.Meta().(*conns.AWSClient).BackupConn
		resp, err := conn.DescribeGlobalSettings(&backup.DescribeGlobalSettingsInput{})
		if err != nil {
			return err
		}

		*settings = *resp

		return nil
	}
}

func testAccBackupGlobalSettingsConfig(setting string) string {
	return fmt.Sprintf(`
resource "aws_backup_global_settings" "test" {
  global_settings = {
    "isCrossAccountBackupEnabled" = %[1]q
  }
}
`, setting)
}
