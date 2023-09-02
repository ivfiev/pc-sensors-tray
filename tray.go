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
	mQuit := systray.AddMenuItem("Quit", "Quit")
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
				bytes, err := svc.GetIcon(newResult)
				if err != nil {
					log.Fatal(err)
				}
				systray.SetIcon(bytes)

				if len(mEntries) == 0 && len(newResult.Lines()) > 0 {
					for _, line := range newResult.Lines() {
						mEntries = append(mEntries, systray.AddMenuItem(line, line))
					}
				} else if len(newResult.Lines()) > 0 {
					for i, entry := range mEntries {
						entry.SetTitle(newResult.Lines()[i])
						entry.SetTooltip(newResult.Lines()[i])
					}
				}

				lastResult = newResult
				lastUpdate = time.Now()
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {

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
