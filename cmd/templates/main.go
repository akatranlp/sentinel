package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/akatranlp/sentinel/web"
)

type StringWriter interface {
	io.Writer
	io.StringWriter
}

func checkIfDir(path string) error {
	p, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !p.IsDir() {
		return fmt.Errorf("path is not a dir")
	}
	return nil
}

func spreadType(type_ reflect.Type, sb StringWriter, indent int) {
	for i := range type_.NumField() {
		field := type_.Field(i)
		if field.Anonymous {
			continue
		}
		appendField(field, false, sb, indent)
	}
}

func appendType(type_ reflect.Type, nullable bool, sb StringWriter, indent int) {
	// indentStr := strings.Repeat(" ", indent)
	switch type_.Kind() {
	case reflect.Pointer:
		appendType(type_.Elem(), true, sb, indent)
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		appendType(type_.Elem(), false, sb, indent)
		sb.WriteString("[] | null")
	case reflect.Struct:
		sb.WriteString(type_.Name())
		// sb.WriteString("{\n")
		// spreadType(type_, sb, indent+2)
		// sb.WriteString(indentStr)
		// sb.WriteString("}")
	case reflect.Bool:
		sb.WriteString("boolean")
	case reflect.String:
		sb.WriteString(type_.Name())
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		sb.WriteString("number")
	}
	if nullable {
		sb.WriteString(" | null")
	}
}

func appendField(field reflect.StructField, nullable bool, sb StringWriter, indent int) {
	fieldName := field.Name

	if value, ok := field.Tag.Lookup("json"); ok {
		fieldName = value
	}

	indentStr := strings.Repeat(" ", indent)
	sb.WriteString(indentStr)
	sb.WriteString(fieldName)
	sb.WriteString(": ")

	appendType(field.Type, nullable, sb, indent)

	sb.WriteString(";\n")
}

func generateType[T any](sb StringWriter) error {
	type_ := reflect.TypeFor[T]()

	typeName := type_.Name()

	if typeName == "SentinelCtx" {
		typeName = "CommonSentinelCtx"
	}

	fmt.Fprintf(sb, "export type %s = {\n", typeName)

	spreadType(type_, sb, 2)

	sb.WriteString("}\n")

	return nil
}

func generateExtendedType[T any](sb StringWriter) error {
	type_ := reflect.TypeFor[T]()

	fmt.Fprintf(sb, "export type %s = CommonSentinelCtx & {\n", type_.Name())

	pageID, ok := pageIDContextMap[type_.Name()]
	if !ok {
		panic(type_.Name())
	}

	fmt.Fprintf(sb, "  pageId: \"%s\"\n", pageID)

	spreadType(type_, sb, 2)

	sb.WriteString("}\n")

	return nil
}

func generateEnum[T fmt.Stringer](sb StringWriter, values []T) error {
	type_ := reflect.TypeFor[T]()

	typeName := type_.Name()
	constName := []byte(typeName)
	constName[0] += 32

	fmt.Fprintf(sb, "export const %s = [", string(constName))

	for i, v := range values {
		if i > 0 {
			sb.WriteString(", ")
		}
		fmt.Fprintf(sb, "\"%s\"", v.String())
	}

	sb.WriteString("] as const;\n")
	fmt.Fprintf(sb, "export type %s = typeof %s[number]\n", typeName, constName)

	return nil
}

func generateSentinelCtx(sb StringWriter) error {
	sb.WriteString("export type SentinelCtx = \n  ")

	var i int
	for k := range pageIDContextMap {
		if i > 0 {
			sb.WriteString(" |\n  ")
		}
		i++
		sb.WriteString(k)
	}
	sb.WriteString("\n")
	sb.WriteString("\n")

	sb.WriteString("export type ExtractSentinelCtx<T extends PageID> = Extract<SentinelCtx, { pageId: T }>\n")
	sb.WriteString("\n")

	sb.WriteString("export type Prettify<T> = {\n  [K in keyof T]: T[K]\n}\n")
	sb.WriteString("\n")

	return nil
}

func generateTypes(folderPath string) error {
	var sb bytes.Buffer
	generateEnum(&sb, web.PageIDValues())
	sb.WriteString("\n")
	generateType[web.SentinelCtx](&sb)
	sb.WriteString("\n")
	generateExtendedType[web.LoginSentinelCtx](&sb)
	sb.WriteString("\n")
	generateExtendedType[web.FormRedirectSentinelCtx](&sb)
	sb.WriteString("\n")
	generateExtendedType[web.FormPostSentinelCtx](&sb)
	sb.WriteString("\n")
	generateExtendedType[web.InfoSentinelCtx](&sb)
	sb.WriteString("\n")
	generateExtendedType[web.ErrorSentinelCtx](&sb)
	sb.WriteString("\n")

	generateEnum(&sb, web.MessageTypeValues())
	sb.WriteString("\n")
	generateType[web.Message](&sb)
	sb.WriteString("\n")

	generateType[web.Provider](&sb)
	sb.WriteString("\n")
	generateType[web.CSRF](&sb)
	sb.WriteString("\n")
	generateType[web.URLs](&sb)
	sb.WriteString("\n")

	generateSentinelCtx(&sb)

	f, err := os.Create(folderPath + "/types.ts")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, &sb)
	return err
}

var pageIDContextMap = map[string]string{
	"LoginSentinelCtx":        string(web.PageIDLogintmpl),
	"ErrorSentinelCtx":        string(web.PageIDErrortmpl),
	"InfoSentinelCtx":         string(web.PageIDInfotmpl),
	"FormPostSentinelCtx":     string(web.PageIDFormPosttmpl),
	"FormRedirectSentinelCtx": string(web.PageIDFormRedirecttmpl),
}

func run(_ context.Context) error {
	var err error
	pathOption := flag.String("project", "web", "Sentinel Frontend Folder")
	flag.Parse()

	webPath := *pathOption
	if err = checkIfDir(webPath); err != nil {
		return err
	}

	typesFolder := webPath + "/src/context"
	if err = checkIfDir(typesFolder); err != nil {
		return fmt.Errorf("context folder is not there please rebuild")
	}

	if err = generateTypes(typesFolder); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
