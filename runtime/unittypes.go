import (
	"runtime/types"
)

var (
	UnitFlare = UnitType{
		name: "flare",
		speed: 3,
		accel: 0.08,
		drag: 0.01,
		flying: true,
		health: 75,
		engineOffset: 5.5,
		urange: 140,
		targetAir: false,
		//as default AI, flares are not very useful in core rushes, they attack nothing in the way
		//playerTargetFlags: new BlockFlag[]{null};
		//targetFlags: new BlockFlag[]{BlockFlag.generator, null};
		commandLimit: 4,
		circleTarget: true,
		hitSize: 7,
