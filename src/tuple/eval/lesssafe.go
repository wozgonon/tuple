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

import "reflect"
import "fmt"
import "log"
import "os/exec"
import "os"
import "bytes"

// These functions are potentially not harmless since they can access resources out of the sandbox
func NewLessSafeSymbolTable() SymbolTable {
	table := NewSafeSymbolTable(&ExecIfNotFound{})
	AddOperatingSystemFunctions(&table)
	return table
}

/////////////////////////////////////////////////////////////////////////////
// Operating system access
/////////////////////////////////////////////////////////////////////////////

//
//  Add shell specific commands
//  These are typically not a 'safe' in that they allow access to the file system
// 
func Pipe(writer string, reader string) bool {

	log.Printf("pipe '%s', '%s'", writer, reader)
	cmdw := exec.Command(writer)
	stdoutw, err := cmdw.StdoutPipe()

	err = cmdw.Start()
	if err != nil {
		log.Printf("Start 1 finished with error: %v", err)
		return false
	}

	cmdr := exec.Command(reader)
	cmdr.Stdin = stdoutw
	cmdr.Stdout = os.Stdout

	if err != nil {
		log.Printf("StdoutPipe finished with error: %v", err)
		log.Fatal(err)
	}
	err = cmdr.Start()
	if err != nil {
		log.Printf("Start 2 finished with error: %v", err)
		log.Fatal(err)
	}
	err = cmdw.Wait()
	if err != nil {
		log.Printf("Wait 1 finished with error: %v", err)
		return false
	}
	err = cmdr.Wait()
	if ; err != nil {
		log.Printf("Wait 2 finished with error: %v", err)
		log.Fatal(err)
	}

	return true
}

func executeProcess(arg string, args... string) bool {
	cmd := exec.Command(arg, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return false
	}
	outStr := string(stdout.Bytes())
	fmt.Print(outStr)
	return true
}

// TODO sort out variadic arguments
func executeProcess0(arg string) bool {
	return executeProcess(arg)
}


func spawnProcess (arg string) bool {
	cmd := exec.Command(arg)
	//var stdout bytes.Buffer
	cmd.Stdout = os.Stdout //&stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		return false
	}
	//outStr := string(stdout.Bytes())
	//fmt.Print(outStr)
	return true
}


// These functions access the operating system and so an actor
// could use them to harm the computer.
func AddOperatingSystemFunctions(table * SymbolTable) {
	//
	//  Add shell specific commands
	//  These are typically not a 'safe' in that they allow access to the file system
	// 
	table.Add("exec", executeProcess0)
	table.Add("spawn", spawnProcess)
	table.Add("pipe", Pipe)

	table.Add("echo", func (context EvalContext, values... Value) bool {
		for k,_:= range values {
			evaluated := Eval(context, values[k])  // TODO This should use table from context so that it uses scope
			fmt.Print(toString(context, evaluated))
		}
		return true
	})
	table.Add("eval", func (context EvalContext, value Value) Value {
		return Eval(context, value)
	})
}

/////////////////////////////////////////////////////////////////////////////

type ExecIfNotFound struct {}

func (exec * ExecIfNotFound) Find (context EvalContext, name Tag, args [] Value) (*SymbolTable, reflect.Value) {

	return nil, reflect.ValueOf(func(context EvalContext, args... Value) bool {
		return executeProcess(name.Name, EvalToStrings(context, args)...)
	})
}

