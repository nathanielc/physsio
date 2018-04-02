import gym
from gym import error, spaces, utils
from gym.utils import seeding
import numpy as np

numMoves  = 3
numBeasts = 5

actionCount = numMoves + numBeasts

activeCol = 0
healthCol = 1

moveSpeedIdx  = 0
moveDamageIdx = 1
moveUsesIdx   = 2

moveSize       = 3 # speed + damage + uses
beastHealthIdx = 0
beastSize      = 1 + numMoves*moveSize # health + moves

numPlayers      = 2
playerAStartRow = 0
playerBStartRow = numBeasts

stateR = numPlayers * numBeasts
stateC = beastSize + 1 # beast + active beast
stateShape = (stateR, stateC)
stateSize = stateR*stateC



class PhyssioEnv(gym.Env):
    metadata = {'render.modes': ['human']}
    _initial_state = np.array([
        # Player A
        # Active, Health, Move0 , Move1, Move2
        1, 100, 5, 10, 10, 5, 1, 100, 5, 2, 50,
        0, 101, 5, 10, 10, 5, 1, 100, 5, 2, 50,
        0, 102, 5, 10, 10, 5, 1, 100, 5, 2, 50,
        0, 103, 5, 10, 10, 5, 1, 100, 5, 2, 50,
        0, 104, 5, 10, 10, 5, 1, 100, 5, 2, 50,
        # Player B
        # Active, Health, Move0 , Move1, Move2
        0, 110, 6, 10, 10, 6, 1, 100, 6, 2, 50,
        0, 111, 6, 10, 10, 6, 1, 100, 6, 2, 50,
        1, 112, 6, 10, 10, 6, 1, 100, 6, 2, 50,
        0, 113, 6, 10, 10, 6, 1, 100, 6, 2, 50,
        0, 114, 6, 10, 10, 6, 1, 100, 6, 2, 50,
    ]).reshape(stateShape)
    StateSize = stateSize
    ActionSize = actionCount

    def __init__(self):
        self.action_space = spaces.Discrete(actionCount)
        #self.observation_space = spaces.Box(low=finfo.min, high=finfo.max, dtype=np.float32, shape=stateShape)

        self.reset()

    def step(self, actionA, actionB):
        self._steps += 1.0

        changeActiveBeast(actionA, self._playerA)
        changeActiveBeast(actionB, self._playerB)

        beastA, moveA = getMove(actionA, self._playerA)
        beastB, moveB = getMove(actionB, self._playerB)

        if moveA[0,moveSpeedIdx] > moveB[0,moveSpeedIdx]:
            applyMove(moveA, beastA, beastB)
            applyMove(moveB, beastB, beastA)
        else:
            applyMove(moveB, beastB, beastA)
            applyMove(moveA, beastA, beastB)

        game_over = False
        healthA = playerHealth(self._playerA)
        healthB = playerHealth(self._playerB)

        game_over = healthA <= 0 or healthB <= 0
        rewardA = healthA / (healthA+healthB)
        if healthB == 0:
            rewardA += 10
        rewardB = healthB / (healthA+healthB)
        if healthA == 0:
            rewardB += 10
        return np.copy(self._state).reshape((1,stateSize)), rewardA, game_over, {}


    def reset(self):
        self._steps = 0.0
        self._state = np.matrix(np.zeros(stateShape))
        np.copyto(self._state, self._initial_state)

        self._playerA = self._state[playerAStartRow:playerBStartRow,:]
        self._playerB = self._state[playerBStartRow:,:]

        return np.copy(self._state).reshape((1,stateSize))

    def render(self, mode='human'):
        print(self._state)

def changeActiveBeast(action, player):
    if action < numMoves:
        return
    for i in range(numBeasts):
        player[i,activeCol] = 0
    player[action-numMoves,activeCol] = 1


noopMove = np.zeros((1,moveSize))

def getMove(a, player):
    beast, changed = getActiveBeast(player)
    if changed or a >= numMoves:
        return beast, noopMove
    return beast, beast[0,1+moveSize*a: 1+moveSize*(a+1)]

def getActiveBeast(player):
    activeIdx = -1
    for i in range(numBeasts):
        if player[i,activeCol] > 0:
            activeIdx = i
            break
    beast = player[activeIdx, 1: beastSize+1]
    changed = False
    if beast[0,beastHealthIdx] <= 0:
        changed = True
        player[activeIdx, activeCol] =  0
        for i in range(numBeasts):
            if player[i, 1+beastHealthIdx] > 0:
                player[i, activeCol] = 1
                beast = player[i, 1:beastSize+1]
                return beast, changed
    return beast, changed

def applyMove(m, src, dst):
    health = src[0,beastHealthIdx]
    uses = m[0,moveUsesIdx] - 1
    if health > 0 and uses >= 0:
        m[0,moveUsesIdx] =  uses
        dst[0,beastHealthIdx] = max(0,dst[0,beastHealthIdx] - m[0,moveDamageIdx])

def playerHealth(player):
    health = 0
    for i in range(numBeasts):
        health += player[i, healthCol]
    return health
