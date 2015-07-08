#!/bin/bash
export GAE_LOCAL_VM_RUNTIME=1
gcloud --verbosity debug preview app run app.yaml --enable-mvm-logs
