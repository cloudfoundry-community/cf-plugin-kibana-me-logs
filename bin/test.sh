#!/bin/bash

set +e
cf delete-orphaned-routes -f

cf delete-space -f cf-plugin-kibana-me-logs-test
set -e

cf create-space cf-plugin-kibana-me-logs-test
cf target -s cf-plugin-kibana-me-logs-test

cf cs logstash14 free logstash-one
cf cs logstash14 free logstash-two


if [[ ! -d tmp/kibana-me-logs ]]; then
  mkdir -p tmp
  git clone https://github.com/cloudfoundry-community/kibana-me-logs tmp/kibana-me-logs
fi
cd tmp/kibana-me-logs

cf push kibana-one --no-start
cf bs kibana-one logstash-one
cf start kibana-one

cf push kibana-two --no-start
cf bs kibana-two logstash-two
cf start kibana-two

cd ../..

if [[ ! -d tmp/cf-env ]]; then
  mkdir -p tmp
  git clone https://github.com/cloudfoundry-community/cf-env tmp/cf-env
  cd ..
fi
cd tmp/cf-env

cf push cf-env-one --no-start
cf bs cf-env-one logstash-one
cf start cf-env-one

cf push cf-env-two --no-start
cf bs cf-env-two logstash-two
cf start cf-env-two

cf kibana-me-logs cf-env-one
cf kibana-me-logs cf-env-two
