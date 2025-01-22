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
	"time"
	"ttit/command"
	"ttit/storage"

	_ "github.com/charmbracelet/bubbletea"
)

func main() {
	time := time.Now()
	tt, err := storage.LoadOrCreateTimetable(time)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(tt)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		entry, err := command.GenerateEntryFromCommand(input, tt.FilledUpTo())

		if err != nil {
			log.Println(err)
			continue
		}

		tt.Insert(*entry)

		if err := storage.Save(time, tt); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(tt)
	}
}
