package main

import (
	"github.com/getlantern/systray"
	"log"
	"os"
	"pc-sensors-tray/icon"
	"pc-sensors-tray/sensors"
	"pc-sensors-tray/types"
	"strconv"
	"time"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	args := os.Args[1:]
	reader := getSensorReader(args)
	svc := icon.NewIconService()
	mRefresh := systray.AddMenuItem("Refresh", "Refresh")
	mEntries := make([]*systray.MenuItem, 0)
	ticker := time.NewTicker(3 * time.Second)
	var lastResult types.Result
	var lastUpdate time.Time

	for {
		select {
		case <-ticker.C:
			newResult := reader()
			changeSignificant := lastResult == nil || !lastResult.IsClose(newResult)
			noUpdateForAWhile := time.Now().UnixMilli()-lastUpdate.UnixMilli() > 60_000
			needsUpdate := changeSignificant || noUpdateForAWhile

			if needsUpdate {
				update(svc, newResult, &mEntries)
				lastResult = newResult
				lastUpdate = time.Now()
			}
		case <-mRefresh.ClickedCh:
			newResult := reader()
			update(svc, newResult, &mEntries)
		}
	}
}

func onExit() {

}

func update(svc icon.IconService, result types.Result, mEntries *[]*systray.MenuItem) {
	bytes, err := svc.GetIcon(result)
	if err != nil {
		log.Fatal(err)
	}
	systray.SetIcon(bytes)

	lines := result.Lines()
	if len(*mEntries) == 0 && len(lines) > 0 {
		for _, line := range lines {
			*mEntries = append(*mEntries, systray.AddMenuItem(line, line))
		}
	} else if len(lines) > 0 {
		for i, entry := range *mEntries {
			entry.SetTitle(lines[i])
			entry.SetTooltip(lines[i])
		}
	}
}

func getSensorReader(args []string) func() types.Result {
	switch args[0] {
	case "cpu-freq":
		{
			minFreq, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				log.Fatal("Bad min frequency!")
			}
			maxFreq, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				log.Fatal("Bad max frequency!")
			}
			return func() types.Result {
				return sensors.GetCpuFrequencies(minFreq, maxFreq)
			}
		}
	default:
		return func() types.Result {
			return sensors.GetTemps(args[0], args[1:])
		}
	}
}
