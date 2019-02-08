package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/abiosoft/readline"
	"github.com/gliderlabs/ssh"
	"github.com/gtfierro/xboswave/ingester/types"
	"github.com/olekukonko/tablewriter"
	logrus "github.com/sirupsen/logrus"
)

func parseFilterFromArgs(args []string) (*RequestFilter, error) {
	if len(args) == 0 {
		return nil, nil
	}
	filter := &RequestFilter{}
	for _, arg := range args {
		if arg == "enabled" {
			filter.Enabled = &_TRUE
			continue
		}
		if arg == "hasError" {
			filter.HasError = &_TRUE
			continue
		}
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("filters need to be of form param=value")
		}
		switch parts[0] {
		case "id":
			_id, err := strconv.Atoi(parts[1])
			filter.Id = &_id
			if err != nil {
				return nil, err
			}
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
	return filter, nil
}

func (ingest *Ingester) shell(cfg Config) {

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

				c.Println("delete id=<id> schema=<schema> plugin=<plugin> namespace=<namespace> resource=<resource>")
				filter, err := parseFilterFromArgs(c.Args)
				if err != nil {
					c.Println(err.Error())
					return
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

		// enable archive requests
		shell.AddCmd(&ishell.Cmd{
			Name: "enable",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				c.Println("enable id=<id> schema=<schema> plugin=<plugin> namespace=<namespace> resource=<resource>")
				filter, err := parseFilterFromArgs(c.Args)
				if err != nil {
					c.Println(err.Error())
					return
				}

				reqs, err := ingest.cfgmgr.List(filter)
				if err != nil {
					c.Err(err)
					return
				}

				for _, req := range reqs {
					c.Printf("Enabling req %d...\n", req.Id)
					if err := ingest.enableArchiveRequest(&req); err != nil {
						logrus.Error(err)
						c.Err(err)
						return
					}
				}
			},
		})

		shell.AddCmd(&ishell.Cmd{
			Name: "disable",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				c.Println("disable id=<id> schema=<schema> plugin=<plugin> namespace=<namespace> resource=<resource>")
				filter, err := parseFilterFromArgs(c.Args)
				if err != nil {
					c.Println(err.Error())
					return
				}

				reqs, err := ingest.cfgmgr.List(filter)
				if err != nil {
					c.Err(err)
					return
				}

				for _, req := range reqs {
					c.Printf("Disabling req %d...\n", req.Id)
					if err := ingest.disableArchiveRequest(&req); err != nil {
						logrus.Error(err)
						c.Err(err)
						return
					}
				}
			},
		})

		shell.AddCmd(&ishell.Cmd{
			Name: "list",
			Func: func(c *ishell.Context) {
				c.ShowPrompt(false)
				defer c.ShowPrompt(true)

				// build filter
				filter, err := parseFilterFromArgs(c.Args)
				if err != nil {
					c.Err(err)
					return
				}

				reqs, err := ingest.cfgmgr.List(filter)
				if err != nil {
					c.Err(err)
					return
				}

				table := tablewriter.NewWriter(s)
				table.SetHeader([]string{"id", "enabled?", "namespace", "resource", "schema", "plugin", "created", "error", "error time"})
				table.SetColumnColor(tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{0},
					tablewriter.Colors{tablewriter.FgHiRedColor},
					tablewriter.Colors{tablewriter.FgHiRedColor},
				)

				for _, req := range reqs {
					var enabledStr string
					if req.Enabled {
						enabledStr = "1"
					} else {
						enabledStr = "0"
					}
					row := []string{fmt.Sprintf("%d", req.Id), enabledStr, req.URI.Namespace, req.URI.Resource, req.Schema, req.Plugin, req.Inserted.Format(time.RFC3339), req.LastError}
					if req.ErrorTimestamp.UnixNano() <= 0 {
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
	if cfg.IngesterShell.PasswordLogin {
		if err := ssh.ListenAndServe("localhost:2222", nil, ssh.HostKeyFile(cfg.IngesterShell.SshHostKey), ssh.PasswordAuth(func(ctx ssh.Context, pass string) bool { return pass == cfg.IngesterShell.Password })); err != nil {
			logrus.Fatal(err)
		}
	} else if err := ssh.ListenAndServe("localhost:2222", nil, ssh.HostKeyFile(cfg.IngesterShell.SshHostKey)); err != nil {
		logrus.Fatal(err)
	}
}
