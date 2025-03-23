module github.com/hurricanehrndz/examples/routemonitor

go 1.23.5

require (
	github.com/hurricanehrndz/examples/rtmprint v0.0.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/net v0.37.0 //v0.32.0
	golang.org/x/sys v0.31.0
)

replace github.com/hurricanehrndz/examples/rtmprint => ../rtmprint/
