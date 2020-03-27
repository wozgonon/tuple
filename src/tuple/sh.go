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

// TODO json xml postfix
