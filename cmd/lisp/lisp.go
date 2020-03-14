package main

import "fmt"

import (
	"bufio"
	"io"
	//"io/ioutil"
	"log"
	"os"
	"unicode"
	//"unicode/utf8"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// "Code point" Go introduces a shorter term for the concept: rune and means exactly the same as "code point", with one interesting addition.
//  Go language defines the word rune as an alias for the type int32, so programs can be clear when an integer value represents a code point.
//  Moreover, what you might think of as a character constant is called a rune constant in Go. 

func readString (r io.RuneScanner, token string, keepLast bool, test func(r rune) bool) (string, error) {
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			log.Printf("ERROR missing close \"")
			return token, nil
		} else if err != nil {
			//log.Printf("ERROR nil")
			//return ""
		} else if ! test(ch) {
			if keepLast {
				r.UnreadRune()
			}
			return token, nil
		} else {
			// TODO not efficient
			token = token + string(ch)
		}
	}

}

func isArithmetic(ch rune) bool {
	switch ch {
		case '+': return true
		case '-': return true
		case '/': return true
		case '*': return true
		case '^': return true
		default: return false
	}
}

func next(r io.RuneScanner) (string, error) {

	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if err == io.EOF {
			return "", nil
		}
		if ! unicode.IsSpace(ch) {
			if ch == '(' {
				return "(", nil
			} else if ch == ')' {
				return ")", nil
			} else if ch == '"' {
				return readString(r, "", false, func(r rune) bool { return r != '"' })  // Handle escape character
			} else if unicode.IsNumber(ch) {
				return readString(r, string(ch), true, func(r rune) bool { return unicode.IsNumber(r) || r== '.' })  // TODO multiple . in number
			} else if ch == '.' {
				return readString(r, ".", true, func(r rune) bool { return unicode.IsNumber(r) })
			//} else if unicode.IsPunct(ch) {
			//	return readString(r, string(ch), func(r rune) bool { return unicode.IsPunct(r) })
			} else if isArithmetic(ch) {
				return readString(r, string(ch), true, func(r rune) bool { return unicode.IsSymbol(r) })
			} else if unicode.IsLetter(ch) {
				return readString(r, string(ch), true, func(r rune) bool { return unicode.IsLetter(r) })
			} else if unicode.IsGraphic(ch) {
				log.Printf("Error graphic character not recognised '%s'", string(ch))
			} else if unicode.IsControl(ch) {
				log.Printf("Error control character not recognised '%d'", ch)
			} else  {
				log.Printf("Error character not recognised '%d'", ch)
			}
		}
	}
}

type Tuple struct {
	list []interface{}
}

func (tuple Tuple) prettyPrintList() string {
	space := ""
	result := "("
	for _, token := range tuple.list {
		result = result + space + token
		space = " "
	}
	result = result + ")"
	return result
}


func (tuple *Tuple) Append(token string) {
	tuple.list = append(tuple.list, token)

}

func NewTuple() tuple {
	return tuple{make([]string, 0)}
}

func parse(reader io.RuneScanner) (tuple, error) {

	tuple := NewTuple()
	
	for {
		token, err := next(reader)
		//fmt.Printf ("%s\n", token)
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
			pretty := subList.prettyPrintList ()
			tuple.Append(pretty)
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
			fmt.Printf("%s", list.prettyPrintList())
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
				fmt.Printf ("%s\n", list.prettyPrintList())
			}
			file.Close()
		}
	}
}
