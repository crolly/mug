# Deployment

The deployment uses the **serverless** framework. If you require some help setting up serverless, [follow the instructions here](https://serverless.com/blog/anatomy-of-a-serverless-app/#setup).

## Deploy to AWS

To deploy your application just run `mug deploy`. This will do the following:

1. Generate a Makefile for the resource or function group.
2. Build binaries in the `bin` folder using `make build` defined in the `Makefile`.
3. Generate a `serverless.yml`.
4. Run `sls deploy` from each directory of resources or function groups. 

**This will deploy your app to AWS and you can now develop against your new serverless API! Yeah!**

::: tip
Just like `mug debug` you can define a list of resources/ function groups, you wish to deploy, in case you do not want to deploy all of them. Just set the `-l` **list** flag and provide a comma separated list (e.g. `mug deploy -l "user,course"` which will only deploy the **user** and **course** resources).
:::

## Deploying with secrets

In case you have environment variables, you want to have added to your `serverless.yml` especially for those, you may not want to share in your git repository, you can easily create a `secrets.yml` file for that resource/ function group (where `serverless.yml` file is), which will be parsed during creation/ update of the `serverless.yml`.
The environments will then be loaded to the environments section.

For example:
```
API_KEY = Example
COGNITO_USERPOOL_ARN = 'arn:aws:cognito-idp:eu-central-1:XXXXXXXXXXXX:userpool/eu-central-1_XXXXXXXXX'
```
will result in an `serverless.yml` like this:
```yaml
service: example

provider:
  name: aws
  runtime: go1.x
  region: "eu-central-1"
  stage: ${opt:stage, 'dev'}
  environment:
      API_KEY: ${file(./secrets.yml):API_KEY}
      COGNITO_USERPOOL_ARN: ${file(./secrets.yml):COGNITO_USERPOOL_ARN}
```

::: tip
Simply adding your `secrets.yml` to your `.gitignore` would be a good practice to pass environment variables to your lambda function without exposing them to the world.
:::