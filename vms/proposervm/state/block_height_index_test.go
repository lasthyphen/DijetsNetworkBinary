// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lasthyphen/DijetsNetworkBinary/database/memdb"
	"github.com/lasthyphen/DijetsNetworkBinary/database/versiondb"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/logging"
)

func TestHasIndexReset(t *testing.T) {
	a := assert.New(t)

	db := memdb.New()
	vdb := versiondb.New(db)
	s := New(vdb)
	wasReset, err := s.HasIndexReset()
	a.NoError(err)
	a.False(wasReset)
	err = s.ResetHeightIndex(logging.NoLog{}, vdb)
	a.NoError(err)
	wasReset, err = s.HasIndexReset()
	a.NoError(err)
	a.True(wasReset)
}
