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
//  An operator grammar
/////////////////////////////////////////////////////////////////////////////

// https://en.wikipedia.org/wiki/Shunting-yard_algorithm
// https://en.wikipedia.org/wiki/Operator-precedence_grammar
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
	(*stack).operatorStack = append((*stack).operatorStack, token)
}

// Remove top from operator stack
func (stack * OperatorGrammar) popOperator() {
	lo := len(stack.operatorStack) // This could be passed in for efficiency
	(*stack).operatorStack = (*stack).operatorStack[:lo-1]
}

func (stack * OperatorGrammar) PushValue(value interface{}) {
	if ! (*stack).wasOperator {
		stack.context.Error("Unexpected value: %s\n", value)
		// TODO handle this situation, flush current contents or add a comma operator
		return
	}
	stack.context.Verbose("PushValue: %s\n", value)
	(*stack).Values.Append(value)
	(*stack).wasOperator = false
}

func (stack * OperatorGrammar) OpenBracket(token Atom) {

	// TODO NewTuple()

	//lv := stack.Values.Length()
	//if ! (*stack).wasOperator && lv > 0 {
	//	atom, ok := stack.Values.List[lv -1].(Atom)
		
	//} else {
	stack.pushOperator(token)
	//}
	(*stack).wasOperator = true
}

func (stack * OperatorGrammar) CloseBracket(token Atom) {
	// TODO should this return an error
	lo := len(stack.operatorStack)
	if lo == 0 {
		stack.context.UnexpectedCloseBracketError (token.Name)
		return
	}

	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		stack.popOperator()

		lv := stack.Values.Length()
		stack.context.Verbose("index=%d Values: %d OpStack: %d op=%s\n", index, lv, lo, top)
		values := &((*stack).Values.List)
		if stack.operators.IsOpenBracket(top) {
			if ! stack.operators.MatchBrackets(top, token) {
				stack.context.Error("Expected close bracket '%s' but found '%s'", top.Name, token.Name)
				return
			}
			val1 := (*values) [lv - 1]
			(*stack).Values.List = append((*values)[:lv-1], val1)
			return
		} else {
			// Replace top of value stack with an expression
			//eval(token, 2)
			val1 := (*values) [lv - 2]
			val2 := (*values) [lv - 1]
			(*stack).Values.List = append((*values)[:lv-2], NewTuple(top, val1, val2))
		}
	}
}

/*func (stack * OperatorGrammar) eval(token string, arity int) {
	values := &((*stack).Values.List)
	lv := stack.Values.Length()
	args := (*values) [lv-arity-1:]
	(*stack).Values.List = append((*values)[:lv-arity-1], NewTuple(token, args...))
}*/

// Signal end of input
func (stack * OperatorGrammar) EOF(next Next) {
	lo := len(stack.operatorStack)
	stack.context.Verbose("OpStack Len=%d\n", lo)
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		stack.popOperator()

		lv := stack.Values.Length()
		stack.context.Verbose("index=%d Values: %d OpStack: %d op=%s\n", index, lv, lo, top)
		values := &((*stack).Values.List)
		if stack.operators.IsOpenBracket(top) {
			val1 := (*values) [lv - 1]
			(*stack).Values.List = append((*values)[:lv-1], val1)
		} else {
			// Replace top of value stack with an expression
			//eval(token, 2)
			val1 := (*values) [lv - 2]
			val2 := (*values) [lv - 1]
			(*stack).Values.List = append((*values)[:lv-2], NewTuple(top, val1, val2))
		}
	}
	//assert len(stack.values) == 1
	// TODO this is a hack to handle space separated expressions: 1+2 3*4 5
	for _, value := range (*stack).Values.List {
		next (value)
	}
	
}

func (stack * OperatorGrammar) PushOperator(atom Atom) {
	values := &((*stack).Values.List)
	lv := stack.Values.Length()
	if lv == 0 || (*stack).wasOperator {
		if stack.operators.IsUnaryPrefix(atom.Name) {
			// TODO treat plus as a no-op
			//eval(atom, 1)
			val1 := (*values) [lv - 1]
			(*stack).Values.List = append((*values)[:lv-2], NewTuple(atom, val1))
		}
	} else {
		atomPrecedence := stack.operators.Precedence(atom)
		lo := len(stack.operatorStack)
		for index := lo-1 ; index >= 0; index -= 1 {
			top := stack.operatorStack[index]
			if stack.operators.IsOpenBracket(top) {
				val1 := (*values) [lv - 1]
				(*stack).Values.List = append((*values)[:lv-1], val1)

			} else if stack.operators.Precedence(top) > atomPrecedence {
				stack.popOperator()
				// Replace top of value stack with an expression
				//eval(atom, 2)
				lv := stack.Values.Length()
				stack.context.Verbose(" -- index=%d Values: %d OpStack: %d op=%s\n", index, lv, lo, top)
				val1 := (*values) [lv - 2]
				val2 := (*values) [lv - 1]
				(*stack).Values.List = append((*values)[:lv-2], NewTuple(top, val1, val2))
			} else {
				break
			}
		}
		stack.pushOperator(atom)
	}
	// TODO postfix
	(*stack).wasOperator = true
}


// A table of operators
type Operators struct {
	precedence map[string]int

}

func NewOperators() Operators {
	return Operators{make(map[string]int, 0)}
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
	return token == "(" || token == "{" || token == "["
}

func (operators *Operators) IsCloseBracket(atom Atom) bool {
	token := atom.Name
	return token == ")" || token == "}" || token == "]"
}

func (operators *Operators) MatchBrackets(open Atom, close Atom) bool {
	switch open.Name {
	case  "(": return close.Name == ")"
	case  "[": return close.Name == "]"
	case  "{": return close.Name == "}"
	default: return false
	}
}

// TODO generalize
func (operators *Operators) IsUnaryPrefix(token string) bool {
	return token == "+" || token == "-"
}

func (operators *Operators) AddStandardCOperators() {
	operators.Add("^", 10)
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
}
