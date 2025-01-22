package main

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"ttit/command"
	"ttit/timetable"

	_ "github.com/charmbracelet/bubbletea"
)

func main() {
	tt := timetable.New()
	fmt.Println(tt)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		entry, err := command.GenerateEntryFromCommand(input, tt.FilledUpTo())
		if err == nil {
			tt.Insert(*entry)
			fmt.Println(tt)
		} else {
			log.Println(err)
		}
	}
}
