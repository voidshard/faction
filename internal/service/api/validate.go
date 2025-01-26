package api

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/uuid"
)

const (
	// "1-255 chars: a-z, A-Z, 0-9 and _"
	validQueueEx = "^[0-9a-zA-Z_]{1,255}$"

	// "1-255 chars: a-z, A-Z, 0-9, space - _ ."
	validUniqueName = "^[0-9a-zA-Z -_.]{1,255}$"
)

var (
	validQueueRe  = regexp.MustCompile(validQueueEx)
	validUniqueRe = regexp.MustCompile(validUniqueName)

	relaxedIdValidation = map[string]bool{
		kindWorld:   true,
		kindRace:    true,
		kindCulture: true,
	}
)

type hasIDs interface {
	GetIds() []string
}

func validID(kind, id string) error {
	lax, _ := relaxedIdValidation[kind]
	if lax {
		if !validUniqueRe.MatchString(id) {
			return fmt.Errorf("%w id %s", ErrInvalid, id)
		}
	} else {
		if !uuid.IsValidUUID(id) {
			return fmt.Errorf("%w id %s", ErrInvalid, id)
		}
	}
	return nil
}

func validIDs(kind string, in hasIDs) error {
	if in == nil {
		return fmt.Errorf("%w nil", ErrInvalid)
	}
	for _, id := range in.GetIds() {
		err := validID(kind, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func validObject(o structs.Object) error {
	if o == nil {
		return fmt.Errorf("%w object was nil", ErrInvalid)
	}

	// ID format is either a unique name or a UUID
	// If the ID is a UUID and is empty, we will set it to a random UUID
	lax, _ := relaxedIdValidation[o.Kind()]
	if o.GetId() == "" {
		if lax {
			// the user should set a name
			return fmt.Errorf("%w id must be set (%s)", ErrInvalid, validUniqueName)
		} else if o.GetEtag() == "" {
			// we will set a random UUID
			o.SetId(uuid.NewID().String()) // set a valid ID
		} else {
			return fmt.Errorf("%w id must be set in order to update", ErrInvalid)
		}
	} else {
		if lax {
			if !validUniqueRe.MatchString(o.GetId()) {
				return fmt.Errorf("%w id %s want (%s)", ErrInvalid, o.GetId(), validUniqueName)
			}
		} else {
			if !uuid.IsValidUUID(o.GetId()) {
				return fmt.Errorf("%w id %s want uuid4", ErrInvalid, o.GetId())
			}
		}
	}

	// Etag should be empty or a valid UUID
	if o.GetEtag() != "" && !uuid.IsValidUUID(o.GetEtag()) {
		return fmt.Errorf("%w etag %s want uuid4", ErrInvalid, o.GetEtag())
	}

	// for everything that isn't a world, we need a world set
	if o.Kind() != kindWorld && !validUniqueRe.MatchString(o.GetWorld()) {
		return fmt.Errorf("%w require valid world set %s", ErrInvalid, o.GetWorld())
	}
	return nil
}

func validDeferChange(in *structs.DeferChangeRequest) error {
	if !structs.IsValidAPIKind(in.Data.Key) {
		return fmt.Errorf("%w key %s", ErrInvalid, in.Data.Key)
	}
	if in.Data.World == "" || in.Data.Id == "" || in.Data.Key == "" {
		return fmt.Errorf("%w world %s id %s key %d", ErrInvalid, in.Data.World, in.Data.Id, in.Data.Key)
	} else if !uuid.IsValidUUID(in.Data.World) {
		return fmt.Errorf("%w world %s", ErrInvalid, in.Data.World)
	} else if !uuid.IsValidUUID(in.Data.Id) {
		return fmt.Errorf("%w id %s", ErrInvalid, in.Data.Id)
	} else if in.Data.Type != "" && !validUniqueRe.MatchString(in.Data.Type) {
		return fmt.Errorf("%w type %s", ErrInvalid, in.Data.Type)
	} else if *in.ToTick <= 0 && *in.ByTick <= 0 {
		return fmt.Errorf("%w one of ToTick (%d) or ByTick (%d) must be set", ErrInvalid, *in.ToTick)
	} else if *in.ToTick > 0 && *in.ByTick > 0 {
		return fmt.Errorf("%w only one of ToTick (%d) or ByTick (%d) can be set", ErrInvalid, *in.ToTick, *in.ByTick)
	}
	return nil
}

func validOnChange(in *structs.OnChangeRequest) error {
	if in.Data.Key != "" && !structs.IsValidAPIKind(in.Data.Key) {
		return fmt.Errorf("%w key %s", ErrInvalid, in.Data.Key)
	}
	if in.Data.World != "" && !uuid.IsValidUUID(in.Data.World) {
		return fmt.Errorf("%w world %s", ErrInvalid, in.Data.World)
	} else if in.Data.Type != "" && !validUniqueRe.MatchString(in.Data.Type) {
		return fmt.Errorf("%w type %s", ErrInvalid, in.Data.Type)
	} else if in.Data.Id != "" && !uuid.IsValidUUID(in.Data.Id) {
		return fmt.Errorf("%w id %s", ErrInvalid, in.Data.Id)
	} else if in.Queue != "" && !validQueueRe.MatchString(in.Queue) {
		return fmt.Errorf("%w queue %s should match %s", ErrInvalid, in.Queue, validQueueEx)
	} else if strings.HasPrefix(in.Queue, "internal.") {
		return fmt.Errorf("%w queue %s may not start with internal.", ErrInvalid, in.Queue)
	}
	return nil
}
