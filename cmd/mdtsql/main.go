package main

import "github.com/noborus/mdtsql/cmd"

// version represents the version
var version = "v0.0.3"

// revision set "git rev-parse --short HEAD"
var revision = "HEAD"

func main() {
	cmd.Execute(version, revision)
}
