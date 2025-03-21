package kind

import v1 "github.com/voidshard/faction/pkg/structs/v1"

type kindBuilder struct {
	o     v1.Object
	short string
	doc   string

	allow_alphanumeric_ids bool
	is_global              bool
	searchable             bool
}

func NewKind(obj v1.Object) *kindBuilder {
	return &kindBuilder{o: obj, searchable: true}
}

func (kb *kindBuilder) DisableSearch() *kindBuilder {
	kb.searchable = false
	return kb
}

func (kb *kindBuilder) SetIsGlobal() *kindBuilder {
	kb.is_global = true
	return kb
}

func (kb *kindBuilder) AllowAlphanumericIds() *kindBuilder {
	kb.allow_alphanumeric_ids = true
	return kb
}

func (kb *kindBuilder) Short(short string) *kindBuilder {
	kb.short = short
	return kb
}

func (kb *kindBuilder) Doc(doc string) *kindBuilder {
	kb.doc = doc
	return kb
}
