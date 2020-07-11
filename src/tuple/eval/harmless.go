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

import "math"
import "strings"
import "reflect"
import "errors"
import "tuple"

/////////////////////////////////////////////////////////////////////////////

// These functions are harmless
// One can execute them from a script without any worry they will access something they ought not to or use up resources.
func AddHarmlessFunctions(table LocalScope) {
	AddBooleanAndArithmeticFunctions(table)
	AddHarmlessStringFunctions(table)
	AddHarmlessArrayFunctions(table)
	AddHarmlessFinanceFunctions(table)
}
/////////////////////////////////////////////////////////////////////////////

// These functions are harmless in the sense that they just do basic functions and do not provide any access to resources
// such as operating system or even allocating memory.
func AddBooleanAndArithmeticFunctions(table LocalScope) {

	table.Add("exp", math.Exp)
	table.Add("sqrt", math.Sqrt)
	table.Add("log", math.Log)
	table.Add("sin", math.Sin)
	table.Add("cos", math.Cos)
	table.Add("tan", math.Tan)
	table.Add("acos", math.Acos)
	table.Add("asin", math.Asin)
	table.Add("atan", math.Atan)
	table.Add("atan2", math.Atan2)
	table.Add("round", math.Round)
	table.Add("round2", func(value float64) float64 { return math.Round(value*100)/100 }) // Not needed
	table.Add("**", math.Pow)
	table.Add("%", func (aa float64) float64 { return aa/100.0 })
	table.Add("++", func (aa float64) float64 { return aa+1 })
	table.Add("+", func (aa float64) float64 { return aa })
	table.Add("-", func (aa float64) float64 { return -aa })
	table.Add("+", func (aa float64, bb float64) float64 { return aa+bb })
	table.Add("-", func (aa float64, bb float64) float64 { return aa-bb })
	table.Add("*", func (aa float64, bb float64) float64 { return aa*bb })
	table.Add("/", func (aa float64, bb float64) float64 { return aa/bb })
/*	table.Add("++", func (aa int64) int64 { return aa+1 })
	table.Add("+", func (aa int64) int64 { return aa })
	table.Add("-", func (aa int64) int64 { return -aa })
	table.Add("+", func (aa int64, bb int64) int64 { return aa+bb })
	table.Add("-", func (aa int64, bb int64) int64 { return aa-bb })
	table.Add("*", func (aa int64, bb int64) int64 { return aa*bb })
	table.Add("/", func (aa int64, bb int64) (int64, error) {
		if bb == 0 {
			return 0, errors.New("Divide by zero")
		}
		return aa/bb, nil
	})*/
	table.Add("%", func (context EvalContext, aa int64, bb int64) (int64, error) {
		if bb == 0 {
			return 0, errors.New("Divide by zero")
			//context.Log("ERROR", "divide by zero in: %d %% %d", aa, bb)
		}
		return aa%bb, nil
	})
	table.Add("streq", func (aa string, bb string) bool { return aa==bb })  // Should take Value as argument
	table.Add("==", func (aa float64, bb float64) bool { return aa==bb })  // Should take Value as argument
	table.Add("!=", func (aa float64, bb float64) bool { return aa!=bb })
	table.Add(">=", func (aa float64, bb float64) bool { return aa>=bb })
	table.Add("<=", func (aa float64, bb float64) bool { return aa<=bb })
	table.Add(">", func (aa float64, bb float64) bool { return aa>bb })
	table.Add("<", func (aa float64, bb float64) bool { return aa<bb })
	table.Add("&&", func (aa bool, bb bool) bool { return aa&&bb })
	table.Add("||", func (aa bool, bb bool) bool { return aa||bb })
	table.Add("!", func (aa bool) bool { return ! aa })
	table.Add("PI", func () float64 { return math.Pi })
	table.Add("PHI", func () float64 { return math.Phi })
	table.Add("E", func () float64 { return math.E })
	table.Add("true", func () bool { return true })
	table.Add("false", func () bool { return false })

}

// String functions that do not allocate any memory
func AddHarmlessStringFunctions(table LocalScope) {
	table.Add("len", func(value string) int64 { return int64(len(value)) })  // TODO arity
	table.Add("lower", strings.ToLower)
	table.Add("upper", strings.ToUpper)
}

// String functions that do not allocate any memory
func AddHarmlessFinanceFunctions(table LocalScope) {

    v := func(i float64) float64 { return 1/(1+i) }
	d := func(i float64) float64 {
    	 return i/(1+i)
    }
    pv := func(i float64, rent_payment tuple.Array) Value {
        values := make([]Value, rent_payment.Arity())
        vv := v(i)
        _ = tuple.ForallKeyValuesInArray(rent_payment, func(index int, value Value) error {
                ival, iok := value.(Int64)
                fval, fok := value.(Float64) // TODO use int64 everywhere
                var val float64
                if iok {
                    val = float64(int64(ival))
                } else if fok  {
                    val = float64(fval)
                } else {
                    // TODO error
                    val = 1.0
                }
                values[index] = Float64(val * vv)
                vv *= vv
                return nil
            })
            return tuple.NewTuple(values...)
        }
	table.Add("pv", pv)
	table.Add("v", v)
	table.Add("d", d)
    // Annuity certain

    // Annuity immediate incur interest immediately and pay at end of period a with a bar
    table.Add("annuityImmediate", func(i float64, n float64) float64 {
          	 return (1-math.Pow(v(i), n))/i
          })
    // Payments at start of period  a with two dots
	table.Add("annuityDue", func(i float64, n float64) float64 {
    	 return (1-math.Pow(v(i), n))/d(i)
    })
	table.Add("perpetuity", func(rent_payment float64, i float64) float64 {
    	 return rent_payment/i
    })
	table.Add("amortization", func(rent_payment float64, debt float64, i float64, n float64) float64 {
        // The amount owed after 'n' payments
        ri := rent_payment/i
    	 return ri - math.Pow(1+i, n) * (ri - debt)
    })

}

// Array functions that do not allocate any memory
func AddHarmlessArrayFunctions(table LocalScope)  {

	table.Add("quote", func(context EvalContext, quoted Quoted) Value {
		// TODO get correct quote
		return tuple.NewTuple(Tag{"quote"}, quoted.Value())
	})
	table.Add("arity", func(context EvalContext, value Value) int64 {
		return int64(value.Arity())
	})

	table.Add("nth", func(context EvalContext, index64 int64, value tuple.Array) Value {
		index := int(index64) // TODO use int64 everywhere
		return value.Get(index)
	})

	table.Add("istuple", func (context EvalContext, value Value) bool {
		_, ok := value.(Tuple)
		return ok
	})
	table.Add("ismap", func (context EvalContext, value Value) bool {
		_, ok := value.(tuple.Map)
		return ok
	})
	table.Add("typeof", func (context EvalContext, value Value) string {
		return reflect.TypeOf(value).Name()
	})
	table.Add("eq", func (context EvalContext, aa Value, bb Value) bool {
		if aa.Arity() != bb.Arity() {
			return false
		}
		//TODO eq ("a1" b2")   1
		//TODO eq ("a1" b2")   ("a1" b2") 
		return reflect.DeepEqual(aa, bb)
	})
}

/////////////////////////////////////////////////////////////////////////////

type ErrorIfFunctionNotFound struct {}

func NewErrorIfFunctionNotFound() Finder {
	finder := ErrorIfFunctionNotFound{}
	return &finder
}

func (function * ErrorIfFunctionNotFound) Find(context EvalContext, name Tag, args [] Value) (LocalScope, reflect.Value) {
	return nil, reflect.ValueOf(func(args... Value) bool {
		context.Log("ERROR", "function not found: '%s' %s\n", name.Name, args)
		return false
	})
}

