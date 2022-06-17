// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"time"

	"github.com/lasthyphen/DijetsNetworkBinary/ids"
	"github.com/lasthyphen/DijetsNetworkBinary/message"
	"github.com/lasthyphen/DijetsNetworkBinary/network/throttling"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/networking/router"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/networking/tracker"
	"github.com/lasthyphen/DijetsNetworkBinary/snow/validators"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/logging"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/timer/mockable"
	"github.com/lasthyphen/DijetsNetworkBinary/version"
)

type Config struct {
	// Size, in bytes, of the buffer this peer reads messages into
	ReadBufferSize int
	// Size, in bytes, of the buffer this peer writes messages into
	WriteBufferSize      int
	Clock                mockable.Clock
	Metrics              *Metrics
	MessageCreator       message.Creator
	Log                  logging.Logger
	InboundMsgThrottler  throttling.InboundMsgThrottler
	Network              Network
	Router               router.InboundHandler
	VersionCompatibility version.Compatibility
	VersionParser        version.ApplicationParser
	MySubnets            ids.Set
	Beacons              validators.Set
	NetworkID            uint32
	PingFrequency        time.Duration
	PongTimeout          time.Duration
	MaxClockDifference   time.Duration

	// Unix time of the last message sent and received respectively
	// Must only be accessed atomically
	LastSent, LastReceived int64

	// Tracks CPU/disk usage caused by each peer.
	ResourceTracker tracker.ResourceTracker

	PingMessage message.OutboundMessage
}
