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

func (parser * ParserContext) ReadRune() (rune, error) {
	ch, _, err := parser.scanner.ReadRune()
	switch {
	case err != nil: return ch, err
	case ch == '\n':
		parser.line ++
		parser.column = 0
	default:
		parser.column ++
	}
	return ch, nil
}

func (parser * ParserContext) UnreadRune() {
	parser.scanner.UnreadRune()
	if parser.column == 0 {
		parser.line --
	} else {
		parser.column --
	}
}

func (parser * ParserContext) Error(format string, args ...interface{}) {

	prefix := fmt.Sprintf("ERROR at %d, %d in '%s': ", parser.line, parser.column, parser.sourceName)
	suffix := fmt.Sprintf(format, args)
	log.Print(prefix + suffix)
}

