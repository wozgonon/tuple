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

import "tuple"
import "reflect"
import "fmt"

/////////////////////////////////////////////////////////////////////////////

type Runner struct {
	locationLogger LocationLogger
	symbols SymbolTable
	root tuple.TagValueMap
}

func NewRunner(notFound Finder, logger LocationLogger) Runner {
	symbols :=  NewSymbolTable(notFound)
	runner := Runner{logger, symbols, tuple.NewTagValueMap()}

	runner.AddToRoot(Tag{"funcs"}, &symbols)
	return runner
}

func (scope * Runner) GlobalScope() GlobalScope {
	return scope
}

func (runner * Runner) Root() Value {
	return runner.root
}

func (runner * Runner) AddToRoot(key Tag, value Value) {
	runner.root.Add(key, value)
}

func (runner * Runner) LocationLogger() LocationLogger {
	return runner.locationLogger
}

func (runner * Runner) Log(level string, format string, args ...interface{}) {
	location := tuple.NewLocation("<eval>", 0, 0, 0)
	runner.locationLogger(location, level, fmt.Sprintf(format, args...))
}

func (runner * Runner) Find(context EvalContext, name Tag, args [] Value) (LocalScope, reflect.Value) {
	return runner.symbols.Find(context, name, args)
}

func (runner * Runner) Add(name string, function interface{}) {
	runner.symbols.Add(name, function)
}

func (runner * Runner) NewLocalScope() EvalContext {
	symbols :=  NewSymbolTable(runner)
	scope := RunnerLocalScope{runner,symbols}
	return &scope
}

/////////////////////////////////////////////////////////////////////////////
// TODO need a function scope and a loop scope, assign should never work below function scope

type RunnerLocalScope struct {
	global * Runner
	symbols SymbolTable
}

func (scope * RunnerLocalScope) NewLocalScope() EvalContext {
	symbols :=  NewSymbolTable(scope)
	newScope := RunnerLocalScope{scope.global,symbols}
	return &newScope
}

func (scope * RunnerLocalScope) GlobalScope() GlobalScope {
	return scope.global
}

func (scope * RunnerLocalScope) Root() Value {
	return scope.global.Root()
}

func (scope * RunnerLocalScope) Log(level string, format string, args ...interface{}) {
	location := tuple.NewLocation("<eval>", 0, 0, 0)
	scope.global.locationLogger(location, level, fmt.Sprintf(format, args...))
}

func (scope * RunnerLocalScope) Find(context EvalContext, name Tag, args [] Value) (LocalScope, reflect.Value) {
	return scope.symbols.Find(context, name, args)
}

func (scope * RunnerLocalScope) Add(name string, function interface{}) {
	scope.symbols.Add(name, function)
}

