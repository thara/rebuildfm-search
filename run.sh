#!/usr/bin/env bash
glide install -s -v
go-wrapper install
go-wrapper run $@
