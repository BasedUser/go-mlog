package runtime

import (
	"strings"
)

type ExecutionContext struct {
	Variables   map[string]*Variable
	PrintBuffer strings.Builder
	DrawBuffer  []DrawStatement
	Objects     map[string]interface{}
	Metrics     map[int64]*Metrics
}

type Metrics struct {
	Executions uint64
}

type Variable struct {
	Value    interface{}
	Constant bool
}

type MLOGLine struct {
	Instruction []string
	Comment     string
	SourceLine  int
}

type Operation struct {
	Line     MLOGLine
	Executor OperationExecutor
}

type OperationExecutor func(ctx *ExecutionContext)

type OperationSetup func(args []string) (OperationExecutor, error)

type Message interface {
	PrintFlush(buffer string)
}

type UnitController struct {
	parent interface{}
	name   string
	ctype  int
}

type UnitType struct {
	/** If true, the unit is always at elevation 1. */
	name                  string
	flying                bool
	targetAir             bool
	targetGround          bool
	faceTarget            bool
	rotateShooting        bool
	isCounted             bool
	lowAltitude           bool
	circleTarget          bool
	canBoost              bool
	destructibleWreck     bool
	hovering              bool
	omniMovement          bool
	allowLegStep          bool
	canDrown              bool
	speed                 float64
	boostMultiplier       float64
	rotateSpeed           float64
	baseRotateSpeed       float64
	drag                  float64
	accel                 float64
	landShake             float64
	rippleScale           float64
	riseSpeed             float64
	fallSpeed             float64
	health                float64
	urange                float64
	miningRange           float64
	armor                 float64
	maxRange              float64
	crashDamageMultiplier float64
	aimDst                float64
	commandRadius         float64
	buildSpeed            float64
	mineSpeed             float64
	strafePenalty         float64
	hitSize               float64
	commandLimit          int
	payloadCapacity       int
	itemCapacity          int
	ammoCapacity          int
	clipSize              int
	ammoType              int
	mineTier              int
	//abilities = []interface{};
	weapons []interface{} //TODO make actual weapon type
}

type ItemStack struct {
	item   int
	amount int
}

type Player struct {
	name   string
	unit   *Unit
	mouseX float64 // tiles
	mouseY float64
}

type Unit struct {
	utype      *UnitType
	formation  []Unit // what you are commanding
	controller *UnitController
	payloads   []interface{}
	vel        Vec2
	stack      ItemStack

	hovering      bool
	dead          bool
	disarmed      bool
	spawnedByCore bool // despawn if no player
	walked        bool
	mining        bool
	shooting      bool
	building      bool

	x                 float64 // tiles, not """block units"""
	y                 float64
	rotation          float64
	elevation         float64
	maxHealth         float64
	drag              float64 // how much we divide speed in a frame
	armor             float64 // how much damage is reduced from each bullet
	hitSize           float64
	health            float64
	ammo              float64
	minFormationSpeed float64
	flag              float64
	dragMultiplier    float64
	strafePenalty     float64
	speedMultiplier   float64
	shootX            float64
	shootY            float64
	resupplyTime      float64

	team  int
	id    int
	mineX int // TODO: replace with Tile?
	mineY int

	buildX int
	buildY int
}

type DrawAction string

const (
	DrawActionClear    = DrawAction("clear")
	DrawActionColor    = DrawAction("color")
	DrawActionStroke   = DrawAction("stroke")
	DrawActionLine     = DrawAction("line")
	DrawActionRect     = DrawAction("rect")
	DrawActionLineRect = DrawAction("lineRect")
	DrawActionPoly     = DrawAction("poly")
	DrawActionLinePoly = DrawAction("linePoly")
	DrawActionTriangle = DrawAction("triangle")
	DrawActionImage    = DrawAction("image")
)

type DrawStatement struct {
	Action    DrawAction
	Arguments []interface{}
}

type Display interface {
	DrawFlush(buffer []DrawStatement)
}

type Memory interface {
	Write(value float64, position int64)
	Read(position int64) float64
}

type PostExecute interface {
	PostExecute()
}
