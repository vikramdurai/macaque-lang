/*
	Macaque is a simple interpreter (which will be based off the Python
	general purpose programming language) made in Golang.

	Right now, all it has is a `print` statement, and variables
	using the `var` keyword.
*/
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Command struct {
	Name      string
	Arguments []interface{}
}

var vars = make(map[string]interface{})

/*
	addArgument is convenience method for adding arguments
	to a command.
*/
func (c *Command) addArgument(arg interface{}) {
	c.Arguments = append(c.Arguments, arg)
}

/*
	Execute executes a command, by using what's
	implied by the command name.
*/
func Execute(c *Command) {
	switch c.Name {
	case "print":
		for _, v := range c.Arguments {
			fmt.Print(v)
		}
		fmt.Print("\n")
	}
}

func fail(reason error) {
	fmt.Errorf("macaque: %v\n", reason)
	os.Exit(1)
}

/*
	Lex is a tokenizer that depends on
	whitespace for hints
*/
func Lex(linum int, line string) {
	tokens := strings.Split(line, " ")
	switch tokens[0] {
	case "print":
		cmd := &Command{Name: "print"}
		var argument string
		if strings.HasPrefix(tokens[1], "'") && strings.HasSuffix(tokens[1], "'") ||
			strings.HasPrefix(tokens[1], "\"") && strings.HasSuffix(tokens[1], "\"") {
			argument = strings.Join(tokens[1:len(tokens)-1], " ")
		}
		if v, ok := vars[tokens[1]]; ok == true {
			switch v.(type) {
			case string:
				argument = v.(string)
			case int:
				argument = strconv.Itoa(v.(int))
			default:
				err := fmt.Sprintf("line %d: variable %s has unknown type or not declared",
					linum, tokens[1])
				fail(errors.New(err))
			}
		}
		cmd.addArgument(argument)
		Execute(cmd)
	case "var":
		vars[tokens[1]] = strings.Join(tokens[3:], " ")
	default:
		err := fmt.Sprintf("line %d: unknown identifer %s", linum, tokens[0])
		fail(errors.New(err))
	}
}

func main() {
	if len(os.Args) == 1 {
		fail(errors.New("program name not provided"))
	}
	progname := os.Args[1]
	file, err := ioutil.ReadFile(progname)
	if err != nil {
		fail(err)
	}
	program := strings.Split(string(file), "\n")
	// iterate over the lines in the program
	for linenum, line := range program {
		Lex(linenum, line)
	}
}
