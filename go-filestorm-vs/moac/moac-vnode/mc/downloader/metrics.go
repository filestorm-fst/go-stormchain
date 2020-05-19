// Copyright 2015 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/filestorm/go-filestorm/moac/moac-lib/metrics"
)

var (
	headerInMeter      = metrics.NewMeter("mc/downloader/headers/in")
	headerReqTimer     = metrics.NewTimer("mc/downloader/headers/req")
	headerDropMeter    = metrics.NewMeter("mc/downloader/headers/drop")
	headerTimeoutMeter = metrics.NewMeter("mc/downloader/headers/timeout")

	bodyInMeter      = metrics.NewMeter("mc/downloader/bodies/in")
	bodyReqTimer     = metrics.NewTimer("mc/downloader/bodies/req")
	bodyDropMeter    = metrics.NewMeter("mc/downloader/bodies/drop")
	bodyTimeoutMeter = metrics.NewMeter("mc/downloader/bodies/timeout")

	receiptInMeter      = metrics.NewMeter("mc/downloader/receipts/in")
	receiptReqTimer     = metrics.NewTimer("mc/downloader/receipts/req")
	receiptDropMeter    = metrics.NewMeter("mc/downloader/receipts/drop")
	receiptTimeoutMeter = metrics.NewMeter("mc/downloader/receipts/timeout")

	stateInMeter   = metrics.NewMeter("mc/downloader/states/in")
	stateDropMeter = metrics.NewMeter("mc/downloader/states/drop")
)
