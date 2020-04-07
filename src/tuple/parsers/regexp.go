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
package parsers

import "tuple"
import "io"
import "errors"

// To provide an example of how the OperatorGrammar can be used for a variety of purposes,
// in this case to parse regular expressions.
// TODO finish off and provide a full set of test cases.
func NewRegexpGrammar(context Context) OperatorGrammar {

	style := NewStyle("", "", "  ",
		OPEN_BRACKET, CLOSE_BRACKET, OPEN_BRACE, CLOSE_BRACE, ":",
		",", "\n", "true", "false", '#', "")

	operators := NewOperators(style)

	operators.AddBracket(OPEN_BRACKET, CLOSE_BRACKET)
	operators.AddBracket(OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET)
	operators.AddBracket(OPEN_BRACE, CLOSE_BRACE)
	operators.AddPostfix("*", 105)
	operators.AddPostfix("+", 105)
	operators.AddPostfix("?", 105)
	operators.AddInfix("-", 100)
	operators.AddInfix("|", 5)
	operators.AddInfix(" ", 50)

	operatorGrammar := NewOperatorGrammar(context, &operators)
	return operatorGrammar;
}

//  Parses a stream of 'runes' that represent an Regular Expression
//  and return a parse tree/AST
func ParseRegexp(context Context) Value {

	grammar := NewRegexpGrammar(context)

	var result Value = tuple.EMPTY
	for {  // TODO read runes
		v, err := context.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			Error(context, "on reading: %s", err)
			break
		}
		s := string(v)
		switch v {
		case '\\': 
			v, err = context.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				Error(context, "on reading: %s", err)
				break
			}
			s := string(v)
			grammar.PushValue(String(s))
		case '[', '(', '{': grammar.OpenBracket(Tag{s})
		case ']', ')', '}': grammar.CloseBracket(Tag{s})
		case '-', '|', '*', '+', '?', '^', '$': grammar.PushOperator(Tag{s})
		default: grammar.PushValue(String(s))
		}
	}
	grammar.EndOfInput(func (value Value) {
		result = value
	})

	return result
}

// A Toy regular expression matcher.
//
// This implementation serves as an example using a dynamic implementation that is easy to implement but not
// as fast as generating a state transition table.
//
// TODO Allow a callback function for matching  regexp to provide a 'lex' functionality and use this to replace clex.go.
// 
func MatchRegexp (scanner io.RuneScanner, value Value) error {
	return MatchRegexpAndCallback(scanner, value, func(_ string) {})
}

// Executes the given callback when matching certain patterns, to provide 'lex' like functionality
func MatchRegexpAndCallback (scanner io.RuneScanner, value Value, callback func (token string)) error {

	tuple, ok := value.(Tuple)
	if ok {
		head, ok := tuple.Get(0).(Tag)
		if ! ok {
			for _, v := range tuple.List {
				err := MatchRegexp(scanner, v)
				if err != nil {
					return err
				}
			}
			return nil
		} else {
			switch head.Name {
			case "_callback_":
				err := MatchRegexp(scanner, tuple.Get(1))
				if err == nil {
					// TODO matching group
					// TODO callback(token)
				}
				return err
			case "|":
				err := MatchRegexp(scanner, tuple.Get(1))
				if err == nil {
					return nil
				}
				err = MatchRegexp(scanner, tuple.Get(2))
				return err
			case "*":
				for {
					err := MatchRegexp(scanner, tuple.Get(1))
					if err != nil {
						scanner.UnreadRune()
						return nil
					}
				}
			case "+":
				err := MatchRegexp(scanner, tuple.Get(1))
				if err != nil {
					return nil
				}
				for {
					err := MatchRegexp(scanner, tuple.Get(1))
					if err != nil {
						return nil
					}
				}
			case "?":
				err := MatchRegexp(scanner, tuple.Get(1))
				return err
			case "-":
				next, _, err := scanner.ReadRune()
				if err == io.EOF {
					return err
				}
				if err != nil {
					return err
				}
				lower := rune(tuple.Get(1).(String)[0]) // TODO check
				upper := rune(tuple.Get(2).(String)[0])
				if lower <= next && next <= upper {
					return nil
				}
				return errors.New("mismatch")
			// TODO case "^":
			// TODO case "$":
			default:
				return errors.New("Unexpected: " + head.Name)
			}
		}
	} else {
		next, _, err := scanner.ReadRune()
		if err == io.EOF {
			return err
		}
		if err != nil {
			return err
		}
		tag := value.(String)
		expected := rune(tag[0])
		if expected == '.' {
			return nil
		}
		if next == expected {
			return nil
		}
		return errors.New("mismatch")
	}
}
