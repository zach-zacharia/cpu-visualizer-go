package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/cpu"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	// Create a paragraph widget for displaying CPU percentages
	p := widgets.NewParagraph()
	p.Title = "CPU Usage"
	p.SetRect(0, 0, 50, 10)
	p.TextStyle.Fg = termui.ColorWhite
	p.BorderStyle.Fg = termui.ColorCyan

	// Update CPU usage every second
	updateInterval := time.Second
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	uiEvents := termui.PollEvents()

	for {
		select {
		case <-ticker.C:
			percentages, err := cpu.Percent(0, true)
			if err != nil {
				log.Fatalf("failed to get CPU usage: %v", err)
			}

			// Build the display text with percentages
			text := ""
			for i, perc := range percentages {
				text += fmt.Sprintf("CPU%d: %.1f%%\n", i, perc)
			}

			// Update the paragraph widget with the new text
			p.Text = text

			termui.Render(p)

		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				return
			}
		}
	}
}
