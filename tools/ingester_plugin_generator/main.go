package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/emicklei/proto"
)

var toparse []string
var subfields [][]string
var parsefields = make(map[string][]string)
var parsefieldstypes = make(map[string][]string)
var parsesubfields = make(map[string][]string)
var parseunits = make(map[string]string)
var f *File
var added = make(map[string]struct{})
var d = make(Dict)
var unitsadded = make(map[string]struct{})
var unitsd = make(Dict)
var unitdefre = regexp.MustCompile(`unit:\s*([a-zA-Z0-9_\-/%# ]+)`)

func main() {

	// base file contains the top-level
	basefile := "../../proto/xbos.proto"
	//filename := "../../proto/iot.proto"
	if len(os.Args) < 3 {
		fmt.Println("ingester_plugin_generator <proto file> <output go file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	destfilename := os.Args[2]

	f = NewFile("main")
	f.PackageComment("This file is AUTOGENERATED")

	files := []string{
		basefile,
		filename,
	}
	for _, filename := range files {
		reader, _ := os.Open(filename)
		defer reader.Close()

		parser := proto.NewParser(reader)
		definition, _ := parser.Parse()

		proto.Walk(definition,
			proto.WithMessage(handleMessage),
		)
		proto.Walk(definition,
			proto.WithMessage(handleMessage2),
		)
	}
	fmt.Println(toparse)

	fmt.Println("-------------")
	//for fieldname, subfields := range parsefields {
	//	for _, subfieldname := range subfields {
	//		fmt.Println(fieldname, subfieldname)
	//	}
	//}
	//fmt.Println("-------------")
	//for fieldname, subfields := range parsesubfields {
	//	for _, subfieldname := range subfields {
	//		fmt.Println("sub>", fieldname, subfieldname)
	//	}
	//}

	f.Var().Id("lookup").Op("=").Map(String()).Func().Params(
		Id("msg").Id("xbospb").Dot("XBOS"),
	).Float64().Values(d)
	f.Var().Id("units").Op("=").Map(String()).String().Values(unitsd)

	destfile, err := os.Create(destfilename + ".go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Printf("%#v\n", f)
	if err := f.Render(destfile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func handleService(s *proto.Service) {
	fmt.Println(s.Name)
}
func handleOption(s *proto.Option) {
	fmt.Println(s.Name)
}

var scalartypes = []string{
	"double",
	"float",
	"int32",
	"int64",
	"uint32",
	"uint64",
	"bool",
}

var disallowed = []string{
	"string",
	"bytes",
}

func contains(l []string, v string) int {
	for idx, s := range l {
		if s == v {
			return idx
		}
	}
	return -1
}

func handleMessage(m *proto.Message) {
	//fmt.Println("Message:", m.Name, len(m.Elements), m.Parent.(*proto.Proto).Filename)

	for _, v := range m.Elements {
		f := v.(*proto.NormalField)
		//fmt.Println(f.Name, f.Type, f.Doc())
		if f.Doc() != nil {
			lines := f.Doc().Lines
			// This looks for every comment line above a field definition
			// for a line starting with 'unit:'. Everything after 'unit:'
			// will be considered the engineering units for this field.
			// Valid characters are: a-zA-Z0-9_/%# and '-' and space.
			// Whitespace is stripped from the unit field.
			//
			//   //this is the unit temperature
			//   //unit: celsius
			//   //double temperature = 1;
			//
			parseunits[f.Name] = "unknown"
			for _, l := range lines {
				ud := pullUnitDef(l)
				if ud != "unknown" {
					parseunits[f.Name] = ud
					break

				}
			}
		}

		// keep track of the top-level messages in other files that we want to parse
		if m.Name == "XBOS" {
			//fmt.Println("PARSE PARSE", f.Type)
			toparse = append(toparse, f.Type)
		} else {

			// otherwise, keep track of the field path
			//if contains(toparse, m.Name) {
			if contains(scalartypes, f.Type) >= 0 {
				fmt.Println("extract", m.Name, ".", f.Name, f.Type)
				//subfields = append(subfields, []string{m.Name})
				parsefields[m.Name] = append(parsefields[m.Name], f.Name)
				parsefieldstypes[m.Name] = append(parsefieldstypes[m.Name], f.Type)
			} else if contains(disallowed, f.Type) < 0 {
				//fmt.Println("subtype", m.Name, ".", f.Name, f.Type)
				parsesubfields[f.Name] = append(parsefields[f.Name], f.Type)
			}
		}
		//}

	}

}

func handleMessage2(m *proto.Message) {

	if contains(toparse, m.Name) >= 0 {
		//fmt.Println("Parse fields in", m.Name)

		for _, v := range m.Elements {
			field := v.(*proto.NormalField)
			propname := field.Name // key
			fmt.Println("TYPE", field.Name, field.Type)
			parsefieldidx := contains(parsefields[m.Name], propname)
			if parsefieldidx >= 0 {
				//fmt.Println("  parsing", propname)
				protofield := strings.Replace(strings.Title(strings.Replace(propname, "_", " ", -1)), " ", "", -1)

				// term if the field "<m.Name>.<protofield>" can be converted to float64
				var _if = Return(Float64().Call(Id("msg").Dot(m.Name).Dot(protofield)))
				if parsefieldstypes[m.Name][parsefieldidx] == "bool" {
					_if = If(Id("msg").Dot(m.Name).Dot(protofield)).Block(Return(Lit(1))).Else().Block(Return(Lit(0)))
				}

				_, found := added[propname]
				if propname != "time" && !found {
					d[Lit(propname)] = Func().Params(
						Id("msg").Id("xbospb").Dot("XBOS"),
					).Float64().Block(
						_if,
					)
					added[propname] = struct{}{}

					_, found = unitsadded[propname]
					if !found {
						unitsd[Lit(propname)] = Lit(parseunits[propname])
						unitsadded[propname] = struct{}{}
					}
				}
			} else if contains(parsesubfields[propname], field.Type) >= 0 {
				for _, subfield := range parsefields[field.Type] {
					//fmt.Println("  parsing", propname, ".", subfield)
					protofield := strings.Replace(strings.Title(strings.Replace(propname, "_", " ", -1)), " ", "", -1)
					subprotofield := strings.Replace(strings.Title(strings.Replace(subfield, "_", " ", -1)), " ", "", -1)

					_, found := unitsadded[subfield]
					if !found {
						unitsd[Lit(subfield)] = Lit(parseunits[subfield])
						unitsadded[subfield] = struct{}{}
					}

					idx := contains(parsefields[protofield], subfield)
					var _if = Return(Float64().Call(Id("msg").Dot(m.Name).Dot(protofield).Dot(subprotofield)))
					if len(parsefields[protofield]) > 0 && parsefieldstypes[protofield][idx] == "bool" {
						_if = If(Id("msg").Dot(m.Name).Dot(protofield).Dot(subprotofield)).Block(Return(Lit(1))).Else().Block(Return(Lit(0)))
					}
					_, found = added[subfield]
					if subfield != "time" && !found {
						d[Lit(subfield)] = Func().Params(
							Id("msg").Id("xbospb").Dot("XBOS"),
						).Float64().Block(
							_if, //Float64().Call(Id("msg").Dot(m.Name).Dot(protofield).Dot(subprotofield)),
						)
						added[subfield] = struct{}{}
					}
				}
			}
		}
	} else {
		fmt.Println("NO PARSE", m.Name)
	}
	return
}

func pullUnitDef(s string) string {

	groups := unitdefre.FindStringSubmatch(s)
	if len(groups) > 1 {
		return groups[1]
	}
	return "unknown"
}