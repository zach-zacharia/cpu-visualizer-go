package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	// CPU percentages
	cpu_usage := widgets.NewParagraph()
	cpu_usage.Title = "CPU usage"
	cpu_usage.SetRect(0, 0, 50, 10)
	cpu_usage.TextStyle.Fg = termui.ColorWhite
	cpu_usage.BorderStyle.Fg = termui.ColorCyan

	memory_usage := widgets.NewParagraph()
	memory_usage.Title = "RAM usage"
	memory_usage.SetRect(52, 11, 100, 21)
	memory_usage.TextStyle.Fg = termui.ColorWhite
	memory_usage.BorderStyle.Fg = termui.ColorCyan

	swap_usage := widgets.NewParagraph()
	swap_usage.Title = "Swap usage"
	swap_usage.SetRect(52, 0, 100, 10)
	swap_usage.TextStyle.Fg = termui.ColorWhite
	swap_usage.BorderStyle.Fg = termui.ColorCyan

	memoryUpdateInterval := time.Second
	cpuUpdateInterval := time.Second
	swapUpdateInterval := time.Second
	tickercpu := time.NewTicker(cpuUpdateInterval)
	tickermem := time.NewTicker(memoryUpdateInterval)
	tickerswap := time.NewTicker(swapUpdateInterval)
	defer tickerswap.Stop()
	defer tickermem.Stop()
	defer tickercpu.Stop()

	uiEvents := termui.PollEvents()

	for {
		select {
		case <-tickercpu.C:
			percentages, err := cpu.Percent(0, true)
			if err != nil {
				fmt.Errorf("failed to get CPU usage: %v", err)
			}

			// Build the display text with percentages
			text := ""
			for i, perc := range percentages {
				text += fmt.Sprintf("CPU%d: %.1f%%\n", i, perc)
			}

			// Update the paragraph widget with the new text
			cpu_usage.Text = text

			termui.Render(cpu_usage)

		case <-tickermem.C:
			memory, err := mem.VirtualMemory()
			if err != nil {
				fmt.Errorf("Failed to get memory usage: %v", err)
			}

			memoryUsage := []float64{
				memory.UsedPercent,
				bytesToGB(memory.Total),
				bytesToGB(memory.Available),
				bytesToGB(memory.Used),
			}

			memoryDataType := []string{
				"Used memory (%)",
				"Total memory (GB)",
				"Available memory (GB)",
				"Used memory (GB)",
			}

			isPercentage := []bool{
				true,
				false,
				false,
				false,
			}

			text := ""

			for i, usage := range memoryUsage {
				text += formatMemoryUsage(memoryDataType[i], usage, isPercentage[i])
			}

			memory_usage.Text = text

			termui.Render(memory_usage)

		case <-tickerswap.C:
			{
				swap, err := mem.SwapMemory()
				if err != nil {
					fmt.Errorf("Failed to get swap usage: %v", err)
				}

				swapUsage := []float64{
					swap.UsedPercent,
					bytesToGB(swap.Total),
					bytesToGB(swap.Free),
					bytesToGB(swap.Used),
				}

				swapDataType := []string{
					"Used swap memory (%)",
					"Total swap memory (GB)",
					"Free swap memory (GB)",
					"Used swap memory (GB)",
				}

				isPercentage := []bool{
					true,
					false,
					false,
					false,
				}

				text := ""

				for i, usage := range swapUsage {
					text += formatMemoryUsage(swapDataType[i], usage, isPercentage[i])
				}

				swap_usage.Text = text

				termui.Render(swap_usage)
			}
		case e := <-uiEvents:
			if e.Type == termui.KeyboardEvent {
				return
			}
		}
	}
}

func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024) // 1 GB = 1024^3 bytes
}

func formatMemoryUsage(label string, value float64, isPercentage bool) string {
	if isPercentage {
		return fmt.Sprintf("%s: %.1f%%\n\n", label, value)
	}
	// For raw float values (e.g., in gigabytes), you might want to specify units
	return fmt.Sprintf("%s: %.2f GB\n\n", label, value)
}
