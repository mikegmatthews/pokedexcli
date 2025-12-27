package pokedex

import (
	"sync"

	"github.com/mikegmatthews/pokedexcli/internal/pokeapi"
)

type Pokedex struct {
	lock sync.Mutex
	dex  map[string]*pokeapi.Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		lock: sync.Mutex{},
		dex:  make(map[string]*pokeapi.Pokemon),
	}
}

func (p *Pokedex) Caught(pokemon *pokeapi.Pokemon) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.dex[pokemon.Name] = pokemon
}
