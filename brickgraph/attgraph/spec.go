package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Spec struct {
	Graph    graph
	Edges    []edge
	Policies []policy
}

// returns list of all entities in the spec
func (spec Spec) Entities() []string {
	var ents = make(map[string]struct{})

	// graph namespaces
	for _, ns := range spec.Graph.Namespaces {
		ents[ns] = struct{}{}
	}

	// edges
	for _, edge := range spec.Edges {
		ents[edge.Namespace] = struct{}{}
		ents[edge.From] = struct{}{}
		ents[edge.To] = struct{}{}
	}

	// policies
	for _, pol := range spec.Policies {
		ents[pol.Namespace] = struct{}{}
	}

	// assemble into list and return
	var ret []string
	for ent := range ents {
		if ent != "" {
			ret = append(ret, ent)
		}
	}
	return ret
}

type graph struct {
	Namespaces   []string
	EntityExpiry duration `toml:"entity_expiry"`
	GrantExpiry  duration `toml:"grant_expiry"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	dur, err := ParseDuration(string(text))
	if err != nil {
		return err
	}
	d.Duration = *dur
	return nil
}

type policy struct {
	Name        string
	Namespace   string
	Resource    string
	Edge        string
	Permissions string
	Pset        string
}

type edge struct {
	From        string
	To          string
	Policy      string
	Namespace   string
	Resource    string
	Edge        string
	Permissions string
	Pset        string
	TTL         int
	Expiry      *duration
}

func (e edge) String() string {
	return fmt.Sprintf("%s => %s %s:%s@%s/%s", e.From, e.To, e.Pset, e.Permissions, e.Namespace, e.Resource)
}

//ParseDuration is a little like the existing time.ParseDuration
//but adds days and years because its really annoying not having that
// from Michael Andersen
func ParseDuration(s string) (*time.Duration, error) {
	if s == "" {
		return nil, nil
	}
	pat := regexp.MustCompile(`^(\d+y)?(\d+d)?(\d+h)?(\d+m)?(\d+s)?$`)
	res := pat.FindStringSubmatch(s)
	if res == nil {
		return nil, fmt.Errorf("Invalid duration")
	}
	res = res[1:]
	sec := int64(0)
	for idx, mul := range []int64{365 * 24 * 60 * 60, 24 * 60 * 60, 60 * 60, 60, 1} {
		if res[idx] != "" {
			key := res[idx][:len(res[idx])-1]
			v, e := strconv.ParseInt(key, 10, 64)
			if e != nil { //unlikely
				return nil, e
			}
			sec += v * mul
		}
	}
	rv := time.Duration(sec) * time.Second
	return &rv, nil
}
