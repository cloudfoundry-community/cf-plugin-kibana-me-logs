#!/bin/bash

set -x
set +e
cf delete-orphaned-routes -f

cf delete-space -f cf-plugin-kibana-me-logs-test
set -e

cf create-space cf-plugin-kibana-me-logs-test
cf target -s cf-plugin-kibana-me-logs-test

set +e
cf kibana-me-logs unknown
set -e

if [[ ! -d tmp/kibana-me-logs ]]; then
  mkdir -p tmp
  git clone https://github.com/cloudfoundry-community/kibana-me-logs tmp/kibana-me-logs
fi

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
cf kibana-me-logs one
set -e

cf cs logstash14 free logstash-one
cf cs logstash14 free logstash-two

cf bs one logstash-one
cf restart one
set +e
cf kibana-me-logs one
set -e


cd tmp/kibana-me-logs

cf push kibana-one --no-start
cf bs kibana-one logstash-one
cf start kibana-one

cf kibana-me-logs one


cf push kibana-two --no-start
cf bs kibana-two logstash-two
cf start kibana-two

cd ../..

cd tmp/simple-go-web-app

cf push two --no-start
cf bs two logstash-one
cf set-env two MESSAGE "I am two"
cf start two

cf push three --no-start
cf bs three logstash-one
cf set-env three MESSAGE "I am three"
cf start three

cf push dedicated-logs --no-start
cf bs dedicated-logs logstash-two
cf set-env dedicated-logs MESSAGE "I have a dedicated logstash"
cf start dedicated-logs

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
