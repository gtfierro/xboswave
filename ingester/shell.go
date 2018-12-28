package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	"github.com/fatih/color"
	"github.com/gliderlabs/ssh"
)

func (ingest *Ingester) shell() {

	ssh.Handle(func(s ssh.Session) {
		//io.WriteString(s, fmt.Sprintf("Hello %s\n", s.User()))

		cfg := &readline.Config{
			Prompt:      ">>",
			Stdin:       s,
			StdinWriter: s,
			Stdout:      s,
			Stderr:      s,
		}

		shell := ishell.NewWithConfig(cfg)

		// display info.
		shell.Println("XBOS/WAVE ingester shell")

		// list archive requests
		shell.AddCmd(&ishell.Cmd{
			Name: "list",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				var filter *RequestFilter
				c.Println("list schema=<schema> plugin=<plugin> namespace=<namespace> resource=<resource>")
				if len(c.Args) > 0 {
					filter = &RequestFilter{}
				}
				for _, arg := range c.Args {
					parts := strings.SplitN(arg, "=", 2)
					if len(parts) == 1 {
						c.Println("filters need to be of form param=value")
						return
					}
					switch parts[0] {
					case "schema":
						filter.Schema = &parts[1]
					case "plugin":
						filter.Plugin = &parts[1]
					case "namespace":
						filter.Namespace = &parts[1]
					case "resource":
						filter.Resource = &parts[1]
					}
				}

				reqs, err := ingest.cfgmgr.List(filter)
				if err != nil {
					c.Err(err)
					return
				}
				for _, req := range reqs {
					c.Println(req)
				}
			},
		})

		// handle login.
		shell.AddCmd(&ishell.Cmd{
			Name: "login",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				c.Println("Let's simulate login")

				// prompt for input
				c.Print("Username: ")
				username := c.ReadLine()
				c.Print("Password: ")
				password := c.ReadPassword()

				// do something with username and password
				c.Println("Your inputs were", username, "and", password+".")

			},
			Help: "simulate a login",
		})

		// handle "greet".
		shell.AddCmd(&ishell.Cmd{
			Name:    "greet",
			Aliases: []string{"hello", "welcome"},
			Help:    "greet user",
			Func: func(c *ishell.Context) {
				name := "Stranger"
				if len(c.Args) > 0 {
					name = strings.Join(c.Args, " ")
				}
				c.Println("Hello", name)
			},
		})

		// multiple choice
		shell.AddCmd(&ishell.Cmd{
			Name: "choice",
			Help: "multiple choice prompt",
			Func: func(c *ishell.Context) {
				choice := c.MultiChoice([]string{
					"Golangers",
					"Go programmers",
					"Gophers",
					"Goers",
				}, "What are Go programmers called ?")
				if choice == 2 {
					c.Println("You got it!")
				} else {
					c.Println("Sorry, you're wrong.")
				}
			},
		})

		// multiple choice
		shell.AddCmd(&ishell.Cmd{
			Name: "checklist",
			Help: "checklist prompt",
			Func: func(c *ishell.Context) {
				languages := []string{"Python", "Go", "Haskell", "Rust"}
				choices := c.Checklist(languages,
					"What are your favourite programming languages ?",
					nil)
				out := func() (c []string) {
					for _, v := range choices {
						c = append(c, languages[v])
					}
					return
				}
				c.Println("Your choices are", strings.Join(out(), ", "))
			},
		})

		// progress bars
		{
			// determinate
			shell.AddCmd(&ishell.Cmd{
				Name: "det",
				Help: "determinate progress bar",
				Func: func(c *ishell.Context) {
					c.ProgressBar().Start()
					for i := 0; i < 101; i++ {
						c.ProgressBar().Suffix(fmt.Sprint(" ", i, "%"))
						c.ProgressBar().Progress(i)
						time.Sleep(time.Millisecond * 100)
					}
					c.ProgressBar().Stop()
				},
			})

			// indeterminate
			shell.AddCmd(&ishell.Cmd{
				Name: "ind",
				Help: "indeterminate progress bar",
				Func: func(c *ishell.Context) {
					c.ProgressBar().Indeterminate(true)
					c.ProgressBar().Start()
					time.Sleep(time.Second * 10)
					c.ProgressBar().Stop()
				},
			})
		}

		// subcommands and custom autocomplete.
		{
			var words []string
			autoCmd := &ishell.Cmd{
				Name: "suggest",
				Help: "try auto complete",
				LongHelp: `Try dynamic autocomplete by adding and removing words.
Then view the autocomplete by tabbing after "words" subcommand.
This is an example of a long help.`,
			}
			autoCmd.AddCmd(&ishell.Cmd{
				Name: "add",
				Help: "add words to autocomplete",
				Func: func(c *ishell.Context) {
					if len(c.Args) == 0 {
						c.Err(errors.New("missing word(s)"))
						return
					}
					words = append(words, c.Args...)
				},
			})

			autoCmd.AddCmd(&ishell.Cmd{
				Name: "clear",
				Help: "clear words in autocomplete",
				Func: func(c *ishell.Context) {
					words = nil
				},
			})

			autoCmd.AddCmd(&ishell.Cmd{
				Name: "words",
				Help: "add words with 'suggest add', then tab after typing 'suggest words '",
				Completer: func([]string) []string {
					return words
				},
			})

			shell.AddCmd(autoCmd)
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		boldRed := color.New(color.FgRed, color.Bold).SprintFunc()
		shell.AddCmd(&ishell.Cmd{
			Name: "color",
			Help: "color print",
			Func: func(c *ishell.Context) {
				c.Print(cyan("cyan\n"))
				c.Println(yellow("yellow"))
				c.Printf("%s\n", boldRed("bold red"))
			},
		})

		// when started with "exit" as first argument, assume non-interactive execution
		if len(os.Args) > 1 && os.Args[1] == "exit" {
			shell.Process(os.Args[2:]...)
		} else {
			// start shell
			shell.Run()
			// teardown
			shell.Close()
		}

	})
	go ssh.ListenAndServe(":2222", nil)

}
