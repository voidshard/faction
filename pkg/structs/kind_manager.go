package structs

var kinds = newKindManager()

type KindManager struct {
	api map[string]Object
}

func (k *KindManager) Register(o Object) {
	k.api[o.Kind()] = o
}

func (k *KindManager) ValidAPIKind(kind string) bool {
	_, ok := k.api[kind]
	return ok
}

func newKindManager() *KindManager {
	return &KindManager{
		api: make(map[string]Object),
	}
}

func IsValidAPIKind(kind string) bool {
	return kinds.ValidAPIKind(kind)
}

func init() {
	kinds.Register(&World{})
	kinds.Register(&Race{})
	kinds.Register(&Actor{})
	kinds.Register(&Culture{})
	kinds.Register(&Faction{})
	kinds.Register(&Job{})
}
