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

func AddSafeFunctions(table LocalScope) {
	AddHarmlessFunctions(table)
	AddAllocatingStringFunctions(table)
	AddAllocatingTupleFunctions(table)
	AddSetAndDeclareFunctions(table)
	AddControlStatementFunctions(table)
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
		if value.Arity() == 0 {
			return "()"
		}
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

func toStrings(context EvalContext, value Value) []string {

	result := make([]string, value.Arity())
	kk := 0
	value.ForallValues(func (val Value) error {
		str := toString(context, val)
		result[kk] = str
		kk += 1
		return nil
	})
	return result
}

func AddAllocatingStringFunctions(table LocalScope) {
	// TODO change this to take an array
	table.Add("join", func (context EvalContext, separator string, value Value) string { return strings.Join(toStrings(context, value), separator) })
	table.Add("concat", func (context EvalContext, values... Value) string  { return strings.Join(EvalToStrings(context, values), "") })
}

func AddAllocatingTupleFunctions(table LocalScope)  {
	
	table.Add("yield", func(context EvalContext, value Value)  {
		// TODO Fix this to be like a real yield function ala Sather C# Python
		fmt.Printf("%s", toString(context, value))
	})

	table.Add("keys", func(context EvalContext, value Value) (Value, error) {

		result := tuple.NewTuple()  // TODO not efficient use stream
		if mmap, ok := value.(tuple.Map); ok {
			mmap.ForallKeyValue(func(k Tag, _ Value) {
				result.Append(k)
			})
		} else {
			for k := 0; k < value.Arity(); k += 1 {
				result.Append(tuple.IntToTag(k))
			}
		}
		return result, nil
	})

	table.Add("values", func(context EvalContext, evaluated Value) (Value, error) {
		result := tuple.NewTuple()  // TODO not efficient use stream
		evaluated.ForallValues(func(value Value) error {
			result.Append(value)
			return nil
		})
		return result, nil
	})

	table.Add("list", func(context EvalContext, values... Value) (Value, error) {
		// TODO evaluate
		array := make([]Value, len(values))
		for k,v := range values {
			evaluated, err := Eval(context, v)
			array[k] = evaluated
			if err != nil {
				return nil, err
			}
		}
		return tuple.NewTuple(array...), nil
	})
	// TODO table.Add("quote", func(value Value) Value { return NewTuple("quote", value) })
}

/////////////////////////////////////////////////////////////////////////////
// Declare variables and functions
/////////////////////////////////////////////////////////////////////////////

// This assign might set the value of a global variable if one exists
func Assign (context EvalContext, tag Tag, evaluated Value) (Value, error) {
	if table, _ := context.Find(context, tag, []Value{}); table != nil {
		table.Add(tag.Name, func () Value { return evaluated })
	} else {
		context.Add(tag.Name, func () Value { return evaluated })
	}
	return evaluated, nil
}

// This assign will only set a loca variable in the top most context.
func AssignLocal (context EvalContext, tag Tag, evaluated Value) (Value, error) {
	context.Add(tag.Name, func () Value { return evaluated })
	return evaluated, nil
}

func AddSetAndDeclareFunctions(table LocalScope) {

	table.Add("get", func(context EvalContext, tag Tag) (Value, error) {
		result, err := Call(context, tag, []Value{})
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
					newScope := context1.NewLocalScope()
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
					evaluated, err := Eval(newScope, code)
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

func AddControlStatementFunctions(table LocalScope) {

	// Perhaps this could be moved to harmless.
	table.Add("if", func(context EvalContext, condition bool, trueCode Quoted, falseCode Quoted) (Value, error) {
		var code Value
		if condition {
			code = trueCode.Value()
		} else {
			code = falseCode.Value()
		}
		return Eval(context, code)
	})
	table.Add("for", func(context EvalContext, tag Tag, list Value, code Quoted) Value {
		var iterator Value = nil
		newScope := context.NewLocalScope()
		newScope.Add(tag.Name, func () Value {
			return iterator
		})
		result := tuple.NewFiniteStream(list, func (v Value, next func(v Value) error) error {
			iterator = v
			value, err := Eval(newScope, code.Value())
			if err != nil {
				return err
			}
			return next(value)
		})
		return result
	})
	table.Add("forkv", func(context EvalContext, key Tag, val Tag, mapp tuple.Map, code Quoted) Value {
		var keyIterator Value = nil
		var valIterator Value = nil
		newScope := context.NewLocalScope()
		newScope.Add(key.Name, func () Value {
			return keyIterator
		})
		newScope.Add(val.Name, func () Value {
			return valIterator
		})
		// TODO Ideally for efficiency allow the method to return a callback iterator rather than collect values into a tuple
		result := tuple.NewTuple()
		mapp.ForallKeyValue(func(kk Tag, vv Value) {
			keyIterator = kk
			valIterator = vv
			value, err := Eval(newScope, code.Value())
			if err != nil {
				Error(context, "%s", err)
				return
			}
			result.Append(value)
		})
		return result
	})
	table.Add("while", func(context EvalContext, condition Quoted, code Quoted) (Value, error) {
		var result Value = tuple.EMPTY
		for {
			condition, err := Eval(context, condition.Value())
			if err != nil {
				return nil, err
			}
			conditionResult, ok := condition.(Bool)
			if ! ok || ! bool(conditionResult) {
				return result, nil
			}
			result, err = Eval(context, code.Value())
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
