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
	"testing"
)

func TestPM_Valid(t *testing.T) {
	type fields struct {
		Pm25 float32
		Pm10 float32
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "1", fields: fields{Pm25: -1., Pm10: -1.}, want: false},
		{name: "2", fields: fields{Pm25: -1., Pm10: 0}, want: false},
		{name: "3", fields: fields{Pm25: 0., Pm10: 0}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PM{
				Pm25: tt.fields.Pm25,
				Pm10: tt.fields.Pm10,
			}
			if got := pm.Valid(); got != tt.want {
				t.Errorf("PM.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPM_Aqi(t *testing.T) {
	type fields struct {
		Pm25 float32
		Pm10 float32
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{name: "1", fields: fields{Pm25: 1, Pm10: 5}, want: 5},
		{name: "2", fields: fields{Pm25: 9.3, Pm10: 9.3}, want: 39},
		{name: "3", fields: fields{Pm25: 9.3, Pm10: 15}, want: 39},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PM{
				Pm25: tt.fields.Pm25,
				Pm10: tt.fields.Pm10,
			}
			if got := pm.Aqi(); got != tt.want {
				t.Errorf("PM.Aqi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_iaqi(t *testing.T) {
	type args struct {
		c   float32
		bps []float32
		q   float32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "PM2.5 - 1", args: args{c: 1.5, bps: pm25Bps, q: 0.1}, want: 6},
		{name: "PM2.5 - 2", args: args{c: 9.3, bps: pm25Bps, q: 0.1}, want: 39},
		{name: "PM2.5 - 3", args: args{c: 15, bps: pm25Bps, q: 0.1}, want: 57},
		{name: "PM2.5 - 4", args: args{c: 49.5, bps: pm25Bps, q: 0.1}, want: 135},
		{name: "PM2.5 - 5", args: args{c: 235.4, bps: pm25Bps, q: 0.1}, want: 285},
		{name: "PM2.5 - 6", args: args{c: 505, bps: pm25Bps, q: 0.1}, want: 500},

		{name: "PM10 - 1", args: args{c: 10.5, bps: pm10Bps, q: 1.0}, want: 9},
		{name: "PM10 - 2", args: args{c: 21.3, bps: pm10Bps, q: 1.0}, want: 19},
		{name: "PM10 - 3", args: args{c: 107.1, bps: pm10Bps, q: 1.0}, want: 77},
		{name: "PM10 - 4", args: args{c: 256, bps: pm10Bps, q: 1.0}, want: 151},
		{name: "PM10 - 5", args: args{c: 377.8, bps: pm10Bps, q: 1.0}, want: 233},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := iaqi(tt.args.c, tt.args.bps, tt.args.q); got != tt.want {
				t.Errorf("iaqi() = %v, want %v", got, tt.want)
			}
		})
	}
}
