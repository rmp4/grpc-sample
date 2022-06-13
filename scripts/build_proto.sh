#!/bin/bash
protoc --proto_path=./pkg/pb ./pkg/pb/src/*.proto --go_out=. 