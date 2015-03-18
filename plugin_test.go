package main

import (
	"testing"

	"github.com/cloudfoundry/cli/plugin/fakes"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	cliConn *fakes.FakeCliConnection
)

func TestNoApp(t *testing.T) {
	setup()
	// Convey("checkArgs should not return error with search-logs test", t, func() {
	// 	err := checkArgs(cliConn, []string{"search-logs", "test"})
	// 	So(err, ShouldBeNil)
	// })
	//
	// Convey("checkArgs should return error with search-logs", t, func() {
	// 	err := checkArgs(cliConn, []string{"search-logs"})
	// 	So(err, ShouldNotBeNil)
	// })

}

func TestAppGuid(t *testing.T) {
	Convey("findAppGuid should not return nothing", t, func() {
		err := findAppGuid(cliConn, "test")
		So(err, ShouldNotBeNil)
	})
}

func setup() {
	cliConn = &fakes.FakeCliConnection{}
}
