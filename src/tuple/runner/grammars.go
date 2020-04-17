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
import "path"
import "os"
import "bufio"
import "fmt"


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

/////////////////////////////////////////////////////////////////////////////

func (grammars * Grammars)  RunFile(locationLogger LocationLogger, fileName string, next Next) (Context, error) {
	suffix := path.Ext(fileName)
	grammar, ok := grammars.FindBySuffix(suffix)
	if ! ok {
		return nil, errors.New("Unsupported file suffix: " + suffix)
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	context := NewParserContext(fileName, reader, locationLogger)
	err = grammar.Parse(&context, next)
	file.Close()
	return &context, err
}

func (grammars * Grammars) RunFiles(locationLogger LocationLogger, args []string, next Next) (int64) {
	errors := int64(0)
	if len(args) == 0 {
		 context, err := RunParserOnStdin(locationLogger, grammars.Default(), next)
		if err != nil {
			location := tuple.NewLocation("<stdin>", 0, 0, 0)
			locationLogger(location, "ERROR", fmt.Sprintf("%s", err))
			errors += 1
		}
		errors += context.Errors()
	} else {
		for _, fileName := range args {
			context, err := grammars.RunFile(locationLogger, fileName, next)
			if err != nil {
				location := tuple.NewLocation(fileName, 0, 0, 0)
				locationLogger(location, "ERROR", fmt.Sprintf("%s", err))
				errors += 1
				break
			}
			errors += context.Errors()
		}
	}
	return errors
}

/////////////////////////////////////////////////////////////////////////////

func (grammars * Grammars) AddAllKnownGrammars() {
	grammars.Add(parsers.NewLispWithInfixGrammar())
	grammars.Add(parsers.NewLispGrammar())
	grammars.Add(parsers.NewInfixExpressionGrammar())
	grammars.Add(parsers.NewYamlGrammar())
	grammars.Add(parsers.NewIniGrammar())
	grammars.Add(parsers.NewPropertyGrammar())
	grammars.Add(parsers.NewJSONGrammar())
	grammars.Add(parsers.NewShellGrammar())
}


/////////////////////////////////////////////////////////////////////////////

func (grammars * Grammars) AddSafeGrammarFunctions(table * eval.Runner) {

	table.AddToRoot(tuple.Tag{"grammars"}, grammars)

	//table.Add("help", func (context eval.EvalContext) Value {
	//	return table.AllSymbols()
	//})
	table.Add("ctx", func (context eval.EvalContext) Value {
		return context.GlobalScope().Root()
	})

	// TODO Add to root
	table.Add("grammars", func (context eval.EvalContext) Value {
		return grammars
	})

	table.Add("read", func (context eval.EvalContext, file String) (Value, error) {
		result := tuple.NewTuple()
		_, err := grammars.RunFile(context.GlobalScope().LocationLogger(), string(file), func(in Value) error { result.Append(in); return nil })
		if err != nil {
			return result, err
		}
		if result.Arity() == 1 {
			return result.Get(0), nil
		}
		return result, nil
	})

	//table.Add("ast", func (expression string) tuple.Value { return parsers.ParseString(inputGrammar, expression) })
	//table.Add("expr", func (expression string) tuple.Value { return  runner.ParseAndEval(&table, inputGrammar, expression) })

	table.Add("expr", func (context eval.EvalContext, expression string) (Value, error) {
		grammar := parsers.NewInfixExpressionGrammar()  // Default???
		return ParseAndEval(context, grammar, expression)
	})

	table.Add("ast2", func (context eval.EvalContext, grammarFileSuffix string, expression string) (Value, error) {
		grammar, ok := grammars.FindBySuffix(grammarFileSuffix)
		if ok {
			return parsers.ParseString(context.GlobalScope().LocationLogger(), grammar, expression)
		} else {
			return nil, errors.New(fmt.Sprintf("No such grammar '%s'", grammarFileSuffix))
		}
	})

	table.Add("expr2", func (context eval.EvalContext, grammarFileSuffix string, expression string) (Value, error) {

		grammar, ok := grammars.FindBySuffix(grammarFileSuffix)
		if ok {
			return ParseAndEval(context, grammar, expression)
		} else {
			return nil, errors.New(fmt.Sprintf("No such grammar '%s'", grammarFileSuffix))
		}
	})
}
