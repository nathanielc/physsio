package main

import (
	"fmt"
)

type Game struct {
	Teams [2]Team
}

type Team struct {
	Name   string
	Active int
	Beasts [3]Beast
}

type Beast struct {
	Name   string
	HP     float64
	Energy float64
	Moves  MoveSet
}

type MoveSet [3]Move

type Move struct {
	Name       string
	Damage     float64
	EnergyCost float64
	beast      *Beast
}

func main() {
	a := Team{
		Name: "a",
		Beasts: [3]Beast{
			{
				Name:   "f",
				HP:     100,
				Energy: 1000,
				Moves: MoveSet{
					{
						Name:       "stomp",
						Damage:     10,
						EnergyCost: 100,
					},
					{
						Name:       "punch",
						Damage:     20,
						EnergyCost: 200,
					},
					{
						Name:       "kick",
						Damage:     30,
						EnergyCost: 300,
					},
				},
			},
			{
				Name:   "t",
				HP:     500,
				Energy: 1000,
				Moves: MoveSet{
					{
						Name:       "tickle",
						Damage:     5,
						EnergyCost: 50,
					},
					{
						Name:       "stomp",
						Damage:     10,
						EnergyCost: 100,
					},
				},
			},
			{
				Name:   "o",
				HP:     200,
				Energy: 1000,
				Moves: MoveSet{
					{
						Name:       "stomp",
						Damage:     10,
						EnergyCost: 100,
					},
				},
			},
		},
	}
	b := Team{
		Name: "b",
		Beasts: [3]Beast{
			{
				Name:   "f",
				HP:     100,
				Energy: 1000,
				Moves: MoveSet{
					{
						Name:       "stomp",
						Damage:     10,
						EnergyCost: 100,
					},
					{
						Name:       "punch",
						Damage:     20,
						EnergyCost: 200,
					},
					{
						Name:       "kick",
						Damage:     30,
						EnergyCost: 300,
					},
				},
			},
			{
				Name: "t",
				HP:   500,
				Moves: MoveSet{
					{
						Name:       "tickle",
						Damage:     5,
						EnergyCost: 50,
					},
					{
						Name:       "stomp",
						Damage:     10,
						EnergyCost: 100,
					},
				},
			},
			{
				Name: "o",
				HP:   200,
				Moves: MoveSet{
					{
						Name:       "stomp",
						Damage:     10,
						EnergyCost: 100,
					},
				},
			},
		},
	}

	g := Game{Teams: [2]Team{a, b}}
	g.Play()
}

func (g *Game) init() {
	for i := range g.Teams {
		g.Teams[i].Active = -1
		for j := range g.Teams[i].Beasts {
			g.Teams[i].Beasts[j].init()
		}
	}
}
func (g *Game) Play() {
	g.init()
	var activeBeasts [2]*Beast
	var moves [2]Move
	for {
		for i := range g.Teams {
			b := g.Teams[i].SelectBeast()
			activeBeasts[i] = b
		}
		for i, b := range activeBeasts {
			move := b.SelectMove()
			moves[i] = move
		}

		g.Teams[0].Apply(moves[1])
		g.Teams[1].Apply(moves[0])

		for i := range g.Teams {
			g.Teams[i].Print()
		}
	}
}

func (t *Team) SelectBeast() *Beast {
	for t.Active < 0 || t.Active >= len(t.Beasts) || t.Beasts[t.Active].HP <= 0 {
		fmt.Println("Team", t.Name)
		for i, b := range t.Beasts {
			if b.HP > 0 {
				fmt.Printf("\t%d - %s\n", i, b.Name)
			}
		}
		fmt.Println("Select Beast:")
		fmt.Scanln(&t.Active)
	}
	return &t.Beasts[t.Active]
}

func (b *Beast) init() {
	for i := range b.Moves {
		b.Moves[i].beast = b
	}
}
func (b *Beast) SelectMove() Move {
	m := -1
	for m < 0 || m >= len(b.Moves) || b.Moves[m].Name == "" || b.Moves[m].EnergyCost > b.Energy {
		for i, m := range b.Moves {
			if m.Name != "" {
				fmt.Printf("\t%d - %-15s Damage: %f\n", i, m.Name, m.Damage)
			}
		}
		fmt.Println("Select Move:")
		fmt.Scanln(&m)
	}
	return b.Moves[m]
}

func (t *Team) Apply(m Move) {
	t.Beasts[t.Active].Apply(m)
}
func (b *Beast) Apply(m Move) {
	b.HP -= m.Damage
	m.beast.Energy -= m.EnergyCost
}

func (t *Team) Print() {
	t.Beasts[t.Active].Print()
}
func (b *Beast) Print() {
	fmt.Printf("%s HP:%f Energy: %f\n", b.Name, b.HP, b.Energy)
}
