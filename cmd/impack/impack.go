package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/javiercbk/impack"
)

func main() {
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	linter, err := impack.NewLinter(compiler, arch)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for _, folder := range folders {
		err = linter.Lint(ctx, folder)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
