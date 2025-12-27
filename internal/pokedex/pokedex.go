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

func (p *Pokedex) Inspect(name string) *pokeapi.Pokemon {
	p.lock.Lock()
	defer p.lock.Unlock()

	poke, ok := p.dex[name]
	if ok {
		return poke
	} else {
		return nil
	}
}

func (p *Pokedex) GetAllPokemonNames() []string {
	pokemon := make([]string, 0)

	p.lock.Lock()
	defer p.lock.Unlock()

	for name := range p.dex {
		pokemon = append(pokemon, name)
	}

	return pokemon
}
