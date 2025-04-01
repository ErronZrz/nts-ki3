package clock

import (
	"math"
	"math/rand"
)

// KMeans function to perform K-means clustering on 1D data
func KMeans(values []float64, k int) []float64 {
	// 1. Initialize random centroids
	centroids := make([]float64, k)
	for i := 0; i < k; i++ {
		centroids[i] = values[rand.New(rand.NewSource(44)).Intn(len(values))]
	}

	// Repeat the process of assignment and update until convergence
	// Convergence criteria: when centroids don't change between iterations
	var previousCentroids []float64
	for {
		// 2. Assign each data point to the nearest centroid
		clusters := make([][]float64, k)
		for _, value := range values {
			closestCentroidIndex := findClosestCentroid(value, centroids)
			clusters[closestCentroidIndex] = append(clusters[closestCentroidIndex], value)
		}

		// 3. Recalculate the centroids
		for i := 0; i < k; i++ {
			centroids[i] = calculateMean(clusters[i])
		}

		// 4. Check if centroids have converged (i.e., they don't change)
		if hasConverged(previousCentroids, centroids) {
			break
		}

		// Save current centroids for the next iteration check
		previousCentroids = append([]float64{}, centroids...)
	}

	return centroids
}

// Helper function to calculate the mean of a slice of float64 values
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

// Helper function to find the index of the closest centroid for a given data point
func findClosestCentroid(value float64, centroids []float64) int {
	closestIndex := 0
	minDistance := math.Abs(value - centroids[0])
	for i := 1; i < len(centroids); i++ {
		distance := math.Abs(value - centroids[i])
		if distance < minDistance {
			closestIndex = i
			minDistance = distance
		}
	}
	return closestIndex
}

// Helper function to check if centroids have converged
func hasConverged(oldCentroids, newCentroids []float64) bool {
	if len(oldCentroids) != len(newCentroids) {
		return false
	}
	for i := 0; i < len(oldCentroids); i++ {
		if oldCentroids[i] != newCentroids[i] {
			return false
		}
	}
	return true
}
