// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package common

import (
	"time"

	"github.com/lasthyphen/DijetsNetworkBinary/ids"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/constants"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/logging"
)

var (
	_ StateSummaryFrontierHandler = &noOpStateSummaryFrontierHandler{}
	_ AcceptedStateSummaryHandler = &noOpAcceptedStateSummaryHandler{}
	_ AcceptedFrontierHandler     = &noOpAcceptedFrontierHandler{}
	_ AcceptedHandler             = &noOpAcceptedHandler{}
	_ AncestorsHandler            = &noOpAncestorsHandler{}
	_ PutHandler                  = &noOpPutHandler{}
	_ QueryHandler                = &noOpQueryHandler{}
	_ ChitsHandler                = &noOpChitsHandler{}
	_ AppHandler                  = &noOpAppHandler{}
)

type noOpStateSummaryFrontierHandler struct {
	log logging.Logger
}

func NewNoOpStateSummaryFrontierHandler(log logging.Logger) StateSummaryFrontierHandler {
	return &noOpStateSummaryFrontierHandler{log: log}
}

func (nop *noOpStateSummaryFrontierHandler) StateSummaryFrontier(validatorID ids.NodeID, requestID uint32, summary []byte) error {
	nop.log.Debug("StateSummaryFrontier(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

func (nop *noOpStateSummaryFrontierHandler) GetStateSummaryFrontierFailed(validatorID ids.NodeID, requestID uint32) error {
	nop.log.Debug("GetStateSummaryFrontierFailed(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

type noOpAcceptedStateSummaryHandler struct {
	log logging.Logger
}

func NewNoOpAcceptedStateSummaryHandler(log logging.Logger) AcceptedStateSummaryHandler {
	return &noOpAcceptedStateSummaryHandler{log: log}
}

func (nop *noOpAcceptedStateSummaryHandler) AcceptedStateSummary(validatorID ids.NodeID, requestID uint32, summaryIDs []ids.ID) error {
	nop.log.Debug("AcceptedStateSummary(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

func (nop *noOpAcceptedStateSummaryHandler) GetAcceptedStateSummaryFailed(validatorID ids.NodeID, requestID uint32) error {
	nop.log.Debug("GetAcceptedStateSummaryFailed(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

type noOpAcceptedFrontierHandler struct {
	log logging.Logger
}

func NewNoOpAcceptedFrontierHandler(log logging.Logger) AcceptedFrontierHandler {
	return &noOpAcceptedFrontierHandler{log: log}
}

func (nop *noOpAcceptedFrontierHandler) AcceptedFrontier(validatorID ids.NodeID, requestID uint32, containerIDs []ids.ID) error {
	nop.log.Debug("AcceptedFrontier(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

func (nop *noOpAcceptedFrontierHandler) GetAcceptedFrontierFailed(validatorID ids.NodeID, requestID uint32) error {
	nop.log.Debug("GetAcceptedFrontierFailed(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

type noOpAcceptedHandler struct {
	log logging.Logger
}

func NewNoOpAcceptedHandler(log logging.Logger) AcceptedHandler {
	return &noOpAcceptedHandler{log: log}
}

func (nop *noOpAcceptedHandler) Accepted(validatorID ids.NodeID, requestID uint32, containerIDs []ids.ID) error {
	nop.log.Debug("Accepted(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

func (nop *noOpAcceptedHandler) GetAcceptedFailed(validatorID ids.NodeID, requestID uint32) error {
	nop.log.Debug("GetAcceptedFailed(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

type noOpAncestorsHandler struct {
	log logging.Logger
}

func NewNoOpAncestorsHandler(log logging.Logger) AncestorsHandler {
	return &noOpAncestorsHandler{log: log}
}

func (nop *noOpAncestorsHandler) Ancestors(validatorID ids.NodeID, requestID uint32, containers [][]byte) error {
	nop.log.Debug("Ancestors(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

func (nop *noOpAncestorsHandler) GetAncestorsFailed(validatorID ids.NodeID, requestID uint32) error {
	nop.log.Debug("GetAncestorsFailed(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

type noOpPutHandler struct {
	log logging.Logger
}

func NewNoOpPutHandler(log logging.Logger) PutHandler {
	return &noOpPutHandler{log: log}
}

func (nop *noOpPutHandler) Put(vdr ids.NodeID, requestID uint32, blkBytes []byte) error {
	if requestID == constants.GossipMsgRequestID {
		nop.log.Verbo("Gossip Put(%s, %d) unhandled by this gear. Dropped.", vdr, requestID)
	} else {
		nop.log.Debug("Put(%s, %d) unhandled by this gear. Dropped.", vdr, requestID)
	}
	return nil
}

func (nop *noOpPutHandler) GetFailed(validatorID ids.NodeID, requestID uint32) error {
	nop.log.Debug("GetFailed(%s, %d) unhandled by this gear. Dropped.", validatorID, requestID)
	return nil
}

type noOpQueryHandler struct {
	log logging.Logger
}

func NewNoOpQueryHandler(log logging.Logger) QueryHandler {
	return &noOpQueryHandler{log: log}
}

func (nop *noOpQueryHandler) PullQuery(vdr ids.NodeID, requestID uint32, blkID ids.ID) error {
	nop.log.Debug("PullQuery(%s, %d, %s) unhandled by this gear. Dropped.", vdr, requestID, blkID)
	return nil
}

func (nop *noOpQueryHandler) PushQuery(vdr ids.NodeID, requestID uint32, blkBytes []byte) error {
	nop.log.Debug("PushQuery(%s, %d) unhandled by this gear. Dropped.", vdr, requestID)
	return nil
}

type noOpChitsHandler struct {
	log logging.Logger
}

func NewNoOpChitsHandler(log logging.Logger) ChitsHandler {
	return &noOpChitsHandler{log: log}
}

func (nop *noOpChitsHandler) Chits(vdr ids.NodeID, requestID uint32, votes []ids.ID) error {
	nop.log.Debug("Chits(%s, %d) unhandled by this gear. Dropped.", vdr, requestID)
	return nil
}

func (nop *noOpChitsHandler) QueryFailed(vdr ids.NodeID, requestID uint32) error {
	nop.log.Debug("QueryFailed(%s, %d) unhandled by this gear. Dropped.", vdr, requestID)
	return nil
}

type noOpAppHandler struct {
	log logging.Logger
}

func NewNoOpAppHandler(log logging.Logger) AppHandler {
	return &noOpAppHandler{log: log}
}

func (nop *noOpAppHandler) AppRequest(nodeID ids.NodeID, requestID uint32, deadline time.Time, request []byte) error {
	nop.log.Debug("AppRequest(%s, %d) unhandled by this gear. Dropped.", nodeID, requestID)
	return nil
}

func (nop *noOpAppHandler) AppRequestFailed(nodeID ids.NodeID, requestID uint32) error {
	nop.log.Debug("AppRequestFailed(%s, %d) unhandled by this gear. Dropped.", nodeID, requestID)
	return nil
}

func (nop *noOpAppHandler) AppResponse(nodeID ids.NodeID, requestID uint32, response []byte) error {
	nop.log.Debug("AppResponse(%s, %d) unhandled by this gear. Dropped.", nodeID, requestID)
	return nil
}

func (nop *noOpAppHandler) AppGossip(nodeID ids.NodeID, msg []byte) error {
	nop.log.Debug("AppGossip(%s) unhandled by this gear. Dropped.", nodeID)
	return nil
}