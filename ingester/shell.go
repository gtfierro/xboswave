package main

import (
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	"github.com/gliderlabs/ssh"
	"github.com/gtfierro/xboswave/ingester/types"
	"github.com/olekukonko/tablewriter"
	logrus "github.com/sirupsen/logrus"
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

		shell.Interrupt(func(c *ishell.Context, count int, input string) {
			c.Println("Use 'exit' or ctl-d to disconnect")
		})

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

				table := tablewriter.NewWriter(s)
				table.SetHeader([]string{"namespace", "resource", "plugin", "schema"})

				for _, req := range reqs {
					table.Append([]string{req.URI.Namespace, req.URI.Resource, req.Plugin, req.Schema})
				}
				table.Render()
			},
		})

		shell.AddCmd(&ishell.Cmd{
			Name: "add",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				if len(c.Args) != 4 {
					c.Println("add <schema> <plugin> <namespace> <resource>")
					return
				}

				req := &ArchiveRequest{
					Schema: c.Args[0],
					Plugin: c.Args[1],
					URI: types.SubscriptionURI{
						Namespace: c.Args[2],
						Resource:  c.Args[3],
					},
				}
				if err := ingest.addArchiveRequest(req); err != nil {
					logrus.Error(err)
					c.Err(err)
					return
				}
				c.Println("Successfully requested archival")
				c.Println(c.Args)
			},
		})

		// del archive requests
		shell.AddCmd(&ishell.Cmd{
			Name: "delete",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				c.Println("delete schema=<schema> plugin=<plugin> namespace=<namespace> resource=<resource>")
				if len(c.Args) == 0 {
					return
				}
				filter := &RequestFilter{}
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
						logrus.Error(parts)
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

				table := tablewriter.NewWriter(s)
				table.SetAutoMergeCells(true)
				table.SetRowLine(true)
				table.SetHeader([]string{"plugin", "namespace", "resource", "schema"})

				for _, req := range reqs {
					if err := ingest.delArchiveRequest(&req); err != nil {
						logrus.Error(err)
						c.Err(err)
						return
					}
					table.Append([]string{req.Plugin, req.URI.Namespace, req.URI.Resource, req.Schema})
				}
				table.Render()
			},
		})

		shell.AddCmd(&ishell.Cmd{
			Name: "status",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				reqs, err := ingest.cfgmgr.Status()
				if err != nil {
					c.Err(err)
					return
				}

				table := tablewriter.NewWriter(s)
				table.SetHeader([]string{"namespace", "resource", "schema", "plugin", "created", "error", "error time"})
				table.SetColumnColor(tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{tablewriter.FgHiRedColor},
					tablewriter.Colors{tablewriter.FgHiRedColor},
				)

				for _, req := range reqs {
					row := []string{req.URI.Namespace, req.URI.Resource, req.Schema, req.Plugin, req.Inserted.Format(time.RFC3339), req.LastError}
					if req.ErrorTimestamp.UnixNano() == 0 {
						row = append(row, "")
					} else {
						row = append(row, req.ErrorTimestamp.Format(time.RFC3339))
					}
					table.Append(row)
				}
				table.Render()
			},
		})

		// start shell
		shell.Run()
		// teardown
		shell.Close()
	})
	if err := ssh.ListenAndServe(":2222", nil, ssh.HostKeyFile("sshkey")); err != nil {
		logrus.Error(err)
	}
}
