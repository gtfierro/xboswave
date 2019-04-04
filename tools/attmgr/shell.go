package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/abiosoft/ishell"
	//"github.com/olekukonko/tablewriter"
)

func parseFilterFromArgs(args []string) (*filter, error) {
	if len(args) == 0 {
		return nil, nil
	}
	filter := &filter{}
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 1 {
			return nil, fmt.Errorf("filters need to be of form param=value")
		}
		switch parts[0] {
		case "attid":
			_id, err := strconv.Atoi(parts[1])
			filter.attid = &_id
			if err != nil {
				return nil, err
			}
		case "polid":
			_id, err := strconv.Atoi(parts[1])
			filter.polid = &_id
			if err != nil {
				return nil, err
			}
		case "hash":
			filter.hash = &parts[1]
		case "policy":
			_pol, err := strconv.Atoi(parts[1])
			filter.policy = &_pol
			if err != nil {
				return nil, err
			}
		case "namespace":
			filter.namespace = &parts[1]
		case "resource":
			filter.resource = &parts[1]
		case "pset":
			filter.pset = &parts[1]
			//case "permissions":
			//	filter.permissions = &parts[1]
		}
	}
	return filter, nil
}

func (db *DB) setupShell() {

	db.shell = ishell.New()
	db.shell.SetHomeHistoryPath(".attmgr_history")

	// list policies that meet the given filters
	db.shell.AddCmd(&ishell.Cmd{
		Name: "listatt",
		Help: "List attestations",
		Func: func(c *ishell.Context) {
			filter, err := parseFilterFromArgs(c.Args)
			if err != nil {
				c.Err(err)
				return
			}
			atts, err := db.listAttestation(filter)
			if err != nil {
				c.Err(err)
				return
			}
			for _, a := range atts {
				fmt.Printf("%+v\n", a)
			}
		},
	})

	// list policies that meet the given filters
	db.shell.AddCmd(&ishell.Cmd{
		Name: "listpol",
		Help: "List policies",
		Func: func(c *ishell.Context) {
			filter, err := parseFilterFromArgs(c.Args)
			if err != nil {
				c.Err(err)
				return
			}
			pols, err := db.listPolicy(filter)
			if err != nil {
				c.Err(err)
				return
			}
			for _, p := range pols {
				fmt.Printf("%+v\n", p)
			}
		},
	})

	db.shell.AddCmd(&ishell.Cmd{
		Name: "list",
		Help: "List unique properteis of policies, etc",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Println("Must be one of [namespace attestation policy resource pset permissions]")
				return
			}
			switch c.Args[0] {
			case "ns", "namespace", "namespaces":
				ns, err := db.getUniqueFromPolicy("namespace", nil)
				if err != nil {
					c.Err(err)
					return
				}
				for _, n := range ns {
					c.Println(n)
				}
			case "resource":
				ns, err := db.getUniqueFromPolicy("resource", nil)
				if err != nil {
					c.Err(err)
					return
				}
				for _, n := range ns {
					c.Println(n)
				}
			case "pset":
				ns, err := db.getUniqueFromPolicy("pset", nil)
				if err != nil {
					c.Err(err)
					return
				}
				for _, n := range ns {
					c.Println(n)
				}
			case "permissions":
				ns, err := db.getUniqueFromPolicy("permissions", nil)
				if err != nil {
					c.Err(err)
					return
				}
				for _, n := range ns {
					c.Println(n)
				}
			case "att", "attestation", "attestations":
			}
		},
	})
}
