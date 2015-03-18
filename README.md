Overview
========

Launches the Kibana UI (from [kibana-me-logs](https://github.com/cloudfoundry-community/kibana-me-logs)\) for an application.

Installation
------------

```
$ go get github.com/cloudfoundry-community/cf-plugin-kibana-me-logs
$ cf install-plugin $GOPATH/bin/cf-plugin-kibana-me-logs
```

Usage
-----

```
$ cf kibana-me-logs <kibana-app> <app>
```

Will launch the Kibana UI and show logs for the requested app.

It assumes that `<kibana-app>` is the [kibana-me-logs](https://github.com/cloudfoundry-community/kibana-me-logs)\), and both are bound to the same `logstash14` logstash service.

Development
-----------

```
cf uninstall-plugin kibana-me-logs
go get ./...
cf install-plugin $GOPATH/bin/cf-plugin-kibana-me-logs
```

Or a one-liner:

```
cf uninstall-plugin kibana-me-logs; go get ./...; cf install-plugin $GOPATH/bin/cf-plugin-kibana-me-logs
```
