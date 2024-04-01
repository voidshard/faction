package structs

import "fmt"

const (
	MaxTuple = 10000
	MinTuple = -10000
)

func (t *Tuple) ObjectID() string {
	return fmt.Sprintf("%s-%s", t.Subject, t.Object)
}
