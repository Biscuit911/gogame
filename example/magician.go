package main

import (
	"github.com/Nightgunner5/gogame/effect"
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
	"math/rand"
)

type (
	Magician interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.Resourcer
		entity.Thinker
		spell.Caster

		effect.EffectAdder
		Cast(spell spell.Spell)

		magician()
	}

	magician struct {
		entity.EntityID
		entity.Healther
		entity.Resourcer
		spell.SpellCaster

		effect.BasicEffectAdder

		x, y, z float64
	}
)

func NewMagician(x, y, z float64) Magician {
	const (
		maxHealth = 100
		maxMana   = 160
	)

	m := &magician{
		x: x,
		y: y,
		z: z,
	}

	m.Healther = entity.BaseHealth(m, maxHealth)
	m.Resourcer = entity.BaseResource(m, maxMana)

	entity.Spawn(m)

	return m
}

func (m *magician) Parent() entity.Entity {
	return entity.World
}

func (m *magician) Position() (x, y, z float64) {
	return m.x, m.y, m.z
}

func (m *magician) Think(delta float64) {
	const (
		manaPerSecond = 10
	)

	if m.Health() <= 0 {
		entity.Despawn(m)
		return
	}

	m.EffectThink(delta)

	if m.CasterThink(delta) {
		// currently casting spell
		return
	}

	m.UseResource(-delta * manaPerSecond)
}

func summonImp(target, caster entity.Entity) {
	m := caster.(Magician)
	x, y, z := m.Position()

	x += rand.Float64()*2 - 1
	y += rand.Float64()*2 - 1
	z += rand.Float64()*2 - 1

	NewImp(m, x, y, z)
}

func summonShield(target, caster entity.Entity) {
	m := caster.(Magician)

	m.AddEffect(effect.AbsorbDamage(20), 5)
}

func (magician) magician() {}
