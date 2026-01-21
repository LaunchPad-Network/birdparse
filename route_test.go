package birdparse

import (
	"reflect"
	"testing"
)

func TestParseIPv4Routes(t *testing.T) {
	data := `BIRD 2.17.1 ready.
Access restricted
Table master4:
40.0.0.0/14          unicast [us_44324_4 2026-01-19 from 1.1.1.1] * (100) [AS4249i]
        via 10.151.104.1 on eth0
        Type: BGP univ
        BGP.origin: IGP
        BGP.as_path: 44324 216211 30058 2914 32787 4249
        BGP.next_hop: 10.151.104.1
        BGP.local_pref: 100
        BGP.atomic_aggr: 
        BGP.aggregator: 40.15.254.191 AS4249
        BGP.community: (2914,410) (2914,1408) (2914,2401) (2914,3400) (32787,64015) (32787,65522)
        BGP.large_community: (44324, 10000, 52) (44324, 10001, 392) (215172, 0, 100) (215172, 3, 1392) (215172, 3, 52) (215172, 3, 3003) (215172, 5, 2) (215172, 6, 44324)
                     unicast [us_1234_4 12:16:59.123] (100) [AS4249i]
        via 1.2.3.4 on eth1
        Type: BGP univ
        BGP.origin: IGP
        BGP.as_path: 1234 4249
        BGP.next_hop: 1.2.3.4
        BGP.local_pref: 100
        BGP.atomic_aggr: 
        BGP.aggregator: 40.15.254.191 AS4249
        BGP.community: (2914,410) (2914,1408) (2914,2401) (2914,3400) (32787,64015) (32787,65522)
        BGP.large_community: (44324, 10000, 52) (44324, 10001, 392) (215172, 0, 100) (215172, 3, 1392) (215172, 3, 52) (215172, 3, 3003) (215172, 5, 2) (215172, 6, 44324)`

	expected := []Route{
		{
			Network:      "40.0.0.0/14",
			Gateway:      "10.151.104.1",
			Interface:    "eth0",
			FromProtocol: "us_44324_4",
			FromAddress:  "1.1.1.1",
			Primary:      true,
			Metric:       100,
			Type:         []string{"BGP", "univ"},
			BGP: &RouteBGPInfo{
				Origin:     "IGP",
				ASPath:     []int{44324, 216211, 30058, 2914, 32787, 4249},
				NextHop:    "10.151.104.1",
				LocalPref:  100,
				AtomicAggr: "",
				Aggregator: "40.15.254.191 AS4249",
				Communities: [][]int{
					{2914, 410},
					{2914, 1408},
					{2914, 2401},
					{2914, 3400},
					{32787, 64015},
					{32787, 65522},
				},
				LargeCommunities: [][]int{
					{44324, 10000, 52},
					{44324, 10001, 392},
					{215172, 0, 100},
					{215172, 3, 1392},
					{215172, 3, 52},
					{215172, 3, 3003},
					{215172, 5, 2},
					{215172, 6, 44324},
				},
			},
		},
		{
			Network:      "40.0.0.0/14",
			Gateway:      "1.2.3.4",
			Interface:    "eth1",
			FromProtocol: "us_1234_4",
			FromAddress:  "",
			Primary:      false,
			Metric:       100,
			Type:         []string{"BGP", "univ"},
			BGP: &RouteBGPInfo{
				Origin:     "IGP",
				ASPath:     []int{1234, 4249},
				NextHop:    "1.2.3.4",
				LocalPref:  100,
				AtomicAggr: "",
				Aggregator: "40.15.254.191 AS4249",
				Communities: [][]int{
					{2914, 410},
					{2914, 1408},
					{2914, 2401},
					{2914, 3400},
					{32787, 64015},
					{32787, 65522},
				},
				LargeCommunities: [][]int{
					{44324, 10000, 52},
					{44324, 10001, 392},
					{215172, 0, 100},
					{215172, 3, 1392},
					{215172, 3, 52},
					{215172, 3, 3003},
					{215172, 5, 2},
					{215172, 6, 44324},
				},
			},
		},
	}

	result := ParseRoutes(data)

	if len(result) != len(expected) {
		t.Fatalf("Expected %d routes, got %d", len(expected), len(result))
	}

	for i, route := range result {
		if !reflect.DeepEqual(route, expected[i]) {
			t.Errorf("Route %d mismatch:\nGot: %+v\nExpected: %+v", i, route, expected[i])
		}
	}
}

func TestParseIPv6Routes(t *testing.T) {
	data := `BIRD 2.17.1 ready.
Access restricted
Table master6:
2a0a:2c0:1a::/48     unicast [rr_tyom10 23:41:27.768 from 2001:678:11a4::2] * (100/145) [AS210092i]
        via fe80::5efe:a64:bfe on tyom10
        Type: BGP univ
        BGP.origin: IGP
        BGP.as_path: 50263 48648 210092
        BGP.next_hop: 2001:678:11a4::12
        BGP.local_pref: 205
        BGP.community: (0,3255) (0,3326) (0,6768) (0,8647) (0,12883) (0,12963) (0,13249) (0,13335) (0,14061) (0,15169) (0,15895) (0,20764) (0,20940) (0,21497) (0,22697) (0,2
5133) (0,25229) (0,25521) (0,29632) (0,30058) (0,34927) (0,35297) (0,35320) (0,39737) (0,41820) (0,43668) (0,44600) (0,44854) (0,47787) (0,48011) (0,48919) (0,49824) (0,5058
1) (0,54994) (0,61049) (0,62041) (0,64289) (48648,6197)
        BGP.originator_id: 118.91.186.99
        BGP.cluster_list: 0.0.0.1
        BGP.ext_community: (rt, 48648, 3)
        BGP.large_community: (50263, 1910, 437) (50263, 1911, 32) (50263, 1912, 101) (50263, 1913, 804) (50263, 1914, 150) (50263, 1915, 1) (215172, 0, 200) (215172, 3, 52) 
(215172, 3, 1085) (215172, 3, 3012) (215172, 5, 1) (215172, 5, 5) (215172, 6, 47498) (215172, 7, 3756)
                     unicast [us_44324_6 23:41:27.768] (100) [AS210092i]
        via fc00:230::1 on eth0
        Type: BGP univ
        BGP.origin: IGP
        BGP.as_path: 44324 216211 6939 35297 48648 210092
        BGP.next_hop: fc00:230::1 fe80::fcb2:c6ff:fe2a:691
        BGP.local_pref: 100
        BGP.community: (23640,65012) (65101,30) (65102,392) (65103,3921) (65104,3921)
        BGP.large_community: (44324, 10000, 52) (44324, 10001, 392) (215172, 0, 100) (215172, 3, 1392) (215172, 3, 52) (215172, 3, 3003) (215172, 5, 1) (215172, 6, 44324)
2001:44b8:4040::/48  unicast [us_44324_6 2026-01-15] * (100) [AS7545i]
        via fc00:230::1 on eth0
        Type: BGP univ
        BGP.origin: IGP
        BGP.as_path: 44324 216211 3491 3491 6453 7545
        BGP.next_hop: fc00:230::1 fe80::fcb2:c6ff:fe2a:691
        BGP.local_pref: 100
        BGP.community: (3491,4000) (3491,4019) (3491,9002)
        BGP.large_community: (44324, 10000, 52) (44324, 10001, 392) (215172, 0, 100) (215172, 3, 1392) (215172, 3, 52) (215172, 3, 3003) (215172, 5, 2) (215172, 6, 44324)`

	expected := []Route{
		{
			Network:      "2a0a:2c0:1a::/48",
			Gateway:      "fe80::5efe:a64:bfe",
			Interface:    "tyom10",
			FromProtocol: "rr_tyom10",
			FromAddress:  "2001:678:11a4::2",
			Primary:      true,
			Metric:       100,
			IGPMetric:    145,
			Type:         []string{"BGP", "univ"},
			BGP: &RouteBGPInfo{
				Origin:    "IGP",
				ASPath:    []int{50263, 48648, 210092},
				NextHop:   "2001:678:11a4::12",
				LocalPref: 205,
				Communities: [][]int{
					{0, 3255}, {0, 3326}, {0, 6768}, {0, 8647}, {0, 12883},
					{0, 12963}, {0, 13249}, {0, 13335}, {0, 14061}, {0, 15169},
					{0, 15895}, {0, 20764}, {0, 20940}, {0, 21497}, {0, 22697},
					{0, 25133}, {0, 25229}, {0, 25521}, {0, 29632}, {0, 30058},
					{0, 34927}, {0, 35297}, {0, 35320}, {0, 39737}, {0, 41820},
					{0, 43668}, {0, 44600}, {0, 44854}, {0, 47787}, {0, 48011},
					{0, 48919}, {0, 49824}, {0, 50581}, {0, 54994}, {0, 61049},
					{0, 62041}, {0, 64289}, {48648, 6197},
				},
				LargeCommunities: [][]int{
					{50263, 1910, 437},
					{50263, 1911, 32},
					{50263, 1912, 101},
					{50263, 1913, 804},
					{50263, 1914, 150},
					{50263, 1915, 1},
					{215172, 0, 200},
					{215172, 3, 52},
					{215172, 3, 1085},
					{215172, 3, 3012},
					{215172, 5, 1},
					{215172, 5, 5},
					{215172, 6, 47498},
					{215172, 7, 3756},
				},
			},
		},
		{
			Network:      "2a0a:2c0:1a::/48",
			Gateway:      "fc00:230::1",
			Interface:    "eth0",
			FromProtocol: "us_44324_6",
			FromAddress:  "",
			Primary:      false,
			Metric:       100,
			IGPMetric:    0,
			Type:         []string{"BGP", "univ"},
			BGP: &RouteBGPInfo{
				Origin:    "IGP",
				ASPath:    []int{44324, 216211, 6939, 35297, 48648, 210092},
				NextHop:   "fc00:230::1",
				LocalPref: 100,
				Communities: [][]int{
					{23640, 65012},
					{65101, 30},
					{65102, 392},
					{65103, 3921},
					{65104, 3921},
				},
				LargeCommunities: [][]int{
					{44324, 10000, 52},
					{44324, 10001, 392},
					{215172, 0, 100},
					{215172, 3, 1392},
					{215172, 3, 52},
					{215172, 3, 3003},
					{215172, 5, 1},
					{215172, 6, 44324},
				},
			},
		},
		{
			Network:      "2001:44b8:4040::/48",
			Gateway:      "fc00:230::1",
			Interface:    "eth0",
			FromProtocol: "us_44324_6",
			FromAddress:  "",
			Primary:      true,
			Metric:       100,
			IGPMetric:    0,
			Type:         []string{"BGP", "univ"},
			BGP: &RouteBGPInfo{
				Origin:    "IGP",
				ASPath:    []int{44324, 216211, 3491, 3491, 6453, 7545},
				NextHop:   "fc00:230::1",
				LocalPref: 100,
				Communities: [][]int{
					{3491, 4000},
					{3491, 4019},
					{3491, 9002},
				},
				LargeCommunities: [][]int{
					{44324, 10000, 52},
					{44324, 10001, 392},
					{215172, 0, 100},
					{215172, 3, 1392},
					{215172, 3, 52},
					{215172, 3, 3003},
					{215172, 5, 2},
					{215172, 6, 44324},
				},
			},
		},
	}

	result := ParseRoutes(data)

	if len(result) != len(expected) {
		t.Fatalf("Expected %d routes, got %d", len(expected), len(result))
	}

	for i, route := range result {
		if !reflect.DeepEqual(route, expected[i]) {
			t.Errorf("Route %d mismatch:\nGot: %+v\nExpected: %+v", i, route, expected[i])
		}
	}
}
