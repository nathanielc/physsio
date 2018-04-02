package physsio

import (
	"gonum.org/v1/gonum/mat"
)

const (
	numMoves  = 3
	numBeasts = 5

	actionCount = numMoves + numBeasts

	activeCol = 0
	healthCol = 1

	moveSpeedIdx  = 0
	moveDamageIdx = 1
	moveUsesIdx   = 2

	moveSize       = 3 // speed + damage + uses
	beastHealthIdx = 0
	beastSize      = 1 + numMoves*moveSize // health + moves

	numPlayers      = 2
	playerAStartRow = 0
	playerBStartRow = numBeasts

	stateR = numPlayers * numBeasts
	stateC = beastSize + 1 // beast + active beast
)

type Env struct {
	state *State

	playerA,
	playerB *mat.Dense
}

type State struct {
	mat.Dense
}

func NewState(s []float64) *State {
	if len(s) != stateR*stateC {
		panic("invalid initial state")
	}
	return &State{Dense: *mat.NewDense(stateR, stateC, s)}
}

func (s *State) Equal(o *State) bool {
	self := s.RawMatrix().Data
	other := o.RawMatrix().Data

	for i := range self {
		if self[i] != other[i] {
			return false
		}
	}
	return true
}

func NewEnv(initialState *State) *Env {
	e := new(Env)

	e.state = &State{Dense: *mat.DenseCopyOf(initialState)}
	e.playerA = e.state.Slice(playerAStartRow, playerBStartRow, 0, stateC).(*mat.Dense)
	e.playerB = e.state.Slice(playerBStartRow, stateR, 0, stateC).(*mat.Dense)

	return e
}

func (e *Env) State() *State {
	return e.state
}
func (e *Env) Step(actionA, actionB int) (rewardA, rewardB float64, over bool) {
	//log.Printf("State:\n%v", mat.Formatted(e.state))

	changeActiveBeast(actionA, e.playerA)
	changeActiveBeast(actionB, e.playerB)

	beastA, moveA := getMove(actionA, e.playerA)
	beastB, moveB := getMove(actionB, e.playerB)

	if moveA.AtVec(moveSpeedIdx) > moveB.AtVec(moveSpeedIdx) {
		applyMove(moveA, beastA, beastB)
		applyMove(moveB, beastB, beastA)
	} else {
		applyMove(moveB, beastB, beastA)
		applyMove(moveA, beastA, beastB)
	}
	if hasPlayerLost(e.playerA) {
		rewardB = 1
		over = true
	}
	if hasPlayerLost(e.playerB) {
		rewardA = 1
		over = true
	}
	//log.Printf("State:\n%v", mat.Formatted(e.state))
	return
}

func changeActiveBeast(a int, player *mat.Dense) {
	if a < numMoves {
		return
	}
	for i := 0; i < numBeasts; i++ {
		player.Set(i, activeCol, 0)
	}
	player.Set(a-numMoves, activeCol, 1)
}

var noopMove = mat.NewVecDense(moveSize, nil)

func getMove(a int, player *mat.Dense) (beast, move *mat.VecDense) {
	beast, changed := getActiveBeast(player)
	if changed || a >= numMoves {
		return beast, noopMove
	}
	return beast, beast.SliceVec(1+moveSize*a, 1+moveSize*(a+1)).(*mat.VecDense)
}

func getActiveBeast(player *mat.Dense) (beast *mat.VecDense, changed bool) {
	activeIdx := -1
	for i := 0; i < numBeasts; i++ {
		if player.At(i, activeCol) > 0 {
			activeIdx = i
			break
		}
	}
	beast = player.Slice(activeIdx, activeIdx+1, 1, beastSize+1).(*mat.Dense).RowView(0).(*mat.VecDense)
	if beast.AtVec(beastHealthIdx) <= 0 {
		changed = true
		player.Set(activeIdx, activeCol, 0)
		for i := 0; i < numBeasts; i++ {
			if player.At(i, 1+beastHealthIdx) > 0 {
				player.Set(i, activeCol, 1)
				beast = player.Slice(i, i+1, 1, beastSize+1).(*mat.Dense).RowView(0).(*mat.VecDense)
				return
			}
		}
	}
	return
}

func applyMove(m *mat.VecDense, src, dst *mat.VecDense) {
	health := src.AtVec(beastHealthIdx)
	uses := m.AtVec(moveUsesIdx) - 1
	if health > 0 && uses >= 0 {
		m.SetVec(moveUsesIdx, uses)
		dst.SetVec(beastHealthIdx, dst.AtVec(beastHealthIdx)-m.AtVec(moveDamageIdx))
	}
}
func hasPlayerLost(player *mat.Dense) bool {
	for i := 0; i < numBeasts; i++ {
		if player.At(i, healthCol) > 0 {
			return false
		}
	}
	return true
}
