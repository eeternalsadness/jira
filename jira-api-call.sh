#!/usr/bin/env bash

if [[ $# -eq 1 ]]; then
  method="GET"
  api_path="$1"
elif [[ $# -eq 2 ]]; then
  method="$1"
  api_path="$2"
else
  echo "Must have 1 or 2 arguments!"
  exit 1
fi

url="https://${JIRA_DOMAIN}/rest/api/3/${api_path}"
auth=$(echo -n "${JIRA_EMAIL}:${JIRA_API_TOKEN}" | base64)
echo "$url"

curl -fsSL -X "$method" \
  -H "Authorization: Basic $auth" \
  -H "Accept: application/json" \
  "$url" |
  jq
