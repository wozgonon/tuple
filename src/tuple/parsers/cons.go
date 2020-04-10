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
import "fmt"
import "tuple"
import "errors"

func consFilterFinal(expression Value) (Value,error) {// TODO add context for errors and make a Tuple

	//fmt.Printf("Tuple: %s\n", expression)
	if isCons(expression) {
		mapp := tuple.NewTagValueMap()
		err := addConsToMap(mapp, expression)
		if err != nil {
			return nil, err
		}
		//context.Log("VERBOSE", "Map len '%d'", mapp.Arity())
		return mapp, nil
	}
	return expression, nil
}

func consFilter(expression Tuple) (Value,error) {// TODO add context for errors and make a Tuple

	arity := expression.Arity()
	if arity == 0 {
		return expression, nil
	}
	if isCons(expression.Get(0)) {
		mapp := tuple.NewTagValueMap()
		err := tuple.ForallInArray(expression, func (value Value) error {
			return addConsToMap(mapp, value)
		})
		//context.Log("VERBOSE", "Len of map '%d'", mapp.Arity())
		if err != nil {
			return nil, err
		}
		return mapp, nil
	}
	for k := 1; k < arity; k+=1 {
		element := expression.Get(k)
		if isCons(element) {
			mapp := tuple.NewTagValueMap()
			addConsToMap(mapp, element)
			expression.Set(k, mapp)
		}

	}
	return expression, nil
}

func addConsToMap(mapp tuple.TagValueMap, value Value) error {
	if isCons(value) {
		array := value.(tuple.Array)
		head := array.Get(1)
		var key Tag
		switch head.(type) {
		case tuple.Tag: key = head.(Tag)
		case tuple.Int64: key = Tag{tuple.Int64ToString(head.(Int64))}
		case tuple.String: key = Tag{string(head.(String))}
		default: 
			message := fmt.Sprintf("Got '%s' expected tag,int or string, got", head)
			return errors.New(message)
		}
		//context.Log("VERBOSE", "Add '%s:%s'", key, array.Get(2))
		mapp.Add(key, array.Get(2))
		return nil
	} else {
		message := fmt.Sprintf("Got '%s' expect CONS cell / Key-Value pair", value)
		return errors.New(message)
	}
}

func isCons(value Value) bool {
	if value.Arity() > 0 {
		array, ok := value.(tuple.Array)
		if ok {
			head := array.Get(0)
			tag, ok := head.(Tag)
			if ok && tag == tuple.CONS_ATOM {
				return true
			}
		}
	}
	return false
}

func isConsInTuple(value Value) bool {
	if array, ok := value.(tuple.Array); ok {
		return value.Arity() > 0 && isCons(array.Get(0))
	}
	return false
}

