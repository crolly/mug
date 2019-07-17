# Introduction

<p align="center"><img src="../logo.svg" width="480px" /></p>

## Idea

**mug** lets you create AWS Lambda for **go** projects directly integration **DynamoDB** if you wish so.

It boilerplates the project structure with a serverless framework configuration. Additionally a resource definition for Dynamo DB is generated, which creates a table for each resource to be added.

## How It Works

In general, there are three units that can be added to a project:

* Resource
* Function Group
* Function

### Resource

A resource is defined as a model which implements basic CRUDL (Create, Read, Update, Delete, List) functions. Each resource will have its own DynamoDB table to persist the data. The model itself will be a go struct and can have nested structs as well, if defined upon creation.

The resource files will be stored in a directory named after the resource in the project's **functions** directory. It will look something like this e.g.:

```
-rw-r--r-- course.json
-rw-r--r-- course.go
-rw-r--r-- serverless.yml
-rw-r--r-- create/main.go
-rw-r--r-- delete/main.go
-rw-r--r-- list/main.go
-rw-r--r-- read/main.go
-rw-r--r-- update/main.go
```

* `course.json` holds the resource configuration (e.g. attributes, key schema etc.)
* `course.go` is the model class containing the go struct and the **CRUDL** functions interacting with DynamoDB
* `serverless.yml` has the configuration for the deployment with serverless framework
* `main.go` files in the subdirectories are the Lambda handler functions calling the model`s methods

### Function Group

A function group is - as the name might suggest - a group of functions. Unlike resource, it does not define a model or **CRUDL** functions but basically just helps to organize smaller functions. **mug** will generate a `serverless.yml` file for the deployment.

### Function

A function is a Lambda handler which can be added to either a *resource* or a *function group*.
**mug** will be able to deploy either a function group or a resource, so in case, you just want to deploy a single function, **mug** might not be the tool you are looking for, however, you would still accomplish this by first adding a function group and then adding a function to that group.