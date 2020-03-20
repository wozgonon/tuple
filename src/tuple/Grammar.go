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
//import "fmt"

// The Grammar interface represents a particular language Grammar or Grammar or File Format.
//
// The print and parse method ought to be inverse functions of each other
// so the output of parse can be passed to print which in principle should be parsable by the parse function.
//
type Grammar interface {
	// A friendly name for the syntax
	Name() string

	// A standard suffix for source files.
	FileSuffix() string
	
	// Parses an input stream of characters into an internal representation (AST)
	// The output ought to be printable by the 'print' method.
	Parse(context * ParserContext) // , next func(tuple Tuple)) (interface{}, error)
	
	// Pretty prints the objects in the given syntax.
	// The output ought to be parsable by the 'parse' method.
	Print(token interface{}, next func(value string))
}

// A set of Grammars
type Grammars struct {
	all map[string]Grammar
}

// Returns a new empty set of syntaxes
func NewGrammars() Grammars{
	return Grammars{make(map[string]Grammar)}
}

func (syntaxes * Grammars) Add(syntax Grammar) {
	suffix := syntax.FileSuffix()
	syntaxes.all[suffix] = syntax
}

func (syntaxes * Grammars) FindBySuffix(suffix string) (*Grammar, bool) {
	if ! strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	syntax, ok := syntaxes.all[suffix]
	return &syntax, ok
}

func (syntaxes * Grammars) FindBySuffixOrPanic(suffix string) *Grammar {
	syntax, ok := syntaxes.FindBySuffix(suffix)
	if ! ok {
		panic("Unsupported file suffix: '" + suffix + "'")
	}
	return syntax
}

/////////////////////////////////////////////////////////////////////////////
//  An operator grammar
/////////////////////////////////////////////////////////////////////////////

// https://en.wikipedia.org/wiki/Shunting-yard_algorithm
type OperatorStack struct {
	context * ParserContext
	operators * Operators
	operatorStack []Atom
	Values Tuple
	wasOperator bool
}

func NewOperatorStack(context * ParserContext, operators * Operators) OperatorStack {
	return OperatorStack{context, operators, make([]Atom, 0), NewTuple(),true}
}

func (stack * OperatorStack) pushOperator(token Atom) {
	(*stack).operatorStack = append((*stack).operatorStack, token)
}

// Remove top from operator stack
func (stack * OperatorStack) popOperator() {
	lo := len(stack.operatorStack) // This could be passed in for efficiency
	(*stack).operatorStack = (*stack).operatorStack[:lo-1]
}

func (stack * OperatorStack) PushValue(value interface{}) {
	stack.context.Verbose("PushValue: %s\n", value)
	(*stack).Values.Append(value)
	(*stack).wasOperator = false
}

func (stack * OperatorStack) OpenBracket(token Atom) {

	// TODO NewTuple()

	//lv := stack.Values.Length()
	//if ! (*stack).wasOperator && lv > 0 {
	//	atom, ok := stack.Values.List[lv -1].(Atom)
		
	//} else {
	stack.pushOperator(token)
	//}
	(*stack).wasOperator = true
}

func (stack * OperatorStack) CloseBracket(token Atom) {
	// TODO should this return an error
	lo := len(stack.operatorStack)
	if lo == 0 {
		stack.context.UnexpectedCloseBracketError (token.Name)
		return
	}
	top := stack.operatorStack[lo-1]
	if top != token {
		stack.context.Error("Expected close bracket '%s' but found '%s'", top.Name, token.Name)
		return
	}
	stack.popOperator()
}

/*func (stack * OperatorStack) eval(token string, arity int) {
	values := &((*stack).Values.List)
	lv := stack.Values.Length()
	args := (*values) [lv-arity-1:]
	(*stack).Values.List = append((*values)[:lv-arity-1], NewTuple(token, args...))
}*/

// Signal end of input
func (stack * OperatorStack) EOF() interface{} {
	lo := len(stack.operatorStack)
	stack.context.Verbose("OpStack Len=%d\n", lo)
	for index := lo-1 ; index >= 0; index -= 1 {
		top := stack.operatorStack[index]
		stack.context.Verbose("  OpStack index=%d op=%s\n", index, top)
		stack.popOperator()
		// Replace top of value stack with an expression
		//eval(token, 2)
		lv := stack.Values.Length()
		values := &((*stack).Values.List)
		val1 := (*values) [lv - 2]
		val2 := (*values) [lv - 1]
		(*stack).Values.List = append((*values)[:lv-2], NewTuple(top, val1, val2))
	}
	//assert len(stack.values) == 1
	return (*stack).Values.List[0]
}

func (stack * OperatorStack) PushOperator(atom Atom) {
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
			if stack.operators.IsOpenBracket(top) || stack.operators.Precedence(top) > atomPrecedence {
				stack.popOperator()
				// Replace top of value stack with an expression
				//eval(atom, 2)
				lv := stack.Values.Length()
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
