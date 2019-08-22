package main

import (
	flag "flag"
	fmt "fmt"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	cpu "github.com/shirou/gopsutil/cpu"
	host "github.com/shirou/gopsutil/host"
	load "github.com/shirou/gopsutil/load"
	log "log"
	time "time"
)

// GetCPUInfo gets cpu information
func GetCPUInfo() string {
	info, _ := cpu.Info()
	return info[0].ModelName
}

// GetTemperature return cpu temp
func GetTemperature() int {
	temps, _ := host.SensorsTemperatures()
	for _, sensor := range temps {
		if sensor.SensorKey == "coretemp_packageid0_input" {
			return int(sensor.Temperature)
		}
	}
	return 0
}

// GetColorByTemp gets temp and returns appropriate color scheme
func GetColorByTemp(temp int) termui.Color {
	var color termui.Color
	switch {
	case temp <= 50:
		color = termui.ColorGreen
	case temp > 50 && temp <= 70:
		color = termui.ColorYellow
	case temp > 70:
		color = termui.ColorRed
	}
	return color
}

// GetLoad1 returns load for the past 1 minute
func GetLoad1() float64 {
	info, _ := load.Avg()
	return info.Load1
}

// LocalMonitor start local measurement of temperature
func LocalMonitor() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	cpuInfo := GetCPUInfo()

	var barLength int

	if 2*len(cpuInfo) < 60 {
		barLength = 60
	} else {
		barLength = 2 * len(cpuInfo)
	}

	g := widgets.NewGauge()
	g.SetRect(0, 0, barLength, 3)
	g.BarColor = termui.ColorGreen
	g.BorderStyle.Fg = termui.ColorWhite
	g.TitleStyle.Fg = termui.ColorWhite

	draw := func(count int) {
		temp := GetTemperature()
		g.Percent = temp
		g.Label = fmt.Sprintf("%v C", temp)
		g.Title = fmt.Sprintf("%v Temperature @ Load: %.2f", cpuInfo, GetLoad1())
		g.BarColor = GetColorByTemp(temp)
		termui.Render(g)
	}

	tickerCount := 1
	draw(tickerCount)
	tickerCount++
	uiEvents := termui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			draw(tickerCount)
			tickerCount++
		}
	}
}

func main() {
	monitor := flag.String("monitor", "", "URI of server EX: 127.0.0.1:123")
	flag.Parse()

	if len(*monitor) < 9 {
		LocalMonitor()
	} else {
		fmt.Println("Not implemented")
	}
}
