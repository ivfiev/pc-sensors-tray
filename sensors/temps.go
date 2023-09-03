package sensors

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log"
	"math"
	"os/exec"
	"pc-sensors-tray/types"
	"slices"
	"strings"
)

type TemperaturesResult struct {
	deviceId     string
	temperatures map[string]float64
	watts        map[string]float64
	fanRpm       map[string]int
}

func (res TemperaturesResult) Value() float64 {
	return lo.Max(lo.Values(res.temperatures))
}

func (res TemperaturesResult) IsClose(result types.Result) bool {
	return math.Abs(res.Value()-result.Value()) < 2
}

func (res TemperaturesResult) Icon() string {
	return fmt.Sprintf("%d°", int(res.Value()))
}

func (res TemperaturesResult) Lines() []string {
	lines := []string{res.deviceId}
	temps := lo.Map(lo.Keys(res.temperatures), func(key string, i int) string {
		return fmt.Sprintf("%s -> %d°", key, int(res.temperatures[key]))
	})
	slices.Sort(temps)
	watts := lo.Map(lo.Keys(res.watts), func(key string, i int) string {
		return fmt.Sprintf("%s -> %dW", key, int(res.watts[key]))
	})
	slices.Sort(watts)
	fanRpm := lo.Map(lo.Keys(res.fanRpm), func(key string, i int) string {
		return fmt.Sprintf("%s -> %dRPM", key, int(res.fanRpm[key]))
	})
	slices.Sort(fanRpm)
	lines = append(lines, temps...)
	lines = append(lines, watts...)
	lines = append(lines, fanRpm...)
	return lines
}

func (res TemperaturesResult) Colour() string {
	return "red"
}

func (res TemperaturesResult) Intensity() float64 {
	return lo.Clamp((res.Value()-20)/(95-20), 0, 1)
}

func GetTemps(deviceId string, paths []string) TemperaturesResult {
	result, err := exec.Command("sensors", "-j", deviceId).Output()
	if err != nil {
		log.Fatal(err)
	}
	obj := make(map[string]interface{})
	err = json.Unmarshal(result, &obj)
	if err != nil {
		log.Fatal(err)
	}
	obj = obj[deviceId].(map[string]interface{})
	temps := make(map[string]float64)
	watts := make(map[string]float64)
	fanRpm := make(map[string]int)
	for _, path := range paths {
		steps := strings.Split(path, ".")
		val := math.Round(getPathTemp(obj, steps))
		if strings.Contains(path, "power") {
			watts[steps[0]] = val
		} else if strings.Contains(path, "fan") {
			fanRpm[steps[0]] = int(val)
		} else {
			temps[steps[0]] = val
		}
	}
	return TemperaturesResult{deviceId, temps, watts, fanRpm}
}

func getPathTemp(obj map[string]interface{}, steps []string) float64 {
	for i := 0; i < len(steps)-1; i++ {
		obj = obj[steps[i]].(map[string]interface{})
	}
	return obj[steps[len(steps)-1]].(float64)
}
