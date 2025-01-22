package command

import (
	"fmt"
	"strings"
	"testing"
	sdl "ttit/schedule"
)

func TestGeneratEntryFromCommand(t *testing.T) {
	thenMinutes := sdl.Minute(60)

	tests := map[string]sdl.Entry{
		"then work for 15m":        {Name: "work", From: thenMinutes, To: thenMinutes + sdl.Minute(15)},
		"then work until 15:00":    {Name: "work", From: thenMinutes, To: sdl.Minute(15 * 60)},
		"work from 15:00 to 16:30": {Name: "work", From: sdl.Minute(15 * 60), To: sdl.Minute(16*60 + 30)},
		"work from 15:00 for 45m":  {Name: "work", From: sdl.Minute(15 * 60), To: sdl.Minute(15*60 + 45)},
	}

	for input, want := range tests {
		got, err := GenerateEntryFromCommand(input, thenMinutes)
		if err != nil {
			t.Errorf("%s: %s", input, err)
			continue
		}

		if *got != want {
			t.Errorf("%s: got %+v, wanted %+v", input, got, want)
		}
	}
}

// TODO
func TestParseCommand(t *testing.T) {
	tests := []string{
		"then work for 15m",
		"then work until 15:00",
		"work from 15:00 to 16:00",
		"work from 15:00 for 45m",
	}

	for _, input := range tests {
		parsed, err := parseCommand(strings.Fields(input))

		if err != nil {
			t.Errorf("%s: %s", input, err)
			continue
		}

		fmt.Print("[")
		for i, part := range parsed {
			if i != 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%+v", part)
		}
		fmt.Println("]")
	}
}

func TestParseValidTime(t *testing.T) {
	tests := map[string]sdl.Minute{
		"00:00": sdl.Minute(0),
		"1":     sdl.Minute(60),
		"01":    sdl.Minute(60),
		"17":    sdl.Minute(17 * 60),
		"15:00": sdl.Minute(15 * 60),
		"15:30": sdl.Minute(15*60 + 30),
		"24:00": sdl.Minute(24 * 60),
		"00:30": sdl.Minute(30),
	}

	for input, want := range tests {
		got, err := minutesFromTimeString(input)
		if err != nil {
			t.Errorf("%s: %s", input, err)
			continue
		}

		if got != want {
			t.Errorf("%s: got %d, wanted %d", input, got, want)
		}
	}
}

func TestParseInvalidTime(t *testing.T) {
	tests := []string{"", ":", "-1", "25", "24:30", "00:89", "15:-2"}

	for _, input := range tests {
		_, err := minutesFromTimeString(input)
		if err == nil {
			t.Errorf("Expected error for invalid input %s, but got none", input)
		}
	}
}

func TestParseValidDuration(t *testing.T) {
	tests := map[string]sdl.Minute{
		"15m":    sdl.Minute(15),
		"1h":     sdl.Minute(60),
		"10h":    sdl.Minute(10 * 60),
		"1m":     sdl.Minute(1),
		"5h5m":   sdl.Minute(5*60 + 5),
		"10h10m": sdl.Minute(10*60 + 10),
	}

	for input, want := range tests {
		got, err := minutesFromDurationString(input)
		if err != nil {
			t.Errorf("%s: %s", input, err)
			continue
		}

		if got != want {
			t.Errorf("%s: got %d, wanted %d", input, got, want)
		}
	}
}

func TestParseInvalidDuration(t *testing.T) {
	tests := []string{"", "-1h", "home", "-30m", "0h0m", "10", "10h-2m"}

	for _, input := range tests {
		_, err := minutesFromDurationString(input)
		if err == nil {
			t.Errorf("Expected error for invalid input %s, but got none", input)
		}
	}
}

func TestIsKeyword(t *testing.T) {
	kwd := &CmdPart{Type: CptKeyword, Value: KwdFrom}
	if !kwd.isKeyword(KwdFrom) {
		t.Error("got false, wanted true")
	}
}
