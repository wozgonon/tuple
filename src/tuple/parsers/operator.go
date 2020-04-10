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
import "fmt"

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
	operatorStack []Tag
	Values Tuple
	wasOperator bool
}

func NewOperatorGrammar(context Context, operators * Operators) OperatorGrammar {
	return OperatorGrammar{context, operators, make([]Tag, 0), NewTuple(),true}
}

func (stack * OperatorGrammar) pushOperator(token Tag) {
	Verbose(stack.context,"PUSH OPERATOR\t '%s'", token.Name)
	stack.operatorStack = append(stack.operatorStack, token)
}

// Remove top from operator stack
func (stack * OperatorGrammar) popOperator() {
	lo := len(stack.operatorStack) // This could be passed in for efficiency
	Verbose(stack.context,"POP OPERATOR\t '%s'", stack.operatorStack[lo-1].Name)
	stack.operatorStack = stack.operatorStack[:lo-1]
}

// Replace top of value stack with an expression
func (stack * OperatorGrammar) reduceOperatorExpression(top Tag) int {

	values := &(stack.Values.List)
	lv := stack.Values.Arity()
	name := stack.operators.Map(top)
	stack.popOperator()
	popped := 0
	index := 0
	tuple := NewTuple()
	if isPrefix(top) {
		val1 := (*values) [lv - 1]
		popped = 1
		tuple = NewTuple(name, val1)
		Verbose(stack.context," REDUCE:\t%s\t'%s'\n", name.Name, val1)
	} else {
		popped = 2
		if stack.operators.IsReduceAllRepeats(top) {
			// TODO Could in principle generalize to make any binary operator n-ary
			// TODO this would work like python 3>2>1
			Verbose(stack.context,"** REDUCE\t'%s'\n", top.Name)
			for {
				ll := len(stack.operatorStack)
				if ll == 0 {
					break
				}
				nextTop := stack.operatorStack[ll - 1]
				if nextTop != top {
					break
				}
				stack.popOperator()
				popped += 1
			}
			// TODO the following is not efficient and should be replaced with a slice: tuple := NewTuple(args...)
			args := (*values) [lv-popped:]
			for _,v := range args {
				tuple.Append(v)
			}
			Verbose(stack.context," REDUCE:\t'SPACE'\t'%s'\t'%s'\t...%d...   \n", tuple.List[0], tuple.List[1], tuple.Arity()) //, tuple.List==(*values))
		} else {
			// TODO this can fail, for instance if no brackets are registered
			val1 := (*values) [lv - 2]
			val2 := (*values) [lv - 1]
			tuple = NewTuple(name, val1, val2)
			Verbose(stack.context," REDUCE:\t'%s'\t'%s'\t'%s'\n", name, val1, val2)
		}
		index = popped - 2
	}

	value, err := consFilter(tuple)  // TODO generalize this
	if err != nil {
		panic(fmt.Sprintf("TODO handle err: %s", err))
		// TODO
	}
	if value == nil {
		panic("Unexpected nil")
	}
	stack.Values.List = append((*values)[:lv-popped], value)
	return index
}


func (stack * OperatorGrammar) PushValueWithoutInsertingMissingSepator(value Value) {
	Verbose(stack.context,"PUSH VALUE\t'%s'\n", value)
	if value == nil {
		panic("Unexpected nil")
	}
	stack.Values.Append(value)
	stack.wasOperator = false
}

// If a list of space separated values are entered such as (1 2 3 4) without any separating operator
// such as (1,2,3,4) then a separator is inserted automatically.
// A strict grammar might insist of having comma separators and less strict would be happy with space
// separated values.
func (stack * OperatorGrammar) PushValue(value Value) {
	if ! stack.wasOperator {
		stack.PushOperator(SPACE_ATOM)  // TODO should this just add a comma rather than a space
	}
	stack.PushValueWithoutInsertingMissingSepator(value)
}

func (stack * OperatorGrammar) OpenBracket(token Tag) {

	Verbose(stack.context,"OPEN '%s'", token.Name)
	if ! stack.wasOperator {
		stack.PushOperator(SPACE_ATOM)
	}
	stack.pushOperator(token)
	stack.wasOperator = true
}

func (stack * OperatorGrammar) CloseBracket(token Tag) {
	Verbose(stack.context,"CLOSE '%s'", token.Name)

	stack.postfix()

	lo := len(stack.operatorStack)
	if lo == 0 || stack.wasOperator {
		if lo > 0 && stack.operators.IsOpenBracket(stack.operatorStack[lo-1]) {  // '()'  Empty list, is this always okay
			stack.wasOperator = false
			stack.popOperator()
			//lv := stack.Values.Arity()
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

func (stack * OperatorGrammar) postfix() {
	lo := len(stack.operatorStack)
	if stack.wasOperator && lo > 0 {
		if stack.operators.IsReduceAllRepeats(stack.operatorStack[lo-1]) {
			stack.popOperator()
			stack.wasOperator = false
		}
	}
}

// Signal end of input
func (stack * OperatorGrammar) EndOfInput(next Next) error {
	Verbose(stack.context,"EOF")

	lo := len(stack.operatorStack)
	empty := stack.Values.Arity() == 0 && lo == 0
	if empty {
		return nil
	}

	stack.postfix()
	
	lo = len(stack.operatorStack)
	if stack.wasOperator {
		UnexpectedEndOfInputErrorBracketError(stack.context)
	} else {
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
			result, err := consFilterFinal(stack.Values.Get(0))
			if err != nil {
				stack.flush()
				return err
			}
			err = next(result)
			if err != nil {
				stack.flush()
				return err
			}
		} else {
			err := next(stack.Values)
			if err != nil {
				stack.flush()
				return err
			}
		}
	}
	stack.flush()
	return nil
}

func (stack * OperatorGrammar) flush() {
	stack.Values = NewTuple()
	stack.operatorStack = make([]Tag, 0)
	stack.wasOperator = true
}

func (stack * OperatorGrammar) PushOperator(operator Tag) {
	_, ok := stack.operators.postfix[operator.Name]
	operatorIsPostfix := ! stack.wasOperator && ok
	prefixOperator, ok := stack.operators.prefix[operator.Name]
	operatorIsPrefix := stack.wasOperator && ok

	Verbose(stack.context,"*PushOperator '%s' isPrefix=%s  (%s)", operator.Name, ok, prefixOperator)
	// TODO postfix

	if operatorIsPrefix {
		stack.pushOperator(prefixOperator.tag)
	} else if operatorIsPostfix {  // TODO move this code to the postfix method, which will allow % as a postfix operator
		values := &(stack.Values.List)
		lv := len(*values)
		val1 := (*values) [lv - 1]
		name := stack.operators.Map(operator)
		tuple := NewTuple(name, val1)
		stack.Values.List = append((*values)[:lv-1], tuple)
		Verbose(stack.context," REDUCE POSTFIX:\t%s\t'%s'\n", name.Name, val1)
		stack.wasOperator = false
		return
	} else {
		tagPrecedence := stack.operators.Precedence(operator)
		lo := len(stack.operatorStack)
		for index := lo-1 ; index >= 0; index -= 1 {
			top := stack.operatorStack[index]
			topIsPrefix := isPrefix(top)
			Verbose(stack.context, "IsPrefix %s %s  precedence=%d", top, topIsPrefix, stack.operators.Precedence(top))
			if operatorIsPrefix {
				break
			} else if topIsPrefix {
				index -= stack.reduceOperatorExpression(top)
			} else if stack.operators.IsOpenBracket(top) {
				break
			} else if top == operator && stack.operators.IsReduceAllRepeats(operator) {
				break
			} else if stack.operators.Precedence(top) >= tagPrecedence {
				Verbose(stack.context,"* PushOperator - Reduce '%s'", top)
				index -= stack.reduceOperatorExpression(top)
			} else {
				break
			}
		}
		if ! operatorIsPrefix && stack.wasOperator {
			Error(stack.context,"Unexpected binary operator '%s'", operator.Name)
			return
		}
		stack.pushOperator(operator)
	} 
	stack.wasOperator = true
}


/////////////////////////////////////////////////////////////////////////////
//  Operators
/////////////////////////////////////////////////////////////////////////////

// A table of operators
type Operators struct {
	Style
	precedence map[string]int
	prefix map[string]Operator
	infix map[string]Operator
	postfix map[string]Operator
	brackets map[string]string
	closeBrackets map[string]string
	evalName map[string]Tag
}

// TODO replace some of the maps in Operators with a class
type Operator struct {
	tag Tag
	precedence int
	reduceAllRepeats bool
	evalName Tag
}

func NewOperators(style Style) Operators {
	tags1 := make(map[string]Operator, 0)
	tags2 := make(map[string]Operator, 0)
	tags3 := make(map[string]Operator, 0)
	tags4 := make(map[string]Tag, 0)
	strings1 := make(map[string]string, 0)
	strings2 := make(map[string]string, 0)
	return Operators{style, make(map[string]int, 0), tags1, tags2, tags3, strings1, strings2, tags4}
}

const PREFIX = "_prefix_"

func isPrefix(top Tag) bool {
	// stack.operators.prefix[tag.Name]
	return strings.HasPrefix(top.Name, PREFIX)
}

func (operators *Operators) Map(top Tag) Tag {

	tag, ok := operators.evalName[top.Name]
	if ok {
		return tag
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

var	COMMA_ATOM = Tag{";"}

func (_ *Operators) IsReduceAllRepeats(top Tag) bool {
	return top == SPACE_ATOM || top == COMMA_ATOM
}

func (operators *Operators) AddInfix(operator string, precedence int) {
	operators.AddInfix3(operator, precedence, operator)
}

func (operators *Operators) AddInfix3(operator string, precedence int, name string) {
	operators.precedence[operator] = precedence
	operators.infix[operator] = Operator{Tag{name}, precedence, false,Tag{operator}}
	operators.evalName[name] = Tag{operator}
}

func (operators *Operators) AddPrefix(operator string, precedence int) {
	name := PREFIX + operator
	operators.prefix[operator] = Operator{Tag{name}, precedence, false,Tag{operator}}
	operators.precedence[operator] = precedence
	operators.precedence[name] = precedence
	operators.evalName[name] = Tag{operator}
}

func (operators *Operators) AddPostfix(operator string, precedence int) {
	name := "_postfix" + operator
	operators.postfix[operator] = Operator{Tag{name}, precedence, false,Tag{operator}}
	operators.precedence[operator] = precedence
	operators.precedence[name] = precedence
	operators.evalName[name] = Tag{operator}
}

func (operators *Operators) AddBracket(open string, close string) {
	operators.brackets[open] = close
	operators.closeBrackets[close] = open
}

func (operators *Operators) Precedence(token Tag) int {
	value, ok := operators.precedence[token.Name]
	if ok {
		return value
	}
	return -1
}

// TODO generalize
func (operators *Operators) IsOpenBracket(tag Tag) bool {
	token := tag.Name
	_, ok := operators.brackets[token]
	return ok
}

func (operators *Operators) MatchBrackets(open Tag, close Tag) bool {
	expectedClose, ok := operators.brackets[open.Name]
	return ok && expectedClose == close.Name
}

/////////////////////////////////////////////////////////////////////////////
// Printer
/////////////////////////////////////////////////////////////////////////////

func (printer Operators) PrintNullaryOperator(depth string, tag Tag, out StringFunction) {
	PrintTuple(&printer, depth, NewTuple(tag), out)
}

func (printer Operators) PrintUnaryOperator(depth string, tag Tag, value Value, out StringFunction) {  // Prefix and Postfix???
	PrintTuple(&printer, depth, NewTuple(tag, value), out)
}

func (printer Operators) PrintBinaryOperator(depth string, tag Tag, value1 Value, value2 Value, out StringFunction) {  // TODO binary to infix

	if _, ok := printer.precedence[tag.Name]; ok {
		out(printer.Style.Open)
		newDepth := depth + "  "
		printer.PrintSuffix(newDepth, out)
		
		PrintExpression1(printer, newDepth, value1, out)

		out(" ")
		out(tag.Name)
		out(" ")

		PrintExpression1(printer, newDepth, value2, out)
		printer.PrintSuffix(newDepth, out)
		
		printer.PrintIndent(depth, out)
		out(printer.Style.Close)
	} else {
		PrintTuple(&printer, depth, NewTuple(tag, value1, value2), out)
	}
}
