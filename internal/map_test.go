package internal

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestMap(t *testing.T) {
	m := NewMap(4, 3)
	m.Reset()
	if m.String() != `. . . 
. . . 
. . . 
. . . 
` {
		t.Errorf("map is wrong size, got `%s`", m)
	}

	loc := m.FromRowCol(3, 2)

	row, col := m.FromLocation(loc)
	if row != 3 || col != 2 {
		t.Errorf("conversion broken, got (%v, %v), wanted (3, 2)", row, col)
	}

	loc2 := m.FromRowCol(3, -1)
	if loc2 != loc {
		t.Errorf("from xy broken, got (%v), wanted (%v)", loc2, loc)
	}

	n := m.FromRowCol(2, 2)
	s := m.FromRowCol(4, 2)
	e := m.FromRowCol(3, 3)
	w := m.FromRowCol(3, 1)

	if n != m.Move(loc, North) {
		t.Errorf("Move north is broken")
	}
	if s != m.Move(loc, South) {
		t.Errorf("Move south is broken")
	}
	if e != m.Move(loc, East) {
		t.Errorf("Move east is broken")
	}
	if w != m.Move(loc, West) {
		t.Errorf("Move west is broken")
	}

	m.AddAnt(n, MY_ANT)
	m.AddAnt(s, ANT_1)
	m.AddAnt(w, MY_HILL)
	m.AddAnt(e, MY_OCCUPIED_HILL)
	m.AddHill(m.FromRowCol(0, 0), HILL_1)
	m.AddHill(m.FromRowCol(1, 0), HILL_1)
	m.AddAnt(m.FromRowCol(1, 0), ANT_1)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 1")
	//fmt.Println(m.String())
	//fmt.Println(m.FromLocation(w))
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 2")

	if m.String() != `1 . b 
B . . 
. . a 
A 0 . 
` {
		t.Errorf("map put ants in wrong place, got `%s`", m)
	}
}

func TestMyBot_GetNextStep(t *testing.T) {
	m := NewMap(5, 5)
	food1 := m.FromRowCol(2, 2)
	food2 := m.FromRowCol(3, 3)
	myAnt := m.FromRowCol(0, 0)
	m.AddAnt(myAnt, MY_ANT)
	m.AddFood(food1)
	m.AddFood(food2)

	var s State
	s.Map = m
	bot := NewBot(&s)
	loc1, near1 := bot.GetNearestAnt(food1, &s)
	loc2, near2 := bot.GetNearestAnt(food2, &s)

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Near 1")
	spew.Dump(m.FromLocation(loc1))
	spew.Dump(loc1)
	spew.Dump(near1)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Near 2")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Near 3")
	spew.Dump(m.FromLocation(loc2))
	spew.Dump(loc2)
	spew.Dump(near2)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Near 4")

	bot.FeelMyAntByFood(food1, &s)
	bot.FeelMyAntByFood(food2, &s)

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Near 5")
	spew.Dump(bot)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> Near 6")
}
