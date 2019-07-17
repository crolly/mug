# Getting Started

## Requirements

mug riffs off [nerdguru's go-sls-crudl](https://github.com/nerdguru/go-sls-crudl) combining the [Dynamo DB Golang samples](https://github.com/awsdocs/aws-doc-sdk-examples/tree/master/go/example_code/dynamodb) and the [Serverless Framework Go example](https://serverless.com/blog/framework-example-golang-lambda-support/).
In order to be able to deploy the generated code please make sure to have the following in place/ installed:


* AWS Account (duh!)
* golang (double duh!)
* [serverless framework](https://serverless.com/framework/docs/getting-started/)
* [dep](https://golang.github.io/dep/) (golang dependancy management tool)
* [aws-cli](https://docs.aws.amazon.com/de_de/cli/latest/userguide/cli-chap-welcome.html) wouldn't hurt, but is not necessarily required
* [aws-sam-cli](https://github.com/awslabs/aws-sam-cli) in case you want to locally debug your code

## Installation

To get mug just run
```
go get github.com/crolly/mug
```
This will create the cobra executable under your `$GOPATH/bin` directory. Make sure to have that path in your `$PATH` configuration.

## Create a Project

The create command generates the boilerplate for project.
```
mug create projectname [flags]
```
This will create the directory in case it doesn't already exist. In case you want to overwrite any existing directory, you can add the `-f` flag to forcefully overwrite.

The structure will generally look like this:
* `mug.config.json` holds the project's configuration required for **mug** to work. 
::: warning 
Do not change the `mug.config.json` file, as it may break your project.
:::
* `Gopkg.toml` initializes the project with dependency management through **dep**

