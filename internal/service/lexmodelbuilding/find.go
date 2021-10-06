package lexmodelbuilding

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tflex "github.com/hashicorp/terraform-provider-aws/aws/internal/service/lex"
	"github.com/hashicorp/terraform-provider-aws/aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	tflexmodelbuilding "github.com/hashicorp/terraform-provider-aws/internal/service/lexmodelbuilding"
	tflexmodelbuilding "github.com/hashicorp/terraform-provider-aws/internal/service/lexmodelbuilding"
)

func FindBotVersionByName(conn *lexmodelbuildingservice.LexModelBuildingService, name, version string) (*lexmodelbuildingservice.GetBotOutput, error) {
	input := &lexmodelbuildingservice.GetBotInput{
		Name:           aws.String(name),
		VersionOrAlias: aws.String(version),
	}

	output, err := conn.GetBot(input)

	if tfawserr.ErrCodeEquals(err, lexmodelbuildingservice.ErrCodeNotFoundException) {
		return nil, &resource.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

// FindLatestBotVersionByName returns the latest published version of a bot or $LATEST if the bot has never been published.
// See https://docs.aws.amazon.com/lex/latest/dg/versioning-aliases.html.
func FindLatestBotVersionByName(conn *lexmodelbuildingservice.LexModelBuildingService, name string) (string, error) {
	input := &lexmodelbuildingservice.GetBotVersionsInput{
		Name: aws.String(name),
	}
	var latestVersion int

	err := conn.GetBotVersionsPages(input, func(page *lexmodelbuildingservice.GetBotVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, bot := range page.Bots {
			version := aws.StringValue(bot.Version)

			if version == tflexmodelbuilding.BotVersionLatest {
				continue
			}

			if version, err := strconv.Atoi(version); err != nil {
				continue
			} else if version > latestVersion {
				latestVersion = version
			}
		}

		return !lastPage
	})

	if err != nil {
		return "", err
	}

	if latestVersion == 0 {
		return tflexmodelbuilding.BotVersionLatest, nil
	}

	return strconv.Itoa(latestVersion), nil
}

// FindLatestIntentVersionByName returns the latest published version of an intent or $LATEST if the intent has never been published.
// See https://docs.aws.amazon.com/lex/latest/dg/versioning-aliases.html.
func FindLatestIntentVersionByName(conn *lexmodelbuildingservice.LexModelBuildingService, name string) (string, error) {
	input := &lexmodelbuildingservice.GetIntentVersionsInput{
		Name: aws.String(name),
	}
	var latestVersion int

	err := conn.GetIntentVersionsPages(input, func(page *lexmodelbuildingservice.GetIntentVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, intent := range page.Intents {
			version := aws.StringValue(intent.Version)

			if version == tflexmodelbuilding.IntentVersionLatest {
				continue
			}

			if version, err := strconv.Atoi(version); err != nil {
				continue
			} else if version > latestVersion {
				latestVersion = version
			}
		}

		return !lastPage
	})

	if err != nil {
		return "", err
	}

	if latestVersion == 0 {
		return tflexmodelbuilding.IntentVersionLatest, nil
	}

	return strconv.Itoa(latestVersion), nil
}
