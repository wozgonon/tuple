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

func (grammar LispGrammar) Parse(context Context, next Next) {
	grammar.parser.Parse(context, next)
}

func (grammar LispGrammar) Print(object Value, out func(value string)) {
	PrintExpression(grammar.parser.lexer, "", object, out)
}

func NewLispGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, "", "", CONS_OPERATOR,
		"", "\n", "true", "false", ';', "")
	return LispGrammar{NewSExpressionParser(style)}
}

/////////////////////////////////////////////////////////////////////////////
// Lisp with an infix notation Grammar
/////////////////////////////////////////////////////////////////////////////

type LispWithInfixGrammar struct {
	style Style
	operators Operators
}

func (grammar LispWithInfixGrammar) Name() string {
	return "Lisp with infix"
}

func (grammar LispWithInfixGrammar) FileSuffix() string {
	return ".infix"
}

func (grammar LispWithInfixGrammar) Parse(context Context, next Next) {

	operators := grammar.operators
	operatorGrammar := NewOperatorGrammar(context, &operators)
	for {
		err := grammar.style.GetNext(context,
			func (open string) {
				operatorGrammar.OpenBracket(Atom{open})
			},
			func (close string) {
				operatorGrammar.CloseBracket(Atom{close})
			},
			func (atom Atom) {
				if operators.Precedence(atom) != -1 {
					operatorGrammar.PushOperator(atom)
				} else {
					operatorGrammar.PushValue(atom)
				}
			},
			func (literal Value) {
				operatorGrammar.PushValue(literal)
			})
		
		if err == io.EOF {
			operatorGrammar.EOF(next)
			break
		}
	}
}

func (grammar LispWithInfixGrammar) Print(token Value, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

var LispWithInfixStyle Style =NewStyle("", "", "  ",
	OPEN_BRACKET, CLOSE_BRACKET, "", "", CONS_OPERATOR, 
	"", "\n", "true", "false", ';', "")

func NewLispWithInfixGrammar() Grammar {
	style := LispWithInfixStyle
	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.Add(CONS_OPERATOR, 105) // CONS Operator
	return LispWithInfixGrammar{style, operators}
}

/////////////////////////////////////////////////////////////////////////////
// Conventional arithmetic expression grammar,
// with an infix notation and conventional function call syntax a(b, c, ...)
/////////////////////////////////////////////////////////////////////////////

type InfixExpressionGrammar struct {
	style Style
	operators Operators
}

func (grammar InfixExpressionGrammar) Name() string {
	return "Expression with infix operators"
}

func (grammar InfixExpressionGrammar) FileSuffix() string {
	return ".expr"
}

func (grammar InfixExpressionGrammar) Parse(context Context, next Next) {

	open := grammar.style.Open
	operators := grammar.operators
	operatorGrammar := NewOperatorGrammar(context, &operators)
	for {
		err := grammar.style.GetNext(context,
			func (open string) {
				operatorGrammar.OpenBracket(Atom{open})
			},
			func (close string) {
				operatorGrammar.CloseBracket(Atom{close})
			},
			func(atom Atom) {
				if operators.Precedence(atom) != -1 {
					operatorGrammar.PushOperator(atom)
				} else {
					ch, err := context.ReadRune()
					if err == io.EOF {
						operatorGrammar.PushValue(atom)
						//operatorGrammar.EOF(next)
						return // TODO eof
					}
					if err != nil {
						// TODO
						return
					}
					if ch == grammar.style.openChar {
						operatorGrammar.OpenBracket(Atom{open})
					} else {
						context.UnreadRune()
					}
					operatorGrammar.PushValue(atom)
				}
			},
			func (literal Value) {
				operatorGrammar.PushValue(literal)
			})
	
		if err == io.EOF {
			operatorGrammar.EOF(next)
			break
		}
		if err != nil {
			// TODO 
			break
		}
	}
}

func (grammar InfixExpressionGrammar) Print(token Value, next func(value string)) {
	PrintExpression(&(grammar.operators), "", token, next)
}

func NewInfixExpressionGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, "", "", ":",
		",", "\n", "true", "false", '%', "") // prolog, sql '--' for 

	operators := NewOperators(style)
	AddStandardCOperators(&operators)
	operators.Add(CONS_OPERATOR, 105) // CONS Operator

	return InfixExpressionGrammar{style, operators}
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
// It would be nice not to have to put quotes around strings, in particular working interactice on the command line (CLI) with a shell such as bash or
// dos command line.  It is very nice not to have to include any boiler plate when just entering commands but this quickly becomes very awkward to add any more complex syntax.

type Tcl struct {
	style Style
}

func (grammar Tcl) Name() string {
	return "Tcl"
}

func (grammar Tcl) FileSuffix() string {
	return ".tcl"
}

func (grammar Tcl ) readCommandString(context Context, token string) (string, error) {
	return ReadString(context, token, true, func (ch rune) bool {
		return ! unicode.IsSpace(ch) && string(ch) != grammar.style.Close && string(ch) != grammar.style.Open && ch != '$'
	})
}

func (grammar Tcl) getNextCommandShell(context Context) (Value, error) {

	style := grammar.style
	for {
		ch, err := readRune(context, grammar.style)
		switch {
		case err != nil: return nil, err
		case err == io.EOF: return nil, nil
		case ch == NEWLINE: return String(NEWLINE), nil
		case unicode.IsSpace(ch): break
		case ch == style.OneLineComment:
			// TODO ignore for now
			//return string(ch), nil
		case ch == style.openChar : return String(style.Open), nil  // Atom???
		case ch == style.closeChar : return String(style.Close), nil
		case ch == '"' :  return ReadCLanguageString(context)
		case ch == '.' || unicode.IsNumber(ch): return ReadNumber(context, string(ch))    // TODO minus
		case ch == '$':
			value, err := grammar.readCommandString(context, "")
			if err != nil {
				return nil, err
			}
			return Atom{value}, nil
		case unicode.IsGraphic(ch):
			value, err := grammar.readCommandString(context, string(ch))
			if err != nil {
				return nil, err
			}
			return String(value), nil
		case unicode.IsControl(ch): Error(context, "Error control character not recognised '%d'", ch)
		default: Error(context, "Error character not recognised '%d'", ch)
		}
	}
}

func (grammar Tcl) parseCommandShellTuple(context Context, tuple *Tuple) (error) {

	style := grammar.style
	for {
		token, err := grammar.getNextCommandShell(context)
		switch {
		case err != nil:
			Error(context, "parsing %s", err);
			return err /// ??? Any need to return
		case token == String(style.Close):
			return nil
		case token == String(style.Open):
			subTuple := NewTuple()
			err := grammar.parseCommandShellTuple(context, &subTuple)
			if err == io.EOF {
				Error(context,"Missing close bracket")
				return err
			}
			if err != nil {
				return err
			}
			tuple.Append(subTuple)
		case token == String(string(NEWLINE)):
		default:
			tuple.Append(token)
		}
	}
}

func (grammar Tcl) printObject1(depth string, token Value, out func(value string)) {

	style := grammar.style
	if tuple, ok := token.(Tuple); ok {

		len := len(tuple.List)
		out(depth)
		if len == 0 {
			out(style.Open)
			out(style.Close)
			return
		}
		newDepth := depth + style.Indent
		head := tuple.List[0]
		atom, ok := head.(Atom)
		first := ok // && style.indentOnly()
		if first {
			out(atom.Name)
		} else if tuple.IsCons() {
			grammar.printObject1(depth, tuple.List[1], out)
			if _, ok = tuple.List[2].(Tuple); ok {
				out (" ")
				out(style.KeyValueSeparator)
				out (style.LineBreak)
				grammar.printObject1(newDepth, tuple.List[2], out)
			} else {
				out (" ")
				out(style.KeyValueSeparator)
				out (" ")
				PrintScalar(style, "", token, out)
				//grammar.style.printScalar(tuple.List[2], out)
			}
			return
		}
		tuple1, ok := tuple.List[0].(Tuple)
		cons := ok && tuple1.IsCons()
		if cons {
			// TODO Need a way to differentiate between [] and {}
			out(style.Open2)
		} else {
			out(style.Open)
		}
		out(style.LineBreak)
		for k, token := range tuple.List {
			grammar.printObject1(newDepth, token, out)
			if ! first && k < len-1 {
				out(style.Separator)
				out(style.LineBreak)
			}
		}
		out(style.LineBreak)
		out(depth)
		if cons {
			out(style.Close2)
		} else {
			out(style.Close)
		}
	} else {
		out(depth)
		out(style.ScalarPrefix)
		PrintScalar (grammar.style, "", token, out)
	}
}

func (grammar Tcl) Parse(context Context, next Next) {

	style := grammar.style

	resultTuple := NewTuple()
	for {
		token, err := grammar.getNextCommandShell(context)
		switch {
		case err == io.EOF:
			return
		case err != nil:
			Error(context,"'%s'", err)
			return
		case token == String(string(NEWLINE)):
			l := resultTuple.Length()
			Verbose(context,  "Newline length of tuple=%d", l)
			switch l {
			case 0: // Ignore
			case 1:
				first := resultTuple.List[0]
				if _, ok := first.(Atom); ok {
					next(resultTuple)
				} else {
					next(token)
				}
			default:
				next(resultTuple)
			}
			resultTuple = NewTuple()
		case token == String(string(style.OneLineComment)):
			comment, err := ReadUntilEndOfLine(context)
			if err != nil {
				return
			}
			next(comment)
		case token == String(style.Close):
			UnexpectedCloseBracketError(context,style.Close)
		case token == String(style.Open):
			subTuple := NewTuple()
			err := grammar.parseCommandShellTuple(context, &subTuple)
			if err != nil {
				return // tuple, err
			}
			resultTuple.Append(subTuple)
		default:
			Verbose(context, "Add token: '%s'", token)
			resultTuple.Append(token) 
		}
	}
}

func (grammar Tcl) Print(token Value, out func(value string)) {

	//PrintExpression(grammar.style, "", object, out)  // TODO Use Printer
	
	style := grammar.style
	if tuple, ok := token.(Tuple); ok {
		len := len(tuple.List)
		for k, token := range tuple.List {
			grammar.printObject1("", token, out)
			if k < len-1 {
				out(style.Indent)
				out(style.Separator)
			}
		}
	} else {
		grammar.printObject1("", token, out)
	}
	out (string(NEWLINE))
}

func NewTclGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_BRACE, CLOSE_BRACE, OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, ":",
		"", "\n", "true", "false", '#', "")
	return Tcl{style}
}

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

// TODO json xml postfix

/////////////////////////////////////////////////////////////////////////////
// JSON Grammar
/////////////////////////////////////////////////////////////////////////////

type JSONGrammar struct {
	Style
}

func (grammar JSONGrammar) Name() string {
	return "JSON"
}

func (grammar JSONGrammar) FileSuffix() string {
	return ".json"
}

func (grammar JSONGrammar) Parse(context Context, next Next) {
	parser := NewSExpressionParser(grammar.Style)
	parser.Parse(context, next)
}

func (grammar JSONGrammar) Print(object Value, out func(value string)) {
	PrintExpression(grammar, "", object, out)  // TODO Use Printer
}

func NewJSONGrammar() Grammar {
	style := NewStyle("", "", "  ",
		OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET, OPEN_BRACE, CLOSE_BRACE, ":",
		",", "\n", "true", "false", '%', "") // prolog, sql '--' for 
	return JSONGrammar{style}
}

func (printer JSONGrammar) PrintNullaryOperator(depth string, atom Atom, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom), out)
}

func (printer JSONGrammar) PrintUnaryOperator(depth string, atom Atom, value Value, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom, value), out)
}

func (printer JSONGrammar) PrintSeparator(depth string, out StringFunction) {
	out(printer.Style.Separator)
}

func (printer JSONGrammar) PrintBinaryOperator(depth string, atom Atom, value1 Value, value2 Value, out StringFunction) {

	if atom == CONS_ATOM {
		newDepth := depth + "  "
		printer.PrintIndent(depth, out)
		PrintExpression1(printer, newDepth, value1, out)
		out(" :")
		if _, ok := value2.(Tuple); ok {
			printer.PrintSuffix(newDepth, out)
			printer.PrintIndent(newDepth, out)
		} else {
			out (" ")
		}
		PrintExpression1(printer, newDepth, value2, out)
	} else {
		PrintTuple(&printer, depth, NewTuple(atom, value1, value2), out)
	}
}
