package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type EntryView struct {
	entries *[]Entry
	cursor  int
	name    string
}

func initEntryView(entries *[]Entry, name string) EntryView {
	return EntryView{
		entries: entries,
		name:    name,
	}
}

func (ev EntryView) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (ev EntryView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return ev, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if ev.cursor > 0 {
				ev.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if ev.cursor < len(*ev.entries)-1 {
				ev.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			if ok, sub := (*ev.entries)[ev.cursor].IsFolder(); ok {
				subView := tea.NewProgram(initEntryView(&sub, (*ev.entries)[ev.cursor].Name()))
				if _, err := subView.Run(); err != nil {
					fmt.Printf("error in subview: %v", err)
					os.Exit(1)
				}

				(*ev.entries)[ev.cursor] = NewFolderEntry((*ev.entries)[ev.cursor].Name(), sub)
			} else if (*ev.entries)[ev.cursor].Name() != "" {
				displayEntry((*ev.entries)[ev.cursor])
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return ev, nil
}

func (ev EntryView) View() string {
	// The header
	s := ev.name + "\n\n"

	// Iterate over our choices
	for i, entry := range *ev.entries {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if ev.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, entry.Format())
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
