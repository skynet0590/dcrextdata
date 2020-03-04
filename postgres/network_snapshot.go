package postgres

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/raedahgroup/dcrextdata/netsnapshot"
	"github.com/raedahgroup/dcrextdata/postgres/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (pg PgDb) SaveSnapshot(ctx context.Context, snapshot netsnapshot.SnapShot) error {
	if snapshot.NodeCount == 0 {
		log.Critical("this cannot be")
	}
	existingSnapshot, err := models.FindNetworkSnapshot(ctx, pg.db, snapshot.Timestamp)
	if err == nil {
		existingSnapshot.Height = snapshot.Height
		existingSnapshot.NodeCount = snapshot.NodeCount
		_, err = existingSnapshot.Update(ctx, pg.db, boil.Infer())
		return err
	}

	snapshotModel := models.NetworkSnapshot{
		Timestamp: snapshot.Timestamp, Height: snapshot.Height, NodeCount: snapshot.NodeCount,
	}
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

func (pg PgDb) Snapshots(ctx context.Context) ([]netsnapshot.SnapShot, error) {
	snapshotSlice, err := models.NetworkSnapshots(qm.OrderBy("timestamp")).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	snapshots := make([]netsnapshot.SnapShot, len(snapshotSlice))
	for i, m := range snapshotSlice {
		snapshots[i] = netsnapshot.SnapShot {
			Timestamp: m.Timestamp,
			Height: m.Height,
			NodeCount: m.NodeCount,
		}
	}

	return snapshots, nil
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

func (pg PgDb) SaveHeartbeat(ctx context.Context, heartbeat netsnapshot.Heartbeat) error {
	heartbeatModel, err := models.Heartbeats(
		models.HeartbeatWhere.NodeID.EQ(heartbeat.Address),
		models.HeartbeatWhere.Timestamp.EQ(heartbeat.Timestamp)).One(ctx, pg.db)
	if err == nil {
		if heartbeat.CurrentHeight > 0 {
			heartbeatModel.CurrentHeight = heartbeat.CurrentHeight
		}

		if heartbeat.Latency > 0 {
			heartbeatModel.Latency = heartbeat.Latency
		}

		if heartbeat.LastSeen > 0 {
			heartbeatModel.LastSeen = heartbeat.LastSeen
		}

		if _, err = heartbeatModel.Update(ctx, pg.db, boil.Infer()); err != nil {
			return fmt.Errorf("error in saving heartbeatModel, %s", err.Error())
		}
		return nil
	}

	newHeartbeat := models.Heartbeat{
		Timestamp:     heartbeat.Timestamp,
		NodeID:        heartbeat.Address,
		LastSeen:      heartbeat.LastSeen,
		Latency:       heartbeat.Latency,
		CurrentHeight: heartbeat.CurrentHeight,
	}

	if err = newHeartbeat.Insert(ctx, pg.db, boil.Infer()); err != nil {
		return fmt.Errorf("error in saving hearbeat, %s", err.Error())
	}
	return nil
}

func (pg PgDb) NodeExists(ctx context.Context, address string) (bool, error) {
	return models.NodeExists(ctx, pg.db, address)
}

func (pg PgDb) SaveNode(ctx context.Context, peer netsnapshot.NetworkPeer) error  {
	newNode := models.Node{
		Address:         peer.Address,
		IPVersion:       peer.IPVersion,
		Country:         peer.CountryName,
		Region:          peer.RegionName,
		City:            peer.City,
		Zip: 			 peer.Zip,
		LastAttempt:     peer.LastSeen,
		LastSeen:        peer.LastSeen,
		LastSuccess: 	 peer.LastSuccess,
		ConnectionTime:  peer.ConnectionTime,
		ProtocolVersion: int(peer.ProtocolVersion),
		UserAgent:       peer.UserAgent,
		Services:        peer.Services,
		StartingHeight:  peer.StartingHeight,
		CurrentHeight:   peer.CurrentHeight,
		IsDead:          false,
	}
	err := newNode.Insert(ctx, pg.db, boil.Infer())
	return err
}

func (pg PgDb) UpdateNode(ctx context.Context, peer netsnapshot.NetworkPeer) error {
	var cols = models.M{
		models.NodeColumns.LastAttempt: peer.LastAttempt,
		models.NodeColumns.LastSeen: peer.LastSeen,
		models.NodeColumns.LastSuccess: peer.LastSuccess,
		models.NodeColumns.Services: peer.Services,
		models.NodeColumns.StartingHeight: peer.StartingHeight,
		models.NodeColumns.UserAgent: peer.UserAgent,
		models.NodeColumns.CurrentHeight: peer.CurrentHeight,
		models.NodeColumns.ConnectionTime: peer.ConnectionTime,
	}
	 _, err := models.Nodes(models.NodeWhere.Address.EQ(peer.Address)).UpdateAll(ctx, pg.db, cols)
	 return err
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
			INNER JOIN node on node.address = heartbeat.node_id WHERE ` + where +
			fmt.Sprintf(" ORDER BY node.last_seen DESC LIMIT %d OFFSET %d", limit, offset)

	var peerSlice models.NodeSlice
	err := models.NewQuery(qm.SQL(sql)).Bind(ctx, pg.db, &peerSlice)
	if err != nil {
		return nil, 0, fmt.Errorf("error %s, on query %s", err.Error(), sql)
	}

	var peers []netsnapshot.NetworkPeer
	for _, node := range peerSlice {
		peer := netsnapshot.NetworkPeer{
			Address:         node.Address,
			LastSeen:        node.LastSeen,
			ConnectionTime:  node.ConnectionTime,
			ProtocolVersion: uint32(node.ProtocolVersion),
			UserAgent:       node.UserAgent,
			StartingHeight:  node.StartingHeight,
			CurrentHeight:   node.CurrentHeight,
			Services:        node.Services,
		}

		peer.IPInfo = netsnapshot.IPInfo{
			CountryName: node.Country,
			RegionName:  node.Region,
			City:        node.City,
			Zip:         node.Zip,
		}
		peers = append(peers, peer)
	}

	sql = "SELECT COUNT(heartbeat.node_id) as total FROM heartbeat INNER JOIN node on node.address = heartbeat.node_id WHERE " + where
	var countResult struct{Total int64}
	err = models.NewQuery(qm.SQL(sql)).Bind(ctx, pg.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	return peers, countResult.Total, nil
}

func (pg PgDb) GetAvailableNodes(ctx context.Context) ([]net.IP, error) {
	peerSlice, err := models.Nodes(models.NodeWhere.IsDead.EQ(false), qm.Select(models.NodeColumns.Address)).All(ctx, pg.db)
	if err != nil {
		return nil, err
	}

	var peers []net.IP
	for _, node := range peerSlice {
		peer := net.ParseIP(node.Address)
		peers = append(peers, peer)
	}

	return peers, nil
}

func (pg PgDb) NetworkPeer(ctx context.Context, address string) (*netsnapshot.NetworkPeer, error) {
	node, err := models.FindNode(ctx, pg.db, address)
	if err != nil {
		return nil, err
	}
	peer := netsnapshot.NetworkPeer{
		Address:         node.Address,
		LastSeen:        node.LastSeen,
		ConnectionTime:  node.ConnectionTime,
		ProtocolVersion: uint32(node.ProtocolVersion),
		UserAgent:       node.UserAgent,
		StartingHeight:  node.StartingHeight,
		CurrentHeight:   node.CurrentHeight,
		Services:        node.Services,
	}

	peer.IPInfo = netsnapshot.IPInfo{
		CountryName: node.Country,
		RegionName:  node.Region,
		City:        node.City,
		Zip:         node.Zip,
	}

	return &peer, nil
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

func (pg PgDb) SeenNodesByTimestamp(ctx context.Context) ([]netsnapshot.NodeCount, error) {
	var result []netsnapshot.NodeCount
	err := models.NewQuery(
		qm.SQL("SELECT heartbeat.timestamp, COUNT(*) FROM heartbeat group by heartbeat.timestamp order by timestamp"),
	).Bind(ctx, pg.db, &result)
	return result, err
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
