# Infra as code on Azure using Pulumi and Go

This is an example of how to create the Azure infrastructure services to deploy a backend for Slack app. The details of the app have been covered in this blog - https://itnext.io/build-a-serverless-app-using-go-and-azure-functions-c4475398f4ab

Instead of Azure Functions (as covered in the aforementioned blog), we deploy the same backend code to Azure App Service as a Docker container - https://hub.docker.com/repository/docker/abhirockzz/funcy-go

It uses the [Pulumi Microsoft Azure provider](https://www.pulumi.com/docs/intro/cloud-providers/azure/) with [Go SDK](https://www.pulumi.com/docs/intro/languages/go/) support.

> Before you begin, install and configure Pulumi - https://www.pulumi.com/docs/get-started/azure/begin/

## Configure

Store Slack Signing Secret and GIPHY API Key as secrets (in Pulumi config)

```bash
pulumi config set --secret giphyapikey <api key>
pulumi config set --secret slacksecret <slack secret>
```

After that:

```bash
pulumi up
```

You should have the application public endpoint. Configure that in Slack and continue using it as specified in this tutorial - https://itnext.io/build-a-serverless-app-using-go-and-azure-functions-c4475398f4ab

To clean up:

```bash
pulumi destroy
```
