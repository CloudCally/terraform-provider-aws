package cognitoidentity_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccAWSCognitoIdentityPool_basic(t *testing.T) {
	name := sdkacctest.RandString(10)
	updatedName := sdkacctest.RandString(10)
	resourceName := "aws_cognito_identity_pool.main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "cognito-identity", regexp.MustCompile(`identitypool/.+`)),
					resource.TestCheckResourceAttr(resourceName, "allow_unauthenticated_identities", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolConfig_basic(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", updatedName)),
				),
			},
		},
	})
}

func TestAccAWSCognitoIdentityPool_supportedLoginProviders(t *testing.T) {
	name := sdkacctest.RandString(10)
	resourceName := "aws_cognito_identity_pool.main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolConfig_supportedLoginProviders(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "supported_login_providers.graph.facebook.com", "7346241598935555"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolConfig_supportedLoginProvidersModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "supported_login_providers.graph.facebook.com", "7346241598935552"),
					resource.TestCheckResourceAttr(resourceName, "supported_login_providers.accounts.google.com", "123456789012.apps.googleusercontent.com"),
				),
			},
			{
				Config: testAccPoolConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
				),
			},
		},
	})
}

func TestAccAWSCognitoIdentityPool_openidConnectProviderArns(t *testing.T) {
	name := sdkacctest.RandString(10)
	resourceName := "aws_cognito_identity_pool.main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolConfig_openidConnectProviderARNs(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_provider_arns.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolConfig_openidConnectProviderARNsModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_provider_arns.#", "2"),
				),
			},
			{
				Config: testAccPoolConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
				),
			},
		},
	})
}

func TestAccAWSCognitoIdentityPool_samlProviderArns(t *testing.T) {
	name := sdkacctest.RandString(10)
	resourceName := "aws_cognito_identity_pool.main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolConfig_samlProviderARNs(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "saml_provider_arns.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolConfig_samlProviderARNsModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "saml_provider_arns.#", "1"),
				),
			},
			{
				Config: testAccPoolConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "saml_provider_arns.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSCognitoIdentityPool_cognitoIdentityProviders(t *testing.T) {
	name := sdkacctest.RandString(10)
	resourceName := "aws_cognito_identity_pool.main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolConfig_cognitoIdentityProviders(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cognito_identity_providers.*", map[string]string{
						"client_id":               "7lhlkkfbfb4q5kpp90urffao",
						"provider_name":           fmt.Sprintf("cognito-idp.%[1]s.%[2]s/%[1]s_Zr231apJu", acctest.Region(), acctest.PartitionDNSSuffix()),
						"server_side_token_check": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cognito_identity_providers.*", map[string]string{
						"client_id":               "7lhlkkfbfb4q5kpp90urffao",
						"provider_name":           fmt.Sprintf("cognito-idp.%[1]s.%[2]s/%[1]s_Ab129faBb", acctest.Region(), acctest.PartitionDNSSuffix()),
						"server_side_token_check": "false",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolConfig_cognitoIdentityProvidersModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cognito_identity_providers.*", map[string]string{
						"client_id":               "6lhlkkfbfb4q5kpp90urffae",
						"provider_name":           fmt.Sprintf("cognito-idp.%[1]s.%[2]s/%[1]s_Zr231apJu", acctest.Region(), acctest.PartitionDNSSuffix()),
						"server_side_token_check": "false",
					}),
				),
			},
			{
				Config: testAccPoolConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
				),
			},
		},
	})
}

func TestAccAWSCognitoIdentityPool_addingNewProviderKeepsOldProvider(t *testing.T) {
	name := sdkacctest.RandString(10)
	resourceName := "aws_cognito_identity_pool.main"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolConfig_cognitoIdentityProviders(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "cognito_identity_providers.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolConfig_cognitoIdentityProvidersAndOpenidConnectProviderARNs(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "cognito_identity_providers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_provider_arns.#", "1"),
				),
			},
			{
				Config: testAccPoolConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "identity_pool_name", fmt.Sprintf("identity pool %s", name)),
					resource.TestCheckResourceAttr(resourceName, "cognito_identity_providers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "openid_connect_provider_arns.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSCognitoIdentityPool_tags(t *testing.T) {
	name := sdkacctest.RandString(10)
	resourceName := "aws_cognito_identity_pool.main"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); testAccPreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, cognitoidentity.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolConfig_Tags1(name, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPoolConfig_Tags2(name, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccPoolConfig_Tags1(name, "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckPoolExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Cognito Identity Pool ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CognitoIdentityConn

		_, err := conn.DescribeIdentityPool(&cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(rs.Primary.ID),
		})

		return err
	}
}

func testAccCheckPoolDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CognitoIdentityConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cognito_identity_pool" {
			continue
		}

		_, err := conn.DescribeIdentityPool(&cognitoidentity.DescribeIdentityPoolInput{
			IdentityPoolId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if wserr, ok := err.(awserr.Error); ok && wserr.Code() == cognitoidentity.ErrCodeResourceNotFoundException {
				return nil
			}
			return err
		}
	}

	return nil
}

func testAccPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CognitoIdentityConn

	input := &cognitoidentity.ListIdentityPoolsInput{
		MaxResults: aws.Int64(int64(1)),
	}

	_, err := conn.ListIdentityPools(input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccPoolConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false
  developer_provider_name          = "my.developer"
}
`, name)
}

func testAccPoolConfig_supportedLoginProviders(name string) string {
	return fmt.Sprintf(`
resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false

  supported_login_providers = {
    "graph.facebook.com" = "7346241598935555"
  }
}
`, name)
}

func testAccPoolConfig_supportedLoginProvidersModified(name string) string {
	return fmt.Sprintf(`
resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false

  supported_login_providers = {
    "graph.facebook.com"  = "7346241598935552"
    "accounts.google.com" = "123456789012.apps.googleusercontent.com"
  }
}
`, name)
}

func testAccPoolConfig_openidConnectProviderARNs(name string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false

  openid_connect_provider_arns = ["arn:${data.aws_partition.current.partition}:iam::123456789012:oidc-provider/server.example.com"]
}
`, name)
}

func testAccPoolConfig_openidConnectProviderARNsModified(name string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false

  openid_connect_provider_arns = ["arn:${data.aws_partition.current.partition}:iam::123456789012:oidc-provider/foo.example.com", "arn:${data.aws_partition.current.partition}:iam::123456789012:oidc-provider/bar.example.com"]
}
`, name)
}

func testAccPoolConfig_samlProviderARNs(name string) string {
	return fmt.Sprintf(`
resource "aws_iam_saml_provider" "default" {
  name                   = "myprovider-%[1]s"
  saml_metadata_document = file("./test-fixtures/saml-metadata.xml")
}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %[1]s"
  allow_unauthenticated_identities = false

  saml_provider_arns = [aws_iam_saml_provider.default.arn]
}
`, name)
}

func testAccPoolConfig_samlProviderARNsModified(name string) string {
	return fmt.Sprintf(`
resource "aws_iam_saml_provider" "default" {
  name                   = "default-%[1]s"
  saml_metadata_document = file("./test-fixtures/saml-metadata.xml")
}

resource "aws_iam_saml_provider" "secondary" {
  name                   = "secondary-%[1]s"
  saml_metadata_document = file("./test-fixtures/saml-metadata.xml")
}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %[1]s"
  allow_unauthenticated_identities = false

  saml_provider_arns = [aws_iam_saml_provider.secondary.arn]
}
`, name)
}

func testAccPoolConfig_cognitoIdentityProviders(name string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false

  cognito_identity_providers {
    client_id               = "7lhlkkfbfb4q5kpp90urffao"
    provider_name           = "cognito-idp.${data.aws_region.current.name}.${data.aws_partition.current.dns_suffix}/${data.aws_region.current.name}_Ab129faBb"
    server_side_token_check = false
  }

  cognito_identity_providers {
    client_id               = "7lhlkkfbfb4q5kpp90urffao"
    provider_name           = "cognito-idp.${data.aws_region.current.name}.${data.aws_partition.current.dns_suffix}/${data.aws_region.current.name}_Zr231apJu"
    server_side_token_check = false
  }
}
`, name)
}

func testAccPoolConfig_cognitoIdentityProvidersModified(name string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false

  cognito_identity_providers {
    client_id               = "6lhlkkfbfb4q5kpp90urffae"
    provider_name           = "cognito-idp.${data.aws_region.current.name}.${data.aws_partition.current.dns_suffix}/${data.aws_region.current.name}_Zr231apJu"
    server_side_token_check = false
  }
}
`, name)
}

func testAccPoolConfig_cognitoIdentityProvidersAndOpenidConnectProviderARNs(name string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_region" "current" {}

resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "identity pool %s"
  allow_unauthenticated_identities = false

  cognito_identity_providers {
    client_id               = "7lhlkkfbfb4q5kpp90urffao"
    provider_name           = "cognito-idp.${data.aws_region.current.name}.${data.aws_partition.current.dns_suffix}/${data.aws_region.current.name}_Ab129faBb"
    server_side_token_check = false
  }

  cognito_identity_providers {
    client_id               = "7lhlkkfbfb4q5kpp90urffao"
    provider_name           = "cognito-idp.${data.aws_region.current.name}.${data.aws_partition.current.dns_suffix}/${data.aws_region.current.name}_Zr231apJu"
    server_side_token_check = false
  }

  openid_connect_provider_arns = ["arn:${data.aws_partition.current.partition}:iam::123456789012:oidc-provider/server.example.com"]
}
`, name)
}

func testAccPoolConfig_Tags1(name, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = %q
  allow_unauthenticated_identities = false

  tags = {
    %q = %q
  }
}
`, name, tagKey1, tagValue1)
}

func testAccPoolConfig_Tags2(name, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = %q
  allow_unauthenticated_identities = false

  tags = {
    %q = %q
    %q = %q
  }
}
`, name, tagKey1, tagValue1, tagKey2, tagValue2)
}
