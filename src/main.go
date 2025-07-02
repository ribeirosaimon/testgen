package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ribeirosaimon/testgen/file"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	ctx := context.Background()

	source := flag.String("source", "", "Caminho de origem")
	flag.Parse()
	fmt.Println("Source:", *source)

	methods := file.New(ctx, *source)

	tmpl := template.Must(template.ParseFiles("test_template.tmpl"))

	replace := strings.Replace(*source, ".go", "_test.go", 1)
	outFile, err := os.Create(replace)
	if err != nil {
		log.Fatal(err)
	}
	if err = tmpl.Execute(outFile, methods); err != nil {
		log.Fatal(err)
	}
	log.Println("Deu tudo certo com seu file:", *source)
}
