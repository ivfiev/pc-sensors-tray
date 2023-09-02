package sensors

import (
	"fmt"
	"github.com/samber/lo"
	"log"
	"math"
	"os/exec"
	"pc-sensors-tray/types"
	"strconv"
	"strings"
)

type CpuFrequencyResult struct {
	minFreq     float64
	maxFreq     float64
	frequencies []float64
}

func (res CpuFrequencyResult) Value() float64 {
	return lo.Max(res.frequencies)
}

func (res CpuFrequencyResult) IsClose(result types.Result) bool {
	return math.Abs(res.Value()-result.Value()) < 100
}

func (res CpuFrequencyResult) Icon() string {
	return fmt.Sprintf("%.1f", res.Value()/1000)
}

func (res CpuFrequencyResult) Lines() []string {
	return lo.Map(res.frequencies, func(freq float64, i int) string {
		return fmt.Sprintf("Core %d -> %.2f GHz", i+1, freq/1000)
	})
}

func (res CpuFrequencyResult) Colour() string {
	return "blue"
}

func (res CpuFrequencyResult) Intensity() float64 {
	return lo.Clamp((res.Value()-res.minFreq)/(res.maxFreq-res.minFreq), 0, 1)
}

func GetCpuFrequencies(minFreq, maxFreq float64) types.Result {
	cmdResult, err := exec.Command("bash", "-c", "cat /proc/cpuinfo | grep MHz").Output()
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(cmdResult), "\n")
	freqs := lo.Map(lines[:24], func(line string, i int) float64 {
		words := strings.Split(line, " ")
		freq, err := strconv.ParseFloat(words[len(words)-1], 64)
		if err != nil {
			log.Fatal(err)
		}
		return freq
	})
	result := CpuFrequencyResult{minFreq, maxFreq, freqs}
	return result
}
