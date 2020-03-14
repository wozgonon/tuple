package tuple

import 	"io"
import 	"log"
import 	"fmt"

type ParserContext struct {
	sourceName string
	line int
	column int
	scanner io.RuneScanner
}

func NewParserContext(sourceName string, scanner io.RuneScanner) ParserContext {
	return ParserContext{sourceName, 1, 0, scanner}
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

