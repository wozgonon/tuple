package tuple

import 	"io"
import 	"log"
import 	"fmt"
import 	"bufio"
import 	"os"
import 	"path"

type ParserContext struct {
	SourceName string
	line int
	column int
	scanner io.RuneScanner
}

func NewParserContext(sourceName string, scanner io.RuneScanner) ParserContext {
	return ParserContext{sourceName, 1, 0, scanner}
}

func (context * ParserContext) Suffix() string {
	return path.Ext(context.SourceName)
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
	prefix := fmt.Sprintf("ERROR at %d, %d in '%s': ", context.line, context.column, context.SourceName)
	suffix := fmt.Sprintf(format, args)
	log.Print(prefix + suffix)
}

func RunParser(args []string, parse func (context * ParserContext)) {

	if len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		context := NewParserContext("<stdin>", reader)
		parse(&context)
	} else {
		for _, fileName := range args {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewParserContext(fileName, reader)
			parse(&context)
			file.Close()
		}
	}
}

