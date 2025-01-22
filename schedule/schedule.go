package schedule

import (
	"bytes"
	"fmt"
)

type Minute int

type Entry struct {
	From Minute `json:"from"`
	To   Minute `json:"to"`
	Name string `json:"name"`
}

type Schedule struct {
	Entries []Entry `json:"entries"`
}

func (m Minute) String() string {
	hours := m / 60
	mins := m % 60

	return fmt.Sprintf("%02d:%02d", hours, mins)
}

func (sdl Schedule) View() string {
	var buffer bytes.Buffer

	if len(sdl.Entries) == 0 || len(sdl.Entries) == 1 && sdl.Entries[0].Name == "" {
		return "Schedule is empty\n"
	}

	for i, entry := range sdl.Entries {
		if i == 0 && entry.Name == "" {
			continue
		}

		buffer.WriteString(entry.From.String())
		buffer.WriteString(" : ")

		name := entry.Name
		if name == "" {
			name = "---"
		}
		buffer.WriteString(name)

		buffer.WriteString(" (")
		buffer.WriteString(entry.durationString())
		buffer.WriteString(")\n")
	}

	buffer.WriteString("24:00 : DAY END\n")

	return buffer.String()
}

func (e Entry) durationString() string {
	durationMinutes := e.To - e.From

	hours := durationMinutes / 60
	minutes := durationMinutes % 60

	var buffer bytes.Buffer
	if hours > 0 {
		buffer.WriteString(fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		buffer.WriteString(fmt.Sprintf("%dm", minutes))
	}

	return buffer.String()
}

// TODO: merge entries with same name?
func (sdl *Schedule) Insert(newEntry Entry) {
	newEntry.From = max(newEntry.From, 0)
	newEntry.To = min(newEntry.To, 24*60)

	if newEntry.To <= newEntry.From {
		return
	}

	cropEnd := -1
	cropStart := -1
	deleteFrom := -1
	deleteToInclusive := -1

	for i, entry := range sdl.Entries {
		if entry.From >= newEntry.From && entry.To <= newEntry.To {
			if deleteFrom == -1 {
				deleteFrom = i
			}
			deleteToInclusive = i
		} else {
			if entry.From < newEntry.From && entry.To > newEntry.From {
				cropEnd = i
			}

			if entry.To > newEntry.To && entry.From < newEntry.To {
				cropStart = i
			}
		}
	}

	// fmt.Printf("cropEnd: %d, cropStart: %d, deleteFrom: %d, deleteTo: %d\n", cropEnd, cropStart, deleteFrom, deleteToInclusive)

	// split
	if cropEnd == cropStart && cropEnd != -1 {
		splitIndex := cropEnd
		newEntries := make([]Entry, 0)

		newEntries = append(newEntries, sdl.Entries[:splitIndex]...)

		entryToSplit := sdl.Entries[splitIndex]
		entryToSplit.To = newEntry.From
		newEntries = append(newEntries, entryToSplit)

		newEntries = append(newEntries, newEntry)

		entryToSplit = sdl.Entries[splitIndex]
		entryToSplit.From = newEntry.To
		newEntries = append(newEntries, entryToSplit)

		newEntries = append(newEntries, sdl.Entries[splitIndex+1:]...)
		sdl.Entries = newEntries
		return
	}

	if cropEnd != -1 {
		toCrop := &(sdl.Entries[cropEnd])
		toCrop.To = newEntry.From
	}

	if cropStart != -1 {
		toCrop := &(sdl.Entries[cropStart])
		toCrop.From = newEntry.To
	}

	includeToInclusive := 0
	switch {
	case cropEnd != -1:
		includeToInclusive = cropEnd
	case deleteFrom != -1:
		includeToInclusive = deleteFrom - 1
	case cropStart != -1:
		includeToInclusive = cropStart - 1
	}

	includeFromInclusive := 0
	switch {
	case cropStart != -1:
		includeFromInclusive = cropStart
	case deleteToInclusive != -1:
		includeFromInclusive = deleteToInclusive + 1
	case cropEnd != -1:
		includeFromInclusive = cropEnd + 1
	}

	// fmt.Printf("to: %d, from: %d\n", includeToInclusive, includeFromInclusive)

	newEntries := make([]Entry, 0)
	if includeToInclusive >= 0 && includeToInclusive < len(sdl.Entries) {
		newEntries = append(newEntries, sdl.Entries[0:includeToInclusive+1]...)
	}

	newEntries = append(newEntries, newEntry)

	if includeFromInclusive >= 0 && includeFromInclusive < len(sdl.Entries) {
		newEntries = append(newEntries, sdl.Entries[includeFromInclusive:len(sdl.Entries)]...)
	}
	sdl.Entries = newEntries
}

func (sdl *Schedule) FilledUpTo() Minute {
	for i := len(sdl.Entries) - 1; i > 0; i-- {
		entry := sdl.Entries[i]
		if entry.Name != "" {
			return entry.To
		}
	}

	return Minute(0)
}

func New() Schedule {
	return Schedule{
		Entries: []Entry{newEmptyEntry(Minute(0), Minute(24*60))},
	}
}

func newEmptyEntry(from, to Minute) Entry {
	return Entry{Name: "", From: from, To: to}
}
