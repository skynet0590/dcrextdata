package netsnapshot

type NetworkPeer struct {
	ID              int    `json:"id"`
	Address         string `json:"address"`
	LastReceiveTime int64  `json:"last_receive_time"`
	LastSendTime    int64  `json:"last_send_time"`
	ConnectionTime  int64  `json:"connection_time"`
	ProtocolVersion int    `json:"protocol_version"`
	UserAgent       string `json:"user_agent"`
	StartingHeight  int    `json:"starting_height"`
	CurrentHeight   int    `json:"current_height"`
}
