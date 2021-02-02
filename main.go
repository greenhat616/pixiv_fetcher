package main

func init() {
	parseFlag()
	initLog()
	initConfig()
}

func main() {
	RunHTTPServer()
}
