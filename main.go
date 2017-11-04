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
func execute(c *Command) {
	switch c.Name {
	case "print":
		for _, v := range c.Arguments {
			fmt.Print(v)
		}
		fmt.Print("\n")
	case "if":
		arg := c.Arguments[0]
		test := arg.(func() bool)
		if ok := test(); ok == true {
			fmt.Println("=> true")
		} else {
			fmt.Println("=> false")
		}
	}
}

func fail(reason error) {
	err := fmt.Errorf("macaque: %v\n", reason)
	fmt.Printf("%v", err)
	os.Exit(1)
}

/*
	Lex is a tokenizer that depends on
	whitespace for hints
*/
func lex(linum int, line string) {
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
				err := fmt.Sprintf("line %v: variable %s has unknown type or not declared",
					linum, tokens[1])
				fail(errors.New(err))
			}
		}
		cmd.addArgument(argument)
		execute(cmd)

	case "var":
		arg := strings.Join(tokens[3:], " ")
		name := tokens[1]
		// check if value is an int
		new_arg, err := strconv.Atoi(arg)
		if err != nil {
			// int conversion failed, let's try bool
			new_arg, err := strconv.ParseBool(arg)
			if err != nil {
				// bool conversion failed, so value is a string
				vars[name] = arg
			} else {
				vars[name] = new_arg
			}
		} else {
			vars[name] = new_arg
		}

	case "if":
		if strings.HasPrefix(tokens[1], "(") {
			if strings.HasSuffix(tokens[1], ")") {
				cmd := &Command{Name: "if"}
				// the test is usually about a variable
				if len(strings.TrimSuffix(strings.TrimPrefix(tokens[1], "("), ")")) == 0 ||
					strings.TrimSuffix(strings.TrimPrefix(tokens[1], "("), ")") == "true" ||
					strings.TrimSuffix(strings.TrimPrefix(tokens[1], "("), ")") == "false" {
					fail(errors.New(fmt.Sprintf("line %d: pointless to evaluate the \"if\"", linum)))
				}
				test := strings.TrimSuffix(strings.TrimPrefix(tokens[1], "("), ")")
				v := vars[test]
				switch v.(type) {
				case int:
					cmd.addArgument(func() bool {
						if v.(int) >= 1 {
							return true
						} else {
							return false
						}
					})
				case string:
					cmd.addArgument(func() bool {
						if len(v.(string)) == 2 {
							return false
						} else {
							return true
						}
					})
				case bool:
					cmd.addArgument(func() bool {
						return v.(bool)
					})
				}
				execute(cmd)
			} else {
				err := fmt.Sprintf("line %v: missing ending parenthesis", linum)
				fail(errors.New(err))
			}
		} else {
			err := fmt.Sprintf("line %v: test unparenthesized", linum)
			fail(errors.New(err))
		}
	default:
		err := fmt.Sprintf("line %v: unknown identifer %s", linum, tokens[0])
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
		lex(linenum, line)
	}
}
