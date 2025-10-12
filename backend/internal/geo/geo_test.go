package geo

import "testing"

func TestDistanceBetween(t *testing.T) {
	tests := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
		expected   float64
	}{
		{"same point", 0, 0, 0, 0, 0},
		{"london to paris", 51.5074, -0.1278, 48.8566, 2.3522, 343},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := DistanceBetween(tc.lat1, tc.lon1, tc.lat2, tc.lon2)
			if diff := mathAbs(d - tc.expected); diff > 5 {
				t.Fatalf("distance mismatch: got %.1f want %.1f", d, tc.expected)
			}
		})
	}
}

func TestWithinRadius(t *testing.T) {
	tests := []struct {
		name                 string
		originLat, originLon float64
		targetLat, targetLon float64
		radius               float64
		expected             bool
	}{
		{"inside", 0, 0, 0.1, 0.1, 20, true},
		{"outside", 0, 0, 10, 10, 100, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := WithinRadius(tc.originLat, tc.originLon, tc.targetLat, tc.targetLon, tc.radius)
			if got != tc.expected {
				t.Fatalf("expected %v got %v", tc.expected, got)
			}
		})
	}
}

func mathAbs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
