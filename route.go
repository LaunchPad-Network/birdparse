package birdparse

import (
	"regexp"
	"strconv"
	"strings"
)

type collectorType int

const (
	collectorNone collectorType = iota
	collectorGateway
	collectorTypeSource
	collectorBGPPrefix
	collectorBGPASPath
	collectorBGPNextHop
	collectorBGPLocalPref
	collectorBGPMED
	collectorBGPAtomicAggr
	collectorBGPAggregator
	collectorBGPCommunity
	collectorBGPLargeCommunity
)

func ParseRoutes(data string) []Route {
	routes := []Route{}
	var currentRoute Route
	lines := strings.Split(data, "\n")

	var currentCollector collectorType
	var collectorLines []string

	resetCollector := func() {
		currentCollector = collectorNone
		collectorLines = nil
	}

	processCollector := func() {
		if len(collectorLines) == 0 {
			return
		}

		fullLine := strings.Join(collectorLines, "")
		fullLine = strings.TrimSpace(fullLine)

		switch currentCollector {
		case collectorGateway:
			if matches := regexp.MustCompile(`^via\s+([0-9a-f\.\:]+)\s+on\s+([a-zA-Z0-9_\.\-\/]+).*$`).FindStringSubmatch(fullLine); matches != nil {
				currentRoute.Gateway = matches[1]
				currentRoute.Interface = matches[2]
			}
		case collectorTypeSource:
			if matches := regexp.MustCompile(`^(?:Type|source):\s+(.*)$`).FindStringSubmatch(fullLine); matches != nil {
				currentRoute.Type = strings.Fields(strings.TrimSpace(matches[1]))
			}
		case collectorBGPCommunity:
			if matches := regexp.MustCompile(`^BGP\.community:\s+(.+)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				communities := parseCommunities(matches[1])
				for _, community := range communities {
					if !containsCommunity(currentRoute.BGP.Communities, community) {
						currentRoute.BGP.Communities = append(currentRoute.BGP.Communities, community)
					}
				}
			}
		case collectorBGPLargeCommunity:
			if matches := regexp.MustCompile(`^BGP\.large_community:\s+(.+)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				largeCommunities := parseLargeCommunities(matches[1])
				for _, community := range largeCommunities {
					if !containsLargeCommunity(currentRoute.BGP.LargeCommunities, community) {
						currentRoute.BGP.LargeCommunities = append(currentRoute.BGP.LargeCommunities, community)
					}
				}
			}
		case collectorBGPASPath:
			if matches := regexp.MustCompile(`^(?:BGP\.as_path|bgp_path):\s+(.*)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				currentRoute.BGP.ASPath, _ = parseASPath(strings.TrimSpace(matches[1]))
			}
		case collectorBGPNextHop:
			if matches := regexp.MustCompile(`^BGP\.next_hop:\s+([0-9a-f\.\:]+)(?:\s+[^\s]+)*$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				currentRoute.BGP.NextHop = matches[1]
			}
		case collectorBGPLocalPref:
			if matches := regexp.MustCompile(`^BGP\.local_pref:\s+(\w+)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				currentRoute.BGP.LocalPref = atoi(matches[1])
			}
		case collectorBGPMED:
			if matches := regexp.MustCompile(`^BGP\.med:\s+(\d+)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				currentRoute.BGP.MED = atoi(matches[1])
			}
		case collectorBGPAtomicAggr:
			if matches := regexp.MustCompile(`^BGP\.atomic_aggr:(.*)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				currentRoute.BGP.AtomicAggr = strings.TrimSpace(matches[1])
			}
		case collectorBGPAggregator:
			if matches := regexp.MustCompile(`^BGP\.aggregator:(.*)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				currentRoute.BGP.Aggregator = strings.TrimSpace(matches[1])
			}
		case collectorBGPPrefix:
			if matches := regexp.MustCompile(`^BGP\.origin:\s+(\w+)$`).FindStringSubmatch(fullLine); matches != nil {
				if currentRoute.BGP == nil {
					currentRoute.BGP = &RouteBGPInfo{}
				}
				currentRoute.BGP.Origin = matches[1]
			}
		}

		resetCollector()
	}

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if strings.HasPrefix(line, "BIRD") ||
			strings.HasPrefix(line, "Access restricted") ||
			strings.HasPrefix(line, "Table ") {
			continue
		}

		if matches := regexp.MustCompile(`^([0-9a-f.:\/]+)\s+((?:via\s+([0-9a-f.:]+)\s+on\s+([a-zA-Z0-9_.\-\/]+))|\w+)\s+\[(\w+)\s+([0-9]{4}-[0-9]{1,2}-[0-9]{1,2}|[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}(?:\.[0-9]+)?)(?:\s+from\s+([0-9a-f.:\/]+))?\](?:\s+(\*))?\s+\((\d+)(?:\/(\-?\d+))?\).*$`).FindStringSubmatch(line); matches != nil {
			processCollector()
			if currentRoute.Network != "" {
				routes = append(routes, currentRoute)
			}
			currentRoute = mainRouteDetail(matches)
			resetCollector()
			continue
		} else if matches := regexp.MustCompile(`^\s+((?:via\s+([0-9a-f.:]+)\s+on\s+([a-zA-Z0-9_.\-\/]+))|\w+)\s+\[(\w+)\s+([0-9]{4}-[0-9]{1,2}-[0-9]{1,2}|[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}(?:\.[0-9]+)?)(?:\s+from\s+([0-9a-f.:\/]+))?\](?:\s+(\*))?\s+\((\d+)(?:\/(\-?\d+))?\).*$`).FindStringSubmatch(line); matches != nil {
			processCollector()
			if currentRoute.Network != "" {
				routes = append(routes, currentRoute)
				network := currentRoute.Network
				fullMatch := matches[0]
				matches = matches[1:]
				matches = append([]string{network}, matches...)
				matches = append([]string{fullMatch}, matches...)
				currentRoute = mainRouteDetail(matches)
				resetCollector()
			}
			continue
		}

		var detectedCollector collectorType
		trimmedLine := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmedLine, "via "):
			detectedCollector = collectorGateway
		case strings.HasPrefix(trimmedLine, "Type:") || strings.HasPrefix(trimmedLine, "source:"):
			detectedCollector = collectorTypeSource
		case strings.HasPrefix(trimmedLine, "BGP.community:"):
			detectedCollector = collectorBGPCommunity
		case strings.HasPrefix(trimmedLine, "BGP.large_community:"):
			detectedCollector = collectorBGPLargeCommunity
		case strings.HasPrefix(trimmedLine, "BGP.as_path:") || strings.HasPrefix(trimmedLine, "bgp_path:"):
			detectedCollector = collectorBGPASPath
		case strings.HasPrefix(trimmedLine, "BGP.next_hop:"):
			detectedCollector = collectorBGPNextHop
		case strings.HasPrefix(trimmedLine, "BGP.local_pref:"):
			detectedCollector = collectorBGPLocalPref
		case strings.HasPrefix(trimmedLine, "BGP.med:"):
			detectedCollector = collectorBGPMED
		case strings.HasPrefix(trimmedLine, "BGP.atomic_aggr:"):
			detectedCollector = collectorBGPAtomicAggr
		case strings.HasPrefix(trimmedLine, "BGP.aggregator:"):
			detectedCollector = collectorBGPAggregator
		case strings.HasPrefix(trimmedLine, "BGP.origin:"):
			detectedCollector = collectorBGPPrefix
		default:
			if currentCollector != collectorNone && trimmedLine != "" {
				collectorLines = append(collectorLines, line)
				continue
			}
		}

		if detectedCollector != collectorNone {
			processCollector()
			currentCollector = detectedCollector
			collectorLines = []string{line}
		}
	}

	processCollector()
	if currentRoute.Network != "" {
		routes = append(routes, currentRoute)
	}

	return routes
}

func mainRouteDetail(matches []string) Route {
	var r Route
	if len(matches) < 9 {
		return r
	}

	r.Network = matches[1]
	hasVia := len(matches) >= 4 && matches[3] != ""

	if hasVia {
		r.Gateway = matches[3]
		if len(matches) >= 5 && matches[4] != "" {
			r.Interface = matches[4]
		} else {
			r.Interface = "none"
		}
	} else {
		r.Gateway = matches[2]
		r.Interface = matches[2]
	}

	if len(matches) >= 6 && matches[5] != "" {
		r.FromProtocol = matches[5]
	}

	if len(matches) >= 8 && matches[7] != "" {
		r.FromAddress = matches[7]
	}

	if len(matches) >= 9 && matches[8] == "*" {
		r.Primary = true
	}

	if len(matches) >= 10 && matches[9] != "" {
		if metric, err := strconv.Atoi(matches[9]); err == nil {
			r.Metric = metric
		}
	}

	if len(matches) >= 11 && matches[10] != "" {
		if igpMetric, err := strconv.Atoi(matches[10]); err == nil {
			r.IGPMetric = igpMetric
		}
	}

	return r
}

func parseCommunities(communityStr string) [][]int {
	communities := [][]int{}
	trimmed := strings.TrimSpace(communityStr)

	for _, community := range strings.Fields(trimmed) {
		if matches := regexp.MustCompile(`^\((\d+),(\d+)\)`).FindStringSubmatch(community); matches != nil {
			if asn1, err1 := strconv.Atoi(matches[1]); err1 == nil {
				if asn2, err2 := strconv.Atoi(matches[2]); err2 == nil {
					communities = append(communities, []int{asn1, asn2})
				}
			}
		}
	}

	return communities
}

func parseLargeCommunities(largeCommunityStr string) [][]int {
	communities := [][]int{}
	trimmed := strings.TrimSpace(largeCommunityStr)
	trimmed = strings.Trim(trimmed, "()")
	pairs := strings.Split(trimmed, ") (")

	for _, pair := range pairs {
		if matches := regexp.MustCompile(`^(\d+),\s*(\d+),\s*(\d+)`).FindStringSubmatch(pair); matches != nil {
			if val1, err1 := strconv.Atoi(matches[1]); err1 == nil {
				if val2, err2 := strconv.Atoi(matches[2]); err2 == nil {
					if val3, err3 := strconv.Atoi(matches[3]); err3 == nil {
						communities = append(communities, []int{val1, val2, val3})
					}
				}
			}
		}
	}

	return communities
}

func containsCommunity(communities [][]int, community []int) bool {
	for _, c := range communities {
		if len(c) == 2 && len(community) == 2 && c[0] == community[0] && c[1] == community[1] {
			return true
		}
	}
	return false
}

func containsLargeCommunity(communities [][]int, community []int) bool {
	for _, c := range communities {
		if len(c) == 3 && len(community) == 3 &&
			c[0] == community[0] && c[1] == community[1] && c[2] == community[2] {
			return true
		}
	}
	return false
}

func parseASPath(asPathStr string) ([]int, error) {
	fields := strings.Fields(asPathStr)
	asPath := make([]int, 0, len(fields))

	for _, field := range fields {
		if asn, err := strconv.Atoi(field); err == nil {
			asPath = append(asPath, asn)
		} else {
			re := regexp.MustCompile(`(\d+)`)
			if matches := re.FindStringSubmatch(field); matches != nil {
				if asn, err := strconv.Atoi(matches[1]); err == nil {
					asPath = append(asPath, asn)
				}
			}
		}
	}

	return asPath, nil
}
