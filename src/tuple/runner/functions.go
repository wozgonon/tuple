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

#  TODO replace concat and join with a yield statement
func print_sexp a {
   if (ismap a) {
      concat "(" {
        join " " {
     	   forkv k v a {
               concat k "." (print_sexp v)
      	   }
        }} ")" 
     
   } {
     if ((arity a) == 0) {
     	a
     } {
       concat "(" {
        join " " {
     	   for v a {
               print_sexp v
      	   }
        }} ")" 
     }
  }
}

func print_yaml a indent {
   if (ismap a) {
      concat  {
        join " " {
     	   forkv k v a {
               concat indent k ":\n" (print_yaml v (concat indent "  "))
      	   }
        }}
     
   } {
     if ((arity a) == 0) {
     	concat indent "- " a "\n"
     } {
       concat  {
        join " " {
     	   for v a {
               print_yaml v (concat indent "  ")
      	   }
        }}  
     }
  }
}


`)

	//func reduce f t { progn c=1 accumulator=first(t) (for v t { accumulator = f(accumulator v))  accumulator}

}

