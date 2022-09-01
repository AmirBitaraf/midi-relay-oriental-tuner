package main

import "testing"

func TestTamperamentTuning(t *testing.T) {
	t1 := newTamperament("3", "C", 250)
	if t1.Tuning != 200 {
		t.Errorf("Tuning is not set correctly")
	}

	t2 := newTamperament("3", "C", -250)
	if t2.Tuning != -200 {
		t.Errorf("Tuning is not set correctly")
	}
}

func TestKeyToOctaveNote(t *testing.T) {
	if octave, note := keyToOctaveAndNote("C8"); octave != "8" || note != "C" {
		t.Errorf("Key to octave note is not working correctly")
	}

	if octave, note := keyToOctaveAndNote("A0"); octave != "0" || note != "A" {
		t.Errorf("Key to octave note is not working correctly")
	}
}

func TestScaleInitial(t *testing.T) {
	s := newScale()

	if _, ok := s.PitchMap["B2"]; !ok {
		t.Errorf("Pitch not set correctly")
	}

	if len(s.ChannelMap) != 1 {
		t.Errorf("Channel map not set correctly")
	}

	if s.ChannelMap[0] != 0 {
		t.Errorf("Channel pitch not set correctly")
	}
}

func TestScaleWithTamperament(t *testing.T) {
	s := newScale()

	s.addTamperament("*", "Bb", -50)

	if s.PitchMap["B3"] != 0 {
		t.Errorf("Pitch not set correctly")
	}

	if s.PitchMap["Bb3"] != -50 {
		t.Errorf("Pitch not set correctly")
	}

	if len(s.ChannelMap) != 2 {
		t.Errorf("Channel map not set correctly")
	}

	if s.ChannelMap[-50] != 1 {
		t.Errorf("Channel pitch not set correctly")
	}
}
