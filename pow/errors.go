// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package pow

type NilClientError struct{}

func (err NilClientError) Error() string {
	return "Cannot use a nil http client"
}
