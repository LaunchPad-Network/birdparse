package birdparse

type BgpProtocolBgpRoutes struct {
	Imported  int `json:"imported"`
	Filtered  int `json:"filtered"`
	Exported  int `json:"exported"`
	Preferred int `json:"preferred"`
}

type BgpProtocolRouteChangeDetail struct {
	Received int `json:"received"`
	Rejected int `json:"rejected"`
	Filtered int `json:"filtered"`
	Ignored  int `json:"ignored"`
	Accepted int `json:"accepted"`
}

type BgpProtocolRouteChanges struct {
	ImportUpdates   *BgpProtocolRouteChangeDetail `json:"import_updates"`
	ImportWithdraws *BgpProtocolRouteChangeDetail `json:"import_withdraws"`
	ExportUpdates   *BgpProtocolRouteChangeDetail `json:"export_updates"`
	ExportWithdraws *BgpProtocolRouteChangeDetail `json:"export_withdraws"`
}

type BgpProtocol struct {
	Protocol         string                   `json:"protocol"`
	Table            string                   `json:"table"`
	State            string                   `json:"state"`
	Connection       string                   `json:"connection"`
	Description      string                   `json:"description"`
	DescriptionShort string                   `json:"description_short"`
	Preference       int                      `json:"preference"`
	InputFilter      string                   `json:"input_filter"`
	OutputFilter     string                   `json:"output_filter"`
	ImportLimit      int                      `json:"import_limit"`
	LimitAction      string                   `json:"limit_action"`
	Routes           *BgpProtocolBgpRoutes    `json:"routes"`
	RouteChanges     *BgpProtocolRouteChanges `json:"route_changes"`
	BgpState         string                   `json:"bgp_state"`
	NeighborAddress  string                   `json:"neighbor_address"`
	NeighborAS       int                      `json:"neighbor_as"`
	NeighborID       string                   `json:"neighbor_id"`
	BgpSession       []string                 `json:"bgp_session"`
	SourceAddress    string                   `json:"source_address"`
	RouteLimitAt     int                      `json:"route_limit_at"`
	HoldTimer        int                      `json:"hold_timer"`
	HoldTimerNow     int                      `json:"hold_timer_now"`
	Keepalive        int                      `json:"keepalive"`
	KeepaliveNow     int                      `json:"keepalive_now"`
}

func (p BgpProtocol) IsValid() bool {
	if p.Protocol == "" &&
		p.Table == "" {
		return false
	}

	return true
}
