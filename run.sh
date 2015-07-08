#!/bin/bash
go build -work -x -v -o coduno
gcloud preview app run app.yaml
