package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/olekukonko/tablewriter"
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

	db.shell.AddCmd(&ishell.Cmd{
		Name: "newp",
		Help: "Create policy",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 5 {
				c.Println("newp namespace pset indirections perm1[,perm2,...]<resource")
				return
			}
			namespace := c.Args[0]
			pset := c.Args[1]
			indir, err := strconv.Atoi(c.Args[2])
			if err != nil {
				c.Err(err)
				return
			}
			perms := strings.Split(c.Args[3], ",")
			resource := c.Args[4]
			if err := db.CreatePolicy(namespace, resource, pset, indir, perms); err != nil {
				if err != nil {
					c.Err(err)
					return
				}
			}
		},
	})

	db.shell.AddCmd(&ishell.Cmd{
		Name: "grant",
		Help: "Grant",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 3 {
				c.Println("grant subject expiry policyid1[,policyid2,...]")
				return
			}
			subject := c.Args[0]
			expiry, err := ParseDuration(c.Args[1])
			if err != nil {
				c.Err(err)
				return
			}
			_policyids := strings.Split(c.Args[2], ",")
			var policyids []int
			for _, _pid := range _policyids {
				pid, err := strconv.Atoi(_pid)
				if err != nil {
					c.Err(err)
					return
				}
				policyids = append(policyids, pid)
			}

			policies, err := db.getPoliciesById(policyids)
			if err != nil {
				c.Err(err)
				return
			}

			fmt.Println(subject, expiry, policies)

			err = db.CreateAttestation(subject, time.Now().Add(*expiry), policies)
			if err != nil {
				c.Err(err)
				return
			}

		},
	})

	// list policies that meet the given filters
	db.shell.AddCmd(&ishell.Cmd{
		Name: "attestations",
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
			table := tablewriter.NewWriter(os.Stdout)
			table.SetRowLine(true)
			table.SetHeader([]string{"hash", "subject", "valid until", "expires in", "policies"})
			for _, a := range atts {
				var _policies []string
				for _, pol := range a.PolicyStatements {
					_policies = append(_policies, fmt.Sprintf("%d", pol.id))
				}
				expires := time.Until(a.ValidUntil)
				expires = time.Second * time.Duration(expires.Seconds())
				table.Append([]string{
					a.Hash,
					a.Subject,
					a.ValidUntil.Format(time.RFC822Z),
					fmt.Sprintf("%s", expires),
					strings.Join(_policies, ", "),
				})
			}
			table.Render()
		},
	})

	// list policies that meet the given filters
	db.shell.AddCmd(&ishell.Cmd{
		Name: "policies",
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

			table := tablewriter.NewWriter(os.Stdout)
			table.SetRowLine(true)
			table.SetHeader([]string{"id", "namespace", "resource", "indir", "pset", "perms"})

			for _, p := range pols {
				_id := fmt.Sprintf("%d", p.id)
				_indir := fmt.Sprintf("%d", p.Indirections)
				_perms := fmt.Sprintf("%s", p.Permissions)
				table.Append([]string{_id, p.Namespace, p.Resource, _indir, p.PermissionSet, _perms})
			}
			table.Render()
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
