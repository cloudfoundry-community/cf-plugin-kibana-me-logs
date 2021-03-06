#!/bin/bash

set -x
set +e
cf delete-orphaned-routes -f

cf delete-space -f cf-plugin-kibana-me-logs-test
set -e

cf create-space cf-plugin-kibana-me-logs-test
cf target -s cf-plugin-kibana-me-logs-test

set +e
# use run.sh first time to rebuild/reinstall plugin
./bin/run.sh one
set -e

if [[ ! -d tmp/simple-go-web-app ]]; then
  mkdir -p tmp
  git clone https://github.com/cloudfoundry-community/simple-go-web-app tmp/simple-go-web-app
fi

cd tmp/simple-go-web-app
cf push one --no-start
cf set-env one MESSAGE "I am one"
cf start one
cd ../..

set +e
echo "Try to open kibana; except app not bound to logstash"
cf kibana-me-logs one
set -e

echo Bind app to logstash
cf cs logstash14 free logstash-one

cf bs one logstash-one
cf restart one

echo "Should auto-deploy kibana UI since it isn't already running"
cf kibana-me-logs one


# Try a 2nd logstash/kibana
cf cs logstash14 free logstash-two

cd tmp/simple-go-web-app
cf push dedicated-logs --no-start
cf bs dedicated-logs logstash-two
cf set-env dedicated-logs MESSAGE "I have a dedicated logstash"
cf start dedicated-logs

echo "Auto-deploy kibana for 2nd logstash14 (via the app)"
cf kibana-me-logs dedicated-logs


cf push two --no-start
cf bs two logstash-one
cf set-env two MESSAGE "I am two"
cf start two

cf push three --no-start
cf bs three logstash-one
cf set-env three MESSAGE "I am three"
cf start three

set +e
# install `open` plugin
cf open one
cf open two
cf kibana-me-logs two
cf open three
cf kibana-me-logs three
cf open dedicated-logs
cf kibana-me-logs dedicated-logs

cd ../..
