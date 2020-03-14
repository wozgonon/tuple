package main

import "fmt"

import (
	"tuple"
	"bufio"
	"io"
	"log"
	"os"
	"unicode"
)

func next(r io.RuneScanner) (interface{}, error) {

	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if err == io.EOF {
			return "", nil
		}
		if ! unicode.IsSpace(ch) {
			switch {
			case ch == '(' :  return "(", nil
			case ch == ')' :  return ")", nil
			case ch == '"' :  return tuple.ReadCLanguageString(r)
			case ch == '.' || unicode.IsNumber(ch): return tuple.ReadNumber(r, string(ch))    // TODO minus
			case tuple.IsArithmetic(ch): return tuple.ReadAtom(r, string(ch), func(r rune) bool { return tuple.IsArithmetic(r) })
			case unicode.IsLetter(ch):  return tuple.ReadAtom(r, string(ch), func(r rune) bool { return unicode.IsLetter(r) })
			case unicode.IsGraphic(ch): log.Printf("Error graphic character not recognised '%s'", string(ch))
			case unicode.IsControl(ch): log.Printf("Error control character not recognised '%d'", ch)
			default: log.Printf("Error character not recognised '%d'", ch)
			}
		}
	}
}

func parse(reader io.RuneScanner) (tuple.Tuple, error) {

	tuple := tuple.NewTuple()
	
	for {
		token, err := next(reader)
		if err == io.EOF {
			// TODO missing brackets?
			return tuple, nil
		}
		if err != nil {
			return tuple, err
		}
		if token == ")" {
			return tuple, nil
		}
		if token == "(" {
			subTuple, err := parse(reader)
			if err != nil {
				return tuple, err
			}
			tuple.Append(subTuple)
		} else {
			tuple.Append(token)
		}
	}
}

func main() {

	if len(os.Args) == 1 {
		reader := bufio.NewReader(os.Stdin)
		list, err := parse(reader)
		if err != nil {
			log.Print("Error after parsing: %s", err)
		} else {
			fmt.Printf("%s", list.PrettyPrint(""))
		}
		//parse("<stdin>", os.Stdin)
	} else {
		for _, fileName := range os.Args[1:] {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			list, err := parse(reader)
			//err = parse(fileName, file)
			if err != nil {
				log.Print("Error after parsing: %s", err)
			} else {
				fmt.Printf ("%s\n", list.PrettyPrint(""))
			}
			file.Close()
		}
	}
}
