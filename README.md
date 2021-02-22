# Impack

Imperfect memory packer for Go

## Features

`impack` will grab every struct in a package and order the fields by lower size to higher size. When sizes match it will order alphabetically.

This is an imperfect memory packer since it does not maximize memory efficiency by packing structs perfectly, it attempts to minimize memory footprint but trying to remain as human readable as possible.

## Install

```sh
go get -u github.com/javiercbk/impack/cmd/impack
```

## Usage

`impack` will lint a whole package, and re-order every struct's fields in the package.

```sh
impack --compiler gc --arch amd64 /home/user/path/to/go/package/folder
```

* The first unnammed parameter will be considered as the package path. By default it will be the current folder
* `--compiler` defaults to `gc`. Affects the size and the alignment values of types.
* `--arch` defaults to `amd64`. Also affects the size and the alignment values of types.