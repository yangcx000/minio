#!/bin/bash

protoc --proto_path=protos/ --go_out=plugins=grpc:protos --go_opt=paths=source_relative protos/*.proto
