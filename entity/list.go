package entity

import (
	"log"
	"math/rand"
	"sync"
)

type EntityList interface {
	// Returns the entity with the given ID or nil if that entity is not
	// in this list.
	Get(EntityID) Entity
	// Attempt to add the given entity to the list. Always returns true
	// unless the entity is already in the list.
	Add(Entity) bool
	// Remove the entity with the given ID from the list. Returns the
	// removed entity or nil if the ID does not occur on this list.
	Remove(EntityID) Entity
	// Remove the given entity and all entities that reference it as
	// their parent entity, recursively.
	RemoveRecursive(EntityID)

	// Returns the number of entities on this list.
	Count() int
	// Loops through each element of this entity list in order of ID,
	// calling the given function for each of them in sequence. Use
	// this if you need the entities in order or if the given
	// function cannot be run in parallel with itself.
	Each(func(Entity))
	// Calls the given function for each entity in the list without
	// waiting for one call to finish before starting the next.
	// The function returns after all calls have completed.
	All(func(Entity))
}

type entityList []Entity

// Copied from http://golang.org/src/pkg/sort/search.go to avoid closure memory allocation
func (list entityList) search(id EntityID) (i int, found bool) {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	j := len(list)
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if list[i].ID() < id {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	found = i < len(list) && list[i].ID() == id
	return
}

func (list entityList) Get(id EntityID) Entity {
	i, found := list.search(id)
	if found {
		return list[i]
	}
	return nil
}

func (list *entityList) Add(entity Entity) bool {
	i, found := list.search(entity.ID())
	if found {
		return false
	}

	*list = append((*list)[:i], append(entityList{entity}, (*list)[i:]...)...)

	return true
}

func (list *entityList) Remove(id EntityID) Entity {
	i, found := list.search(id)

	var ret Entity
	if found {
		ret = (*list)[i]
		*list = append((*list)[:i], (*list)[i+1:]...)
	}
	return ret
}

func (list *entityList) RemoveRecursive(id EntityID) {
	parent := list.Remove(id)

	if parent == nil {
		return
	}

	var toRemove []EntityID
	for _, ent := range *list {
		if ent.Parent() == parent {
			toRemove = append(toRemove, ent.ID())
		}
	}

	for _, remove := range toRemove {
		list.RemoveRecursive(remove)
	}
}

func (list entityList) Count() int {
	return len(list)
}

func (list entityList) Each(f func(Entity)) {
	for i := range list {
		f(list[i])
	}
}

func (list entityList) All(f func(Entity)) {
	var wg sync.WaitGroup

	wg.Add(list.Count())

	list.Each(func(e Entity) {
		go func() {
			f(e)
			wg.Done()
		}()
	})

	wg.Wait()
}

type concurrentEntityList struct {
	l entityList
	m sync.RWMutex
}

func (c *concurrentEntityList) Get(id EntityID) Entity {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.l.Get(id)
}

func (c *concurrentEntityList) Add(entity Entity) bool {
	c.m.Lock()
	defer c.m.Unlock()

	return c.l.Add(entity)
}

func (c *concurrentEntityList) Remove(id EntityID) Entity {
	c.m.Lock()
	defer c.m.Unlock()

	return c.l.Remove(id)
}

func (c *concurrentEntityList) RemoveRecursive(id EntityID) {
	c.m.Lock()
	defer c.m.Unlock()

	c.l.RemoveRecursive(id)
}

func (c *concurrentEntityList) Count() int {
        c.m.RLock()
        defer c.m.RUnlock()

	return c.l.Count()
}

func (c *concurrentEntityList) Each(f func(Entity)) {
        c.m.RLock()
        defer c.m.RUnlock()

	c.l.Each(f)
}

func (c *concurrentEntityList) All(f func(Entity)) {
        c.m.RLock()

	var wg sync.WaitGroup

	wg.Add(c.l.Count())

	c.l.Each(func(e Entity) {
		go func() {
			f(e)
			wg.Done()
		}()
	})

	c.m.RUnlock()

	wg.Wait()
}

// Creates a new, empty EntityList. This list is not synchronized and should
// only ever be accessed from a single goroutine. Use of the All method is
// discouraged on this type of EntityList.
func NewEntityList(capacity int) EntityList {
	list := make(entityList, 0, capacity)
	return &list
}

// Creates a new, empty EntityList. This list is synchronized and may be
// accessed from multiple goroutines at once. This is the type of EntityList
// used for the global entity list.
func ConcurrentEntityList(capacity int) EntityList {
	list := make(entityList, 0, capacity)
	return &concurrentEntityList{l: list}
}

type readOnlyEntityList struct {
	EntityList
}

func (readOnlyEntityList) Add(Entity) bool {
	return false
}

func (readOnlyEntityList) Remove(EntityID) Entity {
	return nil
}

var globalEntityList EntityList = ConcurrentEntityList(1)

// Returns a read-only reference to the global entity list.
func GlobalEntityList() EntityList {
	return readOnlyEntityList{globalEntityList}
}

// Add an entity to the global entity list. Logs a warning if the entity
// was already spawned. This will also add any parent entities if they
// are not already added.
func Spawn(entity Entity) {
	if globalEntityList.Add(entity) {
		for parent := entity.Parent(); parent != nil; parent = parent.Parent() {
			if !globalEntityList.Add(parent) {
				break
			}
		}
	} else {
		log.Printf("SpawnEntity called twice on entity %v", entity)
	}
}

// Remove an entity from the global entity list, along with any entity with it as its parent, recursively.
func Despawn(entity Entity) {
	globalEntityList.RemoveRecursive(entity.ID())
}

// Returns the entity for a given ID. nil will be returned if no entity is
// spawned with the given ID.
func Get(id EntityID) Entity {
	return globalEntityList.Get(id)
}

// Calls the given function for each currently spawned entity concurrently.
// Returns when all the function calls are complete.
func ForAll(f func(Entity)) {
	globalEntityList.All(f)
}

// Calls the given function for each currently spawned Positioner Entity
// within [distance] of [target].
func ForAllNearby(target Positioner, distance float64, f func(Entity)) {
	sX, sY, sZ := target.Position()
	d2 := distance * distance

	globalEntityList.All(func(e Entity) {
		if p, ok := e.(Positioner); ok && target != p {
			x, y, z := p.Position()
			x, y, z = sX-x, sY-y, sZ-z
			x, y, z = x*x, y*y, z*z

			if x+y+z <= d2 {
				f(e)
			}
		}
	})
}

func ForOneNearby(target Positioner, distance float64, allowed func(Entity) bool, f func(Entity)) {
	sX, sY, sZ := target.Position()
	d2 := distance * distance

	var (
		ent   Entity
		count int
	)

	globalEntityList.Each(func(e Entity) {
		if p, ok := e.(Positioner); ok && target != p && allowed(e) {
			x, y, z := p.Position()
			x, y, z = sX-x, sY-y, sZ-z
			x, y, z = x*x, y*y, z*z

			if x+y+z <= d2 {
				count++
				if ent == nil || rand.Intn(count) == 0 {
					ent = e
				}
			}
		}
	})

	if ent != nil {
		f(ent)
	}
}
