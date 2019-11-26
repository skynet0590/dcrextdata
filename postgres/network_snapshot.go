package postgres

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/raedahgroup/dcrextdata/netsnapshot"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg PgDb) SaveSnapshot(ctx context.Context, snapshot netsnapshot.SnapShot) error {
	existingSnapshot, err := models.FindNetworkSnapshot(ctx, pg.db, snapshot.Timestamp)
	if err == nil {
		existingSnapshot.Nodes = snapshot.Nodes
		existingSnapshot.Height = snapshot.Height
		_, err = existingSnapshot.Update(ctx, pg.db, boil.Infer())
		return err
	}

	snapshotModel := models.NetworkSnapshot{Timestamp: snapshot.Timestamp, Height: snapshot.Height, Nodes: snapshot.Nodes}
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
		Nodes:     snapshotModel.Nodes,
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
	existingNode, err := models.NetworkPeers(models.NetworkPeerWhere.Timestamp.EQ(peer.Timestamp),
		models.NetworkPeerWhere.Address.EQ(peer.Address)).One(ctx, pg.db)
	if err == nil {
		existingNode.LastSeen = peer.LastSeen
		existingNode.ConnectionTime = peer.ConnectionTime
		existingNode.ProtocolVersion = int(peer.ProtocolVersion)
		existingNode.UserAgent = peer.UserAgent
		existingNode.StartingHeight = peer.StartingHeight
		existingNode.CurrentHeight = peer.CurrentHeight

		_, err = existingNode.Update(ctx, pg.db, boil.Infer())

		return err
	}
	peerModel := models.NetworkPeer{
		Timestamp:       peer.Timestamp,
		Address:         peer.Address,
		Country:		 peer.Country,
		IPVersion:       peer.IPVersion,
		LastSeen:        peer.LastSeen,
		ConnectionTime:  peer.ConnectionTime,
		ProtocolVersion: int(peer.ProtocolVersion),
		UserAgent:       peer.UserAgent,
		StartingHeight:  peer.StartingHeight,
		CurrentHeight:   peer.CurrentHeight,
	}

	return peerModel.Insert(ctx, pg.db, boil.Infer())
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
		qm.OrderBy(models.NetworkPeerColumns.LastSeen),
	)
	peerSlice, err := models.NetworkPeers(query...).All(ctx, pg.db)
	if err != nil {
		return nil, 0, err
	}

	var peers []netsnapshot.NetworkPeer
	for _, peerModel := range peerSlice {
		peers = append(peers, netsnapshot.NetworkPeer{
			Timestamp:       peerModel.Timestamp,
			Address:         peerModel.Address,
			Country: 		 peerModel.Country,
			LastSeen:        peerModel.LastSeen,
			ConnectionTime:  peerModel.ConnectionTime,
			ProtocolVersion: uint32(peerModel.ProtocolVersion),
			UserAgent:       peerModel.UserAgent,
			StartingHeight:  peerModel.StartingHeight,
			CurrentHeight:   peerModel.CurrentHeight,
		})
	}

	return peers, totalCount, nil
}

func (pg PgDb) GetIPLocation(ctx context.Context, ip string) (string, error) {
	node, err := models.NetworkPeers(
		models.NetworkPeerWhere.Address.EQ(ip),
	).One(ctx, pg.db)
	if err != nil {
		return "", err
	}

	return node.Country, nil
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

func (pg PgDb) PeerCountByUserAgents(ctx context.Context, timestamp int64, offset int, limit int) (userAgents []netsnapshot.UserAgentInfo, err error) {
	sql := fmt.Sprintf("select %s, count(%s) as number from %s WHERE %s = %d GROUP BY %s ORDER BY number OFFSET %d LIMIT %d",
		models.NetworkPeerColumns.UserAgent, models.NetworkPeerColumns.UserAgent,
		models.TableNames.NetworkPeer, models.NetworkPeerColumns.Timestamp, timestamp,
		models.NetworkPeerColumns.UserAgent, offset, limit)

	var result []struct {
		UserAgent string `json:"user_agent"`
		Number    int64  `json:"number"`
	}

	err = models.NetworkPeers(qm.SQL(sql)).Bind(ctx, pg.db, &result)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, item := range result {
		total += item.Number
	}

	if total == 0 {
		return nil, errors.New("No records found")
	}

	for _, item := range result {
		userAgent := item.UserAgent
		if strings.Trim(userAgent, " ") == "" {
			userAgent = "Unkown"
		}
		userAgents = append(userAgents, netsnapshot.UserAgentInfo{
			UserAgent:  userAgent,
			Nodes:      item.Number,
			Percentage: float64(100 * item.Number / total),
		})
	}

	sort.Slice(userAgents, func(i, j int) bool {
		return userAgents[i].Nodes > userAgents[j].Nodes
	})

	return
}

func (pg PgDb) PeerCountByCountries(ctx context.Context, timestamp int64, offset int, limit int) (countries []netsnapshot.CountryInfo, err error) {
	sql := fmt.Sprintf("select %s, count(%s) as number from %s WHERE %s = %d GROUP BY %s ORDER BY number OFFSET %d LIMIT %d",
		models.NetworkPeerColumns.Country, models.NetworkPeerColumns.Country,
		models.TableNames.NetworkPeer, models.NetworkPeerColumns.Timestamp, timestamp,
		models.NetworkPeerColumns.Country, offset, limit)

	var result []struct {
		Country string `json:"country"`
		Number    int64  `json:"number"`
	}

	err = models.NetworkPeers(qm.SQL(sql)).Bind(ctx, pg.db, &result)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, item := range result {
		total += item.Number
	}

	if total == 0 {
		return nil, errors.New("No records found")
	}

	for _, item := range result {
		country := item.Country
		if strings.Trim(country, " ") == "" {
			country = "Unkown"
		}
		countries = append(countries, netsnapshot.CountryInfo{
			Country:  item.Country,
			Nodes:      item.Number,
			Percentage: float64(100 * item.Number / total),
		})
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Nodes > countries[j].Nodes
	})

	return
}

func (pg PgDb) PeerCountByIPVersion(ctx context.Context, timestamp int64, iPVersion int) (int64, error) {
	return models.NetworkPeers(
		models.NetworkPeerWhere.Timestamp.EQ(timestamp), 
		models.NetworkPeerWhere.IPVersion.EQ(iPVersion),
	).Count(ctx, pg.db)
}

func (pg PgDb) LastSnapshotTime(ctx context.Context) (timestamp int64) {
	_ = pg.LastEntry(ctx, models.TableNames.NetworkSnapshot, &timestamp)
	return
}

func (pg PgDb) LastSnapshot(ctx context.Context) (*netsnapshot.SnapShot, error) {
	return pg.FindNetworkSnapshot(ctx, pg.LastSnapshotTime(ctx))
}
