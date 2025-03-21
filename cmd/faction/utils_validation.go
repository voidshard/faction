package main

import (
	"strings"

	"github.com/voidshard/faction/pkg/kind"
)

var (
	help       = map[string]string{} // help text for each object
	shortNames = map[string]string{} // allows short hand names because we're lazy
)

func init() {
	for _, k := range kind.Kinds() {
		help[k] = kind.Doc(k)
		shortNames[kind.ShortName(k)] = k
	}
}

func validKind(name string) string {
	name = strings.ToLower(name)
	longname, ok := shortNames[name]
	if ok {
		name = longname
	}
	if kind.IsValid(name) {
		return name
	}
	return ""
}
