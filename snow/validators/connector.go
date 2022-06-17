// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"github.com/lasthyphen/DijetsNetworkBinary/ids"
	"github.com/lasthyphen/DijetsNetworkBinary/version"
)

// Connector represents a handler that is called when a connection is marked as
// connected or disconnected
type Connector interface {
	Connected(id ids.NodeID, nodeVersion version.Application) error
	Disconnected(id ids.NodeID) error
}
