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

import "tuple"
import 	"strings"
import "tuple/eval"
import "tuple/parsers"
import "errors"



func AddAllKnownGrammars(grammars * Grammars) {
	grammars.Add(parsers.NewLispWithInfixGrammar())
	grammars.Add(parsers.NewLispGrammar())
	grammars.Add(parsers.NewInfixExpressionGrammar())
	grammars.Add(parsers.NewYamlGrammar())
	grammars.Add(parsers.NewIniGrammar())
	grammars.Add(parsers.NewPropertyGrammar())
	grammars.Add(parsers.NewJSONGrammar())
	grammars.Add(parsers.NewShellGrammar())
}

// A set of Grammars
type Grammars struct {
	All map[string]Grammar
	defaultGrammar Grammar
}

// Returns a new empty set of grammars
func NewGrammars(defaultGrammar Grammar) Grammars{
	grammars := Grammars{make(map[string]Grammar),defaultGrammar}
	grammars.Add(defaultGrammar)
	return grammars
}

func (grammars * Grammars) Arity() int {
	return len(grammars.All)
}

func  (grammars * Grammars) ForallValues(next func(value Value) error) error {
	for key, _ := range grammars.All {
		err := next(Tag{key})
		if err != nil {
			return err
		}
	}
	return nil
}

func (grammars * Grammars) Default() Grammar {
	return grammars.defaultGrammar
}

func (grammars * Grammars) Forall(next func(grammar Grammar)) {
	for _, grammar := range grammars.All {
		next (grammar)
	}
}

func (grammars * Grammars) Add(syntax Grammar) {
	suffix := syntax.FileSuffix()
	grammars.All[suffix] = syntax
}

func (grammars * Grammars) FindBySuffix(suffix string) (Grammar, bool) {
	if ! strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}
	syntax, ok := grammars.All[suffix]
	return syntax, ok
}

func (grammars * Grammars) FindBySuffixOrPanic(suffix string) Grammar {
	syntax, ok := grammars.FindBySuffix(suffix)
	if ! ok {
		panic("Unsupported file suffix: '" + suffix + "'")
	}
	return syntax
}

func (grammars * Grammars) AddSafeGrammarFunctions(table * eval.Runner) {

	table.AddToRoot(tuple.Tag{"grammars"}, grammars)

	//table.Add("help", func (context eval.EvalContext) Value {
	//	return table.AllSymbols()
	//})
	table.Add("ctx", func (context eval.EvalContext) Value {
		return context.Root()
	})
	table.Add("query", func (context eval.EvalContext, path string) Value {
		result := tuple.NewTuple()
		query := NewQuery(path)
		query.Match(context.Root(), func (value Value) error {
			result.Append(value)
			return nil
		})
		return result
	})

	// TODO Add to root
	table.Add("grammars", func (context eval.EvalContext) Value {
		tuple := tuple.NewTuple()
		for _,v := range grammars.All {
			tuple.Append(String(v.FileSuffix()))
		}
		return tuple
	})

	//table.Add("ast", func (expression string) tuple.Value { return parsers.ParseString(inputGrammar, expression) })
	//table.Add("expr", func (expression string) tuple.Value { return  runner.ParseAndEval(&table, inputGrammar, expression) })

	table.Add("expr", func (context eval.EvalContext, expression string) (Value, error) {
		grammar := parsers.NewInfixExpressionGrammar()
		return ParseAndEval(context, grammar, expression)
	})

	table.Add("ast2", func (context eval.EvalContext, grammarFileSuffix string, expression string) (Value, error) {
		grammar, ok := grammars.FindBySuffix(grammarFileSuffix)
		if ok {
			return parsers.ParseString(context.LocationLogger(), grammar, expression)
		} else {
			context.Log("ERROR", "No such grammar '%s'", grammarFileSuffix)
			return nil, errors.New("No such grammar")  //tuple.EMPTY // TODO eror
		}
	})

	table.Add("expr2", func (context eval.EvalContext, grammarFileSuffix string, expression string) (Value, error) {

		grammar, ok := grammars.FindBySuffix(grammarFileSuffix)
		if ok {
			return ParseAndEval(context, grammar, expression)
		} else {
			context.Log("ERROR", "No such grammar '%s'", grammarFileSuffix)
			return tuple.EMPTY, nil // TODO error
		}
	})

	//table.Add("grammars", func (context eval.EvalContext, value Value) Value {
	//	return ParseAndEval(grammar, context, value)
	//})

}
