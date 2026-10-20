package core

import "testing"

func TestSegmentStates(t *testing.T) {
	const F, H, E = segmentFull, segmentHalf, segmentEmpty
	tests := []struct {
		name  string
		ratio float64
		want  []segmentState
	}{
		{"full", 1, []segmentState{F, F, F, F}},
		{"empty", 0, []segmentState{E, E, E, E}},
		{"exact half segment", 0.375, []segmentState{F, H, E, E}},
		{"exact whole segments", 0.5, []segmentState{F, F, E, E}},
		{"partial rounds up to half", 0.26, []segmentState{F, H, E, E}},
		{"tiny amount shows half", 0.01, []segmentState{H, E, E, E}},
		{"above max clamps", 1.5, []segmentState{F, F, F, F}},
		{"below zero clamps", -0.2, []segmentState{E, E, E, E}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := segmentStates(tt.ratio, len(tt.want))
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("segmentStates(%v) = %v, want %v", tt.ratio, got, tt.want)
				}
			}
		})
	}
}

// Health is 100 points over 10 hearts, so each heart is 10 points and a half
// heart is 5.
func TestSegmentStatesHealthSteps(t *testing.T) {
	cases := map[int][2]int{ // health -> {full, half}
		100: {10, 0},
		95:  {9, 1},
		91:  {9, 1},
		90:  {9, 0},
		45:  {4, 1},
		1:   {0, 1},
	}
	for hp, want := range cases {
		full, half := 0, 0
		for _, s := range segmentStates(float64(hp)/MaxPlayerHealth, hudSegments) {
			switch s {
			case segmentFull:
				full++
			case segmentHalf:
				half++
			}
		}
		if full != want[0] || half != want[1] {
			t.Errorf("health %d: full=%d half=%d, want %v", hp, full, half, want)
		}
	}
}
