package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/javiercbk/impack"
)

func main() {
	var packagePath string
	var compiler string
	var arch string
	flag.StringVar(&compiler, "compiler", "gc", "the go compiler")
	flag.StringVar(&arch, "arch", "amd64", "the target architecture")
	flag.Parse()
	args := flag.Args()
	folders := make([]string, len(args))
	for i := range args {
		if args[i] == "help" {
			fmt.Println("Usage: impack [--compiler] [--arch] folder1 folder2 folder3")
			flag.PrintDefaults()
			os.Exit(0)
		}
		folders[i] = args[i]
	}
	packagePath = strings.TrimSpace(packagePath)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	linter := impack.NewLinter(compiler, arch)
	err := linter.Lint(ctx, packagePath)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}
