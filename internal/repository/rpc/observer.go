package rpc

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// kHeadsObserverSubscribeTick represents the time between subscription attempts.
const kHeadsObserverSubscribeTick = 5 * time.Second

// kHeadObserverRpcTick represents the time between rpc calls
const kHeadObserverRpcTick = 500 * time.Millisecond

// kObservedHeadChanCapacity represents the capacity of the channel fed with new headers.
const kObservedHeadChanCapacity = 10_000

// ObservedHeadProxy provides a channel fed with new headers.
// It is not guaranteed that there won't be "gaps" between the headers.
// This might happen when subscription fails, and we have to "simulate" it via rpc calls.
func (rpc *OperaRpc) ObservedHeadProxy() <-chan *types.Header {
	// If the channel is nil, initialize it.
	if rpc.headers == nil {
		rpc.sigClose = make(chan struct{}, 1)
		rpc.headers = make(chan *types.Header, kObservedHeadChanCapacity)
		rpc.wg.Add(1)
		go rpc.observeBlocks()
	}
	return rpc.headers
}

// observeBlocks observes new blocks and sends them to the channel.
func (rpc *OperaRpc) observeBlocks() {
	var sub ethereum.Subscription
	defer func() {
		if sub != nil {
			sub.Unsubscribe()
		}
		rpc.wg.Done()
	}()

	sub = rpc.blockSubscription()

	// failed to subscribe, "simulate" subscription via rpc calls
	if sub == nil {
		tm := time.NewTicker(kHeadObserverRpcTick)
		ethClient := ethclient.NewClient(rpc.ftm)
		latestNumber := big.NewInt(0)
		for {
			select {
			case <-rpc.sigClose:
				return
			case <-tm.C:
				h, err := ethClient.HeaderByNumber(context.Background(), nil)
				if err != nil || h.Number.Cmp(latestNumber) <= 0 {
					continue
				}
				latestNumber.Set(h.Number)
				rpc.headers <- h
			}
		}
	}

	// otherwise go standard path
	for {
		// re-subscribe if the subscription ref is not valid
		if sub == nil {
			tm := time.NewTimer(kHeadsObserverSubscribeTick)
			select {
			case <-rpc.sigClose:
				return
			case <-tm.C:
				sub = rpc.blockSubscription()
				continue
			}
		}
		// use the subscriptions
		select {
		case <-rpc.sigClose:
			return
		case <-sub.Err():
			sub = nil
		}
	}
}

// blockSubscription provides a subscription for new blocks received
// by the connected blockchain node.
func (rpc *OperaRpc) blockSubscription() ethereum.Subscription {
	sub, err := rpc.ftm.EthSubscribe(context.Background(), rpc.headers, "newHeads")
	if err != nil {
		return nil
	}
	return sub
}
