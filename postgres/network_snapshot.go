package postgres

import (
	"context"
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
	_, _ = models.Heartbeats(models.HeartbeatWhere.Timestamp.EQ(timestamp)).DeleteAll(ctx, pg.db)
}

func (pg PgDb) SaveNetworkPeer(ctx context.Context, peer netsnapshot.NetworkPeer) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return fmt.Errorf("error in starting db transaction %s", err.Error())
	}
	existingNode, err := models.Nodes(models.NodeWhere.Address.EQ(peer.Address)).One(ctx, tx)
	if err == nil {
		existingNode.LastSeen = peer.LastSeen
		existingNode.LastAttempt = peer.LastSeen
		if existingNode.Services == "" {
			existingNode.Services = peer.Services
		}

		if existingNode.StartingHeight == 0 {
			existingNode.StartingHeight = peer.StartingHeight
		}

		if existingNode.UserAgent == "" {
			existingNode.UserAgent = peer.UserAgent
		}

		if peer.CurrentHeight > 0 {
			existingNode.CurrentHeight = peer.CurrentHeight
		}

		if existingNode.ConnectionTime == 0 {
			existingNode.ConnectionTime = peer.ConnectionTime
		}

		if _, err = existingNode.Update(ctx, tx, boil.Infer()); err != nil {
			return fmt.Errorf("error in updating existing node information %s", err.Error())
		}
	} else {
		newNode := models.Node{
			Address:         peer.Address,
			IPVersion:       peer.IPVersion,
			Country:         peer.Country,
			State:           "",
			City:            "",
			Locality:        "",
			LastAttempt:     peer.LastSeen,
			LastSeen:        peer.LastSeen,
			ConnectionTime:  peer.ConnectionTime,
			ProtocolVersion: int(peer.ProtocolVersion),
			UserAgent:       peer.UserAgent,
			Services:        peer.Services,
			StartingHeight:  peer.StartingHeight,
			CurrentHeight:   peer.CurrentHeight,
			IsDead:          false,
		}
		if err := newNode.Insert(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("cannot save heartbeat, error in saving new node information, %s", err.Error())
		}
	}
	// TODO: don't save heartbeat that is not within this snapshot timestamp
	// TODO: check with @raedah
	heartbeat, err := models.Heartbeats(
		models.HeartbeatWhere.NodeID.EQ(peer.Address),
		models.HeartbeatWhere.Timestamp.EQ(peer.Timestamp)).One(ctx, tx)
	if err == nil {
		if peer.CurrentHeight > 0 {
			heartbeat.CurrentHeight = peer.CurrentHeight
		}

		if peer.Latency > 0 {
			heartbeat.Latency = peer.Latency
		}

		if peer.LastSeen > 0 {
			heartbeat.LastSeen = peer.LastSeen
		}

		if _, err = heartbeat.Update(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error in saving heartbeat, %s", err.Error())
		}
		return nil
	}

	newHeartbeat := models.Heartbeat{
		Timestamp:       peer.Timestamp,
		NodeID:         peer.Address,
		LastSeen:        peer.LastSeen,
		Latency:         peer.Latency,
		CurrentHeight:   peer.CurrentHeight,
	}

	if err = newHeartbeat.Insert(ctx, tx, boil.Infer()); err != nil {
		return fmt.Errorf("error in saving hearbeat, %s", err.Error())
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error in saving heartbeat, %s", err.Error())
	}
	return nil
}

func (pg PgDb) NetworkPeers(ctx context.Context, timestamp int64, q string, offset int, limit int) ([]netsnapshot.NetworkPeer, int64, error) {
	where := fmt.Sprintf("heartbeat.timestamp = %d", timestamp)
	args := []interface{}{timestamp}
	if q != "" {
		where += fmt.Sprintf(" AND (node.address = '%s' OR node.user_agent = '%s' OR node.country = '%s')", q, q, q)
		args = append(args, q, q, q)
	}

	sql := `SELECT node.address, node.country, node.last_seen, node.connection_time, node.protocol_version,
			node.user_agent, node.starting_height, node.current_height, node.services FROM heartbeat 
			INNER JOIN node on node.address = heartbeat.node_id WHERE ` + where + fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	var peerSlice models.NodeSlice
	err := models.NewQuery(qm.SQL(sql)).Bind(ctx, pg.db, &peerSlice)
	if err != nil {
		return nil, 0, fmt.Errorf("error %s, on query %s", err.Error(), sql)
	}

	var peers []netsnapshot.NetworkPeer
	for _, node := range peerSlice {
		peers = append(peers, netsnapshot.NetworkPeer{
			Address:         node.Address,
			Country: 		 node.Country,
			LastSeen:        node.LastSeen,
			ConnectionTime:  node.ConnectionTime,
			ProtocolVersion: uint32(node.ProtocolVersion),
			UserAgent:       node.UserAgent,
			StartingHeight:  node.StartingHeight,
			CurrentHeight:   node.CurrentHeight,
			Services: 		 node.Services,
		})
	}

	sql = "SELECT COUNT(heartbeat.node_id) as total FROM heartbeat INNER JOIN node on node.address = heartbeat.node_id WHERE " + where
	var countResult struct{Total int64}
	err = models.NewQuery(qm.SQL(sql)).Bind(ctx, pg.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	return peers, countResult.Total, nil
}

func (pg PgDb) GetIPLocation(ctx context.Context, ip string) (string, int, error) {
	node, err := models.Nodes(
		models.NodeWhere.Address.EQ(ip),
	).One(ctx, pg.db)
	if err != nil {
		return "", -1, err
	}

	return node.Country, node.IPVersion, nil
}

func (pg PgDb) TotalPeerCount(ctx context.Context, timestamp int64) (int64, error) {
	return models.Heartbeats(models.HeartbeatWhere.Timestamp.EQ(timestamp)).Count(ctx, pg.db)
}

func (pg PgDb) PeerCountByUserAgents(ctx context.Context, timestamp int64) (userAgents []netsnapshot.UserAgentInfo,
	 err error) {

	sql := fmt.Sprintf(`SELECT node.user_agent, COUNT(node.user_agent) AS number from node 
			INNER JOIN heartbeat ON node.address = heartbeat.node_id
		WHERE heartbeat.timestamp = %d GROUP BY node.user_agent ORDER BY number DESC`, timestamp)

	var result []struct {
		UserAgent string `json:"user_agent"`
		Number    int64  `json:"number"`
	}

	err = models.Nodes(qm.SQL(sql)).Bind(ctx, pg.db, &result)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, item := range result {
		total += item.Number
	}

	for _, item := range result {
		userAgent := item.UserAgent
		if strings.Trim(userAgent, " ") == "" {
			userAgent = "Unknown"
		}
		userAgents = append(userAgents, netsnapshot.UserAgentInfo{
			UserAgent:  userAgent,
			Nodes:      item.Number,
			Percentage: int64(100.0 * float64(item.Number) / float64(total)),
		})
	}

	sort.Slice(userAgents, func(i, j int) bool {
		return userAgents[i].Nodes > userAgents[j].Nodes
	})

	return
}

func (pg PgDb) PeerCountByCountries(ctx context.Context, timestamp int64) (countries []netsnapshot.CountryInfo,
	 err error) {

	sql := fmt.Sprintf(`SELECT node.country, COUNT(node.country) AS number from node 
		INNER JOIN heartbeat on heartbeat.node_id = node.address WHERE heartbeat.timestamp = %d 
		GROUP BY node.country ORDER BY number DESC`, timestamp)

	var result []struct {
		Country string `json:"country"`
		Number    int64  `json:"number"`
	}

	err = models.Heartbeats(qm.SQL(sql)).Bind(ctx, pg.db, &result)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, item := range result {
		total += item.Number
	}

	for _, item := range result {
		country := item.Country
		if strings.Trim(country, " ") == "" {
			country = "Unknown"
		}
		countries = append(countries, netsnapshot.CountryInfo{
			Country:  item.Country,
			Nodes:      item.Number,
			Percentage: int64(100.0 * float64(item.Number) / float64(total)),
		})
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Nodes > countries[j].Nodes
	})

	return
}

func (pg PgDb) PeerCountByIPVersion(ctx context.Context, timestamp int64, iPVersion int) (int64, error) {
	var result struct{Total int64}
	err := models.NewQuery(
		qm.Select("COUNT(h.node_id) as total"),
		qm.From(fmt.Sprintf("%s as h", models.TableNames.Heartbeat)),
		qm.InnerJoin(fmt.Sprintf("%s as n on n.address = h.node_id", models.TableNames.Node)),
		qm.Where("h.timestamp = ? and n.ip_version = ?", timestamp, iPVersion)).Bind(ctx, pg.db, &result)

	return result.Total, err
}

func (pg PgDb) LastSnapshotTime(ctx context.Context) (timestamp int64) {
	_ = pg.LastEntry(ctx, models.TableNames.NetworkSnapshot, &timestamp)
	return
}

func (pg PgDb) LastSnapshot(ctx context.Context) (*netsnapshot.SnapShot, error) {
	return pg.FindNetworkSnapshot(ctx, pg.LastSnapshotTime(ctx))
}
