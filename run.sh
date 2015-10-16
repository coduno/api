#!/bin/bash
set -x
dev_appserver.py --host 0.0.0.0 --port 8080 --admin_port 8000 --storage_path .storage --show_mail_body true $@ .
