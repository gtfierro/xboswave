package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	protoparser "github.com/yoheimuta/go-protoparser"
	parser "github.com/yoheimuta/go-protoparser/parser"
)

type Message struct {
	Name      string
	Namespace string
	Class     string
	Fields    []Field
}

type Field struct {
	Name      string
	Namespace string
	Class     string
	Number    int
}

type visitor struct {
	prefix string
}

func newWithPrefix(pfx string) *visitor {
	return &visitor{
		prefix: pfx,
	}
}

type messageVisitor struct {
	prefix   string
	messages []*Message
	*visitor
}

func newMessageVisitor(pfx string) *messageVisitor {
	return &messageVisitor{
		prefix:  pfx,
		visitor: newWithPrefix(pfx),
	}
}
func (v *messageVisitor) VisitMessage(msg *parser.Message) (next bool) {
	prefix := fmt.Sprintf("%s%s.", v.prefix, msg.MessageName)
	//fmt.Printf("message: %s\n", prefix)
	v.messages = append(v.messages, &Message{
		Name: msg.MessageName,
	})
	fieldvisitor := newFieldVisitor(prefix)
	for _, element := range msg.MessageBody {
		element.Accept(fieldvisitor)
	}
	lastmsg := v.messages[len(v.messages)-1]
	for _, field := range fieldvisitor.fields {
		lastmsg.Fields = append(lastmsg.Fields, *field)
	}
	return true
}
func (v *messageVisitor) VisitOption(opt *parser.Option) (next bool) {
	lastidx := len(v.messages) - 1
	if opt.OptionName == "(brick_equip_class).namespace" {
		v.messages[lastidx].Namespace = opt.Constant
	}
	if opt.OptionName == "(brick_equip_class).value" {
		v.messages[lastidx].Class = opt.Constant
	}
	//fmt.Printf("%+v\n", opt)
	return true
}

type fieldVisitor struct {
	prefix string
	fields []*Field
	*visitor
}

func newFieldVisitor(pfx string) *fieldVisitor {
	return &fieldVisitor{
		prefix:  pfx,
		visitor: newWithPrefix(pfx),
	}
}
func (v *fieldVisitor) VisitField(field *parser.Field) (next bool) {
	//fmt.Printf("field: %s%s = %s\n", v.prefix, field.FieldName, field.FieldNumber)
	i, err := strconv.Atoi(field.FieldNumber)
	if err != nil {
		log.Fatal(err)
	}

	myfield := &Field{
		Name:   field.FieldName,
		Number: i,
	}
	for _, opt := range field.FieldOptions {
		if opt.OptionName == "(brick_point_class).namespace" {
			myfield.Namespace = opt.Constant
		}
		if opt.OptionName == "(brick_point_class).value" {
			myfield.Class = opt.Constant
		}
	}
	v.fields = append(v.fields, myfield)
	return true
}

func (v *visitor) VisitComment(*parser.Comment) {
}
func (v *visitor) VisitEmptyStatement(*parser.EmptyStatement) (next bool) {
	return true
}
func (v *visitor) VisitEnum(*parser.Enum) (next bool) {
	return true
}
func (v *visitor) VisitEnumField(*parser.EnumField) (next bool) {
	return true
}
func (v *visitor) VisitExtend(*parser.Extend) (next bool) {
	return true
}
func (v *visitor) VisitField(field *parser.Field) (next bool) {
	return true
}
func (v *visitor) VisitImport(*parser.Import) (next bool) {
	return true
}
func (v *visitor) VisitMapField(*parser.MapField) (next bool) {
	return true
}
func (v *visitor) VisitMessage(msg *parser.Message) (next bool) {
	return true
}
func (v *visitor) VisitOneof(*parser.Oneof) (next bool) {
	return true
}
func (v *visitor) VisitOneofField(*parser.OneofField) (next bool) {
	return true
}
func (v *visitor) VisitOption(*parser.Option) (next bool) {
	return true
}
func (v *visitor) VisitPackage(*parser.Package) (next bool) {
	return true
}
func (v *visitor) VisitReserved(*parser.Reserved) (next bool) {
	return true
}
func (v *visitor) VisitRPC(*parser.RPC) (next bool) {
	return true
}
func (v *visitor) VisitService(*parser.Service) (next bool) {
	return true
}
func (v *visitor) VisitSyntax(*parser.Syntax) (next bool) {
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: protodef_to_json path/to/protofile.proto")
		os.Exit(0)
	}
	filename := os.Args[1]
	//fmt.Println(filename)

	reader, _ := os.Open(filename)
	defer reader.Close()

	got, err := protoparser.Parse(
		reader,
		protoparser.WithDebug(false),
		protoparser.WithPermissive(true),
		protoparser.WithFilename(filename),
	)
	if err != nil {
		log.Fatal(err)
	}

	toplevel := newMessageVisitor(".")
	got.Accept(toplevel)

	gotJSON, err := json.MarshalIndent(toplevel.messages, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}
	fmt.Print(string(gotJSON))
}
