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
	context * ParserContext
	operators * Operators
	operatorStack []Atom
	Values Tuple
	wasOperator bool
}

func NewOperatorGrammar(context * ParserContext, operators * Operators) OperatorGrammar {
	return OperatorGrammar{context, operators, make([]Atom, 0), NewTuple(),true}
}

func (stack * OperatorGrammar) pushOperator(token Atom) {
	stack.context.Verbose("PUSH OPERATOR\t '%s'", token.Name)
	(*stack).operatorStack = append((*stack).operatorStack, token)
}

// Remove top from operator stack
func (stack * OperatorGrammar) popOperator() {
	lo := len(stack.operatorStack) // This could be passed in for efficiency
	stack.context.Verbose("POP OPERATOR\t '%s'", (*stack).operatorStack[lo-1].Name)
	(*stack).operatorStack = (*stack).operatorStack[:lo-1]
}

func (stack * OperatorGrammar) reduceUnary(top Atom) {

	var name Atom
	switch top.Name {
	case "_unary_minus": name = Atom{"-"}
	case "_unary_plus": name = Atom{"+"}
	default: name = top
	}

	values := &((*stack).Values.List)
	lv := stack.Values.Length()
	val1 := (*values) [lv - 1]
	(*stack).Values.List = append((*values)[:lv-1], NewTuple(name, val1))
	stack.context.Verbose(" REDUCE:\t%s\t'%s'\n", name.Name, val1)
}

// Replace top of value stack with an expression
func (stack * OperatorGrammar) reduceOperatorExpression(top Atom) int {

	values := &((*stack).Values.List)
	lv := stack.Values.Length()
	if strings.HasPrefix(top.Name, "_unary_") {
		stack.popOperator()
		stack.reduceUnary(top)
	} else {
		if top == SPACE_ATOM { // TODO Could in principle generalize to make any binary operator n-ary
			stack.context.Verbose("** REDUCE\t'%s'\n", top.Name)
			stack.popOperator()
			count := 2
			for {
				ll := len((*stack).operatorStack)
				if ll == 0 {
					break
				}
				nextTop := (*stack).operatorStack[ll - 1]
				stack.context.Verbose("*** REDUCE\t'%s'\n", top.Name)
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
			stack.context.Verbose(" REDUCE:\t'SPACE'\t'%s'\t'%s'\t...%d...   \n", tuple.List[0], tuple.List[1], tuple.Length()) //, tuple.List==(*values))
			(*stack).Values.List = append((*values)[:lv-count], tuple) // TODO should not need a special case
			return count - 2
		} else {
			val1 := (*values) [lv - 2]
			val2 := (*values) [lv - 1]
			stack.popOperator()
			(*stack).Values.List = append((*values)[:lv-2], NewTuple(top, val1, val2))
			stack.context.Verbose(" REDUCE:\t'%s'\t'%s'\t'%s'\n", top.Name, val1, val2)
		}
	}
	return 0
}

func (stack * OperatorGrammar) PushValue(value interface{}) {
	if ! (*stack).wasOperator {
		stack.PushOperator(SPACE_ATOM)
	}
	stack.context.Verbose("PUSH VALUE\t'%s'\n", value)
	(*stack).Values.Append(value)
	(*stack).wasOperator = false
}

func (stack * OperatorGrammar) OpenBracket(token Atom) {

	stack.context.Verbose("OPEN '%s'", token.Name)
	if ! (*stack).wasOperator {
		stack.PushOperator(SPACE_ATOM)
	}
	stack.pushOperator(token)
	(*stack).wasOperator = true
}

func (stack * OperatorGrammar) CloseBracket(token Atom) {
	stack.context.Verbose("CLOSE '%s'", token.Name)
	lo := len(stack.operatorStack)
	if lo == 0 || (*stack).wasOperator {
		if lo > 0 && stack.operators.IsOpenBracket(stack.operatorStack[lo-1]) {  // '()'  Empty list, is this always okay
			(*stack).wasOperator = false
			stack.popOperator()
			//lv := stack.Values.Length()
			values := (*stack).Values.List
			(*stack).Values.List = append(values, NewTuple())
			stack.context.Verbose(" REDUCE:\t'()'\n")
			return
		} else {
			// TODO this should return an error
			stack.context.UnexpectedCloseBracketError (token.Name)
			return
		}
	}

	(*stack).wasOperator = false
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		if stack.operators.IsOpenBracket(top) {
			if ! stack.operators.MatchBrackets(top, token) {
				stack.context.Error("Expected close bracket '%s' but found '%s'", top.Name, token.Name)
			}
			stack.popOperator()
			return
		} else {
			index -= stack.reduceOperatorExpression(top)
		}
	}
}

// Signal end of input
func (stack * OperatorGrammar) EOF(next Next) {
	stack.context.Verbose("EOF")
	if (*stack).wasOperator {
		stack.context.UnexpectedEndOfInputErrorBracketError()
		return
	}
	lo := len(stack.operatorStack)
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		if stack.operators.IsOpenBracket(top) {
			stack.popOperator()
			break
		} else {
			index -= stack.reduceOperatorExpression(top)
		}
	}
	// TODO this is a hack to handle space separated expressions: 1+2 3*4 5
	if len((*stack).Values.List) == 1 {
		next ((*stack).Values.List[0])
	} else {
		next ((*stack).Values)
	}
}

func (stack * OperatorGrammar) PushOperator(atom Atom) {
	stack.context.Verbose("*PushOperator '%s'", atom.Name)

	unaryOperator, ok := stack.operators.unary[atom.Name]
	unary := (*stack).wasOperator && ok
	if atom != SPACE_ATOM {
		if unary {
			atom = unaryOperator
		}
		atomPrecedence := stack.operators.Precedence(atom)
		lo := len(stack.operatorStack)
		for index := lo-1 ; index >= 0; index -= 1 {
			top := stack.operatorStack[index]
			if !unary && stack.operators.IsOpenBracket(top) {
				break
			} else if stack.operators.Precedence(top) >= atomPrecedence {
				stack.context.Verbose("* PushOperator - Reduce '%s'", top)
				index -= stack.reduceOperatorExpression(top)
			} else {
				break
			}
		}
	}
	if ! unary && (*stack).wasOperator {
		stack.context.Error("Unexpected binary operator '%s'", atom.Name)
	} else {
		if atom.Name == "." {
			atom = CONS_ATOM
		}
		stack.pushOperator(atom)
		// TODO postfix
		(*stack).wasOperator = true
	}
}


/////////////////////////////////////////////////////////////////////////////
//  Operators
/////////////////////////////////////////////////////////////////////////////

// A table of operators
type Operators struct {
	Style
	precedence map[string]int
	unary map[string]Atom
	brackets map[string]string
	closeBrackets map[string]string
}

func NewOperators(style Style) Operators {
	return Operators{style, make(map[string]int, 0), make(map[string]Atom, 0), make(map[string]string, 0), make(map[string]string, 0)}
}

func (operators *Operators) Add(operator string, precedence int) {
	(*operators).precedence[operator] = precedence
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

func (operators *Operators) IsCloseBracket(atom Atom) bool {
	token := atom.Name
	_, ok := operators.closeBrackets[token]
	return ok
}

func (operators *Operators) MatchBrackets(open Atom, close Atom) bool {
	expectedClose, ok := operators.brackets[open.Name]
	return ok && expectedClose == close.Name
}

// TODO generalize
func (operators *Operators) IsUnaryPrefix(token string) bool {
	_, ok := operators.brackets[token]
	return ok
}

func (operators *Operators) AddStandardCOperators() {
	operators.unary["-"] = Atom{"_unary_minus"}
	operators.unary["+"] = Atom{"_unary_plus"}
	operators.AddBracket(OPEN_BRACKET, CLOSE_BRACKET)
	operators.AddBracket(OPEN_SQUARE_BRACKET, CLOSE_SQUARE_BRACKET)
	operators.AddBracket(OPEN_BRACE, CLOSE_BRACE)
	operators.Add("_unary_plus", 110)
	operators.Add("_unary_minus", 110)
	operators.Add("^", 100)
	operators.Add("*", 90)
	operators.Add("/", 90)
	operators.Add("+", 80)
	operators.Add("-", 80)
	operators.Add("<", 60)
	operators.Add(">", 60)
	operators.Add("<=", 60)
	operators.Add(">=", 60)
	operators.Add("==", 60)
	operators.Add("!=", 60)
	operators.Add("&&", 50)
	operators.Add("||", 50)
	//operators.Add(",", 40)
	//operators.Add(";", 30)
	operators.Add(SPACE_ATOM.Name, 10)  // TODO space???
}

/////////////////////////////////////////////////////////////////////////////
// Printer
/////////////////////////////////////////////////////////////////////////////

func (printer Operators) PrintNullaryOperator(depth string, atom Atom, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom), out)
}

func (printer Operators) PrintUnaryOperator(depth string, atom Atom, value interface{}, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(atom, value), out)
}

func (printer Operators) PrintBinaryOperator(depth string, atom Atom, value1 interface{}, value2 interface{}, out StringFunction) {

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
