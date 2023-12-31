Fantom Explorer
===============

This is a simple explorer for Fantom Opera network. It scans the network and stores the latest blocks into
buffer. It provides a graphql interface to query the data.

# Building and Running

## Requirements

For building/running the project, the following tools are required:
* Go: version 1.20 or later; we recommend to use your system's package manager; alternatively, you can follow Go's
[installation manual](https://go.dev/doc/install) or; if you need to maintain multiple versions,
[this tutorial](https://go.dev/doc/manage-install) describes how to do so.
* MongoDB: version 6.0 or later; to install and run it, follow the instructions on the
[MongoDB website](https://docs.mongodb.com/manual/installation/) or; run it in a Docker container as described
below.

## Building

To build the project, run
```
make
```

To run tests, use
```
make test
```
To clean up a build, use `make clean`.

## Running

To run Explorer, you can run the `demonet-explorer` executable created by the build process:
```
build/demonet-explorer <cmd> <args...>
```
To list the available commands, run
```
build/demonet-explorer
```

## Example config
```
{
  "explorer": {
    "blockBufferSize": 10000,
    "isPersisted": false,
    "maxTxsCount": 10000000
  },
  "faucet": {
    "claimLimitSeconds": 86400,
    "claimTokensAmount": 2.5,
    "walletPrivateKey": "904d5dea0bdffb09d78a81c15f0b3b893f504679eb8cd1de585309cad58e6285",
    "claimsPerDay": 5,
    "erc20MintAmountHex": "0x8ac7230489e80000",
    "Erc20sPath": "erc20s.json"
  },
  "maze": {
    "visibilityRange": 3,
    "configPaths": ["path1.json", "path2.json"]
  },
  "metaFetcher": {
    "numberOfAccountsUrl": "number-of-accounts-url",
    "diskSizePer100MTxsUrl": "disk-size-url",
    "diskSizePrunedPer100MTxsUrl": "disk-size-pruned-url",
    "timeToFinalityUrl": "time-to-finality-url",
    "isIdleStatusUrl": "is-idle-status-url"
  },
  "rpc": {
    "operaRpcUrl": "https://rpcapi.fantom.network",
    "sfcAddress": "0xFC00FACE00000000000000000000000000000000"
  },
  "api": {
    "readTimeout": 2,
    "writeTimeout": 15,
    "idleTimeout": 1,
    "headerTimeout": 1,
    "resolverTimeout": 30,
    "bindAddress": "localhost:16761",
    "domainAddress": "localhost:16761",
    "corsOrigin": ["*"]
  },
  "logger": {
    "loggingLevel": 4,
    "logFormat": "%{color}%{level:-8s} %{shortpkg}/%{shortfunc}%{color:reset}: %{message}"
  },
  "mongodb": {
    "host": "localhost",
    "port": 27017,
    "db": "demonet-explorer",
    "user": null,
    "password": null
  }
}
```

## Local development
In case you need to run mongodb in docker container locally for development, you can use the following command:
```
docker run --name demonet-explorer-mongo -p 27017:27017 -e MONGO_INITDB_DATABASE=demonet-explorer mongo:6.0
```
This configuration won't persist data. You need to add volume to the container to do so.