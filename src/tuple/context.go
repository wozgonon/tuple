package tuple

import 	"io"
import 	"log"
import 	"fmt"
import 	"bufio"
import 	"os"
import 	"path"

type ParserContext struct {
	SourceName string
	line int64
	column int64
	scanner io.RuneScanner
	logStyle Style
}

func NewParserContext(sourceName string, scanner io.RuneScanner, logStyle Style) ParserContext {
	return ParserContext{sourceName, 1, 0, scanner, logStyle}
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

	tuple := NewTuple()
	tuple.Append("ERROR")
	tuple.Append(context.line)
	tuple.Append(context.column)
	tuple.Append(context.SourceName)
	tuple.Append(suffix)
	context.logStyle.PrettyPrint(tuple, func (value string) { fmt.Print(value) })

}

func RunParser(args []string, logStyle Style, parse func (context * ParserContext)) {

	if len(args) == 0 {
		reader := bufio.NewReader(os.Stdin)
		context := NewParserContext("<stdin>.l", reader, logStyle)
		parse(&context)
	} else {
		for _, fileName := range args {
			file, err := os.Open(fileName)
			if err != nil {
				log.Fatal(err)
			}
			reader := bufio.NewReader(file)
			context := NewParserContext(fileName, reader, logStyle)
			parse(&context)
			file.Close()
		}
	}
}

