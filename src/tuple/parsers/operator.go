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

import "strings"

/////////////////////////////////////////////////////////////////////////////
//  An operator grammar
/////////////////////////////////////////////////////////////////////////////

//  A Grammar for handling infix expressions for arithmetic.
//
// See technical details please see:
//   https://en.wikipedia.org/wiki/Operator-precedence_grammar
//
// For implementation see:
//    https://en.wikipedia.org/wiki/Shunting-yard_algorithm
type OperatorGrammar struct {
	context Context
	operators * Operators
	operatorStack []Atom
	Values Tuple
	wasOperator bool
}

func NewOperatorGrammar(context Context, operators * Operators) OperatorGrammar {
	return OperatorGrammar{context, operators, make([]Atom, 0), NewTuple(),true}
}

func (stack * OperatorGrammar) pushOperator(token Atom) {
	Verbose(stack.context,"PUSH OPERATOR\t '%s'", token.Name)
	stack.operatorStack = append(stack.operatorStack, token)
}

// Remove top from operator stack
func (stack * OperatorGrammar) popOperator() {
	lo := len(stack.operatorStack) // This could be passed in for efficiency
	Verbose(stack.context,"POP OPERATOR\t '%s'", stack.operatorStack[lo-1].Name)
	stack.operatorStack = stack.operatorStack[:lo-1]
}

func (stack * OperatorGrammar) reducePrefix(top Atom) {

	var name Atom
	switch top.Name {
	case "_prefix_minus": name = Atom{"-"}
	case "_prefix_plus": name = Atom{"+"}
	default: name = top
	}

	values := &(stack.Values.List)
	lv := stack.Values.Length()
	val1 := (*values) [lv - 1]
	stack.Values.List = append((*values)[:lv-1], NewTuple(name, val1))
	Verbose(stack.context," REDUCE:\t%s\t'%s'\n", name.Name, val1)
}

// Replace top of value stack with an expression
func (stack * OperatorGrammar) reduceOperatorExpression(top Atom) int {

	values := &(stack.Values.List)
	lv := stack.Values.Length()
	if strings.HasPrefix(top.Name, "_prefix_") {
		stack.popOperator()
		stack.reducePrefix(top)
	} else {
		if top == SPACE_ATOM { // TODO Could in principle generalize to make any binary operator n-ary
			Verbose(stack.context,"** REDUCE\t'%s'\n", top.Name)
			stack.popOperator()
			count := 2
			for {
				ll := len(stack.operatorStack)
				if ll == 0 {
					break
				}
				nextTop := stack.operatorStack[ll - 1]
				Verbose(stack.context,"*** REDUCE\t'%s'\n", top.Name)
				if nextTop != SPACE_ATOM {
					break
				}
				stack.popOperator()
				count += 1
			}
			// TODO the following is not efficient and should be replaced with a slice: tuple := NewTuple(args...)
			tuple := NewTuple()
			//if top != SPACE_ATOM {
			//	tuple.Append(top)
			//}
			args := (*values) [lv-count:]
			for _,v := range args {
				tuple.Append(v)
			}
			Verbose(stack.context," REDUCE:\t'SPACE'\t'%s'\t'%s'\t...%d...   \n", tuple.List[0], tuple.List[1], tuple.Length()) //, tuple.List==(*values))
			stack.Values.List = append((*values)[:lv-count], tuple) // TODO should not need a special case
			return count - 2
		} else {
			val1 := (*values) [lv - 2]
			val2 := (*values) [lv - 1]
			stack.popOperator()
			stack.Values.List = append((*values)[:lv-2], NewTuple(top, val1, val2))
			Verbose(stack.context," REDUCE:\t'%s'\t'%s'\t'%s'\n", top.Name, val1, val2)
		}
	}
	return 0
}

func (stack * OperatorGrammar) PushValue(value Value) {
	if ! stack.wasOperator {
		stack.PushOperator(SPACE_ATOM)
	}
	Verbose(stack.context,"PUSH VALUE\t'%s'\n", value)
	stack.Values.Append(value)
	stack.wasOperator = false
}

func (stack * OperatorGrammar) OpenBracket(token Atom) {

	Verbose(stack.context,"OPEN '%s'", token.Name)
	if ! stack.wasOperator {
		stack.PushOperator(SPACE_ATOM)
	}
	stack.pushOperator(token)
	stack.wasOperator = true
}

func (stack * OperatorGrammar) CloseBracket(token Atom) {
	Verbose(stack.context,"CLOSE '%s'", token.Name)
	lo := len(stack.operatorStack)
	if lo == 0 || stack.wasOperator {
		if lo > 0 && stack.operators.IsOpenBracket(stack.operatorStack[lo-1]) {  // '()'  Empty list, is this always okay
			stack.wasOperator = false
			stack.popOperator()
			//lv := stack.Values.Length()
			values := stack.Values.List
			stack.Values.List = append(values, NewTuple())
			Verbose(stack.context," REDUCE:\t'()'\n")
			return
		} else {
			// TODO this should return an error
			UnexpectedCloseBracketError(stack.context,token.Name)
			return
		}
	}

	stack.wasOperator = false
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		if stack.operators.IsOpenBracket(top) {
			if ! stack.operators.MatchBrackets(top, token) {
				Error(stack.context,"Expected close bracket '%s' but found '%s'", top.Name, token.Name)
			}
			stack.popOperator()
			return
		} else {
			index -= stack.reduceOperatorExpression(top)
		}
	}
}

// Signal end of input
func (stack * OperatorGrammar) EndOfInput(next Next) {
	Verbose(stack.context,"EOF")

	empty := stack.Values.Length() == 0 && len(stack.operatorStack) == 0
	if empty {
		return
	}
	if stack.wasOperator {
		UnexpectedEndOfInputErrorBracketError(stack.context)
		return
	}
	lo := len(stack.operatorStack)
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		if stack.operators.IsOpenBracket(top) {
			stack.popOperator()
		} else {
			index -= stack.reduceOperatorExpression(top)
		}
	}
	// TODO this is a hack to handle space separated expressions: 1+2 3*4 5
	if len(stack.Values.List) == 1 {
		next (stack.Values.List[0])
	} else {
		next (stack.Values)
	}
	stack.Values = NewTuple()
	stack.operatorStack = make([]Atom, 0)
	stack.wasOperator = true
}

func (stack * OperatorGrammar) PushOperator(atom Atom) {
	Verbose(stack.context,"*PushOperator '%s'", atom.Name)

	prefixOperator, ok := stack.operators.prefix[atom.Name]
	nextIsPrefix := stack.wasOperator && ok
	if atom != SPACE_ATOM {
		if nextIsPrefix {
			atom = prefixOperator
		}
		atomPrecedence := stack.operators.Precedence(atom)
		lo := len(stack.operatorStack)
		for index := lo-1 ; index >= 0; index -= 1 {
			top := stack.operatorStack[index]
			topIsPrefix := strings.HasPrefix(top.Name, "_prefix_")
			Verbose(stack.context, "IsPrefix %s %s", top, topIsPrefix)
			if !nextIsPrefix && (topIsPrefix || stack.operators.IsOpenBracket(top)) {
				break
			} else if nextIsPrefix && topIsPrefix {
				break
			} else if stack.operators.Precedence(top) >= atomPrecedence {
				Verbose(stack.context,"* PushOperator - Reduce '%s'", top)
				index -= stack.reduceOperatorExpression(top)
			} else {
				break
			}
		}
	}
	if ! nextIsPrefix && stack.wasOperator {
		Error(stack.context,"Unexpected binary operator '%s'", atom.Name)
	} else {
		if atom.Name == "." {
			atom = CONS_ATOM
		}
		stack.pushOperator(atom)
		// TODO postfix
		stack.wasOperator = true
	}
}


/////////////////////////////////////////////////////////////////////////////
//  Operators
/////////////////////////////////////////////////////////////////////////////

// A table of operators
type Operators struct {
	Style
	precedence map[string]int
	prefix map[string]Atom
	postfix map[string]Atom
	brackets map[string]string
	closeBrackets map[string]string
}

func NewOperators(style Style) Operators {
	atoms := make(map[string]Atom, 0)
	strings := make(map[string]string, 0)
	return Operators{style, make(map[string]int, 0), atoms, atoms, strings, strings}
}

// Iterates through all operators, this is mainly for testing
func (operators *Operators) Forall(next func (value string)) {
	for k, _ := range operators.precedence {
		next(k)
	}
	/* TODO for k, v := range operators.brackets {
		next(k)
		next(v)
	}*/
	for k, v := range operators.prefix {
		next(k)
		next(v.Name)
	}
}

func (operators *Operators) AddInfix(operator string, precedence int) {
	(*operators).precedence[operator] = precedence
}

func (operators *Operators) AddPrefix(operator string, name string, precedence int) {
	operators.prefix[operator] = Atom{name}
	operators.precedence[name] = precedence
}

func (operators *Operators) AddPostfix(operator string, name string, precedence int) {
	operators.postfix[operator] = Atom{name}
	operators.precedence[name] = precedence
}

func (operators *Operators) AddBracket(open string, close string) {
	(*operators).brackets[open] = close
	(*operators).closeBrackets[close] = open
}

func (operators *Operators) Precedence(token Atom) int {
	value, ok := (*operators).precedence[token.Name]
	if ok {
		return value
	}
	return -1
}

// TODO generalize
func (operators *Operators) IsOpenBracket(atom Atom) bool {
	token := atom.Name
	_, ok := operators.brackets[token]
	return ok
}

func (operators *Operators) MatchBrackets(open Atom, close Atom) bool {
	expectedClose, ok := operators.brackets[open.Name]
	return ok && expectedClose == close.Name
}

/////////////////////////////////////////////////////////////////////////////
// Printer
/////////////////////////////////////////////////////////////////////////////

func (printer Operators) PrintNullaryOperator(depth string, atom Atom, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom), out)
}

func (printer Operators) PrintUnaryOperator(depth string, atom Atom, value Value, out StringFunction) {  // Prefix and Postfix???
	PrintTuple(&printer, depth, NewTuple(atom, value), out)
}

func (printer Operators) PrintBinaryOperator(depth string, atom Atom, value1 Value, value2 Value, out StringFunction) {  // TODO binary to infix

	if _, ok := printer.precedence[atom.Name]; ok {
		out(printer.Style.Open)
		newDepth := depth + "  "
		printer.PrintSuffix(newDepth, out)
		
		PrintExpression(printer, newDepth, value1, out)

		printer.PrintIndent(newDepth, out)
		out(atom.Name)
		printer.PrintSuffix(newDepth, out)

		PrintExpression(printer, newDepth, value2, out)

		printer.PrintIndent(depth, out)
		out(printer.Style.Close)
	} else {
		PrintTuple(&printer, depth, NewTuple(atom, value1, value2), out)
	}
}
