# mug remove

Remove resource or function group from your project

### Synopsis

Remove resource or function group from your project

```
mug remove name [flags]
```

### Options

```
  -h, --help   help for remove
```

## mug remove function

Removes a function from a resource/ function group

### Synopsis

Removes a function from a resource/ function group

```
mug remove function functionName [flags]
```

### Options

```
  -h, --help                help for function
  -a, --assignedTo string   Name of the resource or the function group the function
                            was assigned to (default "generic")
```

## mug remove auth

Remove authentication from the given resource or function group

### Synopsis

Remove authentication from the given resource or function group

```
mug remove auth resourceName [flags]
```

### Options

```
  -h, --help   help for auth
```