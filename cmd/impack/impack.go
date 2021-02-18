package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/javiercbk/impack"
)

func main() {
	var projectRoot string
	var compiler string
	var arch string
	flag.StringVar(&projectRoot, "project", "", "the project folder path")
	flag.StringVar(&compiler, "compiler", "gc", "the go compiler (default is gc)")
	flag.StringVar(&arch, "arch", "amd64", "the target architecture (default is amd64)")
	flag.Parse()
	projectRoot = strings.TrimSpace(projectRoot)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	linter := impack.NewLinter(compiler, arch)
	err := linter.Lint(ctx, projectRoot)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
