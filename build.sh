#!/usr/bin/env bash

# Clear
rm -rf ./dist && mkdir ./dist

# Exercise the right to use and modify LGPL licensed software to hide the watermark
# Since there is no better way ATM if one decides to use it in open source software
# sigh...
find="lk != nil \&\& lk"
replace="lk == nil || !lk"
sed -i -e "s/$find/$replace/" ./vendor/github.com/unidoc/unipdf/v3/model/writer.go
find="lk == nil || !lk"
replace="lk != nil \&\& !lk"
sed -i -e "s/$find/$replace/" ./vendor/github.com/unidoc/unipdf/v3/model/writer.go

# Build the Go program
export GO111MODULE=on
go build -mod=vendor ./src/main.go

# Copy files to dist directory
cp backend.yml ./dist/
cp LICENSE ./dist/
cp ./main* ./dist
