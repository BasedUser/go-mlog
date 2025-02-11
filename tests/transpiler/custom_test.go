package transpiler

import (
	"github.com/MarvinJWendt/testza"
	"github.com/Vilsol/go-mlog/transpiler"
	"strings"
	"testing"
)

func TestCustom(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name: "Const",
			input: TestMain(`x := m.Const("@copper")
print(x)`, true, false),
			output: `set _main_x @copper
print _main_x`,
		},
		{
			name: "NestedSelector",
			input: TestMain(`x := m.This.GetX()
print(x)`, true, false),
			output: `sensor _main_x @this @x
print _main_x`,
		},
		{
			name:  "Inline",
			input: TestMain(`print(m.Const("@copper"), m.B("message1"), float64(1))`, true, false),
			output: `print @copper
print message1
print 1`,
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

			test.output = test.output + "\nend"
			testza.AssertEqual(t, test.output, strings.Trim(mlog, "\n"))
		})
	}
}
