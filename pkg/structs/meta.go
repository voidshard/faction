package structs

func AllMetaKeys() []Meta {
	all := []Meta{}
	for _, r := range Meta_value {
		all = append(all, Meta(r))
	}
	return all
}
