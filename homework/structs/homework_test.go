package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		for i := range name {
			person.name[i] = name[i]
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		higherBits := mana & (oneByteMask << fourBitsShift) >> fourBitsShift
		lowerBits := (mana & fourBitsMask) << fourBitsShift
		person.manaHealth[1] |= uint8(lowerBits)
		person.manaHealth[2] |= uint8(higherBits)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		higherBits := (health & (fourBitsMask << oneByteShift)) >> oneByteShift
		lowerBits := health & oneByteMask
		person.manaHealth[1] |= uint8(higherBits)
		person.manaHealth[0] |= uint8(lowerBits)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectExperience |= uint8(respect << fourBitsShift)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.strengthLevel |= uint8(strength << fourBitsShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respectExperience |= uint8(experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.strengthLevel |= uint8(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typePropertyMask |= hasHouse
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typePropertyMask |= hasGun
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.typePropertyMask |= hasFamily
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.typePropertyMask |= 0b0000 << personType
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	fourBitsMask  = 0xF
	oneByteMask   = 0xFF
	fourBitsShift = 4
	oneByteShift  = 8

	typeBuilder    = 1  // 0000 0001
	typeBlacksmith = 2  // 0000 0010
	typeWarrior    = 4  // 0000 0100
	hasHouse       = 8  // 0000 1000
	hasGun         = 16 // 0001 0000
	hasFamily      = 32 // 0010 0000
)

type GamePerson struct {
	x, y, z           int32    // 12 B, mw: [* * * * * * * *] [* * * * _ _ _ _]
	gold              uint32   // 4 B, mw: [* * * * * * * *] [* * * * * * * *]
	manaHealth        [3]uint8 // 3 B, mw: [* * * * * * * *] [* * * * * * * *] [* * * _ _ _ _ _]
	respectExperience uint8    // 1 B, mw: [* * * * * * * *] [* * * * * * * *] [* * * * _ _ _ _]
	strengthLevel     uint8    // 1 B, mw: [* * * * * * * *] [* * * * * * * *] [* * * * * _ _ _]
	name              [42]byte // 42 B, mw: [* * * * * * * *] x7 + [* * * * * * _ _]
	typePropertyMask  uint8    // 1 B, mw: [* * * * * * * *] x8
}

func NewGamePerson(options ...Option) GamePerson {
	gp := GamePerson{}
	for _, option := range options {
		option(&gp)
	}
	return gp
}

func (p *GamePerson) Name() string {
	return unsafe.String(&p.name[0], len(p.name))
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.manaHealth[2])<<fourBitsShift + int(p.manaHealth[1])>>fourBitsShift
}

func (p *GamePerson) Health() int {
	return int(p.manaHealth[1]&fourBitsMask)<<oneByteShift + int(p.manaHealth[0])
}

func (p *GamePerson) Respect() int {
	return int((p.respectExperience & (fourBitsMask << fourBitsShift)) >> fourBitsShift)
}

func (p *GamePerson) Strength() int {
	return int((p.strengthLevel & (fourBitsMask << fourBitsShift)) >> fourBitsShift)
}

func (p *GamePerson) Experience() int {
	return int(p.respectExperience & fourBitsMask)
}

func (p *GamePerson) Level() int {
	return int(p.strengthLevel & fourBitsMask)
}

func (p *GamePerson) HasHouse() bool {
	return p.typePropertyMask&hasHouse != 0
}

func (p *GamePerson) HasGun() bool {
	return p.typePropertyMask&hasGun != 0
}

func (p *GamePerson) HasFamily() bool {
	return p.typePropertyMask&hasFamily != 0
}

func (p *GamePerson) Type() int {
	if p.typePropertyMask&typeBuilder != 0 {
		return BuilderGamePersonType
	}
	if p.typePropertyMask&typeBlacksmith != 0 {
		return BlacksmithGamePersonType
	}
	if p.typePropertyMask&typeWarrior != 0 {
		return WarriorGamePersonType
	}
	return 0
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
