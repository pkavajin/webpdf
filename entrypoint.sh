#!/bin/bash
set -e
set -o pipefail

wkhtmltopdf "${URL}" - | /callback