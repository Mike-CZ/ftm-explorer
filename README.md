Fantom Explorer
===============

This is a simple explorer for Fantom Opera network. It scans the network and stores the latest blocks into
buffer. It provides a graphql interface to query the data.

# Building and Running

## Requirements

For building/running the project, the following tools are required:
* Go: version 1.20 or later; we recommend to use your system's package manager; alternatively, you can follow Go's
[installation manual](https://go.dev/doc/install) or; if you need to maintain multiple versions,
[this tutorial](https://go.dev/doc/manage-install) describes how to do so


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

To run Explorer, you can run the `ftm-explorer` executable created by the build process:
```
build/ftm-explorer <cmd> <args...>
```
To list the available commands, run
```
build/ftm-explorer