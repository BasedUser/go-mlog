package tests

import (
	"github.com/Vilsol/go-mlog/transpiler"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnitControl(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "UnitStop",
			input:  TestMain(`m.UnitStop()`),
			output: `ucontrol stop`,
		},
		{
			name:   "UnitMove",
			input:  TestMain(`m.UnitMove(1, 2)`),
			output: `ucontrol move 1 2`,
		},
		{
			name:   "UnitApproach",
			input:  TestMain(`m.UnitApproach(1, 2, 3)`),
			output: `ucontrol approach 1 2 3`,
		},
		{
			name:   "UnitBoost",
			input:  TestMain(`m.UnitBoost(1)`),
			output: `ucontrol boost 1`,
		},
		{
			name:   "UnitPathfind",
			input:  TestMain(`m.UnitPathfind()`),
			output: `ucontrol pathfind`,
		},
		{
			name:   "UnitTarget",
			input:  TestMain(`m.UnitTarget(1, 2, 3)`),
			output: `ucontrol target 1 2 3`,
		},
		{
			name:   "UnitTargetP",
			input:  TestMain(`m.UnitTargetP(1, 2)`),
			output: `ucontrol targetp 1 2`,
		},
		{
			name:   "UnitItemDrop",
			input:  TestMain(`m.UnitItemDrop(1, 2)`),
			output: `ucontrol itemDrop 1 2`,
		},
		{
			name:   "UnitItemTake",
			input:  TestMain(`m.UnitItemTake(1, "A", 2)`),
			output: `ucontrol itemTake 1 "A" 2`,
		},
		{
			name:   "UnitPayloadDrop",
			input:  TestMain(`m.UnitPayloadDrop()`),
			output: `ucontrol payDrop`,
		},
		{
			name:   "UnitPayloadTake",
			input:  TestMain(`m.UnitPayloadTake(1)`),
			output: `ucontrol payTake 1`,
		},
		{
			name:   "UnitMine",
			input:  TestMain(`m.UnitMine(1, 2)`),
			output: `ucontrol mine 1 2`,
		},
		{
			name:   "UnitFlag",
			input:  TestMain(`m.UnitFlag(1)`),
			output: `ucontrol flag 1`,
		},
		{
			name:   "UnitBuild",
			input:  TestMain(`m.UnitBuild(1, 2, "A", 3, 4)`),
			output: `ucontrol build 1 2 "A" 3 4`,
		},
		{
			name:   "UnitWithin",
			input:  TestMain(`m.UnitWithin(1, 2, 3)`),
			output: `ucontrol within 1 2 3 @return`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mlog, err := transpiler.GolangToMLOG(test.input, transpiler.Options{
				NoStartup: true,
			})

			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
