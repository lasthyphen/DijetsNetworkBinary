//go:build !linux || !amd64 || !rocksdballowed
// +build !linux !amd64 !rocksdballowed

// ^ Only build this file if this computer is not Linux OR it's not AMD64 OR rocksdb is not allowed
// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.
package rocksdb

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/lasthyphen/DijetsNetworkBinary/database"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/logging"
)

var errUnsupportedDatabase = errors.New("database isn't suppported")

// New returns an error.
func New(string, []byte, logging.Logger, string, prometheus.Registerer) (database.Database, error) {
	return nil, errUnsupportedDatabase
}
