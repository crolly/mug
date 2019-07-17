# Debug your Code

## How It Works

You can use the `mug debug` command to locally run the generated functions or whatever modifications you have made. The command will simply do the following:

1. Generate a Makefile for the resource or function group.
2. Build binaries in the `debug` folder using `make debug` defined in the `Makefile`.
3. Generate a `template.yml` later required by **aws-sam-cli** to provide the API Gateway.
4. Create a local docker network for **aws-sam** and **dynamodb** to talk to each other.
5. Start/ Restart a daemon of `amazon/dynamodb-local` container.
6. Create the tables in the database.
7. Start the API with `sam local start-api`. 


## Step through Code

To be able to debug your code step by step, you have to define the `-r` **remoteDebugger** flag. This tells **mug** to initiate a **delve** debug process. This by default runs on port **5986**, however you can overwrite this by defining the `-p` **debugPort** flag.

When you start with the remote debugger enabled you can easily step through your functions using breakpoints etc. Make sure you have a propper launch configuration beforehand. For Visual Studio Code it may look like this:

```json
{
    "version": "0.2.0",
    "configurations": [
    {
        "name": "lambda debug",
        "type": "go",
        "request": "launch",
        "mode": "remote",
        "remotePath": "",
        "port": 5986,
        "host": "127.0.0.1",
        "program": "${file}",
        "env": {},
        "args": [],
      },
    ]
  }
```

### Debugging bigger projects

In case your projects get a little bigger, you may have no interest in generating, building and running a debug process for the entire project, but only for a couple of your resources or function groups.

You can do so by specifying a comma separated list of resources/ function groups with the `-l` **list** flag. **mug** will then only run the generation and build process for those defined.