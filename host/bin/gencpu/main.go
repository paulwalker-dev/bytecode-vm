package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const TOP = 13

func parameters(keywords []string) int {
	i := 0
	for keywords[i] != "[" {
		i++
	}
	ret := 0
	for keywords[i+ret] != "]" {
		ret++
	}
	return ret - 1
}

func def(keywords []string, out io.Writer) {
	ident := strings.Split(keywords[0], ":")
	name := ident[1]

	switch ident[0] {
	case "reg":
		_, err := io.WriteString(out, fmt.Sprintf("%s? equ 0x%s\n", name, keywords[1]))
		if err != nil {
			log.Fatal(err)
		}
	case "core":
		io.WriteString(out, fmt.Sprintf("macro %s? ", name))
		args := make([]string, parameters(keywords))
		for i := range args {
			args[i] = fmt.Sprintf("a%d", i+1)
		}
		io.WriteString(out, fmt.Sprintf("%s\n", strings.Join(args, ",")))
		args = append([]string{fmt.Sprintf("0x%s", keywords[1])}, args...)
		io.WriteString(out, fmt.Sprintf("db %s\n", strings.Join(args, ",")))
		io.WriteString(out, "end macro\n")
	case "common":
		io.WriteString(out, fmt.Sprintf("macro %s? ", name))
		args := make([]string, parameters(keywords))
		for i := range args {
			args[i] = fmt.Sprintf("a%d", i+1)
		}
		io.WriteString(out, fmt.Sprintf("%s\n", strings.Join(args, ",")))
	}
}

func impl(keywords []string, out io.Writer) {
	if keywords[0] == "done" {
		io.WriteString(out, "end macro\n")
		return
	}
	if keywords[0] == "emit" {
		io.WriteString(out, fmt.Sprintf("emit 0x%s: ", keywords[1]))
		keywords = keywords[2:]
	} else {
		io.WriteString(out, fmt.Sprintf("%s ", keywords[0]))
		keywords = keywords[1:]
	}
	for i, b := range keywords {
		b := b
		switch b[0] {
		case '#':
			n, err := strconv.Atoi(string(b[1]))
			if err != nil {
				log.Fatal(err)
			}
			b = strconv.Itoa(TOP - n)
		case '$':
			b = strings.ReplaceAll(b, "$", "a")
		case '%':
			b = b[1:]
		default:
			b = "0x" + b
		}
		keywords[i] = b
	}
	io.WriteString(out, fmt.Sprintf("%s\n", strings.Join(keywords, ",")))
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("Usage: <in> <out>")
	}

	buf, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.OpenFile(args[1], os.O_TRUNC|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	input := string(buf)
	input = strings.ReplaceAll(input, "\r\n", "\n")
	for _, line := range strings.Split(input, "\n") {
		keywords := strings.Split(line, " ")
		if len(keywords) == 0 {
			continue
		}

		switch keywords[0] {
		case "def":
			def(keywords[1:], out)
		case "impl":
			impl(keywords[1:], out)
		}
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}
