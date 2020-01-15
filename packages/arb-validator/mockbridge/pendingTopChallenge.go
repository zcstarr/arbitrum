/*
 * Copyright 2020, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mockbridge

import (
	"context"
	"math/big"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
)

type PendingTopChallenge struct {
	*BisectionChallenge
}

func NewPendingTopChallenge(address common.Address, client arbbridge.ArbClient) (*PendingTopChallenge, error) {
	bisectionChallenge, err := newBisectionChallenge(address, client)
	if err != nil {
		return nil, err
	}
	vm := &PendingTopChallenge{BisectionChallenge: bisectionChallenge}
	err = vm.setupContracts()
	return vm, err
}

func (c *PendingTopChallenge) setupContracts() error {
	//challengeManagerContract, err := pendingtopchallenge.NewPendingTopChallenge(c.address, c.Client)
	//if err != nil {
	//	return errors2.Wrap(err, "Failed to connect to messagesChallenge")
	//}
	//
	//c.challenge = challengeManagerContract
	return nil
}

func (c *PendingTopChallenge) StartConnection(ctx context.Context, outChan chan arbbridge.Notification, errChan chan error) error {
	if err := c.BisectionChallenge.StartConnection(ctx, outChan, errChan); err != nil {
		return err
	}
	if err := c.setupContracts(); err != nil {
		return err
	}
	//header, err := c.Client.HeaderByNumber(ctx, nil)
	//if err != nil {
	//	return err
	//}
	//
	//filter := ethereum.FilterQuery{
	//	Addresses: []common.Address{c.address},
	//	Topics: [][]common.Hash{{
	//		pendingTopBisectedID,
	//		pendingTopOneStepProofCompletedID,
	//	}},
	//}
	//
	//logs, err := c.Client.FilterLogs(ctx, filter)
	//if err != nil {
	//	return err
	//}
	//for _, log := range logs {
	//	if err := c.processEvents(ctx, log, outChan); err != nil {
	//		return err
	//	}
	//}
	//
	//filter.FromBlock = header.Number
	//logChan := make(chan types.Log)
	//logSub, err := c.Client.SubscribeFilterLogs(ctx, filter, logChan)
	//if err != nil {
	//	return err
	//}
	//
	//go func() {
	//	defer logSub.Unsubscribe()
	//
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			break
	//		case log := <-logChan:
	//			if err := c.processEvents(ctx, log, outChan); err != nil {
	//				errChan <- err
	//				return
	//			}
	//		case err := <-logSub.Err():
	//			errChan <- err
	//			return
	//		}
	//	}
	//}()
	return nil
}

//func (c *PendingTopChallenge) processEvents(ctx context.Context, log types.Log, outChan chan arbbridge.Notification) error {
//	event, err := func() (arbbridge.Event, error) {
//		if log.Topics[0] == pendingTopBisectedID {
//			eventVal, err := c.challenge.ParseBisected(log)
//			if err != nil {
//				return nil, err
//			}
//			return arbbridge.PendingTopBisectionEvent{
//				ChainHashes: eventVal.ChainHashes,
//				TotalLength: eventVal.TotalLength,
//				Deadline:    structures.TimeTicks{Val: eventVal.DeadlineTicks},
//			}, nil
//		} else if log.Topics[0] == pendingTopOneStepProofCompletedID {
//			_, err := c.challenge.ParseOneStepProofCompleted(log)
//			if err != nil {
//				return nil, err
//			}
//			return arbbridge.OneStepProofEvent{}, nil
//		}
//		return nil, errors2.New("unknown arbitrum event type")
//	}()
//
//	if err != nil {
//		return err
//	}
//
//	header, err := c.Client.HeaderByHash(ctx, log.BlockHash)
//	if err != nil {
//		return err
//	}
//	outChan <- arbbridge.Notification{
//		Header: header,
//		VMID:   c.address,
//		Event:  event,
//		TxHash: log.TxHash,
//	}
//	return nil
//}

func (c *PendingTopChallenge) Bisect(
	ctx context.Context,
	chainHashes []common.Hash,
	chainLength *big.Int,
) error {
	//c.auth.Context = ctx
	//tx, err := c.challenge.Bisect(
	//	c.auth,
	//	chainHashes,
	//	chainLength,
	//)
	//if err != nil {
	//	return err
	//}
	//return c.waitForReceipt(ctx, tx, "Bisect")
	return nil
}

func (c *PendingTopChallenge) OneStepProof(
	ctx context.Context,
	lowerHashA common.Hash,
	topHashA common.Hash,
	value common.Hash,
) error {
	//c.auth.Context = ctx
	//tx, err := c.challenge.OneStepProof(
	//	c.auth,
	//	lowerHashA,
	//	topHashA,
	//	value,
	//)
	//if err != nil {
	//	return err
	//}
	//return c.waitForReceipt(ctx, tx, "OneStepProof")
	return nil
}

func (c *PendingTopChallenge) ChooseSegment(
	ctx context.Context,
	assertionToChallenge uint16,
	chainHashes []common.Hash,
	chainLength uint32,
) error {
	return nil
}