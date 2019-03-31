// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package postgres

import "time"

func UnixTimeToString(t int64) string {
	return time.Unix(t, 0).UTC().String()
}
