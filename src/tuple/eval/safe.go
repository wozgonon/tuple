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
package eval

import "strings"
import "tuple"

func NewSafeSymbolTable(notFound CallHandler) SymbolTable {
	table := NewHarmlessSymbolTable(notFound)
	AddAllocatingStringFunctions(&table)
	AddAllocatingTupleFunctions(&table)
	AddSetAndDeclareFunctions(&table)
	AddControlStatementFunctions(&table)
	return table
}

/////////////////////////////////////////////////////////////////////////////
// Allocating functions
/////////////////////////////////////////////////////////////////////////////

func AddAllocatingStringFunctions(table * SymbolTable) {
	table.Add("join", func (context EvalContext, separator string, tuple Tuple) string { return strings.Join(EvalToStrings(context, tuple.List), separator) })
	table.Add("concat", func (context EvalContext, values... Value) string  { return strings.Join(EvalToStrings(context, values), "") })
}

func AddAllocatingTupleFunctions(table * SymbolTable)  {

	table.Add("list", func(_ EvalContext, values... Value) Value { return tuple.NewTuple(values...) })
	// TODO table.Add("quote", func(value Value) Value { return NewTuple("quote", value) })
}

/////////////////////////////////////////////////////////////////////////////
// Declare variables and functions
/////////////////////////////////////////////////////////////////////////////

func Assign (context EvalContext, atom Atom, value Value) Value {
	evaluated := Eval(context, value)
	if table, _ := context.Find(context, atom, []Value{}); table != nil {
		table.Add(atom.Name, func () Value { return evaluated })
	} else {
		context.Add(atom.Name, func () Value { return evaluated })
	}
	return evaluated
}

func AddSetAndDeclareFunctions(table * SymbolTable) {

	table.Add("get", func(context EvalContext, atom Atom) Value {
		return context.Call(atom, []Value{})
	})
	table.Add("set", Assign)
	table.Add("func", func(context EvalContext, atom Atom, arg Atom, code Value) Value {
		context.Add(atom.Name, func (context1 EvalContext, argValue Value) Value {
			evaluated := Eval(context1, argValue)
			context.Log("TRACE", "** FUNC %s argName=%s argValue: %s evaluated:%s", atom.Name, arg.Name, argValue, evaluated)
			newScope := NewSymbolTable(context1)
			newScope.Add(arg.Name, func () Value {
				return evaluated
			})
			return Eval(&newScope, code)
		})
		return atom
	})
}

/////////////////////////////////////////////////////////////////////////////
//  These are nearly harmless but do allocate
/////////////////////////////////////////////////////////////////////////////

func AddControlStatementFunctions(table * SymbolTable) {

	// Perhaps this could be moved to harmless.
	table.Add("if", func(context EvalContext, condition bool, trueCode Value, falseCode Value) Value {
		if condition {
			return Eval(context, trueCode)
		} else {
			return Eval(context, falseCode)
		}
	})
	table.Add("for", func(context EvalContext, atom Atom, list Tuple, code Value) Tuple {
		var iterator Value = nil
		newScope := NewSymbolTable(context)
		newScope.Add(atom.Name, func () Value {
			return iterator
		})
		result := tuple.NewTuple()
		for _, v := range list.List {
			iterator = Eval(context, v)
			value := Eval(&newScope, code)
			result.Append(value)
		}
		return result
	})

	// https://www.gnu.org/software/emacs/manual/html_node/eintr/progn.html
	table.Add("progn", func(context EvalContext, values... Value) Value {
		var result Value = tuple.EMPTY
		for _, v := range values {
			result = Eval(context, v)
		}
		return result

	})


}
