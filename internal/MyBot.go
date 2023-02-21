package internal

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"math"
	"math/rand"
	"sort"
)

type MyAnt struct {
	ID   Item     // Идентиффикатор муравья
	Goal Location // Координаты цели, к которой пойдем
	//GoalID     Item     // Идентиффикатор цели, к которой пойдем
	GoalType   Item    // Тип цели, к которой идем
	StepToGoal float64 // Колицество шагов до цели
}

type MyBot struct {
	MyAnts  map[Location]MyAnt // Мапа всех муравьев. Ключи - текущая координата муравья
	MaxRows int
	MaxCell int
}

// NewBot creates a new instance of your bot
func NewBot(s *State) *MyBot {
	//ants := make([]Ant, len(s.Map.Ants) - 1)
	mb := &MyBot{
		MaxCell: s.Map.Cols,
		MaxRows: s.Map.Rows,
		//Ants: ants,
		//do any necessary initialization here
	}
	return mb
}

func (mb *MyBot) ClearMyAnts(s *State) {
	s.LogToFile("MY ANTS CLEAR:", mb.MyAnts)
	// Зачищаем мапу с муравьями
	mb.MyAnts = map[Location]MyAnt{}
}

// DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) error {
	// Если еды больше чем мурашей, и еще не вся еда с целями -
	s.LogToFile("-------------------------------" + fmt.Sprint(s.Turn) + "---------------------------")
	emptyFood := mb.DoTurnByFood(s)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DoTurn рекурсия 0.1")
	/*

	 */
	notEmptyFood := len(s.Map.Food) - emptyFood // Всего еды занятой (должно быть равно кол-ву муравьев)
	spew.Dump(emptyFood)
	spew.Dump(len(s.Map.Ants))
	spew.Dump(len(mb.MyAnts))
	spew.Dump(len(s.Map.Ants) > len(mb.MyAnts))
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DoTurn рекурсия 0.2")

	if emptyFood > 0 && len(mb.MyAnts) > notEmptyFood {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DoTurn рекурсия 1")
		mb.DoTurn(s)
	} else {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DoTurn рекурсия 2")
	}

	//s.LogToFile("-------------------------------" + fmt.Sprint(s.Turn) + "---------------------------")
	s.LogToFile(s.Map.FromLocationMyAnts(*mb))
	s.LogToFile("ANTS:", s.Map.FromLocationMapAnts())
	s.LogToFile("FOOD", s.Map.FromLocationMapFoods())

	mb.NextSteps(s)

	//returning an error will halt the whole program!
	return nil
}

func (mb *MyBot) DoTurnByAnts(s *State) error {
	for loc, ant := range s.Map.Ants {
		// Нашли чужого муравья
		if ant != MY_ANT {
			continue
		}

		// Заполним муравьев
		mb.FeelMyAnt(loc, ant, s)
	}
	return nil
}
func (mb *MyBot) DoTurnByFood(s *State) int {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DoTurnByFood 1")
	spew.Dump(s.Map.Food)
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DoTurnByFood 2")

	foodKeys := s.Map.SortFoodByLocation()
	s.LogToFile("DoTurnByFood foodKeys1.1 :: " + fmt.Sprint(foodKeys))
	s.LogToFile("DoTurnByFood foodKeys1.2 :: ", s.Map.FromLocationMapFoods())
	s.LogToFile("DoTurnByFood foodKeys1.3 :: ", s.Map.FromLocationMapAnts())
	for _, loc := range foodKeys {
		// Заполним муравьев по еде // Возможно имеет смысл отсортировать муравьев
		mb.FeelMyAntByFood(loc, s)
	}

	s.LogToFile("DoTurnByFood foodKeys2 :: " + fmt.Sprint(foodKeys))
	foodKeys = s.Map.SortFoodByLocation()
	s.LogToFile("DoTurnByFood foodKeys3 :: " + fmt.Sprint(foodKeys))
	return len(foodKeys)
}

func (mb *MyBot) FeelMyAntByFood(locFood Location, s *State) {
	// Ближайший доступный муравей
	antLoc, steps := mb.GetNearestAnt(locFood, s)

	// Если не нашли мураша
	if antLoc == Location(0) {
		return
	}

	// Нет муравья
	if steps == float64(0) {
		return
	}

	if mb.MyAnts == nil {
		mb.MyAnts = map[Location]MyAnt{}
	}

	if _, ok := mb.MyAnts[antLoc]; !ok {
		mb.MyAnts[antLoc] = MyAnt{}
	}

	if mb.MyAnts[antLoc].StepToGoal == 0 || mb.MyAnts[antLoc].StepToGoal > steps {
		//if mb.MyAnts[antLoc].StepToGoal == 0 {
		// Если цель ближе
		mb.MyAnts[antLoc] = MyAnt{
			GoalType:   FOOD, // Пока по умолчанию идем за едой
			Goal:       locFood,
			StepToGoal: steps,
		}
	}

}

func (mb *MyBot) FeelMyAnt(loc Location, id Item, s *State) {
	goal, steps := mb.GetNearestGoal(FOOD, loc, s)
	myAnt := MyAnt{
		ID:         id,
		GoalType:   FOOD, // Пока по умолчанию идем за едой
		Goal:       goal,
		StepToGoal: steps,
	}
	mb.MyAnts[loc] = myAnt
}

// GetNearestAnt Возвращает ближайшего муравья (возвращает - локация мураша + дистанция)
func (mb *MyBot) GetNearestAnt(foodLocal Location, s *State) (Location, float64) {
	maxDist := math.Max(float64(s.Map.Rows), float64(s.Map.Cols))
	near := maxDist
	newNear := maxDist
	result := Location(0)

	antKeys := s.Map.SortAntsByLocation()

	for _, antLocation := range antKeys {
		if s.Map.Ants[antLocation] != MY_ANT {
			continue
		}

		newNear = mb.GetDistance(foodLocal, antLocation, s)

		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetNearestAnt 1")
		fRow, fCol := s.Map.FromLocation(foodLocal)
		aRow, aCol := s.Map.FromLocation(antLocation)
		fmt.Println("Координаты еды:", foodLocal, " Row/Col:", fRow, fCol)
		fmt.Println("Проверяемый муравей:", antLocation, " Row/Col:", aRow, aCol)
		if _, ok := mb.MyAnts[antLocation]; ok {
			fmt.Println("Цель у муравья и расстояние:", mb.MyAnts[antLocation].Goal, mb.MyAnts[antLocation].StepToGoal)
		} else {
			fmt.Println("Цель у муравья и расстояние:", 0, 0)
		}
		fmt.Println("Высчитываемая дистанция:", newNear)
		fmt.Println("Лучшая дистанция:", near)
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetNearestAnt 2")

		// Если муравей уже в списке
		if filledAnt, ok := mb.MyAnts[antLocation]; ok {
			// Если у мураша уже есть задача и ему до цели ближе, чем найденная - оставим его
			if filledAnt.Goal != Location(0) {
				if filledAnt.Goal == foodLocal {
					s.Map.Food[filledAnt.Goal] = false
				}
				if filledAnt.StepToGoal <= newNear {
					continue
				} else {
					// Обозначим еду как незахваченную в прицел муравьями
					s.Map.Food[filledAnt.Goal] = true
					fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> filledAnt 1")
					fmt.Println("filledAnt.StepToGoal > newNear : ", filledAnt.StepToGoal, newNear)
					spew.Dump("Переназначение цели у мураша: ", filledAnt)
					fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> filledAnt 2")

				}
			}
		}

		if near > newNear {
			near = newNear
			result = antLocation
		} else {
			//s.LogToFile("NEAR : newNear===" + fmt.Sprint(newNear) + " near====" + fmt.Sprint(near) + "  maxDist===" + fmt.Sprint(maxDist))
		}
	}

	if result != Location(0) {
		// Пометим еду как назначенную
		s.Map.Food[foodLocal] = false
	}

	if near == newNear {
		s.Map.Food[foodLocal] = false
	}

	if near == maxDist {
		near = 0
	}
	return result, near
}

// GetNearestDark Возвращает ближайшую темноту
func (mb *MyBot) GetNearestDark(antLocal Location, s *State) (Location, float64) {

	near := float64(s.Map.Rows)
	newNear := float64(s.Map.Rows)
	result := Location(0)
	for loc, item := range s.Map.itemGrid {
		if item != UNKNOWN {
			continue
		}
		newNear = mb.GetDistance(antLocal, Location(loc), s)
		if near > newNear {
			near = newNear
			result = Location(loc)
		}
	}

	return result, near
}

// GetNearestGoal Возвращает ближайшую цель
func (mb *MyBot) GetNearestGoal(goalType Item, myAntLocal Location, s *State) (Location, float64) {
	// Сортируем еду по ближайшей
	foodKeys := mb.SortByNear(myAntLocal, s)

	near := float64(s.Map.Rows)
	newNear := float64(s.Map.Rows)
	result := Location(0)
	if goalType == FOOD {
		for _, loc := range foodKeys {
			// Если еда уже назначена - пропускаем ее
			//if s.Map.Food[loc] == false {
			//	continue
			//}
			newNear = mb.GetDistance(myAntLocal, loc, s)
			if near > newNear {
				near = newNear
				result = loc
			}
		}
	}

	// Пометим еду как назначеную
	s.Map.Food[result] = false
	return result, near
}

// Неверное определение дистанции!!!!!
func (mb *MyBot) GetDistance(loc1 Location, loc2 Location, s *State) float64 {
	row1, col1 := s.Map.FromLocation(loc1)
	row2, col2 := s.Map.FromLocation(loc2)
	rowForward := math.Abs(float64(row1) - float64(row2))
	colForward := math.Abs(float64(col1) - float64(col2))
	rowBack := float64(s.Map.Rows) - math.Max(float64(row1), float64(row2)) + math.Min(float64(row1), float64(row2))
	colBack := float64(s.Map.Cols) - math.Max(float64(col1), float64(col2)) + math.Min(float64(col1), float64(col2))

	row := math.Min(rowForward, rowBack)
	col := math.Min(colForward, colBack)
	//row := math.Max(float64(row1), float64(row2)) - math.Min(float64(row1), float64(row2))
	//col := math.Max(float64(col1), float64(col1)) - math.Min(float64(col2), float64(col2))
	res := math.Sqrt(math.Pow(row, 2) + math.Pow(col, 2))

	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetDistance 0")
	//spew.Dump(s.Map.Cols - col1 - col2)
	//
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetDistance 1")
	//spew.Dump(row1)
	//spew.Dump(col1)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetDistance 2")
	//spew.Dump(row2)
	//spew.Dump(col2)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetDistance 3")
	//spew.Dump(row)
	//spew.Dump(col)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetDistance 4")
	//spew.Dump(res)
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> GetDistance 5")

	return res
}

func (mb *MyBot) NextSteps(s *State) {
	row := NoMovement
	cel := NoMovement
	//for _, loc := range mb.SortBySteps() {
	for loc, ant := range s.Map.Ants {
		// Нашли чужого муравья
		if ant != MY_ANT {
			continue
		}

		if myAnt, ok := mb.MyAnts[loc]; ok {
			row, cel = mb.PossibleMoveDirections(loc, myAnt, s)
			if row != NoMovement {
				loc2 := s.Map.Move(loc, row)
				if s.Map.SafeDestination(loc2) {
					s.IssueOrderLoc(loc, row)
					continue
				}
			}
			if cel != NoMovement {
				loc2 := s.Map.Move(loc, cel)
				if s.Map.SafeDestination(loc2) {
					s.IssueOrderLoc(loc, cel)
					continue
				}
			}
		}
		// По умолчанию идем в темноту
		//darkLoc, _ := mb.GetNearestDark(loc, s)
		//row, cel = mb.PossibleMoveDirections(loc, MyAnt{Goal: darkLoc}, s)
		//if row != NoMovement {
		//	loc2 := s.Map.Move(loc, row)
		//	if s.Map.SafeDestination(loc2) {
		//		s.IssueOrderLoc(loc, row)
		//		continue
		//	}
		//}

		//continue
		dirs := []Direction{North, East, South, West}
		p := rand.Perm(4)
		for _, i := range p {
			d := dirs[i]

			loc2 := s.Map.Move(loc, d)
			if s.Map.SafeDestination(loc2) {
				s.IssueOrderLoc(loc, d)
				//there's also an s.IssueOrderRowCol if you don't have a Location handy
				break
			}
		}
	}

	return
}

func (mb *MyBot) PossibleMoveDirections(antLoc Location, myAnt MyAnt, s *State) (Direction, Direction) {
	goalLoc := myAnt.Goal
	resRow := NoMovement
	resCol := NoMovement
	if goalLoc == Location(0) {
		return resRow, resCol
	}
	antRow, antCol := s.Map.FromLocation(antLoc)
	goalRow, goalCol := s.Map.FromLocation(goalLoc)

	if antRow > goalRow {
		if float64(antRow-goalRow) < float64(mb.MaxRows/2) {
			resRow = North
		} else {
			resRow = South
		}
	}
	if antRow < goalRow {
		if float64(goalRow-antRow) < float64(mb.MaxRows)/2 {
			resRow = South
		} else {
			resRow = North
		}
	}
	if antCol > goalCol {
		if float64(antCol-goalCol) < float64(mb.MaxCell)/2 {
			resCol = West
		} else {
			resCol = East
		}
	}
	if antCol < goalCol {
		if float64(goalCol-antCol) < float64(mb.MaxCell)/2 {
			resCol = East
		} else {
			resCol = West
		}

	}
	return resRow, resCol
}

func (mb *MyBot) SortBySteps() []Location {
	keys := make([]Location, 0, len(mb.MyAnts))

	for key := range mb.MyAnts {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return mb.MyAnts[keys[i]].StepToGoal < mb.MyAnts[keys[j]].StepToGoal
	})

	return keys
}

func (mb *MyBot) SortByNear(myAntLocal Location, s *State) []Location {
	keys := make([]Location, 0, len(s.Map.Food))
	newNear := float64(s.Map.Rows)
	near := float64(s.Map.Rows)

	for key := range s.Map.Food {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		newNear = mb.GetDistance(myAntLocal, keys[i], s)

		if near > newNear {
			near = newNear
		}

		return mb.GetDistance(myAntLocal, keys[i], s) < mb.GetDistance(myAntLocal, keys[j], s)
	})

	return keys
}

/*
#!/usr/bin/env sh
./playgame.py --player_seed 42 --end_wait=0.25 --verbose --log_dir game_logs --turns 1000 --map_file maps/random_walk/random_walk_p02_08.map "$@" \
	"php sample_bots/php/MyBot.php" \
	"/Users/pavel.sukhorukov/Downloads/tools/sample_bots/go/ants1"

*/
