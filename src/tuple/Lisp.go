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

/////////////////////////////////////////////////////////////////////////////
// Lisp with conventional Prefix Grammar
/////////////////////////////////////////////////////////////////////////////

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

func (grammar LispGrammar) Print(object interface{}, out func(value string)) {
	PrintExpression(&(grammar.parser.style), "", object, out)
}

func NewLispGrammar() Grammar {
	style := Style{"", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, "", "", ".",
		"", "\n", "true", "false", ';', ""}
	return LispGrammar{NewSExpressionParser(style)}
}

/////////////////////////////////////////////////////////////////////////////
// Lisp with an infix notation Grammar
/////////////////////////////////////////////////////////////////////////////

type LispWithInfixGrammar struct {
	parser SExpressionParser
	operators Operators
}

func (grammar LispWithInfixGrammar) Name() string {
	return "Lisp with infix"
}

func (grammar LispWithInfixGrammar) FileSuffix() string {
	return ".infix"
}

func (grammar LispWithInfixGrammar) Parse(context * ParserContext) {

	operators := grammar.operators
	operatorGrammar := NewOperatorGrammar(context, &operators)
	for {
		token, err := grammar.parser.GetNext(context)
		if err == io.EOF {
			operatorGrammar.EOF(context.next)
			break
		}
		if token == "(" {
			operatorGrammar.OpenBracket(Atom{"("})
		} else if token == ")" {
			operatorGrammar.CloseBracket(Atom{")"})

		} else if atom, ok := token.(Atom); ok {
			if operators.Precedence(atom) != -1 {
				operatorGrammar.PushOperator(atom)
			} else if operators.IsOpenBracket(atom) {
				operatorGrammar.OpenBracket(atom)
			} else if operators.IsCloseBracket(atom) {
				operatorGrammar.CloseBracket(atom)
			} else {
				
				operatorGrammar.PushValue(atom)
			}
		} else {
			operatorGrammar.PushValue(token)
		}
	}
}

func (grammar LispWithInfixGrammar) Print(token interface{}, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

func NewLispWithInfixGrammar() Grammar {
	style := Style{"", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, "", "", ".", 
		"", "\n", "true", "false", ';', ""}
	operators := NewOperators(style)
	operators.AddStandardCOperators()
	return LispWithInfixGrammar{NewSExpressionParser(style), operators}
}

/////////////////////////////////////////////////////////////////////////////
// Conventional arithmetic expression grammar Lisp with an infix notation Grammar
/////////////////////////////////////////////////////////////////////////////

type InfixGrammar struct {
	parser SExpressionParser
	operators Operators
}

func (grammar InfixGrammar) Name() string {
	return "Lisp with infix"
}

func (grammar InfixGrammar) FileSuffix() string {
	return ".infix"
}

func (grammar InfixGrammar) Parse(context * ParserContext) {

	operators := grammar.operators
	operatorGrammar := NewOperatorGrammar(context, &operators)
	for {
		token, err := grammar.parser.GetNext(context)
		if err == io.EOF {
			operatorGrammar.EOF(context.next)
			break
		}
		if token == "(" {
			operatorGrammar.OpenBracket(Atom{"("})
		} else if token == ")" {
			operatorGrammar.CloseBracket(Atom{")"})

		} else if atom, ok := token.(Atom); ok {
			if operators.Precedence(atom) != -1 {
				operatorGrammar.PushOperator(atom)
			} else if operators.IsOpenBracket(atom) {
				operatorGrammar.OpenBracket(atom)
			} else if operators.IsCloseBracket(atom) {
				operatorGrammar.CloseBracket(atom)
			} else {
				
				operatorGrammar.PushValue(atom)
			}
		} else {
			operatorGrammar.PushValue(token)
		}
	}
}

func (grammar InfixGrammar) Print(token interface{}, next func(value string)) {
	// TODO
	PrintExpression(&(grammar.operators), "", token, next)
	//grammar.parser.Print(token, next)
}

func NewInfixGrammar() Grammar {
	style := Style{"", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, "", "", ".", 
		"", "\n", "true", "false", ';', ""}
	operators := NewOperators(style)
	operators.AddStandardCOperators()

	return InfixGrammar{NewSExpressionParser(style), operators}
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
		ch, err := readRune(context, grammar.parser)
		switch {
		case err != nil: return "", err
		case err == io.EOF: return "", nil
		case ch == NEWLINE: return string(NEWLINE), nil
		case unicode.IsSpace(ch): break
		case ch == parser.style.OneLineComment:
			// TODO ignore for now
			//return string(ch), nil
		case ch == parser.openChar : return parser.style.Open, nil
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
			context.UnexpectedCloseBracketError (parser.style.Close)
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
			grammar.parser.printObject("", token, out)
			if k < len-1 {
				out(style.Indent)
				out(style.Separator)
			}
		}
	} else {
		grammar.parser.printObject("", token, out)
	}
	out (string(NEWLINE))
}

func NewTclGrammar() Grammar {
	style := Style{"", "", "  ",
		OPEN_BRACE, CLOSE_BRACE, OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, ":",
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
	grammar.parser.Print(token, next)
}

func NewJmlGrammar() Grammar {
	style := Style{"\n", "", "  ", OPEN_BRACE, CLOSE_BRACE, "", "\n", "true", "false", '#'}
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
		token, err := parser.GetNext(context)
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
		token, err := parser.GetNext(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			context.Error ("'%s'", err)
			return
		case token == parser.style.Close:
			context.UnexpectedCloseBracketError (parser.style.Close)
		default:
			if atom,ok := token.(Atom); ok {
				bracket, err := parser.GetNext(context)
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
	grammar.parser.Print(token, next)
}

func NewTupleGrammar() Grammar {
	style := Style{"", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, "", "", ":",
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

func (grammar Yaml) printObject(depth string, token interface{}, out func(value string)) {

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
				grammar.printObject(newDepth, token, out)
				if k < len-1 {
					out(style.LineBreak)
				}
			}
		}

	} else {
		out(depth)
		out(style.ScalarPrefix)
		grammar.parser.printScalar(token, out)
	}
}

func (grammar Yaml) Print(token interface{}, out func(value string)) {
	grammar.printObject("", token, out)
	out (string(NEWLINE))
}

func NewYamlGrammar() Grammar {
	style := Style{"---\n", "...\n", "  ", 
		":", "", OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, "",
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
func (grammar Ini ) printObject(depth string, key string, token interface{}, out func(value string)) {

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
		grammar.parser.printScalar(token, out)
	}
}

func (grammar Ini) Print(token interface{}, out func(value string)) {
	grammar.printObject("", "", token, out)
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

func (grammar PropertyGrammar) printObject(depth string, token interface{}, out func(value string)) {
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
		grammar.parser.printScalar(token, out)
	}
}

func (grammar PropertyGrammar) Print(token interface{}, out func(value string)) {
	grammar.printObject("", token, out)
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
	grammar.parser.Print(token, next)
}

func NewJSONGrammar() Grammar {
	style := Style{"", "", "  ",
		OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, OPEN_BRACE, CLOSE_BRACE, ":",
		",", "\n", "true", "false", '%', ""} // prolog, sql '--' for 
	return JSONGrammar{NewSExpressionParser(style)}
}

