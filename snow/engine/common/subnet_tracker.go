// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.

package common

import (
	"github.com/lasthyphen/DijetsNetworkBinary/ids"
)

// SubnetTracker describes the interface for checking if a node is tracking a
// subnet
type SubnetTracker interface {
	// TracksSubnet returns true if [nodeID] tracks [subnetID]
	TracksSubnet(nodeID ids.NodeID, subnetID ids.ID) bool
}
