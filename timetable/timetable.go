package timetable

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

type Timetable struct {
	Entries []Entry `json:"entries"`
}

func (m Minute) String() string {
	hours := m / 60
	mins := m % 60

	return fmt.Sprintf("%02d:%02d", hours, mins)
}

func (t Timetable) String() string {
	var buffer bytes.Buffer
	for _, entry := range t.Entries {
		buffer.WriteString(entry.From.String())
		buffer.WriteString("-")
		buffer.WriteString(entry.To.String())
		buffer.WriteString(" : ")
		buffer.WriteString(" : ")
		buffer.WriteString(entry.Name)
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// TODO: merge entries with same name?
func (t *Timetable) Insert(newEntry Entry) {
	newEntry.From = max(newEntry.From, 0)
	newEntry.To = min(newEntry.To, 24*60)

	if newEntry.To <= newEntry.From {
		return
	}

	cropEnd := -1
	cropStart := -1
	deleteFrom := -1
	deleteToInclusive := -1

	for i, entry := range t.Entries {
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

		newEntries = append(newEntries, t.Entries[:splitIndex]...)

		entryToSplit := t.Entries[splitIndex]
		entryToSplit.To = newEntry.From
		newEntries = append(newEntries, entryToSplit)

		newEntries = append(newEntries, newEntry)

		entryToSplit = t.Entries[splitIndex]
		entryToSplit.From = newEntry.To
		newEntries = append(newEntries, entryToSplit)

		newEntries = append(newEntries, t.Entries[splitIndex+1:]...)
		t.Entries = newEntries
		return
	}

	if cropEnd != -1 {
		toCrop := &(t.Entries[cropEnd])
		toCrop.To = newEntry.From
	}

	if cropStart != -1 {
		toCrop := &(t.Entries[cropStart])
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
	if includeToInclusive >= 0 && includeToInclusive < len(t.Entries) {
		newEntries = append(newEntries, t.Entries[0:includeToInclusive+1]...)
	}

	newEntries = append(newEntries, newEntry)

	if includeFromInclusive >= 0 && includeFromInclusive < len(t.Entries) {
		newEntries = append(newEntries, t.Entries[includeFromInclusive:len(t.Entries)]...)
	}
	t.Entries = newEntries
}

func (t *Timetable) FilledUpTo() Minute {
	for i := len(t.Entries) - 1; i > 0; i-- {
		entry := t.Entries[i]
		if entry.Name != "" {
			return entry.To
		}
	}

	return Minute(0)
}

func New() Timetable {
	return Timetable{
		Entries: []Entry{newEmptyEntry(Minute(0), Minute(24*60))},
	}
}

func newEmptyEntry(from, to Minute) Entry {
	return Entry{Name: "", From: from, To: to}
}
