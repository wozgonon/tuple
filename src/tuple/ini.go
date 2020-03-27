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

/////////////////////////////////////////////////////////////////////////////
// Yaml Grammar
/////////////////////////////////////////////////////////////////////////////

// http://www.yamllint.com/
type Yaml struct {
	Style
}

func (grammar Yaml) Name() string {
	return "Yaml"
}

func (grammar Yaml) FileSuffix() string {
	return ".yaml"
}

func (grammar Yaml) Parse(context Context, _ Next) {
	Error(context, "Not implemented file suffix: '%s'", grammar.FileSuffix())
}

// TODO Replace with Printer
func (grammar Yaml) printObject(depth string, token Value, out func(value string)) {

	style := grammar.Style
	if tuple, ok := token.(Tuple); ok {

		out(depth)
		out(style.ScalarPrefix)
		len := len(tuple.List)
		if len == 0 {
			out("[]")
			return
		}
		out(style.LineBreak)
		depth = depth + style.Indent
		head := tuple.List[0]
		atom, first := head.(Atom)
		newDepth := depth
		if first {
			out(depth)
			quote(atom.Name, out)
			out(style.Open)
			out(style.LineBreak)
			newDepth = depth + style.Indent
		}
		for k, token := range tuple.List {
			if ! first || k >0  {
				grammar.printObject(newDepth, token, out)
				if k < len-1 {
					out(style.LineBreak)
				}
			}
		}

	} else {
		out(depth)
		out(style.ScalarPrefix)

		switch token.(type) {
		case Atom:
			quote(token.(Atom).Name, out)
		default:
			PrintScalar(style, "", token, out)
		}
	}
}

func (grammar Yaml) Print(object Value, out func(value string)) {
	// TODO PrintExpression(grammar, "", object, out)  // TODO Use Printer
	grammar.printObject("", object, out)
	out (string(NEWLINE))
}

func NewYamlGrammar() Grammar {
	style := NewStyle("---\n", "...\n", "  ", 
		":", "", OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, "",
		"", "\n", "true", "false", '#', "- ")
	return Yaml{style}
}


func (parser Yaml) PrintIndent(depth string, out StringFunction) {
	out(depth)
}

func (parser Yaml) PrintSuffix(depth string, out StringFunction) {
	out(string(NEWLINE))
}

func (parser Yaml) PrintSeparator(depth string, out StringFunction) {}

func (parser Yaml) PrintEmptyTuple(depth string, out StringFunction) {
	out("[]")
}
func (parser Yaml) PrintOpenTuple(depth string, tuple Tuple, out StringFunction) string {
	out("- ")
	return depth + "  "
}

func (parser Yaml) PrintHeadAtom(atom Atom, out StringFunction) {
	quote(atom.Name, out)
	out(": ")
}

func (parser Yaml) PrintCloseTuple(depth string, tuple Tuple, out StringFunction) {}

func (parser Yaml) PrintAtom(depth string, atom Atom, out StringFunction) {
	quote(atom.Name, out)
	//bout(atom.Name)
}

func (parser Yaml) PrintScalarPrefix(depth string, out StringFunction) {
	out ("- ")
}

func (parser Yaml) PrintNullaryOperator(depth string, atom Atom, out StringFunction) {
	PrintTuple(&parser, depth, NewTuple(atom), out)
}

func (parser Yaml) PrintUnaryOperator(depth string, atom Atom, value Value, out StringFunction) {
	PrintTuple(&parser, depth, NewTuple(atom, value), out)
}

func (parser Yaml) PrintBinaryOperator(depth string, atom Atom, value1 Value, value2 Value, out StringFunction) {
	PrintTuple(&parser, depth, NewTuple(atom, value1, value2), out)
}

/////////////////////////////////////////////////////////////////////////////
// Ini Grammar
/////////////////////////////////////////////////////////////////////////////

type Ini struct {
	style Style
}

func (grammar Ini) Name() string {
	return "Ini"
}

func (grammar Ini) FileSuffix() string {
	return ".ini"
}

func (grammar Ini) Parse(context Context, _ Next) {
	Error(context, "Not implemented file suffix: '%s'", grammar.FileSuffix())
}

// TODO 
func (grammar Ini ) printObject(depth string, key string, token Value, out func(value string)) {

	style := grammar.style
	
	if tuple, ok := token.(Tuple); ok {

		len := len(tuple.List)
		if len == 0 {
			out(depth)
			out(style.ScalarPrefix)
			return
		}

		var newDepth string
		head := tuple.List[0]
		atom, ok := head.(Atom)
		first := false

		var prefix string
		if depth == "" {
			prefix = ""
		} else {
			prefix = depth + "."
		}
		if ok {
			key = atom.Name
			newDepth = prefix + atom.Name
			first = true
		} else {
			key = "."
			newDepth = depth
		}
		out(style.LineBreak)
		out(OPEN_SQUARE_BRACKET)
		out(depth)
		out(CLOSE_SQUARE_BRACKET)
		out(style.LineBreak)
		for k, token := range tuple.List {
			if ! first || k >0  {
				grammar.printObject(newDepth, key, token, out)
				if k < len-1 {
					out(style.Separator)
					out(style.LineBreak)
				}
			}
		}

	} else {
		out(key) // TODO just key
		out(style.ScalarPrefix)
		PrintScalar(grammar.style, "", token, out)
	}
}

func (grammar Ini) Print(token Value, out func(value string)) {
	grammar.printObject("", "", token, out)
	out (string(NEWLINE))
}

func NewIniGrammar() Grammar {
	// https://en.wikipedia.org/wiki/INI_file
	style := NewStyle("", "", "",
		"", "", "", "", "",
		"= ", "\n", "true", "false", '#', "=")
	return Ini{style}
}

/////////////////////////////////////////////////////////////////////////////
// PropertyGrammar Grammar
/////////////////////////////////////////////////////////////////////////////

type PropertyGrammar struct {
	style Style
}

func (grammar PropertyGrammar) Name() string {
	return "PropertyGrammar"
}

func (grammar PropertyGrammar) FileSuffix() string {
	return ".properties"
}

func (grammar PropertyGrammar) Parse(context Context, _ Next) {
	Error(context, "Not implemented file suffix: '%s'", grammar.FileSuffix())
}

func (grammar PropertyGrammar) printObject(depth string, token Value, out func(value string)) {
	style := grammar.style
	
	if tuple, ok := token.(Tuple); ok {
		len := len(tuple.List)
		if len == 0 {
			out(depth)
			out(style.ScalarPrefix)
			return
		}
		var newDepth string
		head := tuple.List[0]
		atom, first := head.(Atom)

		var prefix string
		if depth == "" {
			prefix = ""
		} else {
			prefix = depth + "."
		}
		if first {
			newDepth = prefix + atom.Name
		} else {
			newDepth = depth + "."
		}
		for k, token := range tuple.List {
			if ! first || k >0  {
				grammar.printObject(newDepth, token, out)
				if k < len-1 {
					out(style.Separator)
					out(style.LineBreak)
				}
			}
		}

	} else {
		out(depth)
		out(style.ScalarPrefix)
		PrintScalar(grammar.style, "", token, out)
	}
}

func (grammar PropertyGrammar) Print(token Value, out func(value string)) {
	grammar.printObject("", token, out)
	//out (string(NEWLINE))
}

func NewPropertyGrammar() Grammar {
	// https://en.wikipedia.org/wiki/.properties
	style := NewStyle("", "", "",
		"", "", "", "", "",
		" = ", "\n", "true", "false", '#', " = ")
	return PropertyGrammar{style}
}

