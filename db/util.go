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

package db

import "database/sql"

func toNullInt64(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Valid: true, Int64: int64(*i)}
}

func toNullFloat64(f *float32) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Valid: true, Float64: float64(*f)}
}

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: *s}
}

func fromNullInt64(ni64 sql.NullInt64) *int {
	if ni64.Valid {
		v := int(ni64.Int64)
		return &v
	}
	return nil
}

func fromNullFloat64(nf64 sql.NullFloat64) *float32 {
	if nf64.Valid {
		f32 := float32(nf64.Float64)
		return &f32
	}
	return nil
}
