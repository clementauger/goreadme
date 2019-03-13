# goreadme

generate `doc.go` file from the `README.md`.

if a `doc.go` file is encountered bu it is not marked as `// goreadme autogen `, it its not rewriten.

otherwise, if a `README.md` file is found, within a golang `package main`, and that no `doc.go` file is present, of this last file is marked as `// goreadme autogen `, then the documentation file is generated automatically.

# install

```sh
go get github.com/clementauger/goreadme
```

# usage

```sh
$ goreadme -h
  -d	dry run
  -r	recursive lookup
  -v	verbose
```

# examples

```sh
goreadme -d -r
goreadme -d -v
goreadme -d -v -r
```
