// Copyright 2019 The go-filestorm Authors
// This file is part of the go-filestorm library.
//
// The go-filestorm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-filestorm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-filestorm library. If not, see <http://www.gnu.org/licenses/>.

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/filestorm/go-filestorm/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("fst/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("fst/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("fst/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("fst/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("fst/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("fst/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("fst/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("fst/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("fst/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("fst/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("fst/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("fst/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("fst/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("fst/downloader/states/drop", nil)
)
