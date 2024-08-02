package app

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"

	"github.com/hidori/go-genmapper/generator"
	"github.com/pkg/errors"
)

const (
	doNotEdit = "// Code generated by github.com/hidori/go-genmapper/cmd/genmapper DO NOT EDIT."
)

func Run() error {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Printf("usage: %s [OPTION]... <FILE>\n\noption(s):\n", path.Base(os.Args[0]))
		flag.PrintDefaults()

		return nil
	}

	return generate(os.Stdout, args[0])
}

func generate(writer io.Writer, fileName string) error {
	config := &generator.GeneratorConfig{}

	file, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
	if err != nil {
		return errors.Wrap(err, "fail to parser.ParseFile()")
	}

	generator := generator.NewGenerator(config)

	decls, err := generator.Generate(token.NewFileSet(), file)
	if err != nil {
		return errors.Wrap(err, "fail to generator.Generate()")
	}

	buffer := bytes.NewBuffer([]byte{})

	err = format.Node(buffer, token.NewFileSet(), &ast.File{
		Name:  ast.NewIdent(file.Name.Name),
		Decls: decls,
	})
	if err != nil {
		return errors.Wrap(err, "fail to format.Node()")
	}

	cooked := buffer.Bytes()
	// cooked, err := imports.Process("", buffer.Bytes(), &imports.Options{FormatOnly: false})
	// if err != nil {
	// 	return errors.Wrap(err, "fail to imports.Process()")
	// }

	_, _ = fmt.Fprintln(writer, doNotEdit)
	_, _ = writer.Write(cooked)

	return nil
}