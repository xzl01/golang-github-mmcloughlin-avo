//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var output = flag.String("output", "", "path to output file (default stdout)")

func PrintConstType(w io.Writer, name, typ, format string, size int, doc string) {
	r := typ[0]
	fmt.Fprintf(w, "// %s\n", doc)
	fmt.Fprintf(w, "type %s %s\n", name, typ)
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "func (%c %s) Asm() string { return fmt.Sprintf(\"$%s\", %c) }\n", r, name, format, r)
	fmt.Fprintf(w, "func (%c %s) Bytes() int  { return %d }\n", r, name, size)
	fmt.Fprintf(w, "func (%c %s) constant()   {}\n", r, name)
	fmt.Fprintf(w, "\n")
}

func PrintConstTypes(w io.Writer) {
	_, self, _, _ := runtime.Caller(0)
	fmt.Fprintf(w, "// Code generated by %s. DO NOT EDIT.\n\n", filepath.Base(self))
	fmt.Fprintf(w, "package operand\n\n")
	fmt.Fprintf(w, "import \"fmt\"\n\n")
	for n := 1; n <= 8; n *= 2 {
		bits := n * 8
		bs := strconv.Itoa(bits)

		if n >= 4 {
			PrintConstType(w, "F"+bs, "float"+bs, "(%#v)", n, fmt.Sprintf("F%d is a %d-bit floating point constant.", bits, bits))
		}
		PrintConstType(w, "I"+bs, "int"+bs, "%+d", n, fmt.Sprintf("I%d is a %d-bit signed integer constant.", bits, bits))
		PrintConstType(w, "U"+bs, "uint"+bs, "%#0"+strconv.Itoa(2*n)+"x", n, fmt.Sprintf("U%d is a %d-bit unsigned integer constant.", bits, bits))
	}
}

func main() {
	flag.Parse()

	w := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		w = f
	}

	buf := bytes.NewBuffer(nil)
	PrintConstTypes(buf)

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(src)
	if err != nil {
		log.Fatal(err)
	}
}