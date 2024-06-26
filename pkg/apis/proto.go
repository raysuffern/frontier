package apis

const (
	// edge related
	RPCEdgeOnline    = "edge_online"
	RPCEdgeOffline   = "edge_offline"
	RPCEdgeHeartbeat = "edge_heartbeat"

	// service related
	RPCServiceOnline    = "service_online"
	RPCServiceOffline   = "service_offline"
	RPCServiceHeartbeat = "service_heartbeat"

	// frontier related
	RPCFrontierStats = "frontier_stats"
)

type FrontierInstance struct {
	FrontierID string `json:"frontier_id"`
	// in k8s, it should be podIP:port
	AdvertisedServiceboundAddr string `json:"advertised_servicebound_addr"`
	// in k8s, it should be NodeportIP:port
	AdvertisedEdgeboundAddr string `json:"advertised_edgebound_addr"`
}

type FrontierStats struct {
	FrontierID   string `json:"frontier_id"`
	EdgeCount    int    `json:"edge_count"`
	ServiceCount int    `json:"service_count"`
}

// edge protocols
type EdgeOnline struct {
	FrontierID string `json:"frontier_id"`
	EdgeID     uint64 `json:"edge_id"`
	Addr       string `json:"addr"`
}

type EdgeOffline struct {
	FrontierID string `json:"frontier_id"`
	EdgeID     uint64 `json:"edge_id"`
}

type EdgeHeartbeat struct {
	FrontierID string `json:"frontier_id"`
	EdgeID     uint64 `json:"edge_id"`
}

// service protocols
type ServiceOnline struct {
	FrontierID string `json:"frontier_id"`
	ServiceID  uint64 `json:"service_id"`
	Service    string `json:"service"`
	Addr       string `json:"addr"`
}

type ServiceOffline struct {
	FrontierID string `json:"frontier_id"`
	ServiceID  uint64 `json:"service_id"`
}

type ServiceHeartbeat struct {
	FrontierID string `json:"frontier_id"`
	ServiceID  uint64 `json:"service_id"`
}
