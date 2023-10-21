package simutil

type object interface {
	ObjectID() string
}

// unique returns a new slice containing only unique objects by their ID.
// We iterate the list backwards, so the last added duplicates win out.
// This does not merge objects.
func unique[O object](in []O) []O {
	out := []O{}
	seen := map[string]bool{}
	for i := len(in) - 1; i >= 0; i-- {
		o := in[i]

		id := o.ObjectID()
		if seen[id] {
			continue
		}
		out = append(out, o)
		seen[id] = true
	}
	return out
}
