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
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"
	"log"
	"math/big"
	"time"

	"github.com/offchainlabs/arbitrum/packages/arb-util/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
)

var reorgError = errors.New("reorg occured")
var headerRetryDelay = time.Second * 2
var maxFetchAttempts = 5

type MockArbClient struct {
	MockEthClient *mockEthdata
}

func NewEthClient(ethURL string) (*MockArbClient, error) {
	// call to mockEth.go - getMockEth(ethURL)
	return &MockArbClient{getMockEth(ethURL)}, nil
}

func (c *MockArbClient) SubscribeBlockHeaders(ctx context.Context, startBlockId *structures.BlockId) (<-chan arbbridge.MaybeBlockId, error) {
	blockIdChan := make(chan arbbridge.MaybeBlockId, 100)

	blockIdChan <- arbbridge.MaybeBlockId{BlockId: startBlockId}
	prevBlockId := startBlockId
	go func() {
		defer close(blockIdChan)

		for {
			var nextHeader *types.Header
			fetchErrorCount := 0
			for {
				var err error
				//nextHeader, err = c.client.HeaderByNumber(ctx, new(big.Int).Add(prevBlockId.Height.AsInt(), big.NewInt(1)))
				//if err == nil {
				//	// Got next header
				//	break
				//}

				select {
				case <-ctx.Done():
					// Getting header must have failed due to context cancellation
					return
				default:
				}

				if err != nil && err.Error() != ethereum.NotFound.Error() {
					log.Printf("Failed to fetch next header on attempt %v with error: %v", fetchErrorCount, err)
					fetchErrorCount++
				}

				if fetchErrorCount >= maxFetchAttempts {
					blockIdChan <- arbbridge.MaybeBlockId{Err: err}
					return
				}

				// Header was not found so wait before checking again
				time.Sleep(headerRetryDelay)
			}

			if nextHeader.ParentHash != prevBlockId.HeaderHash.ToEthHash() {
				blockIdChan <- arbbridge.MaybeBlockId{Err: reorgError}
				return
			}

			//prevBlockId = getBlockID(nextHeader)
			//blockIdChan <- arbbridge.MaybeBlockId{BlockId: prevBlockId}
		}
	}()

	return blockIdChan, nil
}

func (c *MockArbClient) NewArbFactoryWatcher(address common.Address) (arbbridge.ArbFactoryWatcher, error) {
	return newArbFactoryWatcher(address, c)
}

func (c *MockArbClient) NewRollupWatcher(address common.Address) (arbbridge.ArbRollupWatcher, error) {
	return newRollupWatcher(address, c)
}

func (c *MockArbClient) NewExecutionChallengeWatcher(address common.Address) (arbbridge.ExecutionChallengeWatcher, error) {
	return newExecutionChallengeWatcher(address.ToEthAddress(), c)
}

func (c *MockArbClient) NewMessagesChallengeWatcher(address common.Address) (arbbridge.MessagesChallengeWatcher, error) {
	return newMessagesChallengeWatcher(address.ToEthAddress(), c)
}

func (c *MockArbClient) NewPendingTopChallengeWatcher(address common.Address) (arbbridge.PendingTopChallengeWatcher, error) {
	return newPendingTopChallengeWatcher(address.ToEthAddress(), c)
}

func (c *MockArbClient) NewOneStepProof(address common.Address) (arbbridge.OneStepProof, error) {
	return newOneStepProof(address, c)
}

func (c *MockArbClient) CurrentBlockId(ctx context.Context) (*structures.BlockId, error) {
	return c.MockEthClient.LatestBlock, nil
}

func (c *MockArbClient) BlockIdForHeight(ctx context.Context, height *common.TimeBlocks) (*structures.BlockId, error) {
	block, err := c.MockEthClient.blockNumbers[height]
	if err {
		return nil, errors.New("block height not found")
	}
	return block, nil
}

type TransOpts struct {
	From  common.Address // Ethereum account to send the transaction from
	Nonce *big.Int       // Nonce to use for the transaction execution (nil = use pending state)

	Value    *big.Int // Funds to transfer along along the transaction (nil = 0 = no funds)
	GasPrice *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
	GasLimit uint64   // Gas limit to set for the transaction execution (0 = estimate)
}

type MockArbAuthClient struct {
	*MockArbClient
	auth *TransOpts
}

func NewEthAuthClient(ethURL string, auth *TransOpts) (*MockArbAuthClient, error) {
	client, err := NewEthClient(ethURL)
	if err != nil {
		return nil, err
	}
	return &MockArbAuthClient{
		MockArbClient: client,
		auth:          auth,
	}, nil
}

func (c *MockArbAuthClient) Address() common.Address {
	return c.auth.From
}

func (c *MockArbAuthClient) NewArbFactory(address common.Address) (arbbridge.ArbFactory, error) {
	return newArbFactory(address, c.MockArbClient)
}

func (c *MockArbAuthClient) NewRollup(address common.Address) (arbbridge.ArbRollup, error) {
	return newRollup(address, c)
}

func (c *MockArbAuthClient) NewPendingInbox(address common.Address) (arbbridge.PendingInbox, error) {
	return newPendingInbox(address, c.MockArbClient)
}

func (c *MockArbAuthClient) NewChallengeFactory(address common.Address) (arbbridge.ChallengeFactory, error) {
	return newChallengeFactory(address, c, c.auth)
}

func (c *MockArbAuthClient) NewExecutionChallenge(address common.Address) (arbbridge.ExecutionChallenge, error) {
	return NewExecutionChallenge(address, c)
}

func (c *MockArbAuthClient) NewMessagesChallenge(address common.Address) (arbbridge.MessagesChallenge, error) {
	return newMessagesChallenge(address, c)
}

func (c *MockArbAuthClient) NewPendingTopChallenge(address common.Address) (arbbridge.PendingTopChallenge, error) {
	return NewPendingTopChallenge(address, c)
}

func (c *MockArbAuthClient) DeployChallengeTest(ctx context.Context) (*ChallengeTester, error) {
	//testerAddress, tx, _, err := challengetester.DeployChallengeTester(c.auth, c)
	//if err != nil {
	//	return nil, err
	//}
	//if err := waitForReceipt(
	//	context.Background(),
	//	c.client,
	//	c.auth.From,
	//	tx,
	//	"DeployChallengeTester",
	//); err != nil {
	//	return nil, err
	//}
	tester, err := NewChallengeTester(common.Address{}, c, c.auth)
	if err != nil {
		return nil, err
	}
	return tester, nil
}
