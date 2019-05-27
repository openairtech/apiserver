// Copyright © 2019 Victor Antonovich <victor@antonovich.me>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aqi

import (
	"math"
	"sync"
)

var (
	aqiVals = []float32{0, 51, 101, 151, 201, 301, 401, 500}
	pm25Bps = []float32{0, 12.1, 35.5, 55.5, 150.5, 250.5, 350.5, 500}
	pm10Bps = []float32{0, 55, 155, 255, 355, 425, 505, 605}
)

type PM struct {
	sync.RWMutex
	Pm25, Pm10 float32
}

func (pm *PM) Valid() bool {
	pm.Lock()
	defer pm.Unlock()
	return pm.Pm25 >= 0 && pm.Pm10 >= 0
}

func (pm *PM) Aqi() int {
	pm.Lock()
	defer pm.Unlock()
	iaqi25 := iaqi(pm.Pm25, pm25Bps, 0.1)
	iaqi10 := iaqi(pm.Pm10, pm10Bps, 1.0)
	if iaqi10 > iaqi25 {
		return iaqi10
	}
	return iaqi25
}

func iaqi(c float32, bps []float32, q float32) int {
	c = float32(math.Floor(float64(c/q))) * q
	bp := breakpoint(bps, c)
	bpLo, bpHi, aqiLo, aqiHi := bps[bp], bps[bp+1]-q, aqiVals[bp], aqiVals[bp+1]-1
	return int(math.Round(float64(constrain(linear(c, bpLo, bpHi, aqiLo, aqiHi), aqiVals[0], aqiVals[len(aqiVals)-1]))))
}

func breakpoint(bps []float32, val float32) int {
	var bp int
	for i := range bps[:len(bps)-1] {
		bp = i
		if bps[i] <= val && val < bps[i+1] {
			break
		}
	}
	return bp
}

func linear(value, fromLow, fromHigh, toLow, toHigh float32) float32 {
	return ((value-fromLow)/(fromHigh-fromLow))*(toHigh-toLow) + toLow
}

func constrain(x, a, b float32) float32 {
	if x < a {
		return a
	} else if x > b {
		return b
	}
	return x
}
