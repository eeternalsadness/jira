package util

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

func PrettyPrintStructSlice[T any](headerMap map[string]string, structSlice []T) error {
	if len(structSlice) == 0 {
		panic("structSlice must not be empty!")
	}

	if t := reflect.TypeOf(structSlice[0]).Kind(); t != reflect.Struct {
		panic(fmt.Sprintf("expected struct slice, got %s", t.String()))
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	var builder strings.Builder

	// print out headers
	builder.WriteString("#\t")
	if len(headerMap) > 0 {
		// use header map if it's passed in
		for header := range headerMap {
			_, err := builder.WriteString(fmt.Sprintf("%s\t", header))
			if err != nil {
				return err
			}
		}
	} else {
		// print all struct fields if header map not passed in
		t := reflect.TypeOf(structSlice[0])
		for i := range t.NumField() {
			_, err := builder.WriteString(fmt.Sprintf("%s\t", t.Field(i).Name))
			if err != nil {
				return err
			}
		}
	}
	fmt.Fprintln(w, builder.String())
	builder.Reset()

	// print out struct slice
	for i, st := range structSlice {
		v := reflect.ValueOf(st)

		if len(headerMap) > 0 {
			// use the header map to generate tab-separated string of values
			builder.WriteString(fmt.Sprintf("%d\t", i+1))
			for _, fieldName := range headerMap {
				builder.WriteString(fmt.Sprintf("%s\t", v.FieldByName(fieldName).String()))
			}
		} else {
			// print all values if header map not passed in
			for i := range v.NumField() {
				builder.WriteString(fmt.Sprintf("%s\t", v.Field(i).String()))
			}
		}
		builder.WriteString("\n")
		fmt.Fprintf(w, builder.String())
		builder.Reset()
	}
	w.Flush()

	return nil
}

func PrettyPrintStringSlice(stringSlice []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "#\tName\t")

	// print out string slice
	for i, str := range stringSlice {
		fmt.Fprintf(w, fmt.Sprintf("%d\t%s\t\n", i+1, str))
	}
	w.Flush()

	return nil
}
