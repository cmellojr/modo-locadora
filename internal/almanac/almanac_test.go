package almanac

import (
	"testing"
)

func TestTodaysEphemeride(t *testing.T) {
	entry := TodaysEphemeride()
	if entry == "" {
		t.Fatal("expected TodaysEphemeride to return a non-empty string")
	}
}

func TestEphemeridesDataIntegrity(t *testing.T) {
	if len(genericEphemerides) == 0 {
		t.Fatal("genericEphemerides slice should not be empty")
	}
	for i, text := range genericEphemerides {
		if text == "" {
			t.Errorf("genericEphemerides[%d] is empty", i)
		}
	}

	if len(ephemerides) == 0 {
		t.Fatal("ephemerides map should not be empty")
	}
	for day, text := range ephemerides {
		if day < 1 || day > 366 {
			t.Errorf("ephemerides day %d is out of valid range (1-366)", day)
		}
		if text == "" {
			t.Errorf("ephemerides[%d] is empty", day)
		}
	}
}
