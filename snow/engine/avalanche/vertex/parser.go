// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.

package vertex

import (
	"github.com/lasthyphen/DijetsNetworkBinary/snow/consensus/avalanche"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/hashing"
)

// Parser parses bytes into a vertex.
type Parser interface {
	// Parse a vertex from a slice of bytes
	ParseVtx(vertex []byte) (avalanche.Vertex, error)
}

// Parse parses the provided vertex bytes into a stateless vertex
func Parse(bytes []byte) (StatelessVertex, error) {
	vtx := innerStatelessVertex{}
	version, err := c.Unmarshal(bytes, &vtx)
	if err != nil {
		return nil, err
	}
	vtx.Version = version

	return statelessVertex{
		innerStatelessVertex: vtx,
		id:                   hashing.ComputeHash256Array(bytes),
		bytes:                bytes,
	}, nil
}
