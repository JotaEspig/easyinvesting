package config

import "os"

var BRAPI_TOKEN = ""

func init() {
	BRAPI_TOKEN = os.Getenv("BRAPI_TOKEN")
	if BRAPI_TOKEN == "" {
		panic("BRAPI_TOKEN environment variable is not set")
	}
}
