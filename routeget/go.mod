module github.com/hurricanehrndz/examples/routeget

go 1.23.5

require (
	github.com/hurricanehrndz/examples/rtmprint v0.0.0
	golang.org/x/net v0.34.0
	golang.org/x/sys v0.30.0
)

replace golang.org/x/net => ../../golang-net/
replace github.com/hurricanehrndz/examples/rtmprint => ../rtmprint/
