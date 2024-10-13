package uuid

// Inspired by https://github.com/komuw/yuyuid/blob/master/yuyuid.go

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/martinlindhe/base36"
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

	// checks if we match a base64 encoded no-dashes UUID like
	// 1H2R2Q1G1L2U2U2R1L1G1J1C1G2P2R1H1L1E1F2R2U2T1G1H1L2P2P1K2S1D1H1F
	valudB36ID = regexp.MustCompile("^[0-9A-Z]{65}$")

	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type UUID [16]byte

// Create a new ID string deterministically
func NewID(args ...interface{}) UUID {
	if len(args) == 0 {
		// random
		return newUUID(
			rng.Int(),
			rng.Int(),
			rng.Int(),
			rng.Int(),
			rng.Int(),
		)
	}
	return newUUID(args...)
}

// IsValidUUID returns if the given string represents a UUID
func IsValidUUID(in string) bool {
	if in == "" {
		return false
	}
	return validUUID.MatchString(in)
}

// IsValidB36ID returns if the given string represents a base36 encoded UUID
// We expect: no dashes, all uppercase
func IsValidB36ID(in string) bool {
	if in == "" {
		return false
	}
	return valudB36ID.MatchString(in)
}

// String returns the string form of a UUID
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// StringB36 returns the base36 string form of a UUID
// This string is longer than the standard UUID string, but is all uppercase.
func (u UUID) StringB36() string {
	return base36.EncodeBytes(u[:])
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
