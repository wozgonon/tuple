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

import "io"
import "unicode"
import "errors"
//import "fmt"

/////////////////////////////////////////////////////////////////////////////
//
/////////////////////////////////////////////////////////////////////////////


func (parser SExpressionParser) parseSExpressionTuple(context * ParserContext, tuple *Tuple) (error) {

	//fmt.Printf("parseSExpressionTuple depth=%d, s\n", context.depth)
	
	style := parser.style
	for {
		token, err := parser.getNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == style.Close:
			//fmt.Printf("*** close=%s\n", token)
			return nil
		case token == style.Close2:
			//fmt.Printf("*** close=%s\n", token)
			return nil
		case token == style.Open || token == style.Open2:
			context.Open()
			subTuple := NewTuple()
			err := parser.parseSExpressionTuple(context, &subTuple)
			context.Close()
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			//fmt.Printf("1. s=%s", subTuple)
			tuple.Append(subTuple)
		case token == style.KeyValueSeparator:  // TODO check if it is an operator
			//fmt.Printf("--------------------\n")
			if tuple.Length() == 0 {
				context.Error("Unexpected operator '%s'", style.KeyValueSeparator)
				return errors.New("Unexpected")
			}
			key := tuple.List[tuple.Length()-1]
			//fmt.Printf("** key=%s\n", key)
			value, err := parser.parse(context)
			if err != nil {
				return err
			}
			if value == style.Close || value == style.Close2 {
				context.Error ("Unexpected close bracket '%s'", token)
				return errors.New("Unexpected")

			}
			//fmt.Printf("depth=%d, key=%s value=%s\n", context.depth, key, value)
			tuple.List[tuple.Length() -1] = NewTuple(Atom{"_cons"}, key, value)
		default:
			if _,ok := token.(Comment); ok {
				// TODO Ignore ???
			} else {
				//fmt.Printf("depth=%d, append=%s\n", context.depth, token)
				tuple.Append(token)
			}
		}
	}
}

func (parser SExpressionParser) parse(context * ParserContext) (interface{}, error) {

	style := parser.style
	token, err := parser.getNext(context)
	switch {
	case err == io.EOF:
		return nil, err
	case err != nil:
		context.Error ("'%s'", err)
		return nil, err
	case token == style.Close:
		if context.depth == 0 {
			context.Error ("Unexpected close bracket '%s'", style.Close)
			return nil, errors.New("Unexpected")
		}
		//context.Close()
		return token, nil
	case token == style.Close2:
		context.Error ("Unexpected close bracket '%s'", style.Close2)
		return nil, errors.New("Unexpected")
	case token == style.Open || token == style.Open2:
		//fmt.Printf("!!!Open\n")
		context.Open()
		tuple := NewTuple()
		err := parser.parseSExpressionTuple(context, &tuple)
		context.Close()
		if err != nil {
			return nil, err
		}
		return tuple, nil
	default:
		//if _,ok := token.(Comment); ok {
			// TODO Ignore ???
		//} else {
			return token, nil
		//}
	}
}

func (parser SExpressionParser) Parse(context * ParserContext) {

	for {
		value, err := parser.parse(context)
		if err == nil {
			context.next(value)
		} else {
			return
		}
	}
}


/////////////////////////////////////////////////////////////////////////////
// Lisp Grammar
/////////////////////////////////////////////////////////////////////////////

// A [S-Expression](https://en.wikipedia.org/wiki/S-expression) or symbolic expression is a very old and general notation.
// A nested structure of scalars (atoms and numbers), lists and key-values pairs (called cons cells).
// These are used for the syntax of LISP but also any other language can typically be converted to an S-Expression,
// it is in particular a very useful format for debugging a parser by printing out the Abstract Syntaxt Tree (AST) created by parsing.
type LispGrammar struct {
	parser SExpressionParser
}

func (grammar LispGrammar) Name() string {
	return "Lisp"
}

func (grammar LispGrammar) FileSuffix() string {
	return ".l"
}

func (grammar LispGrammar) Parse(context * ParserContext) {
	grammar.parser.Parse(context)
}

func (grammar LispGrammar) Print(token interface{}, next func(value string)) {
	grammar.parser.style.PrettyPrint(token, next)
}

func NewLispGrammar() Grammar {
	style := Style{"", "", "  ",
		"(", ")", "", "", ".",
		"", "\n", "true", "false", ';', ""}
	return LispGrammar{NewSExpressionParser(style)}
}

/////////////////////////////////////////////////////////////////////////////
// Tcl Grammar
/////////////////////////////////////////////////////////////////////////////

// The basic syntax of command shell parsers is typically very simple:
// everything is a line of strings separated by spaces terminate by a newline.
// Mostly no need for quotes or double quotes or semi-colons.
// This makes it very easy to type command with a few parameters on a command line interface (CLI).
//
// Examples include [TCL](https://en.wikipedia.org/wiki/Tcl), [Bash](https://en.wikipedia.org/wiki/Bash_(Unix_shell)), DOS cmd shell
// TCL and bash use braces { ... } for nesting.
//
// Bash and DOS of course have lots of extra syntax for working with files but the
// basic syntax typically does not understand arithmetic, with infix notation, one has to use a special tool:
// * TCL has a 'expr' function that understand arithmetic with infix notation
// * Bash one can use an external 'expr(1)' tool:  'expr 8.3 + 6'
// * DOS has a special version of the SET command 'SET /a c=a+b'
//
type Tcl struct {
	parser SExpressionParser
}

func (grammar Tcl) Name() string {
	return "Tcl"
}

func (grammar Tcl) FileSuffix() string {
	return ".tcl"
}

func (grammar Tcl ) readCommandString(context * ParserContext, token string) (string, error) {
	parser := grammar.parser
	return ReadString(context, token, true, func (ch rune) bool {
		return ! unicode.IsSpace(ch) && string(ch) != parser.style.Close && string(ch) != parser.style.Open && ch != '$'
	})
}

func (grammar Tcl) getNextCommandShell(context * ParserContext) (interface{}, error) {

	parser := grammar.parser
	for {
		ch, err := context.ReadRune()
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case ch == NEWLINE: return string(NEWLINE), nil
		case unicode.IsSpace(ch): break
		case ch == parser.style.OneLineComment:
			// TODO ignore for now
			//return string(ch), nil
		case ch == parser.openChar :  return parser.style.Open, nil
		case ch == parser.closeChar : return parser.style.Close, nil
		case ch == '"' :  return ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return ReadNumber(context, string(ch))    // TODO minus
		case ch == '$':
			value, err := grammar.readCommandString(context, "")
			if err != nil {
				return nil, err
			}
			return Atom{value}, nil
		case unicode.IsGraphic(ch): return grammar.readCommandString(context, string(ch))
		case unicode.IsControl(ch): context.Error("Error control character not recognised '%d'", ch)
		default: context.Error("Error character not recognised '%d'", ch)
		}
	}
}

func (grammar Tcl) parseCommandShellTuple(context * ParserContext, tuple *Tuple) (error) {

	parser := grammar.parser
	for {
		token, err := grammar.getNextCommandShell(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == parser.style.Close:
			return nil
		case token == parser.style.Open:
			subTuple := NewTuple()
			err := grammar.parseCommandShellTuple(context, &subTuple)
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			tuple.Append(subTuple)
		case token == string(NEWLINE):
		default:
			tuple.Append(token)
		}
	}
}

func (grammar Tcl) Parse(context * ParserContext) {

	parser := grammar.parser

	resultTuple := NewTuple()
	for {
		token, err := grammar.getNextCommandShell(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			context.Error ("'%s'", err)
			return
		case token == string(NEWLINE):
			l := resultTuple.Length()
			context.Verbose ("Newline length of tuple=%d", l)
			switch l {
			case 0: // Ignore
			case 1:
				first := resultTuple.List[0]
				if _, ok := first.(Atom); ok {
					context.next(resultTuple)
				} else {
					context.next(token)
				}
			default:
				context.next(resultTuple)
			}
			resultTuple = NewTuple()
		case token == parser.style.OneLineComment:
			comment, err := ReadUntilEndOfLine(context)
			if err != nil {
				return
			}
			context.next(comment)
		case token == parser.style.Close:
			context.Error ("Unexpected close bracket '%s'", parser.style.Close)
		case token == parser.style.Open:
			subTuple := NewTuple()
			err := grammar.parseCommandShellTuple(context, &subTuple)
			if err != nil {
				return // tuple, err
			}
			resultTuple.Append(subTuple)
		default:
			context.Verbose("Add token: '%s'", token)
			resultTuple.Append(token)
		}
	}
}

func (grammar Tcl) Print(token interface{}, out func(value string)) {
	style := grammar.parser.style
	if tuple, ok := token.(Tuple); ok {
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

func NewTclGrammar() Grammar {
	style := Style{"", "", "  ",
		"{", "}", "[", "]", ":",
		"", "\n", "true", "false", '#', ""}
	return Tcl{NewSExpressionParser(style)}
}

/*/////////////////////////////////////////////////////////////////////////////
// JML Grammar
/////////////////////////////////////////////////////////////////////////////

type Jml struct {
	parser SExpressionParser
}

func (grammar Jml) Name() string {
	return "Jml"
}

func (grammar Jml) FileSuffix() string {
	return ".jml"
}

func (grammar Jml) Parse(context * ParserContext) {
	grammar.parser.ParseSExpression(context)
}

func (grammar Jml) Print(token interface{}, next func(value string)) {
	grammar.parser.style.PrettyPrint(token, next)
}

func NewJmlGrammar() Grammar {
	style := Style{"\n", "", "  ", "{", "}", "", "\n", "true", "false", '#'}
	return Jml{NewSExpressionParser(style)}
}*/

/////////////////////////////////////////////////////////////////////////////
// Tuple Grammar
/////////////////////////////////////////////////////////////////////////////

type TupleGrammar struct {
	parser SExpressionParser
}

func (grammar TupleGrammar) Name() string {
	return "TupleGrammar"
}

func (grammar TupleGrammar) FileSuffix() string {
	return ".tuple"
}

func (grammar TupleGrammar) parseCommaTuple(context * ParserContext, tuple *Tuple) (error) {

	parser := grammar.parser
	// TODO comma and semi-colon
	for {
		token, err := parser.getNext(context)
		switch {
		case err != nil:
			context.Error("parsing %s", err);
			return err /// ??? Any need to return
		case token == parser.style.Close:
			return nil
		case token == parser.style.Open:
			subTuple := NewTuple()
			err := grammar.parseCommaTuple(context, &subTuple)
			if err == io.EOF {
				context.Error ("Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			tuple.Append(subTuple)
		default:
			if _,ok := token.(Comment); ok {
				// TODO Ignore ???
			} else {
				tuple.Append(token)
			}
		}
	}

}

func (grammar TupleGrammar) Parse(context * ParserContext) {

	parser := grammar.parser
	for {
		token, err := parser.getNext(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			context.Error ("'%s'", err)
			return
		case token == parser.style.Close:
			context.Error ("Unexpected close bracket '%s'", parser.style.Close)
		default:
			if atom,ok := token.(Atom); ok {
				bracket, err := parser.getNext(context)
				if err != nil {
					context.Error ("'%s'", err)
					return
				}
				if bracket != parser.style.Open {
					context.Error ("Expected open bracket '%s' after '%s', not '%s'", parser.style.Open, token, bracket)
				} else {
					subTuple := NewTuple()
					subTuple.Append(atom)
					err := grammar.parseCommaTuple(context, &subTuple)
					if err != nil {
						return
					}
					context.next(subTuple)
				}
			} else {
				context.next(token)
			}
		}
	}
}

func (grammar TupleGrammar) Print(token interface{}, next func(value string)) {
	grammar.parser.style.PrettyPrint(token, next)
}

func NewTupleGrammar() Grammar {
	style := Style{"", "", "  ",
		"(", ")", "", "", ":",
		",", "\n", "true", "false", '%', ""} // prolog, sql '--' for 
	return TupleGrammar{NewSExpressionParser(style)}
}

/////////////////////////////////////////////////////////////////////////////
// Yaml Grammar
/////////////////////////////////////////////////////////////////////////////

// http://www.yamllint.com/
type Yaml struct {
	parser SExpressionParser
}

func (grammar Yaml) Name() string {
	return "Yaml"
}

func (grammar Yaml) FileSuffix() string {
	return ".yaml"
}

func (grammar Yaml) Parse(context * ParserContext) {
	context.Error("Not implemented file suffix: '%s'", grammar.FileSuffix())
}

func (grammar Yaml) printToken(depth string, token interface{}, out func(value string)) {

	style := grammar.parser.style
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
				grammar.printToken(newDepth, token, out)
				if k < len-1 {
					out(style.LineBreak)
				}
			}
		}

	} else {
		out(depth)
		out(style.ScalarPrefix)
		style.printScalar(token, out)
	}
}

func (grammar Yaml) Print(token interface{}, out func(value string)) {
	grammar.printToken("", token, out)
	out (string(NEWLINE))
}

func NewYamlGrammar() Grammar {
	style := Style{"---\n", "...\n", "  ", 
		":", "", "[", "]", "",
		"", "\n", "true", "false", '#', "- "}
	return Yaml{NewSExpressionParser(style)}
}

/////////////////////////////////////////////////////////////////////////////
// Ini Grammar
/////////////////////////////////////////////////////////////////////////////

type Ini struct {
	parser SExpressionParser
}

func (grammar Ini) Name() string {
	return "Ini"
}

func (grammar Ini) FileSuffix() string {
	return ".ini"
}

func (grammar Ini) Parse(context * ParserContext) {
	context.Error("Not implemented file suffix: '%s'", grammar.FileSuffix())
}

// TODO 
func (grammar Ini ) printToken(depth string, key string, token interface{}, out func(value string)) {

	style := grammar.parser.style
	
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
		out("[")
		out(depth)
		out("]")
		out(style.LineBreak)
		for k, token := range tuple.List {
			if ! first || k >0  {
				grammar.printToken(newDepth, key, token, out)
				if k < len-1 {
					out(style.Separator)
					out(style.LineBreak)
				}
			}
		}

	} else {
		out(key) // TODO just key
		out(style.ScalarPrefix)
		style.printScalar(token, out)
	}
}

func (grammar Ini) Print(token interface{}, out func(value string)) {
	grammar.printToken("", "", token, out)
	out (string(NEWLINE))
}

func NewIniGrammar() Grammar {
	// https://en.wikipedia.org/wiki/INI_file
	style := Style{"", "", "",
		"", "", "", "", "",
		"= ", "\n", "true", "false", '#', "="}
	return Ini{NewSExpressionParser(style)}
}

/////////////////////////////////////////////////////////////////////////////
// PropertyGrammar Grammar
/////////////////////////////////////////////////////////////////////////////

type PropertyGrammar struct {
	parser SExpressionParser
}

func (grammar PropertyGrammar) Name() string {
	return "PropertyGrammar"
}

func (grammar PropertyGrammar) FileSuffix() string {
	return ".properties"
}

func (grammar PropertyGrammar) Parse(context * ParserContext) {
	context.Error("Not implemented file suffix: '%s'", grammar.FileSuffix())
}

func (grammar PropertyGrammar) printToken(depth string, token interface{}, out func(value string)) {
	style := grammar.parser.style
	
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
				grammar.printToken(newDepth, token, out)
				if k < len-1 {
					out(style.Separator)
					out(style.LineBreak)
				}
			}
		}

	} else {
		out(depth)
		out(style.ScalarPrefix)
		style.printScalar(token, out)
	}
}

func (grammar PropertyGrammar) Print(token interface{}, out func(value string)) {
	grammar.printToken("", token, out)
	out (string(NEWLINE))
}

func NewPropertyGrammar() Grammar {
	// https://en.wikipedia.org/wiki/.properties
	style := Style{"", "", "",
		"", "", "", "", "",
		" = ", "\n", "true", "false", '#', " = "}
	return PropertyGrammar{NewSExpressionParser(style)}
}

// TODO json xml postfix

/////////////////////////////////////////////////////////////////////////////
// JSON Grammar
/////////////////////////////////////////////////////////////////////////////

type JSONGrammar struct {
	parser SExpressionParser
}

func (grammar JSONGrammar) Name() string {
	return "JSON"
}

func (grammar JSONGrammar) FileSuffix() string {
	return ".json"
}

func (grammar JSONGrammar) Parse(context * ParserContext) {
	grammar.parser.Parse(context)
}

func (grammar JSONGrammar) Print(token interface{}, next func(value string)) {
	grammar.parser.style.PrettyPrint(token, next)
}

func NewJSONGrammar() Grammar {
	style := Style{"", "", "  ",
		"[", "]", "{", "}", ":",
		",", "\n", "true", "false", '%', ""} // prolog, sql '--' for 
	return JSONGrammar{NewSExpressionParser(style)}
}

