[![GoDoc](https://godoc.org/github.com/natural/missmolly?status.svg)](https://godoc.org/github.com/natural/missmolly)

# Miss Molly #

Miss Molly is a web application server and library, written in Go, featuring
a friendly YAML config syntax and scritable request processing via an embedded
Lua interpreter.

## Quick Start ##

```sh
$ go get github.com/natural/missmolly
$ ./bin/missmolly
```


## Installation ##

First thing, make sure you've got Go installed.  Then you can `go get`
Miss Molly with a shell command.  For example, you can make a
workspace for a sample project, initialize it, and fetch MM:

```sh
$ mkdir example-workspace
$ cd example-workspace
$ export GOPATH=`pwd`
$ go get github.com/natural/missmolly
```

After the `go get` we've got the server in `./bin/missmolly` but we
need to make a config file before we can run it.  Continuing with the
shell commands above:

```sh
$ cat <<EOF >example-config.yaml
- location: /
  content: >-
    request:write('Hello, world.')
EOF
```

Now we can start the server:

```sh
$ ./bin/missmolly run example-config.yaml
```

(this section (and the rest of this README) are incomplete)
