//go:build !race

package ws

// raceDetectorEnabled is false in normal (non -race) test builds, so the quarantined
// tests still run and provide coverage when the race detector is off.
const raceDetectorEnabled = false
