package tuple

import 	"io"
import 	"log"
import 	"fmt"
import 	"bufio"
import 	"os"
import 	"path"

type ParserContext struct {
	sourceName string
	line int
	column int
	scanner io.RuneScanner
}

func NewParserContext(sourceName string, scanner io.RuneScanner) ParserContext {
	return ParserContext{sourceName, 1, 0, scanner}
}

func (context * ParserContext) Suffix() string {
	return path.Ext(context.sourceName)
}

func (context * ParserContext) ReadRune() (rune, error) {
	ch, _, err := context.scanner.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == '\n':
		context.line ++
		context.column = 0
	default:
		context.column ++
	}
	return ch, nil
}

func (context * ParserContext) UnreadRune() {
	context.scanner.UnreadRune()
	if context.column == 0 {
		context.line --
	} else {
		context.column --
	}
}

func (context * ParserContext) Error(format string, args ...interface{}) {
	prefix := fmt.Sprintf("ERROR at %d, %d in '%s': ", context.line, context.column, context.sourceName)
	suffix := fmt.Sprintf(format, args)
	log.Print(prefix + suffix)
}

func run(context * ParserContext, parse func (context * ParserContext) (Tuple, error), processTuple func(tuple Tuple)) {
	tuple, err := parse(context)
	if err != io.EOF && err != nil {
		context.Error("Failed while parsing: %s", err)
	} else {
		processTuple(tuple)
	}

}

func RunParser(args []string, parse func (context * ParserContext) (Tuple, error), processTuple func(tuple Tuple))  {

	if len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		context := NewParserContext("<stdin>", reader)
		run(&context, parse, processTuple)
	} else {
		for _, fileName := range os.Args[1:] {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewParserContext(fileName, reader)
			run(&context, parse, processTuple)
			file.Close()
		}
	}
}

