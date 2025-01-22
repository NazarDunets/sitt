package command

import (
	"fmt"
	"strings"
	"testing"
	tt "ttit/timetable"
)

func TestGeneratEntryFromCommand(t *testing.T) {
	thenMinutes := tt.Minute(60)

	tests := map[string]tt.Entry{
		"then work for 15m":        {Name: "work", From: thenMinutes, To: thenMinutes + tt.Minute(15)},
		"then work until 15:00":    {Name: "work", From: thenMinutes, To: tt.Minute(15 * 60)},
		"work from 15:00 to 16:30": {Name: "work", From: tt.Minute(15 * 60), To: tt.Minute(16*60 + 30)},
		"work from 15:00 for 45m":  {Name: "work", From: tt.Minute(15 * 60), To: tt.Minute(15*60 + 45)},
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
	tests := map[string]tt.Minute{
		"00:00": tt.Minute(0),
		"1":     tt.Minute(60),
		"01":    tt.Minute(60),
		"17":    tt.Minute(17 * 60),
		"15:00": tt.Minute(15 * 60),
		"15:30": tt.Minute(15*60 + 30),
		"24:00": tt.Minute(24 * 60),
		"00:30": tt.Minute(30),
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
	tests := map[string]tt.Minute{
		"15m":    tt.Minute(15),
		"1h":     tt.Minute(60),
		"10h":    tt.Minute(10 * 60),
		"1m":     tt.Minute(1),
		"5h5m":   tt.Minute(5*60 + 5),
		"10h10m": tt.Minute(10*60 + 10),
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
