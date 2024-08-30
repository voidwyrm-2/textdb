package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type EntryDisplay struct {
	entryName, content string
	ready              bool
	viewport           viewport.Model
}

func (ed EntryDisplay) Init() tea.Cmd {
	return nil
}

func (ed EntryDisplay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return ed, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(ed.headerView())
		footerHeight := lipgloss.Height(ed.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !ed.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			ed.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			ed.viewport.YPosition = headerHeight
			ed.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			ed.viewport.SetContent(ed.content)
			ed.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			ed.viewport.YPosition = headerHeight + 1
		} else {
			ed.viewport.Width = msg.Width
			ed.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(ed.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	ed.viewport, cmd = ed.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return ed, tea.Batch(cmds...)
}

func (ed EntryDisplay) View() string {
	if !ed.ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", ed.headerView(), ed.viewport.View(), ed.footerView())
}

func (ed EntryDisplay) headerView() string {
	title := titleStyle.Render(ed.entryName)
	line := strings.Repeat("─", max(0, ed.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (ed EntryDisplay) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", ed.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, ed.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
