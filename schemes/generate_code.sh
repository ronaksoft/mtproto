#!/bin/sh

go run build_tl_scheme.go < tl-schema-57.json  > ../tl_schema.go
gofmt -w ../tl_schema.go
