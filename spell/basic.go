package spell

import "github.com/Nightgunner5/gogame/entity"

type BasicSpell struct {
	CastTime      float64
	currentTime   float64
	Interruptable bool
	Caster_       entity.EntityID
	Target_       entity.EntityID
	Action        func(target, caster entity.Entity)
}

func (s *BasicSpell) Caster() entity.Entity {
	return entity.Get(s.Caster_)
}

func (s *BasicSpell) Target() entity.Entity {
	return entity.Get(s.Target_)
}

func (s *BasicSpell) Interrupt() bool {
	if !s.Interruptable || s.currentTime >= s.CastTime {
		return false
	}
	s.currentTime = s.CastTime
	return true
}

func (s *BasicSpell) TotalTime() float64 {
	return s.CastTime
}

func (s *BasicSpell) TimeLeft() float64 {
	return s.CastTime - s.currentTime
}

func (s *BasicSpell) Tick(Δtime float64) bool {
	if s.currentTime >= s.CastTime {
		return true
	}
	s.currentTime += Δtime
	if s.currentTime >= s.CastTime {
		s.currentTime = s.CastTime
		target, caster := s.Target(), s.Caster()
		if target == nil || caster == nil {
			return true
		}

		s.Action(target, caster)

		return true
	}
	return false
}

var _ Spell = new(BasicSpell)
