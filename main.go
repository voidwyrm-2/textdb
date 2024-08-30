package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func displayEntry(entry Entry) {
	p := tea.NewProgram(
		EntryDisplay{entryName: entry.Name(), content: entry.Content()},
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}

func main() {
	dbfile := ".database.sds"
	if len(os.Args) > 1 {
		dbfile = os.Args[1]
	}

	content, err := readFile(dbfile)
	if err != nil {
		if err.Error() == "open .database.sds: The system cannot find the file specified." {
			writeFile(".database.sds", "")
		}
		fmt.Println(err.Error())
		return
	}

	if strings.TrimSpace(content) == "" {
		fmt.Println("database file is empty")
		return
	}

	entries, _, err := interpretSDS(content, false, 0)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(entries)

	mainView := tea.NewProgram(initEntryView(&entries, "<MAIN>"))
	if _, err := mainView.Run(); err != nil {
		fmt.Printf("error in main view: %v", err)
		os.Exit(1)
	}
}
