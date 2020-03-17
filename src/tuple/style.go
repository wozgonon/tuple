/*
    This file is part of WOZG.

    WOZG is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    WOZG is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with WOZG.  If not, see <https://www.gnu.org/licenses/>.
*/
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
	StartDoc string
	EndDoc string
	Indent string
	Open string
	Close string
	Separator string
	LineBreak string
	True string
	False string
	OneLineComment rune
}

func quote(value string, out func(value string)) {
	out(DOUBLE_QUOTE)
	out(value)   // TODO Escape
	out(DOUBLE_QUOTE)
}

func (style Style) printToken(depth string, token interface{}, out func(value string)) {
	if tuple, ok := token.(Tuple); ok {
		style.printTuple(depth, tuple, out)
	} else {
		out(depth)
		if style.IsWholePath() {
			out(style.Open)
		} else if style.indentOnly() {
			out("- ")
		}
		switch token.(type) {
		case Atom:
			if style.indentOnly() {
				quote(token.(Atom).Name, out)
			} else {
				out(token.(Atom).Name)
			}
		case string:
			quote(token.(string), out)   // TODO Escape
		case bool:
			if token.(bool) {
				out(style.True)
			} else {
				out(style.False)
			}				
		case Comment:
			out(string(style.OneLineComment))
			out(token.(Comment).Comment)
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

	len := len(tuple.List)
	if len == 0 {

		 if style.indentOnly() {
			out(depth)
			out("- ")
			out("[]")
		} else {
			out(depth)
			out(style.Open)
			//out("$")
			out(style.Close)
		}
		return
	}

	if style.indentOnly() {
		out(depth)
		out("- ")
		out(style.LineBreak)
	}
	var newDepth string
	head := tuple.List[0]
	atom, ok := head.(Atom)
	first := ok && style.indentOnly()

	if style.IsWholePath() {
		var prefix string
		if depth == "" {
			prefix = ""
		} else {
			prefix = depth + "."
		}
		if ok {
			newDepth = prefix + atom.Name
			//out(depth)
			//out(style.Open)
			//out("*")
			//out(style.LineBreak)
			first = true
		} else {
			newDepth = depth
		}	
	} else if first {
		depth = depth + style.Indent
		out(depth)
		quote(atom.Name, out)
		out(style.Open)
		out(style.LineBreak)
		newDepth = depth + style.Indent
	} else if style.indentOnly() {
		depth = depth + style.Indent
		newDepth = depth
	} else {
		out(depth)
		out(style.Open)
		out(style.LineBreak)
		newDepth = depth + style.Indent
	}
	
	for k, token := range tuple.List {
		if ! first || k >0  {
			style.printToken(newDepth, token, out)
			if k < len-1 {
				out(style.Separator)
				out(style.LineBreak)
			}
			//out("@")
			//out("$")
		}
	}
	if style.IsWholePath() {
	} else if !style.indentOnly() {
		out(style.LineBreak)
		out(depth)
		out(style.Close)
	}
	//out("&")
}

func (style Style) IsWholePath() bool {
	return style.Indent == "" // TODO
}

func (style Style) indentOnly() bool {
	// TODO set this as a parameter
	return style.Close == "" && style.Indent != "" // TODO
}

func (style Style) PrettyPrint(token interface{}, out func(value string)) {

	ignoreOuterBrackets := style.Open == string('{')
	if tuple, ok := token.(Tuple); ok && ignoreOuterBrackets {
		len := len(tuple.List)
		for k, token := range tuple.List {
			style.printToken("", token, out)
			if k < len-1 {
				out(style.Indent)
				out(style.Separator)
			}
		}
	} else {
		style.printToken("", token, out)
	}
	out (string(NEWLINE))
}

