package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	server := gin.Default()

	server.LoadHTMLGlob("static/*.html")

	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Route to serve JSON data
	server.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html", nil)
		data, err := getSystemStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	server.Run(":7000")
}

func getSystemStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get CPU usage
	percentages, err := cpu.Percent(0, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %v", err)
	}
	cpuUsage := make([]string, len(percentages))
	for i, perc := range percentages {
		cpuUsage[i] = fmt.Sprintf("CPU%d: %.1f%%", i, perc)
	}
	stats["cpu"] = cpuUsage

	// Get memory usage
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage: %v", err)
	}

	memoryStats := map[string]interface{}{
		"UsedPercent": memory.UsedPercent,
		"TotalGB":     bytesToGB(memory.Total),
		"AvailableGB": bytesToGB(memory.Available),
		"UsedGB":      bytesToGB(memory.Used),
	}
	stats["memory"] = memoryStats

	return stats, nil
}

func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024) // 1 GB = 1024^3 bytes
}
