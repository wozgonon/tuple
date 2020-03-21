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
	values := &((*stack).Values.List)
	lv := stack.Values.Length()
	val1 := (*values) [lv - 1]
	(*stack).Values.List = append((*values)[:lv-1], NewTuple(top, val1))
	stack.context.Verbose(" REDUCE:\t%s\t'%s'\n", top.Name, val1)
}

// Replace top of value stack with an expression
func (stack * OperatorGrammar) reduceOperatorExpression(top Atom) {

	values := &((*stack).Values.List)
	lv := stack.Values.Length()
	if strings.HasPrefix(top.Name, "_unary_") {
		stack.reduceUnary(top)
	} else {
		val1 := (*values) [lv - 2]
		val2 := (*values) [lv - 1]
		if top == SPACE_ATOM {
			// TODO This breaks 'eval'
			(*stack).Values.List = append((*values)[:lv-2], NewTuple(val1, val2)) // TODO should not need a special case
			stack.context.Verbose(" REDUCE:\tSPACE\t'%s'\t'%s'\n", val1, val2)
		} else {
			(*stack).Values.List = append((*values)[:lv-2], NewTuple(top, val1, val2))
			stack.context.Verbose(" REDUCE:\t'%s'\t'%s'\t'%s'\n", top.Name, val1, val2)
		}
	}
}

func (stack * OperatorGrammar) PushValue(value interface{}) {
	if ! (*stack).wasOperator {
		stack.PushOperator(SPACE_ATOM)
		//stack.context.Error("Unexpected value '%s' after value %s\n", value)
		// TODO handle this situation, flush current contents or add a comma operator
		//	return
	}
	stack.context.Verbose("PUSH VALUE\t'%s'\n", value)
	(*stack).Values.Append(value)
	(*stack).wasOperator = false
}

func (stack * OperatorGrammar) OpenBracket(token Atom) {

	if ! (*stack).wasOperator {
		stack.PushOperator(SPACE_ATOM)
	}
	stack.pushOperator(token)
	(*stack).wasOperator = true
}

func (stack * OperatorGrammar) CloseBracket(token Atom) {
	// TODO should this return an error
	lo := len(stack.operatorStack)
	if lo == 0 || (*stack).wasOperator {
		stack.context.UnexpectedCloseBracketError (token.Name)
		return
	}

	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		stack.popOperator()

		if stack.operators.IsOpenBracket(top) {
			if ! stack.operators.MatchBrackets(top, token) {
				stack.context.Error("Expected close bracket '%s' but found '%s'", top.Name, token.Name)
				return
			}
			//stack.reduceUnary(top)
			return
		} else {
			stack.reduceOperatorExpression(top)
		}
	}
}

// Signal end of input
func (stack * OperatorGrammar) EOF(next Next) {
	if (*stack).wasOperator {
		stack.context.UnexpectedEndOfInputErrorBracketError()
		return
	}

	lo := len(stack.operatorStack)
	//stack.context.Verbose("OpStack Len=%d\n", lo)
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		stack.popOperator()

		if stack.operators.IsOpenBracket(top) {
			//stack.reduceUnary(top)
			//val1 := (*values) [lv - 1]
			//(*stack).Values.List = append((*values)[:lv-1], val1)
			break
		} else {
			stack.reduceOperatorExpression(top)
		}
	}
	// TODO this is a hack to handle space separated expressions: 1+2 3*4 5
	if len((*stack).Values.List) == 1 {
		next ((*stack).Values.List[0])
	} else {
		next ((*stack).Values)
	}
	//assert len(stack.values) == 1
	// TODO this is a hack to handle space separated expressions: 1+2 3*4 5
	//for _, value := range (*stack).Values.List {
	//	next (value)
	//}
	
}

func (stack * OperatorGrammar) PushOperator(atom Atom) {
	values := &((*stack).Values.List)

	unaryOperator, ok := stack.operators.unary[atom.Name]
	unary := (*stack).wasOperator && ok
	if unary {
		atom = unaryOperator
	}
	atomPrecedence := stack.operators.Precedence(atom)
	lo := len(stack.operatorStack)
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		stack.context.Verbose("PushOperator\t%s\ttop=%s\t%d", atom, top, len(*values))
		if !unary && stack.operators.IsOpenBracket(top) {
			//lv := stack.Values.Length()
			//val1 := (*values) [lv - 1]
			//(*stack).Values.List = append((*values)[:lv-1], val1)
			break
		} else if stack.operators.Precedence(top) >= atomPrecedence {
			stack.popOperator()
			stack.reduceOperatorExpression(top)
		} else {
			break
		}
	}
	stack.pushOperator(atom)
	
	// TODO postfix
	(*stack).wasOperator = true
}


/////////////////////////////////////////////////////////////////////////////
//  Operators
/////////////////////////////////////////////////////////////////////////////

// A table of operators
type Operators struct {
	precedence map[string]int
	unary map[string]Atom
	brackets map[string]string
}

func NewOperators() Operators {
	return Operators{make(map[string]int, 0), make(map[string]Atom, 0), make(map[string]string, 0)}
}

func (operators *Operators) Add(operator string, precedence int) {
	(*operators).precedence[operator] = precedence
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
	//return token == "(" || token == "{" || token == "["
}

func (operators *Operators) IsCloseBracket(atom Atom) bool {
	token := atom.Name
	return token == ")" || token == "}" || token == "]"
}

func (operators *Operators) MatchBrackets(open Atom, close Atom) bool {
	expectedClose, ok := operators.brackets[open.Name]
	return ok && expectedClose == close.Name
	//switch open.Name {
	//case  "(": return close.Name == ")"
	//case  "[": return close.Name == "]"
	//case  "{": return close.Name == "}"
	//default: return false
	//}
}

// TODO generalize
func (operators *Operators) IsUnaryPrefix(token string) bool {
	return token == "+" || token == "-"
}

func (operators *Operators) AddStandardCOperators() {
	operators.unary["-"] = Atom{"_unary_-"}
	operators.unary["+"] = Atom{"_unary_+"}
	operators.brackets[OPEN_BRACKET] = CLOSE_BRACKET
	operators.brackets[OPEN_SQUARE_BRACKET] = CLOSE_SQUARE_BRACKET
	operators.brackets[OPEN_BRACE] = CLOSE_BRACE
	operators.Add("_unary_+", 110)
	operators.Add("_unary_-", 110)
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
