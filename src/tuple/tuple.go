package tuple

import "fmt"
import "log"
import "strconv"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type Atom struct {
	Name string
}

/////////////////////////////////////////////////////////////////////////////
// Tuple
/////////////////////////////////////////////////////////////////////////////

type Tuple struct {
	list []interface{}
}

func (tuple Tuple) PrettyPrint(depth string) string {
	//space := ""
	result :=  depth + "(\n"
	newDepth := depth + "  "
	for _, token := range tuple.list {
		subTuple, ok := token.(Tuple)
		if ok {
			value := subTuple.PrettyPrint (newDepth);
			result = result + value + "\n"
		} else {
			atom, ok := token.(Atom)
			if ok {
				result = result + newDepth + atom.Name + "\n"
			} else {
				stringValue, ok := token.(string)
				if !ok {
					float, ok := token.(float64)
					if !ok {
						intValue, ok := token.(int64)
						if !ok {
							log.Printf("Type not recognised: %s", token);
						}
						result = result + newDepth + strconv.FormatInt(int64(intValue), 10) + "\n"
					} else {
						result = result + newDepth + fmt.Sprint(float) + "\n"
					}
				} else {
					result = result + newDepth + "\"" + stringValue + "\"" + "\n"  // TODO Escape
				}
			}
		}
	}
	result = result + depth + ")"
	return result
}


func (tuple *Tuple) Append(token interface{}) {
	tuple.list = append(tuple.list, token)
}

func NewTuple() Tuple {
	return Tuple{make([]interface{}, 0)}
}

