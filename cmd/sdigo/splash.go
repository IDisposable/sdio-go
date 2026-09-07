package main

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// splashDuration is how long screenSplash stays up before yielding to
// screenScanning on its own - purely cosmetic, so kept short. A user
// in a hurry can also dismiss it early with any keypress; a scan that
// finishes before either happens skips it entirely (scanDoneMsg sets
// m.screen directly, regardless of what's currently showing).
const splashDuration = 1200 * time.Millisecond

// splashDoneMsg fires once splashDuration elapses.
type splashDoneMsg struct{}

func tickSplashCmd() tea.Cmd {
	return tea.Tick(splashDuration, func(t time.Time) tea.Msg {
		return splashDoneMsg{}
	})
}

// splashBanner is the ASCII-art title screenSplash shows for
// splashDuration on startup - "GO FORTH" in block letters, this
// rewrite's own name (see the About screen and README.md), not a
// reproduction of the original VCL app's own splash/logo. Every line
// uses fixed-width source lines by construction, so the block-letter
// banner stays aligned in a monospace terminal.
const splashBanner = ` ██████╗  ██████╗      ███████╗ ██████╗ ██████╗ ████████╗██╗  ██╗
██╔════╝ ██╔═══██╗     ██╔════╝██╔═══██╗██╔══██╗╚══██╔══╝██║  ██║
██║  ███╗██║   ██║     █████╗  ██║   ██║██████╔╝   ██║   ███████║
██║   ██║██║   ██║     ██╔══╝  ██║   ██║██╔══██╗   ██║   ██╔══██║
╚██████╔╝╚██████╔╝     ██║     ╚██████╔╝██║  ██║   ██║   ██║  ██║
 ╚═════╝  ╚═════╝      ╚═╝      ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝`

// splashStyle gives splashBanner its color - the only screen in the
// TUI that's purely decorative, so it's the one place a splash of
// color doesn't fight with the table/status styling used everywhere
// else.
var splashStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)

// Brand colors from the diamond mark (assets/branding/icon.svg): the
// top facet, the left facet, and the right facet.
var (
	markTopStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ADD8")).Bold(true)
	markLeftStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#CE3262")).Bold(true)
	markRightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00758D")).Bold(true)
)

// diamondMark renders the icon's three-facet diamond as block-character
// ASCII art, split the same way as the SVG: a solid top facet over two
// halves colored left/right. Built row by row instead of as a literal
// so the two triangles and the facet split line up exactly.
func diamondMark() string {
	const n = 4 // half-height in rows; the widest row is 2*n-2 units wide
	var lines []string
	for r := 0; r < 2*n-1; r++ {
		width := n - abs(r-(n-1))
		pad := strings.Repeat("  ", n-width)
		if r < n-1 {
			// Top facet: still widening, full width in one color.
			lines = append(lines, pad+markTopStyle.Render(strings.Repeat("██", width))+pad)
			continue
		}
		left := (width + 1) / 2
		right := width - left
		row := markLeftStyle.Render(strings.Repeat("██", left)) +
			markRightStyle.Render(strings.Repeat("██", right))
		lines = append(lines, pad+row+pad)
	}
	return strings.Join(lines, "\n")
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// splashView renders screenSplash: the diamond mark beside splashBanner,
// plus a subtitle, centered in the terminal - see splashDuration/
// tickSplashCmd for how long it stays up.
func (m model) splashView() string {
	lockup := lipgloss.JoinHorizontal(lipgloss.Center,
		diamondMark(), "   ", splashStyle.Render(splashBanner))
	body := lockup + "\n\n" +
		"Snappy Driver Installer - reimplemented in Go\n\n" +
		"press any key to skip..."
	if m.width <= 0 || m.height <= 0 {
		return body
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
}
