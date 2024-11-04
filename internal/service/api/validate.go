package api

import (
	"fmt"
	"regexp"

	"github.com/voidshard/faction/pkg/structs"
	"github.com/voidshard/faction/pkg/util/uuid"
)

const (
	validQueueEx = "^[0-9a-zA-Z_]{120}$"
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

func validDeferChange(in *structs.DeferChangeRequest) error {
	if in == nil {
		return fmt.Errorf("%w nil", ErrInvalid)
	}
	key, ok := structs.Metakey_name[int32(in.Data.Key)]
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
	key, ok := structs.Metakey_name[int32(in.Data.Key)]
	if !ok {
		return fmt.Errorf("%w key %d", ErrInvalid, in.Data.Key)
	}
	if in.Data.World != "" && !uuid.IsValidUUID(in.Data.World) {
		return fmt.Errorf("%w world %s", ErrInvalid, in.Data.World)
	} else if in.Data.Area != "" && !uuid.IsValidUUID(in.Data.Area) {
		return fmt.Errorf("%w area %s", ErrInvalid, in.Data.Area)
	} else if in.Data.Id != "" && !uuid.IsValidUUID(in.Data.Id) {
		return fmt.Errorf("%w id %s", ErrInvalid, in.Data.Id)
	} else if in.Queue != "" && !valudQueueName.MatchString(in.Queue) {
		return fmt.Errorf("%w queue %s should match %s", ErrInvalid, in.Queue, validQueueEx)
	}
	return nil
}
