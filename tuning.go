package main

import (
	"fmt"
	"regexp"
	"sort"
)

type Tamperament struct {
	Octave string
	Note   string
	Tuning int
}

type Scale struct {
	Tamperaments []Tamperament
	PitchMap     map[string]int
	ChannelMap   map[int]uint8
}

func keyToOctaveAndNote(key string) (string, string) {
	re := regexp.MustCompile(`(\w+)(\d+)`)
	occurences := re.FindAllStringSubmatch(key, 1)
	if len(occurences) == 0 {
		return "", ""
	}
	return occurences[0][2], occurences[0][1]
}

func newTamperament(
	octave string, note string, tuning int,
) Tamperament {
	if tuning < -200 {
		tuning = -200
	}
	if tuning > 200 {
		tuning = 200
	}
	return Tamperament{
		octave, note, tuning,
	}
}

func newScale() *Scale {
	s := Scale{}
	s.ChannelMap = make(map[int]uint8)
	s.PitchMap = make(map[string]int)
	s.Tamperaments = make([]Tamperament, 0)
	s.updateMaps()
	return &s
}

func (s *Scale) addTamperament(octave string, note string, tuning int) {
	s.Tamperaments = append(s.Tamperaments, newTamperament(octave, note, tuning))
	s.updateMaps()
}

func (s *Scale) sortTamperaments() {
	sort.Slice(s.Tamperaments, func(i, j int) bool {
		if s.Tamperaments[i].Octave == s.Tamperaments[j].Octave {
			return s.Tamperaments[i].Note < s.Tamperaments[j].Note
		}
		return s.Tamperaments[i].Octave < s.Tamperaments[j].Octave
	})
}

func (s *Scale) pitchForKey(key string) int {
	octave, note := keyToOctaveAndNote(key)
	pitch := 0
	for _, t := range s.Tamperaments {
		if note == t.Note && (octave == t.Octave ||
			t.Octave == "*") {
			pitch = t.Tuning
		}
	}
	return pitch
}

func (s *Scale) updateMaps() {
	s.sortTamperaments()
	pitches := make(map[int]bool)
	pitches[0] = true
	notes := []string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B"}
	for octave := 0; octave <= 9; octave++ {
		for _, note := range notes {
			key := fmt.Sprintf("%v%d", note, octave)
			pitch := s.pitchForKey(key)
			s.PitchMap[key] = pitch
			pitches[pitch] = true
		}
	}
	s.ChannelMap[0] = 0 //Always set channel-0 to pitch-bend zero
	channel := uint8(1)
	for pitch := range pitches {
		if pitch == 0 {
			continue
		}
		s.ChannelMap[pitch] = channel
		channel++
	}
}
