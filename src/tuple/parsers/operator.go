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
//import "fmt"

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
var	COMMA_ATOM = Atom{";"}


func reduceToTuple(top Atom) bool {
	return top == SPACE_ATOM || top == COMMA_ATOM
}

// Replace top of value stack with an expression
func (stack * OperatorGrammar) reduceOperatorExpression(top Atom) int {

	values := &(stack.Values.List)
	lv := stack.Values.Length()
	name := stack.operators.Map(top)
	stack.popOperator()
	count := 0
	index := 0
	tuple := NewTuple()
	if isPrefix(top) {
		val1 := (*values) [lv - 1]
		count = 1
		tuple = NewTuple(name, val1)
		Verbose(stack.context," REDUCE:\t%s\t'%s'\n", name.Name, val1)
	} else {
		count = 2
		if reduceToTuple(top) {
			// TODO Could in principle generalize to make any binary operator n-ary
			// TODO this would work like python 3>2>1
			Verbose(stack.context,"** REDUCE\t'%s'\n", top.Name)
			for {
				ll := len(stack.operatorStack)
				if ll == 0 {
					break
				}
				nextTop := stack.operatorStack[ll - 1]
				Verbose(stack.context,"*** REDUCE\t'%s'\n", top.Name)
				if ! reduceToTuple(nextTop) {
					break
				}
				stack.popOperator()
				count += 1
			}
			// TODO the following is not efficient and should be replaced with a slice: tuple := NewTuple(args...)
			args := (*values) [lv-count:]
			for _,v := range args {
				tuple.Append(v)
			}
			Verbose(stack.context," REDUCE:\t'SPACE'\t'%s'\t'%s'\t...%d...   \n", tuple.List[0], tuple.List[1], tuple.Length()) //, tuple.List==(*values))
		} else {
			val1 := (*values) [lv - 2]
			val2 := (*values) [lv - 1]
			tuple = NewTuple(name, val1, val2)
			Verbose(stack.context," REDUCE:\t'%s'\t'%s'\t'%s'\n", name, val1, val2)
		}
		index = count - 2
	}
	stack.Values.List = append((*values)[:lv-count], tuple) // TODO should not need a special case
	return index
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

func (stack * OperatorGrammar) PushOperator(operator Atom) {
	prefixOperator, ok := stack.operators.prefix[operator.Name]
	Verbose(stack.context,"*PushOperator '%s' isPrefix=%s  (%s)", operator.Name, ok, prefixOperator)

	nextIsPrefix := stack.wasOperator && ok
	head := operator
	if true { // ! reduceToTuple(atom) {
		if nextIsPrefix {
			head = prefixOperator
		}
		atomPrecedence := stack.operators.Precedence(operator)
		lo := len(stack.operatorStack)
		for index := lo-1 ; index >= 0; index -= 1 {
			top := stack.operatorStack[index]
			topIsPrefix := isPrefix(top)
			Verbose(stack.context, "IsPrefix %s %s  precedence=%d", top, topIsPrefix, stack.operators.Precedence(top))
			if !nextIsPrefix && topIsPrefix {
				index -= stack.reduceOperatorExpression(top)
			} else if !nextIsPrefix && stack.operators.IsOpenBracket(top) {
				break
			} else if nextIsPrefix && topIsPrefix {
				break
			} else if nextIsPrefix && ! topIsPrefix {
				break
			} else if top == head && reduceToTuple(head) {
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
		Error(stack.context,"Unexpected binary operator '%s'", operator.Name)
	} else {
		stack.pushOperator(head)
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
	infix map[string]Atom
	postfix map[string]Atom
	brackets map[string]string
	closeBrackets map[string]string
	evalName map[string]Atom
}

func NewOperators(style Style) Operators {
	atoms1 := make(map[string]Atom, 0)
	atoms2 := make(map[string]Atom, 0)
	atoms3 := make(map[string]Atom, 0)
	atoms4 := make(map[string]Atom, 0)
	strings1 := make(map[string]string, 0)
	strings2 := make(map[string]string, 0)
	return Operators{style, make(map[string]int, 0), atoms1, atoms2, atoms3, strings1, strings2, atoms4}
}

const PREFIX = "_prefix_"

func isPrefix(top Atom) bool {
	// stack.operators.prefix[atom.Name]
	return strings.HasPrefix(top.Name, PREFIX)
}

func (operators *Operators) Map(top Atom) Atom {

	atom, ok := operators.evalName[top.Name]
	if ok {
		return atom
	}
	return top
}

// Iterates through all operators, this is mainly for testing
func (operators *Operators) Forall(next func (value string)) {
	for k, _ := range operators.prefix {
		next(k)
	}
	for k, _ := range operators.infix {
		next(k)
	}
	for k, _ := range operators.postfix {
		next(k)
	}
	/* TODO for k, v := range operators.brackets {
		next(k)
		next(v)
	}*/
	//for k, _ := range operators.prefix {
		//next(k)
		//next(v.Name)
	//}
}

func (operators *Operators) AddInfix(operator string, precedence int) {
	operators.AddInfix3(operator, precedence, operator)
}

func (operators *Operators) AddInfix3(operator string, precedence int, name string) {
	operators.precedence[operator] = precedence
	operators.infix[operator] = Atom{name}
	operators.evalName[name] = Atom{operator}
}

func (operators *Operators) AddPrefix(operator string, precedence int) {
	name := PREFIX + operator
	operators.prefix[operator] = Atom{name}
	operators.precedence[operator] = precedence
	operators.precedence[name] = precedence
	operators.evalName[name] = Atom{operator}
}

func (operators *Operators) AddPostfix(operator string, precedence int) {
	name := "_postfix" + operator
	operators.postfix[operator] = Atom{name}
	operators.precedence[operator] = precedence
	operators.precedence[name] = precedence
	operators.evalName[name] = Atom{operator}
}

func (operators *Operators) AddBracket(open string, close string) {
	operators.brackets[open] = close
	operators.closeBrackets[close] = open
}

func (operators *Operators) Precedence(token Atom) int {
	value, ok := operators.precedence[token.Name]
	//fmt.Printf("Precedence %s %s %s\n", token, value, ok)
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
