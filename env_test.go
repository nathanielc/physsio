package physsio_test

import (
	"testing"

	"github.com/nathanielc/physsio"
	"gonum.org/v1/gonum/mat"
)

var simpleGame = []float64{
	// Player A
	// Active, Health, Move0 , Move1, Move2
	1, 100, 5, 10, 10, 5, 1, 100, 5, 2, 50,
	0, 101, 5, 10, 10, 5, 1, 100, 5, 2, 50,
	0, 102, 5, 10, 10, 5, 1, 100, 5, 2, 50,
	0, 103, 5, 10, 10, 5, 1, 100, 5, 2, 50,
	0, 104, 5, 10, 10, 5, 1, 100, 5, 2, 50,
	// Player B
	// Active, Health, Move0 , Move1, Move2
	0, 110, 6, 10, 10, 6, 1, 100, 6, 2, 50,
	0, 111, 6, 10, 10, 6, 1, 100, 6, 2, 50,
	1, 112, 6, 10, 10, 6, 1, 100, 6, 2, 50,
	0, 113, 6, 10, 10, 6, 1, 100, 6, 2, 50,
	0, 114, 6, 10, 10, 6, 1, 100, 6, 2, 50,
}

func TestEnv(t *testing.T) {
	testCases := []struct {
		name    string
		init    *physsio.State
		actions [][]int
		want    *physsio.State
	}{
		{
			name: "playerB wins",
			init: physsio.NewState(simpleGame),
			actions: [][]int{
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
				{0, 0},
			},
			want: physsio.NewState([]float64{
				0, 0, 5, 10, 1, 5, 1, 100, 5, 2, 50,
				1, 101, 5, 10, 10, 5, 1, 100, 5, 2, 50,
				0, 102, 5, 10, 10, 5, 1, 100, 5, 2, 50,
				0, 103, 5, 10, 10, 5, 1, 100, 5, 2, 50,
				0, 104, 5, 10, 10, 5, 1, 100, 5, 2, 50,
				0, 110, 6, 10, 10, 6, 1, 100, 6, 2, 50,
				0, 111, 6, 10, 10, 6, 1, 100, 6, 2, 50,
				1, 22, 6, 10, 0, 6, 1, 100, 6, 2, 50,
				0, 113, 6, 10, 10, 6, 1, 100, 6, 2, 50,
				0, 114, 6, 10, 10, 6, 1, 100, 6, 2, 50,
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			e := physsio.NewEnv(tc.init)
			for _, a := range tc.actions {
				e.Step(a[0], a[1])
			}
			got := e.State()
			if !got.Equal(tc.want) {
				t.Errorf("unexpected state got:\n%v\nwant:\n%v\n", mat.Formatted(got), mat.Formatted(tc.want))
			}
		})
	}
}
func BenchmarkEnv(b *testing.B) {
	b.ReportAllocs()

	e := physsio.NewEnv(physsio.NewState(simpleGame))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		e.Step(0, 0)
	}
}
