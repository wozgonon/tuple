package main

import "fmt"

import (
	"bufio"
	// "io"
	//"io/ioutil"
	"log"
	"os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func emit_word(str string) {
	fmt.Print(str)
}

func emit(str rune) {
	fmt.Print(string(str))
	//fmt.Print(str)
}

/////////////////////////////////////////////////////////////////////////////
//  Parser
/////////////////////////////////////////////////////////////////////////////

type item struct {
	args []interface{}
}

type parser struct {
	depth int
	lines int
	columns int
	errors int
	stack []item
}

func (p *parser) error(message string) {
	log.Printf("ERROR at %d:%d:  %s", p.lines, p.columns, message)
	p.errors += 1
}

func (p *parser) setColumn(column int) {
	p.columns = column
}

func (p *parser) token(token string) {

	emit_word (token)
	fmt.Print(" ")

	topItem := p.stack[len(p.stack)-1]
	topItem.args = append(topItem.args, token)
}

func (p *parser) open() {

	p.depth += 1
	p.stack = append(p.stack, item{})
	emit('{')
}

func (p *parser) close() {

	if p.depth < 0 {
		p.error("Missing open bracket")
	} else {
		p.depth -= 1
	}


	l := len(p.stack)
	// Pop top item from stack
	topItem := p.stack[l-1]
	p.stack = p.stack[:l-1]
	p.stack[l-2].args = append(p.stack[l-2].args, topItem)
	log.Printf("topItem %d", len(topItem.args))
	// do something with top.args
	//p.stack.args = append(p.stack.args, "top")
	emit('}')
}

func (p *parser) endline() {

	fmt.Println()
	p.lines += 1
	p.columns = 0
	
	//l := len(p.stack)
	//top := p.stack[l-1]
	//top = append(top, token)
}

func print(i item) {

	fmt.Print("(")
	for _, a := range i.args {
		ii, ok := a.(item)
		if !ok {
			print(ii)
		} else {
			ss, ok2 := a.(string)
			if ok2 {
				fmt.Print("%s", ss)
			}
		}
		fmt.Print("%s", a)
	}
	fmt.Print(")")
}

func (p *parser) end() {

	if p.depth > 0 {
		log.Printf("ERROR at %d:%d,  unclosed braces", p.lines, p.columns);
	}

	log.Printf("Parsed %d lines, %d unclosed braces, %d errors", p.lines, p.depth, p.errors)

	print(p.stack[0])
}


/////////////////////////////////////////////////////////////////////////////

func parse(fileName string) {
		
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	p := parser{0, 0, 0, 0,  make([]item, 0)}
	p.open()
	for scanner.Scan() {
		line := scanner.Text()

		var start = 0
		for k, ch := range line {
			p.setColumn(k)
			if k == len(line) - 1 {
			} else {
				switch ch {
				case ' ':
					p.token(line[start:k])
					start = k+1
				case '=':
				case '{':
					p.open();
				case '}':
					p.close();
				default:
				}
			}
		}
		p.token (line[start:len(line)])
		p.endline()
		start = 0
	}
	p.close()
	p.end()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
}

func main() {
	
	for _, fileName := range os.Args[1:] {
		parse(fileName)
	}
	//fileName := "../../site/in.jml"

}

func expression() {
	fmt.Println("Hello World")

	fileName := "../../site/in.jml"

	//dat, err := ioutil.ReadFile(fileName)
	//check(err)
	//fmt.Print(string(dat))


	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)


	for scanner.Scan() {
		line := scanner.Text()

		var start = 0
		for k, ch := range line {

			if k == len(line) - 1 {
			} else {
				switch ch {
				case ' ':
					emit_word (line[start:k])
					start = k+1
					//fmt.Print("%x", ' ')
				case '+':
					emit (ch)
				case '-':
					emit (ch)
				case '*':
					emit (ch)
				case '/':
					emit (ch)
				case '&':
					emit (ch)
				case '|':
					emit (ch)
				case '=':
					emit (ch)
				case '{':
					emit (ch)
				case '}':
					emit (ch)
				default:
					emit (ch)
				}
			}
		}
		emit_word (line[start:len(line)])
		start = 0
		fmt.Println()
		
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
}
