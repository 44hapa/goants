package internal

import (
	"math"
	"math/rand"
	"sort"
)

type MyAnt struct {
	ID   Item     // Идентиффикатор муравья
	Goal Location // Координаты цели, к которой пойдем
	//GoalID     Item     // Идентиффикатор цели, к которой пойдем
	GoalType   Item // Тип цели, к которой идем
	StepToGoal int  // Колицество шагов до цели
}

type MyBot struct {
	MyAnts map[Location]MyAnt // Мапа всех муравьев. Ключи - текущая координата муравья
}

// NewBot creates a new instance of your bot
func NewBot(s *State) Bot {
	//ants := make([]Ant, len(s.Map.Ants) - 1)
	mb := &MyBot{
		//Ants: ants,
		//do any necessary initialization here
	}
	return mb
}

// DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) error {
	// Зачищаем мапу с муравьями
	mb.MyAnts = map[Location]MyAnt{}

	//mb.DoTurnByAnts(s)
	mb.DoTurnByFood(s)

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
func (mb *MyBot) DoTurnByFood(s *State) error {
	for loc := range s.Map.Food {

		// Заполним муравьев по еде
		mb.FeelMyAntByFood(loc, s)
	}
	return nil
}

func (mb *MyBot) FeelMyAntByFood(locFood Location, s *State) {
	antLoc, steps := mb.GetNearestAnt(locFood, s)

	myAnt := MyAnt{
		GoalType:   FOOD, // Пока по умолчанию идем за едой
		Goal:       locFood,
		StepToGoal: steps,
	}
	//TODO: !!!!!!
	if _, ok := mb.MyAnts[antLoc]; !ok {
		mb.MyAnts = map[Location]MyAnt{}
	}
	mb.MyAnts[antLoc] = myAnt
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

// GetNearestAnt Возвращает ближайшего муравья
func (mb *MyBot) GetNearestAnt(foodLocal Location, s *State) (Location, int) {

	near := s.Map.Rows
	newNear := s.Map.Rows
	result := Location(0)
	for loc, ant := range s.Map.Ants {
		if ant != MY_ANT {
			continue
		}
		// Если у мураша уже есть задача
		if filledAnt, ok := mb.MyAnts[loc]; ok {
			if filledAnt.Goal != Location(0) {
				continue
			}
		}
		newNear = mb.GetDistance(foodLocal, loc, s)
		if near > newNear {
			near = newNear
			result = loc
		}
	}
	// Пометим еду как назначеную
	s.Map.Food[result] = false
	return result, near
}

// GetNearestDark Возвращает ближайшую темноту
func (mb *MyBot) GetNearestDark(antLocal Location, s *State) (Location, int) {

	near := s.Map.Rows
	newNear := s.Map.Rows
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
func (mb *MyBot) GetNearestGoal(goalType Item, myAntLocal Location, s *State) (Location, int) {
	// Сортируем еду по ближайшей
	foodKeys := mb.SortByNear(myAntLocal, s)

	near := s.Map.Rows
	newNear := s.Map.Rows
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

func (mb *MyBot) GetDistance(loc1 Location, loc2 Location, s *State) int {
	row1, col1 := s.Map.FromLocation(loc1)
	row2, col2 := s.Map.FromLocation(loc2)
	row := math.Max(float64(row1), float64(row2)) - math.Min(float64(row1), float64(row2))
	col := math.Max(float64(col1), float64(col1)) - math.Min(float64(col2), float64(col2))
	res := math.Sqrt(math.Pow(row, 2) + math.Pow(col, 2))
	return int(res)
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
			row, cel = mb.PossibleMoveDirections(loc, myAnt.Goal, s)
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
		darkLoc, _ := mb.GetNearestDark(loc, s)
		row, cel = mb.PossibleMoveDirections(loc, darkLoc, s)
		if row != NoMovement {
			loc2 := s.Map.Move(loc, row)
			if s.Map.SafeDestination(loc2) {
				s.IssueOrderLoc(loc, row)
				continue
			}
		}

		continue
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

func (mb *MyBot) PossibleMoveDirections(antLoc Location, goalLoc Location, s *State) (Direction, Direction) {
	resRow := NoMovement
	resCol := NoMovement
	if goalLoc == 0 {
		return resRow, resCol
	}
	antRow, antCol := s.Map.FromLocation(antLoc)
	goalRow, goalCol := s.Map.FromLocation(goalLoc)

	if antRow > goalRow {
		resRow = North
	}
	if antRow < goalRow {
		resRow = South
	}
	if antCol > goalCol {
		resCol = West
	}
	if antCol < goalCol {
		resCol = East
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
	newNear := s.Map.Rows
	near := s.Map.Rows

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
