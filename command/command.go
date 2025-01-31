package command

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	sdl "github.com/NazarDunets/sitt/schedule"
)

type CmdPart struct {
	Type  CmdPartType
	Value any
}

type Keyword string
type CmdPartType string

type parserFn func([]string) (*CmdPart, int, error)

// 'until' and 'to' are equal in terms of handling. It's just more natural to write 'from 15:00 to 16:00' and 'do_smth until 16:00'
const (
	KwdThen  Keyword = "then"
	KwdFrom  Keyword = "from"
	KwdUntil Keyword = "until"
	KwdTo    Keyword = "to"
	KwdFor   Keyword = "for"
)

const (
	CptKeyword  CmdPartType = "keyword"
	CptName     CmdPartType = "name"
	CptTime     CmdPartType = "time"
	CptDuration CmdPartType = "duration"
)

const (
	NameClear = "clear"
	TimeNow   = "now"
)

var (
	keywords = map[Keyword]bool{
		KwdThen:  true,
		KwdFrom:  true,
		KwdUntil: true,
		KwdFor:   true,
		KwdTo:    true,
	}

	parsers = []parserFn{
		parseTime,
		parseDuration,
		parseKeyword,
		parseName,
	}
)

func GenerateEntryFromCommand(command string, thenMinutes sdl.Minute) (*sdl.Entry, error) {
	tokens := strings.Fields(command)
	parts, err := parseCommand(tokens)

	if err != nil {
		return nil, err
	}

	entry := new(sdl.Entry)

	nameFound := false
	startTimeFound := false
	endTimeFound := false

	for len(parts) > 0 {
		switch {
		// 'then' or entry name
		case !nameFound && !startTimeFound:
			part := parts[0]
			switch {
			case part.isKeyword(KwdThen):
				entry.From = thenMinutes
				startTimeFound = true

			case part.Type == CptName:
				entry.Name = part.Value.(string)
				nameFound = true

			default:
				return nil, errors.New("expected 'then' or entry name")
			}

			parts = parts[1:]

		// entry name
		case !nameFound:
			part := parts[0]
			if part.Type != CptName {
				return nil, errors.New("expected 'then' or entry name")
			}
			entry.Name = part.Value.(string)
			nameFound = true
			parts = parts[1:]

		// start time (command didn't start with 'then')
		case nameFound && !startTimeFound:
			if len(parts) < 2 {
				return nil, invalidCommandError()
			}

			fromKwd := parts[0]
			fromTime := parts[1]

			if !fromKwd.isKeyword(KwdFrom) || fromTime.Type != CptTime {
				return nil, errors.New("expected start time")
			}

			entry.From = fromTime.Value.(sdl.Minute)
			startTimeFound = true

			parts = parts[2:]

		// end time or duration
		case nameFound && startTimeFound && !endTimeFound:
			if len(parts) < 2 {
				return nil, invalidCommandError()
			}

			kwd := parts[0]
			switch {
			case kwd.isKeyword(KwdFor):
				duration := parts[1]
				if duration.Type != CptDuration {
					return nil, invalidCommandError()
				}
				entry.To = entry.From + duration.Value.(sdl.Minute)
				endTimeFound = true

			case kwd.isKeyword(KwdUntil) || kwd.isKeyword(KwdTo):
				time := parts[1]
				if time.Type != CptTime {
					return nil, invalidCommandError()
				}
				entry.To = time.Value.(sdl.Minute)
				endTimeFound = true

			default:
				return nil, errors.New("expected 'for' or 'until'")
			}
			parts = parts[2:]

		default:
			return nil, invalidCommandError()
		}
	}

	if !nameFound || !startTimeFound || !endTimeFound {
		return nil, invalidCommandError()
	}

	if entry.From >= entry.To {
		return nil, errors.New("entry duration (To - From) must be a positive value")
	}

	return entry, nil
}

func parseCommand(tokens []string) ([]*CmdPart, error) {
	result := make([]*CmdPart, 0)

	for len(tokens) > 0 {
		parsed := false
		for _, parser := range parsers {
			part, consumed, err := parser(tokens)

			if err == nil {
				parsed = true
				result = append(result, part)
				tokens = tokens[consumed:]
				break
			}
		}

		if !parsed {
			return []*CmdPart{}, errors.New(fmt.Sprintf("failed to parse token '%s'", tokens[0]))
		}
	}

	return result, nil
}

func parseName(tokens []string) (*CmdPart, int, error) {
	if len(tokens) < 1 {
		return parseResultTooShort()
	}

	token := tokens[0]
	if isKeyword(token) {
		return nil, 0, errors.New(fmt.Sprintf("expected entry name. Got keyword %s instead", token))
	}

	if token == NameClear {
		token = ""
	}

	return &CmdPart{Type: CptName, Value: token}, 1, nil
}

func parseKeyword(tokens []string) (*CmdPart, int, error) {
	if len(tokens) < 1 {
		return parseResultTooShort()
	}

	token := tokens[0]
	if !isKeyword(token) {
		return nil, 0, errors.New(fmt.Sprintf("expected a keyword. Got %s instead", token))
	}

	return &CmdPart{Type: CptKeyword, Value: Keyword(token)}, 1, nil
}

func parseTime(tokens []string) (*CmdPart, int, error) {
	if len(tokens) < 1 {
		return parseResultTooShort()
	}

	token := tokens[0]

	if token == TimeNow {
		token = time.Now().Format("15:04")
	}

	minutes, err := minutesFromTimeString(token)

	if err != nil {
		return nil, 0, err
	}

	return &CmdPart{Type: CptTime, Value: minutes}, 1, nil
}

func parseDuration(tokens []string) (*CmdPart, int, error) {
	if len(tokens) < 1 {
		return parseResultTooShort()
	}

	token := tokens[0]
	minutes, err := minutesFromDurationString(token)

	if err != nil {
		return nil, 0, err
	}

	return &CmdPart{Type: CptDuration, Value: minutes}, 1, nil
}

func minutesFromDurationString(token string) (sdl.Minute, error) {
	var err error
	hours := 0
	minutes := 0

	beforeH, afterH, hoursFound := strings.Cut(token, "h")
	if hoursFound {
		hours, err = strconv.Atoi(beforeH)
		if err != nil {
			return 0, invalidDurationError(token)
		}
	}

	minutesPart := ""
	if hoursFound {
		minutesPart = afterH
	} else {
		minutesPart = token
	}

	minutesCountString, _, minutesFound := strings.Cut(minutesPart, "m")
	if minutesFound {
		minutes, err = strconv.Atoi(minutesCountString)
		if err != nil {
			return 0, invalidDurationError(token)
		}
	}

	if hours < 0 || minutes < 0 || (hours == 0 && minutes == 0) {
		return 0, invalidDurationError(token)
	}

	return minFromHour(hours) + sdl.Minute(minutes), nil
}

func minutesFromTimeString(token string) (sdl.Minute, error) {
	switch len(token) {
	case 1, 2:
		hours, err := strconv.Atoi(token)

		if err != nil {
			return 0, invalidTimeError(token)
		}

		if hours < 0 || hours > 24 {
			return 0, invalidTimeError(token)
		}

		return minFromHour(hours), nil

	case 4, 5:
		hoursString, minutesString, found := strings.Cut(token, ":")

		if !found {
			return 0, invalidTimeError(token)
		}

		hours, err := strconv.Atoi(hoursString)
		minutes, err1 := strconv.Atoi(minutesString)

		if err != nil || err1 != nil {
			return 0, invalidTimeError(token)
		}

		if hours < 0 || hours > 24 || minutes < 0 || minutes > 59 || (hours == 24 && minutes != 0) {
			return 0, invalidTimeError(token)
		}

		return minFromHour(hours) + sdl.Minute(minutes), nil
	default:
		return 0, invalidTimeError(token)
	}
}

func (p *CmdPart) isKeyword(kwd Keyword) bool {
	return p.Type == CptKeyword && p.Value == kwd
}

func isKeyword(token string) bool {
	return keywords[Keyword(token)]
}

func invalidCommandError() error {
	return errors.New("invalid commad structure")
}

func parseResultTooShort() (*CmdPart, int, error) {
	return nil, 0, errors.New("command to short")
}

func invalidTimeError(got string) error {
	return errors.New(fmt.Sprintf("expected time in format 'HH' or 'HH:MM'. Got %s", got))
}

func invalidDurationError(got string) error {
	return errors.New(fmt.Sprintf("expected time in format {X}h{Y}m/{X}h/{X}m. Got %s", got))
}

func minFromHour(hours int) sdl.Minute {
	return sdl.Minute(hours * 60)
}
