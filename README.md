Overview
========

Launches the Kibana UI (from [kibana-me-logs](https://github.com/cloudfoundry-community/kibana-me-logs)\) for an application.

![demo](http://cl.ly/image/0N0J0h1Q2X2a/kibana-me-logs.gif)

Learn more about using CLI plugins from https://blog.starkandwayne.com/2015/03/04/installing-cloud-foundry-cli-plugins/

Installation
------------

From community plugin repo:

```
$ cf add-plugin-repo community http://plugins.cfapps.io/
$ cf install-plugin kibana-me-logs -r community
```

From source:

```
$ go get github.com/cloudfoundry-community/cf-plugin-kibana-me-logs
$ cf install-plugin $GOPATH/bin/cf-plugin-kibana-me-logs
```

Upgrade
-------

To upgrade you must first uninstall the plugin and then install as above:

```
$ cf uninstall-plugin kibana-me-logs
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

There is a helper to uninstall, build, install, and run the plugin:

```
./bin/run.sh <ARGS>
```

Or manually:

```
cf uninstall-plugin kibana-me-logs
go get ./...
cf install-plugin $GOPATH/bin/cf-plugin-kibana-me-logs
cf kibana-me-logs <ARGS>
```

Or a one-liner:

```
cf uninstall-plugin kibana-me-logs; go get ./...; cf install-plugin $GOPATH/bin/cf-plugin-kibana-me-logs && cf kibana-me-logs <ARGS>
```

### Bump version

There is a helper script to bump the version number `VERSION` file and regenerate the gobindata `version.go` file:

```
./bin/bump_version.sh X.Y.Z
```
