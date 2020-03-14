package main

import "fmt"

import (
	"tuple"
	"strconv"
	"bufio"
	"io"
	//"io/ioutil"
	"log"
	"os"
	"unicode"
	//"unicode/utf8"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////


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
			if ch == '(' {
				return "(", nil
			} else if ch == ')' {
				return ")", nil
			} else if ch == '"' {
				return readString(r, "", false, func(r rune) bool { return r != '"' })  // TODO Handle escape character
			} else if unicode.IsNumber(ch) || ch == '.'{
				number, err := readString(r, string(ch), true, func(r rune) bool { return unicode.IsNumber(r) || r== '.' })  // TODO multiple . in number
				if err != nil {
					return "", err
				}
				if strings.Contains(number, ".") {
					return strconv.ParseFloat(number, 64)
				}
				return strconv.ParseInt(number, 10, 0)
				// TODO "NaN", "+Inf", and "-Inf" 

			} else if isArithmetic(ch) {
				operator, err := readString(r, string(ch), true, func(r rune) bool { return unicode.IsSymbol(r) })
				if err != nil {
					return "", err
				}
				return tuple.Atom{operator}, err
				
			} else if unicode.IsLetter(ch) {
				atom, err := readString(r, string(ch), true, func(r rune) bool { return unicode.IsLetter(r) })
				if err != nil {
					return "", err
				}
				return tuple.Atom{atom}, err
				
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
