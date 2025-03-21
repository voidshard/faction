package kind

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/voidshard/faction/pkg/structs/api"
	v1 "github.com/voidshard/faction/pkg/structs/v1"
	"github.com/voidshard/faction/pkg/util/log"

	"gopkg.in/yaml.v3"
)

var manager = newKindManager()

type KindManager struct {
	kinds map[string]*kindBuilder
}

func newKindManager() *KindManager {
	return &KindManager{
		kinds: make(map[string]*kindBuilder),
	}
}

// New returns an instance of the object with the given input data.
//
// If the kind is not provided, we will try to infer it from the input data.
// This is less efficient, so generally it is recommended to give the kind
// if possible.
func New(kind string, in interface{}) (v1.Object, error) {
	if kind == "" {
		var err error
		data, ok := in.([]byte)
		if !ok {
			data, err = yaml.Marshal(in)
		}
		raw := &struct {
			Kind string `json:"_kind" yaml:"_kind"`
		}{}
		err = yaml.Unmarshal(data, raw)
		if err != nil {
			return nil, err
		}
		kind = raw.Kind
	}
	kb, ok := manager.kinds[kind]
	if !ok {
		return nil, fmt.Errorf("kind %s not registered", kind)
	}
	return kb.o.New(in)
}

func KindOf(in interface{}) string {
	desired := reflect.TypeOf(in).Elem()
	for k, kb := range manager.kinds {
		if reflect.TypeOf(kb.o) == desired {
			return k
		}
	}
	return ""
}

func IsGlobal(kind string) bool {
	kb, ok := manager.kinds[kind]
	if !ok {
		return false
	}
	return kb.is_global
}

func Kinds() []string {
	keys := []string{}
	for k := range manager.kinds {
		keys = append(keys, k)
	}
	return keys
}

func IsValid(kind string) bool {
	_, ok := manager.kinds[kind]
	return ok
}

func IsValidId(kind, id string) bool {
	kb, ok := manager.kinds[kind]
	if !ok {
		return false
	}
	if kb.allow_alphanumeric_ids {
		return isAlphanum(id)
	} else {
		return isUUID4(id)
	}
}

func Validate(kind string, in interface{}) error {
	validate := validator.New()
	validate.RegisterValidation("alphanum-or-empty", ValidateAlphanumOrNone)
	validate.RegisterValidation("uuid4-or-empty", ValidateUUID4OrNone)
	validate.RegisterValidation("alphanumsymbol", ValidateAlphanumSymbol)

	is_global := false
	allow_alphanumeric_ids := false

	kb, ok := manager.kinds[kind]
	if ok {
		is_global = kb.is_global
		allow_alphanumeric_ids = kb.allow_alphanumeric_ids
	}

	if is_global {
		// if not global, require this field to be set
		validate.RegisterValidation("alphanum-if-non-global", ValidateAlphanum)
	} else {
		// otherwise, ignore this validation
		validate.RegisterValidation("alphanum-if-non-global", ValidateNoOp)
	}

	if allow_alphanumeric_ids {
		// enforce that ID is alphanumeric
		validate.RegisterValidation("valid_id", ValidateAlphanumOrNone)
	} else {
		// enforce that ID is a UUID4
		validate.RegisterValidation("valid_id", ValidateUUID4OrNone)
	}

	// If it is an Object, we need to cast it so validate.Struct can
	// pick up everything
	t := reflect.TypeOf(in)
	i := reflect.TypeOf((*v1.Object)(nil)).Elem()
	if t.Implements(i) {
		obj, _ := in.(v1.Object)
		return validate.Struct(obj)
	}

	searchreq, ok := in.(*api.SearchRequest)
	if ok {
		err := validateSearchRequest(searchreq)
		if err != nil {
			return err
		}
	}

	return validate.Struct(in)
}

func validMatchValue(v interface{}) bool {
	switch v.(type) {
	case string, int, float64, bool:
		return true
	default:
		return false
	}
}

func validateSearchRequest(q *api.SearchRequest) error {
	// make sure we don't hit a null
	if q.Filter.All == nil {
		q.Filter.All = []v1.Match{}
	}
	if q.Filter.Any == nil {
		q.Filter.Any = []v1.Match{}
	}
	if q.Filter.Not == nil {
		q.Filter.Not = []v1.Match{}
	}
	if q.Score == nil {
		q.Score = []v1.Score{}
	}

	// validate the query
	for _, m := range q.Filter.All {
		if !validMatchValue(m.Value) {
			return fmt.Errorf("invalid value for match %s", m.Field)
		}
	}
	for _, m := range q.Filter.Any {
		if !validMatchValue(m.Value) {
			return fmt.Errorf("invalid value for match %s", m.Field)
		}
	}
	for _, m := range q.Filter.Not {
		if !validMatchValue(m.Value) {
			return fmt.Errorf("invalid value for match %s", m.Field)
		}
	}
	for _, s := range q.Score {
		if !validMatchValue(s.Value) {
			return fmt.Errorf("invalid value for match %s", s.Field)
		}
	}

	return nil
}

func ShortName(kind string) string {
	kb, ok := manager.kinds[kind]
	if !ok {
		return ""
	}
	return kb.short
}

func Doc(kind string) string {
	kb, ok := manager.kinds[kind]
	if !ok {
		return ""
	}
	return kb.doc
}

func Register(kb *kindBuilder) error {
	if kb.o.GetKind() == "" || kb.o.GetKind() == "event" {
		return fmt.Errorf("kind %s is reserved", kb.o.GetKind())
	}
	_, ok := manager.kinds[kb.o.GetKind()]
	if ok {
		return fmt.Errorf("kind %s already registered", kb.o.GetKind())
	}
	manager.kinds[kb.o.GetKind()] = kb
	return nil
}

func IsSearchable(kind string) bool {
	kb, ok := manager.kinds[kind]
	if !ok {
		return false
	}
	return kb.searchable
}

func init() {
	world := NewKind(&v1.World{Meta: v1.Meta{Kind: "world"}})
	world.AllowAlphanumericIds() // ie. the world id can be "narnia" rather than a UUID
	world.DisableSearch()        // ie. we don't index worlds into the search service
	world.SetIsGlobal()          // ie. operations don't need to specify a world
	world.Short("wo").Doc("The root world object, all other objects are within some world namespace")
	log.Debug().Err(Register(world)).Msg("Registered world kind")

	// TODO: unsure if required
	config := NewKind(&v1.Config{Meta: v1.Meta{Kind: "config"}})
	config.AllowAlphanumericIds()
	config.DisableSearch()
	config.Short("co").Doc("Configuration for the world")
	log.Debug().Err(Register(config)).Msg("Registered config kind")

	race := NewKind(&v1.Race{Meta: v1.Meta{Kind: "race"}})
	race.AllowAlphanumericIds()
	race.DisableSearch()
	race.Short("ra").Doc("Some race in the world")
	log.Debug().Err(Register(race)).Msg("Registered race kind")

	culture := NewKind(&v1.Culture{Meta: v1.Meta{Kind: "culture"}})
	culture.AllowAlphanumericIds()
	culture.DisableSearch()
	culture.Short("cu").Doc("Some culture in the world")
	log.Debug().Err(Register(culture)).Msg("Registered culture kind")

	actor := NewKind(&v1.Actor{Meta: v1.Meta{Kind: "actor"}})
	actor.Short("ac").Doc("An actor in the world")
	log.Debug().Err(Register(actor)).Msg("Registered actor kind")

	faction := NewKind(&v1.Faction{Meta: v1.Meta{Kind: "faction"}})
	faction.Short("fa").Doc("A faction in the world")
	log.Debug().Err(Register(faction)).Msg("Registered faction kind")
}
