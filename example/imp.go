package main

import (
	"github.com/Nightgunner5/gogame/entity"
	"github.com/Nightgunner5/gogame/spell"
)

type (
	Imp interface {
		entity.Entity
		entity.Positioner
		entity.Healther
		entity.Thinker
		spell.Caster

		imp()
	}

	imp struct {
		entity.EntityID
		entity.Positioner
		entity.Healther
		spell.SpellCaster

		master Magician
	}
)

func NewImp(master Magician, x, y, z float64) Imp {
	const (
		maxHealth = 10
	)
	i := &imp{master: master}

	i.Positioner = entity.BasePosition(i, x, y, z)
	i.Healther = entity.BaseHealth(i, maxHealth)

	entity.Spawn(i)

	return i
}

func (i *imp) Parent() entity.Entity {
	return i.master
}

func (m *imp) Tag() string {
	return "imp"
}

func (i *imp) Think(delta float64) {
	const (
		maxCastDistance = 100
		spellCastTime   = 1
		spellDamage     = 5
	)

	if i.Health() <= 0 {
		entity.Despawn(i)
		return
	}

	if i.CasterThink(delta) {
		// currently casting spell
		return
	}

	entity.ForOneNearby(i, maxCastDistance, func(e entity.Entity) bool {
		if o, ok := e.(Magician); ok {
			return o != i.master
		}
		if o, ok := e.(Imp); ok {
			return o.Parent() != i.master
		}
		return false
	}, func(e entity.Entity) {
		i.Cast(spell.DamageSpell(spellCastTime, spellDamage, i, e, false))
	})
}

func (imp) imp() {}
