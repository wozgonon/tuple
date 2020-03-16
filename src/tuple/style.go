package tuple

import "fmt"
import "log"
import "strconv"
import "reflect"
import "math"

/////////////////////////////////////////////////////////////////////////////
// Pretty printer
/////////////////////////////////////////////////////////////////////////////

type Style struct {
	Indent string
	Open string
	Close string
	Separator string
	LineBreak string
	True string
	False string
	OneLineComment rune
}

func (style Style) printToken(depth string, token interface{}, out func(value string)) {
	if tuple, ok := token.(Tuple); ok {
		style.printTuple(depth, tuple, out)
	} else {
		out(depth)
		switch token.(type) {
		case Atom:
			out(token.(Atom).Name)
		case string:
			out(DOUBLE_QUOTE)
			out(token.(string))   // TODO Escape
			out(DOUBLE_QUOTE)
		case bool:
			if token.(bool) {
				out(style.True)
			} else {
				out(style.False)
			}				
		case int64:
			out(strconv.FormatInt(int64(token.(int64)), 10))
		case float64:
			float := token.(float64)
			if math.IsInf(float, 64) {
				out("Inf")  // Do not print +Inf
			} else {
				out(fmt.Sprint(token.(float64)))
			}
		default:
			log.Printf("ERROR type '%s' not recognised: %s", reflect.TypeOf(token), token);
			out(UNKNOWN)
		}
	}
}

func (style Style) printTuple(depth string, tuple Tuple, out func(value string)) {
	out(depth)
	out(style.Open)
	out(style.LineBreak)
	newDepth := depth + style.Indent
	len := len(tuple.List)
	for k, token := range tuple.List {
		style.printToken(newDepth, token, out)
		if k < len-1 {
			out(style.Separator)
		}
		out(style.LineBreak)
	}
	out(depth)
	out(style.Close)
}

func (style Style) PrettyPrint(token interface{}, out func(value string)) {
	style.printToken("", token, out)
}

