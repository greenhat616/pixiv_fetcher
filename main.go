package main

func init() {
	parseFlag()
	initLog()
}

func main() {
	RunHTTPServer()
}
