package main

import (
	flag "flag"
	fmt "fmt"
	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	host "github.com/shirou/gopsutil/host"
	log "log"
	time "time"
)

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

// LocalMonitor start local measurement of temperature
func LocalMonitor() {
	if err := termui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	g := widgets.NewGauge()
	g.Title = "CPU temp"
	g.Percent = 0
	g.SetRect(0, 0, 50, 3)
	g.BarColor = termui.ColorGreen
	g.BorderStyle.Fg = termui.ColorWhite
	g.TitleStyle.Fg = termui.ColorWhite

	draw := func(count int) {
		temp := GetTemperature()
		g.Percent = temp
		g.Label = fmt.Sprintf("%v C", temp)
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
