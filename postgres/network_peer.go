package postgres

import (
	"context"
	"fmt"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/raedahgroup/dcrextdata/netsnapshot"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
)

func (pg PgDb) SaveNetworkPeer(ctx context.Context, peer netsnapshot.NetworkPeer) error {
	peerModel := models.NetworkPeer{
		ID:              0,
		Address:         peer.Address,
		LastReceiveTime: peer.LastReceiveTime,
		LastSendTime:    peer.LastSendTime,
		ConnectionTime:  peer.ConnectionTime,
		ProtocolVersion: peer.ProtocolVersion,
		UserAgent:       peer.UserAgent,
		StartingHeight:  peer.StartingHeight,
		CurrentHeight:   peer.CurrentHeight,
	}

	if existingPeer, err := models.NetworkPeers(models.NetworkPeerWhere.Address.EQ(peer.Address)).One(ctx, pg.db); err == nil {
		if _, err = existingPeer.Delete(ctx, pg.db); err != nil {
			return fmt.Errorf("new peer not added, error in deleting existing duplicate address, %s", err.Error())
		}
	}

	return peerModel.Insert(ctx, pg.db, boil.Infer())
}

func (pg PgDb) NetworkPeers(ctx context.Context, q string, offset int, limit int) ([]netsnapshot.NetworkPeer, int64, error) {
	var query []qm.QueryMod
	if q != "" {
		query = append(query,
			models.NetworkPeerWhere.Address.EQ(q),
			qm.Or2(models.NetworkPeerWhere.UserAgent.EQ(q)),
		)
	}

	totalCount, err := models.NetworkPeers(query...).Count(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	query = append(query, qm.Limit(limit), qm.Offset(offset))
	peerSlice, err := models.NetworkPeers(query...).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	var peers []netsnapshot.NetworkPeer
	for _, peerModel := range peerSlice {
		peers = append(peers, netsnapshot.NetworkPeer{
			ID:              peerModel.ID,
			Address:         peerModel.Address,
			LastReceiveTime: peerModel.LastReceiveTime,
			LastSendTime:    peerModel.LastSendTime,
			ConnectionTime:  peerModel.ConnectionTime,
			ProtocolVersion: peerModel.ProtocolVersion,
			UserAgent:       peerModel.UserAgent,
			StartingHeight:  peerModel.StartingHeight,
			CurrentHeight:   peerModel.CurrentHeight,
		})
	}

	return peers, totalCount, nil
}

func (pg PgDb) TotalPeerCount(ctx context.Context) (int64, error) {
	return models.NetworkPeers().Count(ctx, pg.db)
}

func (pg PgDb) TotalPeerCountByProtocol(ctx context.Context, protocolVersion int) (int64, error) {
	return models.NetworkPeers(models.NetworkPeerWhere.ProtocolVersion.EQ(protocolVersion)).Count(ctx, pg.db)
}

func (pg PgDb) PeerCountByUserAgents(ctx context.Context) (counts map[string]int64, err error) {
	sql := fmt.Sprintf("select %s, count(%s) as no from %s GROUP BY %s",
		models.NetworkPeerColumns.UserAgent, models.NetworkPeerColumns.UserAgent,
		models.TableNames.NetworkPeer, models.NetworkPeerColumns.UserAgent)

	var result []struct {
		UserAgent string `json:"user_agent"`
		Number    int64    `json:"number"`
	}
	err = models.NetworkPeers(qm.SQL(sql)).Bind(ctx, pg.db, &result)
	counts = map[string]int64{}
	for _, item := range result {
		counts[item.UserAgent] = item.Number
	}
	return
}
