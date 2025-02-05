module github.com/hurricanehrndz/examples/routemonitor

go 1.23.5

require (
	github.com/hurricanehrndz/examples/rtmprint v0.0.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/net v0.34.0 //v0.32.0
	golang.org/x/sys v0.30.0
)

replace golang.org/x/net => ../../golang-net/

replace github.com/hurricanehrndz/examples/rtmprint => ../rtmprint/
