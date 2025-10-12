package geo

import "math"

const earthRadiusKm = 6371.0

// DistanceBetween returns the great-circle distance in kilometers between two coordinates.
func DistanceBetween(lat1, lon1, lat2, lon2 float64) float64 {
	// 1.- Convert degrees to radians and calculate deltas.
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)
	lat1Rad := degreesToRadians(lat1)
	lat2Rad := degreesToRadians(lat2)

	// 2.- Apply the haversine formula to compute the arc distance.
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// ScoreByProximity ranks a slice of distances (in kilometers) where a lower score is better.
func ScoreByProximity(distances []float64) []float64 {
	scores := make([]float64, len(distances))
	for i, d := range distances {
		scores[i] = 1 / (1 + d)
	}
	return scores
}

// WithinRadius reports whether the target coordinate lies inside radiusKm of the origin.
func WithinRadius(originLat, originLon, targetLat, targetLon, radiusKm float64) bool {
	return DistanceBetween(originLat, originLon, targetLat, targetLon) <= radiusKm
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
