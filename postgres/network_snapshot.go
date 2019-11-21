package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrextdata/netsnapshot"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg PgDb) SaveSnapshot(ctx context.Context, snapshot netsnapshot.SnapShot) error {
	snapshotModel := models.NetworkSnapshot{Timestamp:snapshot.Timestamp, Height:snapshot.Height, Nodes: snapshot.Nodes}
	if err := snapshotModel.Insert(ctx, pg.db, boil.Infer()); err != nil {
		if !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
			return err
		}
	}

	return nil
}

func (pg PgDb) FindNetworkSnapshot(ctx context.Context, timestamp int64) (*netsnapshot.SnapShot, error) {
	snapshotModel, err := models.FindNetworkSnapshot(ctx, pg.db, timestamp)
	if err != nil {
		return nil, err
	}
	return &netsnapshot.SnapShot{
		Timestamp: snapshotModel.Timestamp,
		Height:    snapshotModel.Height,
	}, nil
}

func (pg PgDb) PreviousSnapshot(ctx context.Context, timestamp int64) (*netsnapshot.SnapShot, error) {
	snapshotModel, err := models.NetworkSnapshots(
		models.NetworkSnapshotWhere.Timestamp.LT(timestamp),
		qm.OrderBy(fmt.Sprintf("%s DESC", models.NetworkSnapshotColumns.Timestamp)),
		qm.Limit(1),
	).One(ctx, pg.db)

	if err != nil {
		return nil, err
	}
	snapshot := netsnapshot.SnapShot{
		Timestamp: snapshotModel.Timestamp,
		Height:    snapshotModel.Height,
	}

	return &snapshot, err
}

func (pg PgDb) NextSnapshot(ctx context.Context, timestamp int64) (*netsnapshot.SnapShot, error) {
	snapshotModel, err := models.NetworkSnapshots(
		models.NetworkSnapshotWhere.Timestamp.GT(timestamp),
		qm.OrderBy(models.NetworkSnapshotColumns.Timestamp),
		qm.Limit(1),
	).One(ctx, pg.db)

	if err != nil {
		return nil, err
	}
	snapshot := netsnapshot.SnapShot{
		Timestamp: snapshotModel.Timestamp,
		Height:    snapshotModel.Height,
	}

	return &snapshot, err
}

func (pg PgDb) DeleteSnapshot(ctx context.Context, timestamp int64) {
	snapshot, err := models.FindNetworkSnapshot(ctx, pg.db, timestamp)
	if err == nil {
		_, _ = snapshot.Delete(ctx, pg.db)
	}
	_, _ = models.NetworkPeers(models.NetworkPeerWhere.Timestamp.EQ(timestamp)).DeleteAll(ctx, pg.db)
}

func (pg PgDb) SaveNetworkPeer(ctx context.Context, peer netsnapshot.NetworkPeer) error {
	peerModel := models.NetworkPeer{
		Timestamp:       peer.Timestamp,
		Address:         peer.Address,
		LastReceiveTime: peer.LastReceiveTime,
		LastSendTime:    peer.LastSendTime,
		ConnectionTime:  peer.ConnectionTime,
		ProtocolVersion: int(peer.ProtocolVersion),
		UserAgent:       peer.UserAgent,
		StartingHeight:  peer.StartingHeight,
		CurrentHeight:   peer.CurrentHeight,
	}

	err := peerModel.Insert(ctx, pg.db, boil.Infer())
	if err != nil && !strings.Contains(err.Error(), "unique constraint") { // Ignore duplicate entries
		return err
	}

	return nil
}

func (pg PgDb) NetworkPeers(ctx context.Context, timestamp int64, q string, offset int, limit int) ([]netsnapshot.NetworkPeer, int64, error) {
	query := []qm.QueryMod{models.NetworkPeerWhere.Timestamp.EQ(timestamp)}
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

	query = append(query,
		qm.Limit(limit),
		qm.Offset(offset),
		qm.OrderBy(models.NetworkPeerColumns.LastReceiveTime),
		qm.OrderBy(models.NetworkPeerColumns.LastSendTime),
	)
	peerSlice, err := models.NetworkPeers(query...).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	var peers []netsnapshot.NetworkPeer
	for _, peerModel := range peerSlice {
		peers = append(peers, netsnapshot.NetworkPeer{
			Timestamp: 		 peerModel.Timestamp,
			Address:         peerModel.Address,
			LastReceiveTime: peerModel.LastReceiveTime,
			LastSendTime:    peerModel.LastSendTime,
			ConnectionTime:  peerModel.ConnectionTime,
			ProtocolVersion: uint32(peerModel.ProtocolVersion),
			UserAgent:       peerModel.UserAgent,
			StartingHeight:  peerModel.StartingHeight,
			CurrentHeight:   peerModel.CurrentHeight,
		})
	}

	return peers, totalCount, nil
}

func (pg PgDb) TotalPeerCount(ctx context.Context, timestamp int64) (int64, error) {
	return models.NetworkPeers(models.NetworkPeerWhere.Timestamp.EQ(timestamp)).Count(ctx, pg.db)
}

func (pg PgDb) TotalPeerCountByProtocol(ctx context.Context, timestamp int64, protocolVersion int) (int64, error) {
	return models.NetworkPeers(
		models.NetworkPeerWhere.Timestamp.EQ(timestamp),
		models.NetworkPeerWhere.ProtocolVersion.EQ(protocolVersion),
	).Count(ctx, pg.db)
}

func (pg PgDb) PeerCountByUserAgents(ctx context.Context, timestamp int64) (counts map[string]int64, err error) {
	sql := fmt.Sprintf("select %s, count(%s) as no from %s WHERE %s = %d GROUP BY %s",
		models.NetworkPeerColumns.UserAgent, models.NetworkPeerColumns.UserAgent,
		models.TableNames.NetworkPeer, models.NetworkPeerColumns.Timestamp, timestamp,
		models.NetworkPeerColumns.UserAgent)

	var result []struct {
		UserAgent string `json:"user_agent"`
		Number    int64  `json:"number"`
	}

	err = models.NetworkPeers(qm.SQL(sql)).Bind(ctx, pg.db, &result)
	counts = map[string]int64{}
	for _, item := range result {
		counts[item.UserAgent] = item.Number
	}
	return
}

func (pg PgDb) LastSnapshotTime(ctx context.Context) (timestamp int64) {
	_ = pg.LastEntry(ctx, models.TableNames.NetworkSnapshot, &timestamp)
	return
}
