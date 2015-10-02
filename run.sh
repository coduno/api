#!/bin/bash
set -x
dev_appserver.py --port 8080 --admin_port 8000 --storage_path .storage --show_mail_body true $@ .
