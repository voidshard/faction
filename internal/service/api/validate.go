package api

import (
	"fmt"
	"regexp"

	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/uuid"
)

const (
	// "At least 3 chars; a-z A-Z 0-9 and underscore. No leading or trailing underscores."
	validQueueEx = "^[0-9a-zA-Z]{1}[0-9a-zA-Z_]{0,118}[0-9a-zA-Z]{1}$"
)

var (
	validQueueRe = regexp.MustCompile(validQueueEx)
)

type hasIDs interface {
	GetIds() []string
}

func validIDs(in hasIDs) error {
	if in == nil {
		return fmt.Errorf("%w nil", ErrInvalid)
	}
	for _, id := range in.GetIds() {
		if !uuid.IsValidUUID(id) {
			return fmt.Errorf("%w id %s", ErrInvalid, id)
		}
	}
	return nil
}

func validObject(o structs.Object) error {
	if o == nil {
		return fmt.Errorf("%w nil", ErrInvalid)
	}
	if o.GetId() == "" && o.GetEtag() == "" { // creating
		// pass
	} else if uuid.IsValidUUID(o.GetId()) && uuid.IsValidUUID(o.GetEtag()) { // updating
		// pass
	} else {
		return fmt.Errorf("%w id %s etag %s, should be valid UUIDs or empty", ErrInvalid, o.GetId(), o.GetEtag())
	}
	// The World ID of a World is it's own ID, so we don't need to set that unless we're updating,
	// which we've already checked above.
	if o.Kind() != kindWorld && !uuid.IsValidUUID(o.GetWorld()) {
		return fmt.Errorf("%w require valid world set %s", ErrInvalid, o.GetWorld())
	}
	return nil
}

func validDeferChange(in *structs.DeferChangeRequest) error {
	if in == nil {
		return fmt.Errorf("%w nil", ErrInvalid)
	}
	_, ok := structs.Metakey_name[int32(in.Data.Key)]
	if !ok {
		return fmt.Errorf("%w key %d", ErrInvalid, in.Data.Key)
	}
	if in.Data.World == "" || in.Data.Id == "" || in.Data.Key == structs.Metakey_KeyNone {
		return fmt.Errorf("%w world %s id %s key %d", ErrInvalid, in.Data.World, in.Data.Id, in.Data.Key)
	} else if !uuid.IsValidUUID(in.Data.World) {
		return fmt.Errorf("%w world %s", ErrInvalid, in.Data.World)
	} else if !uuid.IsValidUUID(in.Data.Id) {
		return fmt.Errorf("%w id %s", ErrInvalid, in.Data.Id)
	} else if in.Data.Area != "" && !uuid.IsValidUUID(in.Data.Area) {
		return fmt.Errorf("%w area %s", ErrInvalid, in.Data.Area)
	} else if *in.ToTick <= 0 && *in.ByTick <= 0 {
		return fmt.Errorf("%w one of ToTick (%d) or ByTick (%d) must be set", ErrInvalid, *in.ToTick)
	} else if *in.ToTick > 0 && *in.ByTick > 0 {
		return fmt.Errorf("%w only one of ToTick (%d) or ByTick (%d) can be set", ErrInvalid, *in.ToTick, *in.ByTick)
	}
	return nil
}

func validOnChange(in *structs.OnChangeRequest) error {
	_, ok := structs.Metakey_name[int32(in.Data.Key)]
	if !ok {
		return fmt.Errorf("%w key %d", ErrInvalid, in.Data.Key)
	}
	if in.Data.World != "" && !uuid.IsValidUUID(in.Data.World) {
		return fmt.Errorf("%w world %s", ErrInvalid, in.Data.World)
	} else if in.Data.Area != "" && !uuid.IsValidUUID(in.Data.Area) {
		return fmt.Errorf("%w area %s", ErrInvalid, in.Data.Area)
	} else if in.Data.Id != "" && !uuid.IsValidUUID(in.Data.Id) {
		return fmt.Errorf("%w id %s", ErrInvalid, in.Data.Id)
	} else if in.Queue != "" && !validQueueRe.MatchString(in.Queue) {
		return fmt.Errorf("%w queue %s should match %s", ErrInvalid, in.Queue, validQueueEx)
	}
	return nil
}
