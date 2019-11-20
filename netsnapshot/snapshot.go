package netsnapshot

import (
	"context"
	"github.com/decred/dcrd/rpcclient"
)

func NewTake(dcrClient *rpcclient.Client, store DataStore) *taker {
	return &taker{
		dcrClient: dcrClient,
		dataStore: store,
	}
}

func (t taker) TakeSnapshot(ctx context.Context) {
	peerInfo, err := t.dcrClient.GetPeerInfo()
	if err != nil {
		log.Errorf("Error in getting peer info, %s", err.Error())
		return
	}

	for _, peer := range peerInfo {
		err := t.dataStore.SaveNetworkPeer(ctx, NetworkPeer{
			Address:         peer.Addr,
			LastReceiveTime: peer.LastRecv,
			LastSendTime:    peer.LastSend,
			ConnectionTime:  peer.ConnTime,
			ProtocolVersion: peer.Version,
			UserAgent:       peer.SubVer,
			StartingHeight:  peer.StartingHeight,
			CurrentHeight:   peer.CurrentHeight,
		})
		if err != nil {
			log.Errorf("Error in saving peer info, %s", err.Error())
		}
	}
}
