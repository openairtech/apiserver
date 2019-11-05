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

package util

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// CloseQuietly closes given closer without error checking.
func CloseQuietly(closer io.Closer) {
	_ = closer.Close()
}

// ParseBBox parses given bounding box in comma-delimited string format, like "44.43,48.65,44.53,48.7".
func ParseBBox(sbb string) ([]float64, error) {
	if sbb == "" {
		return nil, nil
	}

	sbbvs := strings.Split(sbb, ",")

	if len(sbbvs) < 4 {
		return nil, errors.New(fmt.Sprintf("invalid bounding box: [%s]", sbb))
	}

	bbvs := make([]float64, 4)

	for i, sbbv := range sbbvs {
		if i >= 4 {
			break
		}

		bbv, err := strconv.ParseFloat(sbbv, 64)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("invalid bounding box [%s] value: %s", sbb, sbbv))
		}

		bbvs[i] = bbv
	}

	return bbvs, nil
}

// ParseDuration parses given duration string ds into duration.
// It returns parsed duration value or nil for empty string, and error if duration string is invalid.
func ParseDuration(ds string) (*time.Duration, error) {
	if ds == "" {
		return nil, nil
	}

	d, err := time.ParseDuration(ds)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ParseUnixTime parses given Unix time string uts into time.
// It returns parsed time value or nil for empty string, and error if given string can't be parsed as an Unix time.
func ParseUnixTime(uts string) (*time.Time, error) {
	if uts == "" {
		return nil, nil
	}

	ut, err := strconv.ParseInt(uts, 10, 64)
	if err != nil {
		return nil, err
	}

	t := time.Unix(ut, 0)

	return &t, nil
}

// StringInSlice checks string s is in slice list.
func StringInSlice(s string, list []interface{}) bool {
	for _, l := range list {
		if l == s {
			return true
		}
	}
	return false
}
