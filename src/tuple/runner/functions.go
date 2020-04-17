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
package runner

import "tuple/eval"
import "tuple/parsers"

/////////////////////////////////////////////////////////////////////////////

// Provides a set of 'safe' functions that implemented in 'wsh'
func AddTranslatedSafeFunctions(context eval.EvalContext) {
	inputGrammar := parsers.NewShellGrammar()
	ParseAndEval(context, inputGrammar, `

func count  t { progn (c=0) (for v t { c=c+1 }) c }
func first  t { nth 0 t }
func second t { nth 1 t }
func third  t { nth 2 t }

# TODO  Flatten a nested structure
func flatten a {
   if (ismap a) {
   	   for v a {
               flatten v
      	   }
   } {
     if ((arity a) == 0) {
     	a
     } {
     	   for v a {
               flatten v
      	   }
     }
  }
}

# TODO
func print a {
   if (ismap a) {
        join " " {
     	   for v a {
               print v
      	   }
        }
   } {
     if ((arity a) == 0) {
     	a
     } {
        join " " {
     	   for v a {
               print v
      	   }
        }
     }
  }
}


`)

	//func reduce f t { progn c=1 accumulator=first(t) (for v t { accumulator = f(accumulator v))  accumulator}

}

