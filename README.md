# Impack

Imperfect memory packer for Go

## Features

`impack` will grab every struct in a package and order the fields by lower size to higher size. When sizes match it will order alphabetically.

This is an imperfect memory packer since it does not maximize memory efficiency by packing structs perfectly, it attempts to minimize memory footprint but trying to remain as human readable as possible.