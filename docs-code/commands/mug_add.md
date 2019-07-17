# mug add

Add resources, function groups or functions to your project

### Synopsis

Add resources, function groups or functions to your project

### Options

```
  -h, --help   help for add
```

## mug add resource

Adds CRUDL functions for the defined resource

### Synopsis

Adds CRUDL functions for the defined resource

```
mug add resource name [flags]
```

### Options

```
  -h, --help                 help for resource
  -a, --attributes string    attributes of the resource
  -k, --keySchema string     Key Schema definition for the DynamoDB Table Resource
                             [only applied if noID flag is set to true]
                             (default "id:HASH")
  -d, --addDates             automatically add createdAt and updatedAt attributes
  -n, --noID                 disable automatic generation of id attribute
  -s, --softDelete           automatically add deletedAt attribute
  -b, --billingMode string   Choose between 'provisioned' or 'ondemand'
                             (default "provisioned")
  -r, --readUnits int        Set the ReadCapacityUnits if billingMode is set to 
                             ProvisionedThroughput (default 1)
  
  -w, --writeUnits int       Set the WriteCapacityUnits if billingMode is set to
                             ProvisionedThroughput (default 1)
```

## mug add functionGroup

Adds a new function group, you can then add functions to with 'mug add function -r name [flags]'

### Synopsis

Adds a new function group, you can then add functions to with 'mug add function -r name [flags]'

```
mug add functionGroup name [flags]
```

### Options

```
  -h, --help   help for functionGroup
```

## mug add function

Adds a function to a resource

### Synopsis

Adds a function to a resource

```
mug add function functionName [flags]
```

### Options

```
  -h, --help              help for function
  -a, --assign string     Name of the resource or function group the function should be
                          assigned to (default "generic")
  -m, --method string     Method the function will respond to e.g. get
  -p, --path string       Path the function will respond to e.g. /users
```

## mug add auth

Add authentication to a resource or function group

### Synopsis

Add authentication to a resource or function group

```
mug add auth [flags]
```

### Options

```
  -h, --help               help for auth
  -x, --excludes string    list of functions in resource/ function group
                           without authentication
  -p, --user pool string   define the user pool to authenticate against
```