// Copyright 2015 The go-stormchain Authors
// This file is part of the go-stormchain library.
//
// The go-stormchain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-stormchain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-stormchain library. If not, see <http://www.gnu.org/licenses/>.

// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/filestorm-fst/go-stormchain/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("storm/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("storm/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("storm/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("storm/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("storm/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("storm/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("storm/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("storm/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("storm/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("storm/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("storm/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("storm/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("storm/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("storm/downloader/states/drop", nil)
)
