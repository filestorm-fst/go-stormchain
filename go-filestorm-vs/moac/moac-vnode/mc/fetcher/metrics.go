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

// Contains the metrics collected by the fetcher.

package fetcher

import (
	"github.com/filestorm/go-filestorm/moac/moac-lib/metrics"
)

var (
	propAnnounceInMeter   = metrics.NewMeter("mc/fetcher/prop/announces/in")
	propAnnounceOutTimer  = metrics.NewTimer("mc/fetcher/prop/announces/out")
	propAnnounceDropMeter = metrics.NewMeter("mc/fetcher/prop/announces/drop")
	propAnnounceDOSMeter  = metrics.NewMeter("mc/fetcher/prop/announces/dos")

	propBroadcastInMeter   = metrics.NewMeter("mc/fetcher/prop/broadcasts/in")
	propBroadcastOutTimer  = metrics.NewTimer("mc/fetcher/prop/broadcasts/out")
	propBroadcastDropMeter = metrics.NewMeter("mc/fetcher/prop/broadcasts/drop")
	propBroadcastDOSMeter  = metrics.NewMeter("mc/fetcher/prop/broadcasts/dos")

	headerFetchMeter = metrics.NewMeter("mc/fetcher/fetch/headers")
	bodyFetchMeter   = metrics.NewMeter("mc/fetcher/fetch/bodies")

	headerFilterInMeter  = metrics.NewMeter("mc/fetcher/filter/headers/in")
	headerFilterOutMeter = metrics.NewMeter("mc/fetcher/filter/headers/out")
	bodyFilterInMeter    = metrics.NewMeter("mc/fetcher/filter/bodies/in")
	bodyFilterOutMeter   = metrics.NewMeter("mc/fetcher/filter/bodies/out")
)
