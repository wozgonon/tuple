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
import "fmt"
import "errors"

func NewSafeSymbolTable(global Global) SymbolTable {
	table := NewHarmlessSymbolTable(global)
	AddAllocatingStringFunctions(&table)
	AddAllocatingTupleFunctions(&table)
	AddSetAndDeclareFunctions(&table)
	AddControlStatementFunctions(&table)
	return table
}

/////////////////////////////////////////////////////////////////////////////
// Allocating functions
/////////////////////////////////////////////////////////////////////////////

func toString(context EvalContext, value Value) string {
	switch val := value.(type) {
	case Tag: return val.Name
	case String: return string(val)  // Quote ???
	case Float64: return  tuple.FloatToString(float64(val))
	case Int64: return tuple.Int64ToString(val)
	case Bool: return tuple.BoolToString(bool(val))
	default: 
		context.Log("ERROR", "cannot convert '%s' to string", value)
		return "..." // TODO
	}
}

func EvalToStrings(context EvalContext, values []Value) []string {

	result := make([]string, len(values))
	for k,_:= range values {
		value, _ := Eval(context, values[k])
		result[k] = toString(context, value)
	}
	return result
}

func AddAllocatingStringFunctions(table * SymbolTable) {
	// TODO change this to take an array
	table.Add("join", func (context EvalContext, separator string, tuple Tuple) string { return strings.Join(EvalToStrings(context, tuple.List), separator) })
	table.Add("concat", func (context EvalContext, values... Value) string  { return strings.Join(EvalToStrings(context, values), "") })
}

func AddAllocatingTupleFunctions(table * SymbolTable)  {

	table.Add("keys", func(context EvalContext, value Value) (Value, error) {

		evaluated, err := Eval(context, value)
		if err != nil {
			return nil, err
		}
		result := tuple.NewTuple()  // TODO not efficient use stream
		if mmap, ok := evaluated.(tuple.Map); ok {
			mmap.ForallKeyValue(func(k Tag, _ Value) {
				result.Append(k)
			})
		} else {
			for k := 0; k < evaluated.Arity(); k += 1 {
				result.Append(tuple.IntToTag(k))
			}
		}
		return result, nil
	})

	table.Add("values", func(context EvalContext, value Value) (Value, error) {
		evaluated, err := Eval(context, value)
		if err != nil {
			return nil, err
		}
		result := tuple.NewTuple()  // TODO not efficient use stream
		evaluated.ForallValues(func(value Value) error {
			result.Append(value)
			return nil
		})
		return result, nil
	})

	table.Add("list", func(_ EvalContext, values... Value) Value { return tuple.NewTuple(values...) })
	// TODO table.Add("quote", func(value Value) Value { return NewTuple("quote", value) })
}

/////////////////////////////////////////////////////////////////////////////
// Declare variables and functions
/////////////////////////////////////////////////////////////////////////////

// This assign might set the value of a global variable if one exists
func Assign (context EvalContext, tag Tag, value Value) (Value, error) {
	evaluated, err := Eval(context, value)
	if err != nil {
		return nil, err
	}
	if table, _ := context.Find(context, tag, []Value{}); table != nil {
		table.Add(tag.Name, func () Value { return evaluated })
	} else {
		context.Add(tag.Name, func () Value { return evaluated })
	}
	return evaluated, nil
}

// This assign will only set a loca variable in the top most context.
func AssignLocal (context EvalContext, tag Tag, value Value) (Value, error) {
	evaluated, err := Eval(context, value)
	if err != nil {
		return nil, err
	}
	context.Add(tag.Name, func () Value { return evaluated })
	return evaluated, nil
}

func AddSetAndDeclareFunctions(table * SymbolTable) {

	table.Add("get", func(context EvalContext, tag Tag) (Value, error) {
		result, err := context.Call(tag, []Value{})
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	table.Add("set", Assign)  // TODO Assign or AssignLocal
	table.Add("func", func(context EvalContext, values... Value) Value {
		lv := len(values)
		if lv <= 1 {
			context.Log("ERROR", "No name or arguments provided to 'func'")
			return tuple.EMPTY
		} else {
			tag := values[0]
			args := values[1:lv-1]
			code := values[lv-1]
			for k,v := range values[:lv-2] {
				if _, ok := v.(Tag); ! ok {
					context.Log("ERROR", "Expected identifier not '%s' for arg %d", v, k)
					return tuple.EMPTY
				}
			}
			functionName := tag.(Tag).Name
			context.Add(functionName, func (context1 EvalContext, values... Value) (Value, error) {
				if len(values) != len(args) {
					message := fmt.Sprintf( "For '%s' Expected %d arguments not %d", functionName, len(args), len(values))
					return tuple.EMPTY, errors.New(message)
				} else {
					context.Log("TRACE", "** FUNC %s argValue: %s", functionName, values)
					newScope := NewSymbolTable(context1)
					for k,v := range values {
						evaluated, err := Eval(context1, v)
						if err != nil {
							return nil, err
						}
						name := args[k].(Tag).Name
						newScope.Add(name, func () Value {
							return evaluated
						})
					}
					evaluated, err := Eval(&newScope, code)
					if err != nil {
						return nil, err
					}
					return evaluated, nil
				}
			})
			return tag
		}
	})
}

/////////////////////////////////////////////////////////////////////////////
//  These are nearly harmless but do allocate
/////////////////////////////////////////////////////////////////////////////

func AddControlStatementFunctions(table * SymbolTable) {

	// Perhaps this could be moved to harmless.
	table.Add("if", func(context EvalContext, condition bool, trueCode Value, falseCode Value) (Value, error) {
		var code Value
		if condition {
			code = trueCode
		} else {
			code = falseCode
		}
		return Eval(context, code)
	})
	table.Add("for", func(context EvalContext, tag Tag, list Value, code Value) Value {
		var iterator Value = nil
		newScope := NewSymbolTable(context)
		newScope.Add(tag.Name, func () Value {
			return iterator
		})
		// TODO Ideally for efficiency allow the method to return a callback iterator rather than collect values into a tuple

		result := tuple.NewFiniteStream(list, func (v Value, next func(v Value) error) error {
			evaluated, err := Eval(context, v)
			iterator = evaluated
			if err != nil {
				return err
			}
			value, err := Eval(&newScope, code)
			if err != nil {
				return err
			}
			return next(value)
		})
		return result
	})
	table.Add("while", func(context EvalContext, condition Value, code Value) (Value, error) {
		var result Value = tuple.EMPTY
		for {
			condition, err := Eval(context, condition)
			if err != nil {
				return nil, err
			}
			conditionResult, ok := condition.(Bool)
			if ! ok || ! bool(conditionResult) {
				return result, nil
			}
			result, err = Eval(context, code)
			if err != nil {
				return nil, err
			}
		}
	})

	// https://www.gnu.org/software/emacs/manual/html_node/eintr/progn.html
	table.Add("progn", func(context EvalContext, values... Value) (Value, error) {
		var result Value = tuple.EMPTY
		for _, v := range values {
			evaluated, err := Eval(context, v)
			if err != nil {
				return nil, err
			}
			result = evaluated
		}
		return result, nil

	})


}
