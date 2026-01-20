package birdparse

type Route struct {
	Network      string        `json:"network"`
	Gateway      string        `json:"gateway"`
	Interface    string        `json:"interface"`
	FromProtocol string        `json:"from_protocol"`
	FromAddress  string        `json:"from_address"`
	Primary      bool          `json:"primary"`
	Metric       int           `json:"metric"`
	IGPMetric    int           `json:"igp_metric"`
	Type         []string      `json:"type"`
	BGP          *RouteBGPInfo `json:"bgp"`
}

type RouteBGPInfo struct {
	Origin           string  `json:"origin"`
	ASPath           []int   `json:"as_path"`
	NextHop          string  `json:"next_hop"`
	LocalPref        int     `json:"local_pref"`
	MED              int     `json:"med"`
	AtomicAggr       string  `json:"atomic_aggr"`
	Aggregator       string  `json:"aggregator"`
	Communities      [][]int `json:"communities"`
	LargeCommunities [][]int `json:"large_communities"`
}
