Overview
========

Launches the Kibana UI (from [kibana-me-logs](https://github.com/cloudfoundry-community/kibana-me-logs)\) for an application.

![demo](http://cl.ly/image/2H1x2m1B3m0v/cf%20kibana-me-logs%20v0.3.gif)

Learn more about using CLI plugins from https://blog.starkandwayne.com/2015/03/04/installing-cloud-foundry-cli-plugins/

Installation
------------

From community plugin repo:

```
$ cf add-plugin-repo community http://plugins.cfapps.io/
$ cf install-plugin kibana-me-logs -r community
```

From source, see the Development section below.

Upgrade
-------

To upgrade you must first uninstall the plugin and then install as above:

```
$ cf uninstall-plugin kibana-me-logs
```

Usage
-----

```
$ cf kibana-me-logs <kibana-app> <app> [--no-auth]
```

Will launch the Kibana UI and show logs for the requested app.

It assumes that `<kibana-app>` is the [kibana-me-logs](https://github.com/cloudfoundry-community/kibana-me-logs)\), and both are bound to the same `logstash14` logstash service.

`cf-kibana-me-logs` will automatically generate a user/password for you to use in conjunction with your kibana-me-logs app, unless you provide the --no-auth option.

Development
-----------

To build from source, first fetch the Cloud Foundry CLI project and generate its internal go files:

```
go get github.com/cloudfoundry/cli/cf
cd $GOPATH/src/github.com/cloudfoundry/cli
./bin/build
```

Next fetch this project:

```
$ go get github.com/cloudfoundry-community/cf-plugin-kibana-me-logs
$ cf install-plugin $GOPATH/bin/cf-plugin-kibana-me-logs
```

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

### Bump version & release

There is a helper script to bump the version number `VERSION` file and regenerate the gobindata `version.go` file:

```
export VERSION=X.Y.Z
./bin/bump_version.sh $VERSION
git commit -a -m "bump v$VERSION"
git push
./bin/build.sh
./bin/release.sh
./bin/plugin_index.sh
```

Put the output from `plugin_index.sh` into the
[cloudfoundry-incubator/cli-plugin-repo][1] repository and
submit a pull request to update http://plugins.cloudfoundry.org


[1]: https://github.com/cloudfoundry-incubator/cli-plugin-repo
