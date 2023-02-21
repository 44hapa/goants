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

	if m.String() != `1 . b 
B . . 
. . a 
A 0 . 
` {
		t.Errorf("map put ants in wrong place, got `%s`", m)
	}
}

// Две еды, два муравья
func TestMyBot_FeelMyAntByFood_1(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	food1 := m.FromRowCol(2, 2)
	food2 := m.FromRowCol(3, 3)
	myAnt1 := m.FromRowCol(1, 1)
	myAnt2 := m.FromRowCol(4, 4)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddAnt(myAnt2, MY_ANT)
	m.AddFood(food1)
	m.AddFood(food2)

	var s State
	s.Map = m
	bot := NewBot(&s)

	bot.FeelMyAntByFood(food1, &s)
	bot.FeelMyAntByFood(food2, &s)

	if bot.MyAnts[myAnt1].Goal != food1 {
		t.Errorf("муравью1 назначена не ближняя еда, хотел цель `%d` а получил `%d`", bot.MyAnts[myAnt1].Goal, food1)
	}
	if bot.MyAnts[myAnt2].Goal != food2 {
		t.Errorf("муравью2 назначена не ближняя еда, хотел цель `%d` а получил `%d`", bot.MyAnts[myAnt2].Goal, food2)
	}

	////fmt.Println(s.Map.String())
}

// Две еды, один муравей
func TestMyBot_FeelMyAntByFood_2(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(2, 2)
	food2 := m.FromRowCol(3, 3)
	myAnt1 := m.FromRowCol(1, 0)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)
	m.AddFood(food2)

	bot.FeelMyAntByFood(food1, &s)
	bot.FeelMyAntByFood(food2, &s)

	if bot.MyAnts[myAnt1].Goal != food1 {
		t.Errorf("муравью1 назначена не ближняя еде, got `%d`", bot.MyAnts[myAnt1].Goal)
	}

	////fmt.Println(s.Map.String())
}

// Два муравья, одина еда
func TestMyBot_FeelMyAntByFood_3(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(2, 2)
	myAnt1 := m.FromRowCol(0, 1)
	myAnt2 := m.FromRowCol(3, 3)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddAnt(myAnt2, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)

	if bot.MyAnts[myAnt2].Goal != food1 {
		t.Errorf("муравью1 назначена не ближняя еде, got `%d`", bot.MyAnts[myAnt2].Goal)
	}

	////fmt.Println(s.Map.String())
}

// У одного муравья есть цель, значит у него должно быть направление
func TestMyBot_PossibleMoveDirections_1(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	s.Rows = 5
	s.Cols = 5
	bot := NewBot(&s)

	food1 := m.FromRowCol(2, 2)
	myAnt1 := m.FromRowCol(0, 1)
	myAnt2 := m.FromRowCol(3, 3)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddAnt(myAnt2, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)

	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)
	row2, cel2 := bot.PossibleMoveDirections(myAnt2, bot.MyAnts[myAnt2], &s)

	if bot.MyAnts[myAnt1].Goal != Location(0) {
		t.Errorf("у муравья1 не должно быть цели, got `%d`", bot.MyAnts[myAnt1].Goal)
	}
	if row1 != NoMovement {
		t.Errorf("у муравья1 не должно быть движения по вертикали, got `%d`", row1)
	}
	if cel1 != NoMovement {
		t.Errorf("у муравья1 не должно быть движения по горизонтали, got `%d`", cel1)
	}

	if bot.MyAnts[myAnt2].Goal != food1 {
		t.Errorf("у муравья2 должна быть цели - еда1, цель муравья `%d` , еда в ", food1)
	}
	if row2.String() != North.String() {
		t.Errorf("у муравья2 должно быть движения по вертикали на `%s`, got `%s`", North, row2)
	}
	if cel2.String() != West.String() {
		t.Errorf("у муравья2 должно быть движения по горизонтали на `%s`, got `%s`", West, cel2)
	}

	//fmt.Println(s.Map.String())
}

// Один муравей, две еды (не идет по вертикале 1)
func TestMyBot_PossibleMoveDirections_2(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	s.Rows = 5
	s.Cols = 5
	bot := NewBot(&s)

	food1 := m.FromRowCol(1, 1)
	food2 := m.FromRowCol(2, 2)
	myAnt1 := m.FromRowCol(3, 2)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)
	m.AddFood(food2)

	bot.FeelMyAntByFood(food1, &s)
	bot.FeelMyAntByFood(food2, &s)

	/*
		North Direction = iota // Север
		East                   // Восток
		South                  // Юг
		West                   // Запад

		NoMovement

	*/
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)

	if bot.MyAnts[myAnt1].Goal != food2 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food2)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if row1.String() != North.String() {
		t.Errorf("у муравья1 должно быть движения по вертикали на `%s`, got `%s`", North, row1)
	}
	if cel1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по горизонтали на `%s`, got `%s`", NoMovement, cel1)
	}

	fmt.Println(s.Map.StringNum())

}

// Один муравей, две еды (не идет по вертикале 2)
func TestMyBot_PossibleMoveDirections_3(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	s.Rows = 5
	s.Cols = 5
	bot := NewBot(&s)

	food1 := m.FromRowCol(3, 2)
	myAnt1 := m.FromRowCol(2, 2)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)

	/*
		North Direction = iota // Север
		East                   // Восток
		South                  // Юг
		West                   // Запад

		NoMovement

	*/
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)

	if bot.MyAnts[myAnt1].Goal != food1 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food1)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if row1.String() != South.String() {
		t.Errorf("у муравья1 должно быть движения по вертикали на `%s`, got `%s`", South, row1)
	}
	if cel1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по горизонтали на `%s`, got `%s`", NoMovement, cel1)
	}
	////fmt.Println(s.Map.String())
}

// Один муравей, две еды (не идет по горизонтали 1)
func TestMyBot_PossibleMoveDirections_4(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(3, 3)
	myAnt1 := m.FromRowCol(3, 1)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)

	if bot.MyAnts[myAnt1].Goal != food1 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food1)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if row1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по вертикали на `%s`, got `%s`", NoMovement, row1)
	}
	if cel1.String() != East.String() {
		t.Errorf("у муравья1 должно быть движения по горизонтали на `%s`, got `%s`", East, cel1)
	}
	////fmt.Println(s.Map.String())
}

// Один муравей, две еды (не идет по горизонтали 2)
func TestMyBot_PossibleMoveDirections_5(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(3, 1)
	myAnt1 := m.FromRowCol(3, 3)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)

	if bot.MyAnts[myAnt1].Goal != food1 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food1)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if row1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по вертикали на `%s`, got `%s`", NoMovement, row1)
	}
	if cel1.String() != West.String() {
		t.Errorf("у муравья1 должно быть движения по горизонтали на `%s`, got `%s`", West, cel1)
	}
	//fmt.Println(s.Map.String())
}

// Проход через экран на юг
func TestMyBot_PossibleMoveDirections_6(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(0, 2)
	myAnt1 := m.FromRowCol(4, 2)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)
	if bot.MyAnts[myAnt1].Goal != food1 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food1)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if row1.String() != South.String() {
		t.Errorf("у муравья1 должно быть движения по вертикали на `%s`, got `%s`", South, row1)
	}
	if cel1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по горизонтали на `%s`, got `%s`", NoMovement, cel1)
	}
	//fmt.Println(s.Map.StringNum())
}

// Проход через экран на север
func TestMyBot_PossibleMoveDirections_7(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(4, 2)
	myAnt1 := m.FromRowCol(0, 2)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)
	if bot.MyAnts[myAnt1].Goal != food1 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food1)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if row1.String() != North.String() {
		t.Errorf("у муравья1 должно быть движения по вертикали на `%s`, got `%s`", North, row1)
	}
	if cel1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по горизонтали на `%s`, got `%s`", NoMovement, cel1)
	}
	//fmt.Println(s.Map.StringNum())
}

// Проход через экран на запад
func TestMyBot_PossibleMoveDirections_8(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(2, 4)
	myAnt1 := m.FromRowCol(2, 0)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)
	if bot.MyAnts[myAnt1].Goal != food1 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food1)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if cel1.String() != West.String() {
		t.Errorf("у муравья1 должно быть движения по вертикали на `%s`, got `%s`", West, row1)
	}
	if row1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по горизонтали на `%s`, got `%s`", NoMovement, cel1)
	}
	//fmt.Println(s.Map.StringNum())
}

// Проход через экран на восток
func TestMyBot_PossibleMoveDirections_9(t *testing.T) {
	m := NewMap(5, 5)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(2, 0)
	myAnt1 := m.FromRowCol(2, 4)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddFood(food1)

	bot.FeelMyAntByFood(food1, &s)
	row1, cel1 := bot.PossibleMoveDirections(myAnt1, bot.MyAnts[myAnt1], &s)
	if bot.MyAnts[myAnt1].Goal != food1 {
		botRow, botCel := s.Map.FromLocation(bot.MyAnts[myAnt1].Goal)
		foodRow, foodCel := s.Map.FromLocation(food1)
		t.Errorf("у муравья1 должна быть цель row/cel:%d/%d , цель муравья row/cel:%d/%d", foodRow, foodCel, botRow, botCel)
	}
	if cel1.String() != East.String() {
		t.Errorf("у муравья1 должно быть движения по вертикали на `%s`, got `%s`", East, row1)
	}
	if row1.String() != NoMovement.String() {
		t.Errorf("у муравья1 не должно быть движения по горизонтали на `%s`, got `%s`", NoMovement, cel1)
	}
	fmt.Println(s.Map.StringNum())
}

func TestMyBot_DoTurnByFood_1(t *testing.T) {
	m := NewMap(10, 20)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(3, 1)
	food2 := m.FromRowCol(3, 3)
	food3 := m.FromRowCol(5, 6)
	food4 := m.FromRowCol(8, 8)
	food5 := m.FromRowCol(9, 9)
	myAnt1 := m.FromRowCol(3, 4)
	myAnt2 := m.FromRowCol(4, 4)
	myAnt3 := m.FromRowCol(5, 4)
	myAnt4 := m.FromRowCol(6, 4)
	myHill := m.FromRowCol(3, 8)
	m.AddAnt(myAnt1, MY_ANT)
	m.AddAnt(myAnt2, MY_ANT)
	m.AddAnt(myAnt3, MY_ANT)
	m.AddAnt(myAnt4, MY_ANT)
	m.AddHill(myHill, MY_HILL)
	m.AddFood(food1)
	m.AddFood(food2)
	m.AddFood(food3)
	m.AddFood(food4)
	m.AddFood(food5)

	//fmt.Println(s.Map.String())
	//bot.DoTurnByFood(&s)
	//bot.DoTurnByFood(&s)
	bot.DoTurn(&s)
	//fmt.Println(s.Map.String())

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> qwe 1")
	spew.Dump(s.Map.FromLocationMyAnts(*bot))
	spew.Dump(bot.MyAnts)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> qwe 2")
	fmt.Println(s.Map.StringNum())
	fmt.Println(s.Map.Food)

	////fmt.Println(s.Map.String())
	//bot.NextSteps(&s)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> bot 1")
	//spew.Dump(bot)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> bot 2")
	//
	//bot.DoTurnByFood(&s)
	////fmt.Println(s.Map.String())
	//bot.NextSteps(&s)

}

func TestMyBot_DoTurnByFood_2(t *testing.T) {
	m := NewMap(52, 82)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(4, 2)
	food2 := m.FromRowCol(10, 71)
	food3 := m.FromRowCol(14, 75)
	food4 := m.FromRowCol(49, 68)
	food5 := m.FromRowCol(51, 1)

	myAnt1 := m.FromRowCol(11, 78)
	myAnt2 := m.FromRowCol(0, 72)
	myAnt3 := m.FromRowCol(5, 77)
	myAnt4 := m.FromRowCol(7, 77)
	/*
	 (string) (len=42) "antLoc row/col: 0/72  goalLoc row/col: 4/2",
	 (string) (len=44) "antLoc row/col: 7/77  goalLoc row/col: 10/71",
	 (string) (len=45) "antLoc row/col: 11/78  goalLoc row/col: 14/75",
	 (string) (len=44) "antLoc row/col: 5/77  goalLoc row/col: 49/68"
	*/

	myHill := m.FromRowCol(5, 77)

	m.AddAnt(myAnt1, MY_ANT)
	m.AddAnt(myAnt2, MY_ANT)
	m.AddAnt(myAnt3, MY_ANT)
	m.AddAnt(myAnt4, MY_ANT)
	m.AddHill(myHill, MY_HILL)
	m.AddFood(food1)
	m.AddFood(food2)
	m.AddFood(food3)
	m.AddFood(food4)
	m.AddFood(food5)

	//fmt.Println(s.Map.String())
	//bot.DoTurnByFood(&s)
	//bot.DoTurnByFood(&s)
	bot.DoTurnByFood(&s)
	//fmt.Println(s.Map.String())

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> qwe 1")
	spew.Dump(s.Map.FromLocationMyAnts(*bot))
	spew.Dump(bot.MyAnts)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> qwe 2")
	fmt.Println(s.Map.StringNum())
	fmt.Println(s.Map.Food)

	////fmt.Println(s.Map.String())
	//bot.NextSteps(&s)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> bot 1")
	//spew.Dump(bot)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> bot 2")
	//
	//bot.DoTurnByFood(&s)
	////fmt.Println(s.Map.String())
	//bot.NextSteps(&s)

}

func TestMyBot_GetNearestAnt_1(t *testing.T) {
	m := NewMap(52, 82)
	m.Reset()
	var s State
	s.Map = m
	bot := NewBot(&s)

	food1 := m.FromRowCol(4, 2)
	myAnt2 := m.FromRowCol(0, 72)
	/*
	 (string) (len=42) "antLoc row/col: 0/72  goalLoc row/col: 4/2",
	 (string) (len=44) "antLoc row/col: 7/77  goalLoc row/col: 10/71",
	 (string) (len=45) "antLoc row/col: 11/78  goalLoc row/col: 14/75",
	 (string) (len=44) "antLoc row/col: 5/77  goalLoc row/col: 49/68"
	*/

	myHill := m.FromRowCol(5, 77)

	m.AddAnt(myAnt2, MY_ANT)
	m.AddHill(myHill, MY_HILL)
	m.AddFood(food1)

	//fmt.Println(s.Map.String())
	//bot.DoTurnByFood(&s)
	//bot.DoTurnByFood(&s)
	antLoc, stepToGoal := bot.GetNearestAnt(food1, &s)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetNearestAnt 1")
	spew.Dump(antLoc)
	spew.Dump(stepToGoal)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetNearestAnt 2")

	//fmt.Println(s.Map.String())

	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> qwe 1")
	//spew.Dump(s.Map.FromLocationMyAnts(*bot))
	//spew.Dump(bot.MyAnts)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> qwe 2")
	//fmt.Println(s.Map.StringNum())
	//fmt.Println(s.Map.Food)

	////fmt.Println(s.Map.String())
	//bot.NextSteps(&s)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> bot 1")
	//spew.Dump(bot)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> bot 2")
	//
	//bot.DoTurnByFood(&s)
	////fmt.Println(s.Map.String())
	//bot.NextSteps(&s)

}

/*
  (string) (len=47) "foodLoc [loc:330] - row/col: 4/2 - isBool: true",
  (string) (len=50) "foodLoc [loc:891] - row/col: 10/71 - isBool: false",
  (string) (len=51) "foodLoc [loc:1223] - row/col: 14/75 - isBool: false",
  (string) (len=51) "foodLoc [loc:4086] - row/col: 49/68 - isBool: false",
  (string) (len=49) "foodLoc [loc:4183] - row/col: 51/1 - isBool: true"
 }
}
([]interface {}) (len=2 cap=2) {
 (string) (len=28) "DoTurnByFood foodKeys1.3 :: ",
 ([]string) (len=4 cap=4) {
  (string) (len=33) "antLoc [loc:980] - row/col: 11/78",
  (string) (len=31) "antLoc [loc:72] - row/col: 0/72",
  (string) (len=32) "antLoc [loc:487] - row/col: 5/77",
  (string) (len=32) "antLoc [loc:651] - row/col: 7/77"
*/
