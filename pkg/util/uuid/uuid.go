package uuid

// Inspired by https://github.com/komuw/yuyuid/blob/master/yuyuid.go

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

const (
	RESERVED_NCS       byte = 0x80 //Reserved for NCS compatibility
	RFC_4122           byte = 0x40 //Specified in RFC 4122
	RESERVED_MICROSOFT byte = 0x20 //Reserved for Microsoft compatibility
	RESERVED_FUTURE    byte = 0x00 // Reserved for future definition.
)

var (
	// checks if we match a UUID like 123e4567-e89b-12d3-a456-426655440000
	validUUID = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type UUID [16]byte

// Create a new ID string deterministically
func New(args ...interface{}) string {
	if len(args) == 0 {
		// random
		return newUUID(
			rng.Int(),
			rng.Int(),
			rng.Int(),
			rng.Int(),
			rng.Int(),
		).string()
	}
	return newUUID(args...).string()
}

// IsValidUUID returns if the given string represents a UUID
func IsValidUUID(in string) bool {
	if in == "" {
		return false
	}
	return validUUID.MatchString(in)
}

// String returns the string form of a UUID
func (u UUID) string() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func (u *UUID) setVariant(variant byte) {
	switch variant {
	case RESERVED_NCS:
		u[8] &= 0x7F
	case RFC_4122:
		u[8] &= 0x3F
		u[8] |= 0x80
	case RESERVED_MICROSOFT:
		u[8] &= 0x1F
		u[8] |= 0xC0
	case RESERVED_FUTURE:
		u[8] &= 0x1F
		u[8] |= 0xE0
	}
}

func (u *UUID) setVersion(version byte) {
	u[6] = (u[6] & 0x0F) | (version << 4)
}

func newUUID(args ...interface{}) UUID {
	var uuid UUID
	var version byte = 4

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprint(args...)))

	sum := hasher.Sum(nil)
	copy(uuid[:], sum[:len(uuid)])

	uuid.setVariant(RFC_4122)
	uuid.setVersion(version)
	return uuid
}
