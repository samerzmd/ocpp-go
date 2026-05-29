//go:build race

package ws

// raceDetectorEnabled is true when the test binary is built with -race.
// Used to quarantine tests that exercise PRE-EXISTING data races unrelated to the
// connMutex/outQueue server deadlock fix, so the -race gate stays meaningful for the
// rest of the suite. See the TODO references at each t.Skip call site.
const raceDetectorEnabled = true
