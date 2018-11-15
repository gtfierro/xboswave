package main

import (
	"github.com/abiosoft/ishell"
	"strings"
)

func startShell() {
	shell := ishell.New()
	shell.Println("WAVEMQ Data Ingester")

	// register a function for "greet" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "list all archive requests",
		Func: func(c *ishell.Context) {
			// list <all (default) | uripattern >
			c.Println("Hello", strings.Join(c.Args, " "))
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "add",
		Help: "add archive request",
		Func: func(c *ishell.Context) {
			c.Println("Hello", strings.Join(c.Args, " "))
		},
	})

	// run shell
	shell.Run()
}
