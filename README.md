# Bridge Description Language

Specification and interpreter for TriggerMesh's Bridge Description Language.

The Bridge Description Language is a configuration language based on the [HCL syntax][hcl-spec] which purpose is to
offer a user-friendly interface for describing [TriggerMesh Bridges][tm-brg].

## Documentation

Please refer to the [`docs/`](./docs) directory for details about the language specification and other technical
documents.

## Getting started

The interpreter is written in the [Go programming language][go] and leverages the [HCL toolkit][hcl].

After installing the [Go toolchain][go-dl] (version 1.16 or above), the interpreter can be compiled for the current OS
and architecture by executing the following command inside the root of the repository (the `main` Go package):

```console
$ go build .
```

The above command creates an executable called `bridgedl` inside the current directory.

The `-h` or `--help` flag can be appended to any command or subcommand to print some usage instructions about that
command:

```console
$ ./bridgedl --help
```

[tm-brg]: https://www.triggermesh.com/integrations

[go]: https://golang.org/
[go-dl]: https://golang.org/dl/

[hcl]: https://github.com/hashicorp/hcl
[hcl-spec]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md
