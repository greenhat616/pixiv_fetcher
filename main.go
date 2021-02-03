package main

var isTest bool

func init() {
	if !isTest {
		parseFlag()
		initLog()
		initConfig()
	}
}

func main() {
	RunHTTPServer()
}
