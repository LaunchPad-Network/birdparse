package birdparse

import (
	"reflect"
	"testing"
)

func TestParseBGPProtocols(t *testing.T) {
	data := `BIRD 2.18 ready.
Access restricted
rpki_launchpad RPKI       ---        up     2026-01-16    Established
  Cache server:     v4.rpki.launchpadx.top
  Cache port:       8282
  Status:           Established
  Transport:        Unprotected over TCP
  Protocol version: 2
  Session ID:       27605
  Serial number:    6217
  Last update:      before 201.498 s
  Refresh timer   : 42998.501/43200
  Retry timer     : ---
  Expire timer    : 86198.501/86400
  No roa4 channel
  Channel roa6
    State:          UP
    Table:          roa_table_v6
    Preference:     100
    Input filter:   ACCEPT
    Output filter:  REJECT
    Routes:         162155 imported, 0 exported, 162150 preferred
    Route change stats:     received   rejected   filtered    ignored   accepted
      Import updates:         162838          0          0          0     162838
      Import withdraws:          683          0        ---          0        683
      Export updates:              0          0          0        ---          0
      Export withdraws:            0        ---        ---        ---          0
  Channel aspa
    State:          UP
    Table:          aspa_table
    Preference:     100
    Input filter:   ACCEPT
    Output filter:  REJECT
    Routes:         698 imported, 0 exported, 698 preferred
    Route change stats:     received   rejected   filtered    ignored   accepted
      Import updates:            738          0          0          0        738
      Import withdraws:           14          0        ---          0         14
      Export updates:              0          0          0        ---          0
      Export withdraws:            0        ---        ---        ---          0

AS213605_13_V6 BGP        ---        up     23:41:27.768    Established   
  BGP state:          Established
    Neighbor address: 2602:f92a:1315::e
    Neighbor AS:      213605
    Local AS:         203168
    Neighbor ID:      23.151.104.19
    Local capabilities
      Multiprotocol
        AF announced: ipv6
      Route refresh
      Graceful restart
      4-octet AS numbers
      Enhanced refresh
      Long-lived graceful restart
    Neighbor capabilities
      Multiprotocol
        AF announced: ipv6
      Route refresh
      Graceful restart
        Restart time: 120
        AF supported: ipv6
        AF preserved:
      4-octet AS numbers
      Enhanced refresh
      Long-lived graceful restart
    Session:          external route-server AS4
    Source address:   2602:f92a:1315::1
    Hold timer:       211.190/240
    Keepalive timer:  32.592/80
    Send hold timer:  370.477/480
  Channel ipv6
    State:          UP
    Table:          master6
    Preference:     100
    Input filter:   import_filter_test1
    Output filter:  output_filter_test1
    Import limit:   48
      Action:       block
    Export limit:   300
      Action:       block
    Routes:         31 imported, 0 filtered, 14 exported, 31 preferred
    Route change stats:     received   rejected   filtered    ignored   accepted
      Import updates:           1257          0          0        331        926
      Import withdraws:         7473          0        ---       6616        857
      Export updates:           1207       1036          0        ---        171
      Export withdraws:          880        ---        ---        ---        157
    BGP Next hop:   2602:f92a:1315::1 fe80::28a:70ff:fe18:84db

AS151673_16_V6 BGP        ---        start  2026-01-16    Passive       
  BGP state:          Passive
    Neighbor address: 2602:f92a:1315::11
    Neighbor AS:      151673
    Local AS:         203168
  Channel ipv6
    State:          DOWN
    Table:          master6
    Preference:     100
    Input filter:   (unnamed)
    Output filter:  (unnamed)`

	expected := []BgpProtocol{
		{
			Protocol:     "AS213605_13_V6",
			Table:        "master6",
			State:        "up",
			Connection:   "Established",
			Preference:   100,
			InputFilter:  "import_filter_test1",
			OutputFilter: "output_filter_test1",
			ImportLimit:  "48",
			LimitAction:  "block",
			Routes: &BgpProtocolBgpRoutes{
				Imported:  "31",
				Filtered:  "0",
				Exported:  "14",
				Preferred: "31",
			},
			RouteChanges: &BgpProtocolRouteChanges{
				ImportUpdates: &BgpProtocolRouteChangeDetail{
					Received: "1257",
					Rejected: "0",
					Filtered: "0",
					Ignored:  "331",
					Accepted: "926",
				},
				ImportWithdraws: &BgpProtocolRouteChangeDetail{
					Received: "7473",
					Rejected: "0",
					Filtered: "0",
					Ignored:  "6616",
					Accepted: "857",
				},
				ExportUpdates: &BgpProtocolRouteChangeDetail{
					Received: "1207",
					Rejected: "1036",
					Filtered: "0",
					Ignored:  "0",
					Accepted: "171",
				},
				ExportWithdraws: &BgpProtocolRouteChangeDetail{
					Received: "880",
					Rejected: "0",
					Filtered: "0",
					Ignored:  "0",
					Accepted: "157",
				},
			},
			BgpState:        "Established",
			NeighborAddress: "2602:f92a:1315::e",
			NeighborAS:      213605,
			NeighborID:      "23.151.104.19",
			BgpSession: []string{
				"external",
				"route-server",
				"AS4",
			},
			SourceAddress: "2602:f92a:1315::1",
			RouteLimitAt:  "31",
			HoldTimer:     240,
			HoldTimerNow:  211,
			Keepalive:     80,
			KeepaliveNow:  32,
		},
		{
			Protocol:        "AS151673_16_V6",
			Table:           "master6",
			State:           "start",
			Connection:      "Passive",
			Preference:      100,
			InputFilter:     "(unnamed)",
			OutputFilter:    "(unnamed)",
			BgpState:        "Passive",
			NeighborAddress: "2602:f92a:1315::11",
			NeighborAS:      151673,
		},
	}

	result := ParseBGPProtocols(data)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseBGPProtocols() = %v, want %v", result, expected)
	}
}
