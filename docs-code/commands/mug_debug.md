# mug debug

Start Local API for debugging

### Synopsis

This command generates a template.yml for aws-sam-cli and starts a local api to test or debug against

```
mug debug [flags]
```

### Options

```
  -h, --help               help for debug
  -d, --debugPort string   defines the remote port (default "5986")
  -g, --gwPort string      defines the port of local API Gateway (default "3000")
  -l, --list string        list of resources/ function groups to debug (default "all")
  -r, --remoteDebugger     indicates whether you want to run a remote debugger
                           (e.g. step through your code with VSCode)
```