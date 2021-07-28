package main

import (
	"fmt"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/web"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	resourceGroupName           string = "funcy-app-rg"
	plan                        string = "funcy-app-plan"
	planOS                             = "Linux"
	planSKUCode                        = "B1"
	planSKU                            = "Basic"
	dockerImage                        = "abhirockzz/funcy-go"
	appName                            = "funcy-api-backend"
	storageConfigName                  = "WEBSITES_ENABLE_APP_SERVICE_STORAGE"
	giphyAPIKeyAppConfigName           = "GIPHY_API_KEY"
	slackSecretAppConfigName           = "SLACK_SIGNING_SECRET"
	giphyAPIPulumiConfigName           = "giphyapikey"
	slackSecretPulumiConfigName        = "slacksecret"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		resourceGroup, err := resources.NewResourceGroup(ctx, resourceGroupName, nil)
		if err != nil {
			return err
		}

		appSvcPlan, err := web.NewAppServicePlan(ctx, plan, &web.AppServicePlanArgs{
			ResourceGroupName: resourceGroup.Name,
			Kind:              pulumi.String(planOS),
			Reserved:          pulumi.Bool(true),
			Sku: &web.SkuDescriptionArgs{
				Name: pulumi.String(planSKUCode),
				Tier: pulumi.String(planSKU),
			},
		})
		if err != nil {
			return err
		}

		cfg := config.New(ctx, "")
		giphyAPIKey := cfg.RequireSecret(giphyAPIPulumiConfigName)
		slackSecret := cfg.RequireSecret(slackSecretPulumiConfigName)

		helloApp, err := web.NewWebApp(ctx, appName, &web.WebAppArgs{
			ResourceGroupName: resourceGroup.Name,
			ServerFarmId:      appSvcPlan.ID(),
			SiteConfig: &web.SiteConfigArgs{
				AppSettings: web.NameValuePairArray{
					&web.NameValuePairArgs{
						Name:  pulumi.String(storageConfigName),
						Value: pulumi.String("false"),
					},
					&web.NameValuePairArgs{
						Name:  pulumi.String(giphyAPIKeyAppConfigName),
						Value: giphyAPIKey,
					},
					&web.NameValuePairArgs{
						Name:  pulumi.String(slackSecretAppConfigName),
						Value: slackSecret,
					},
				},
				AlwaysOn:       pulumi.Bool(true),
				LinuxFxVersion: pulumi.String(fmt.Sprintf("%v%v", "DOCKER|", dockerImage)),
			},
			HttpsOnly: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		ctx.Export("appurl", helloApp.DefaultHostName.ApplyT(func(defaultHostName string) (string, error) {
			return fmt.Sprintf("%v%v%v", "https://", defaultHostName, "/api/funcy"), nil
		}).(pulumi.StringOutput))

		return nil
	})
}
