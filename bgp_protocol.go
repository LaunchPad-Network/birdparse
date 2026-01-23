package birdparse

import (
	"regexp"
	"strconv"
	"strings"
)

func ParseBGPProtocol(data string) BgpProtocol {
	result := BgpProtocol{}
	seenChannels := make(map[string]bool)

	lines := strings.Split(data, "\n")

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if strings.Contains(line, "Channel ipv4") {
			seenChannels["ipv4"] = true
		}
		if strings.Contains(line, "Channel ipv6") {
			seenChannels["ipv6"] = true
		}

		headerRE := regexp.MustCompile(`^(\S+)\s+BGP\s+([-\w]+|\.{3,}|-+)\s+(\w+)\s+([0-9]{4}-[0-9]{2}-[0-9]{2}|[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}(?:\.[0-9]{1,3})?)\s*(.*)$`)

		if m := headerRE.FindStringSubmatch(line); m != nil {
			result.Protocol = m[1]
			result.Table = m[2]
			result.State = m[3]

			result.Connection = strings.TrimSpace(m[5])
			continue
		}

		if m := regexp.MustCompile(`^\s+Description:\s+(.*)$`).FindStringSubmatch(line); m != nil {
			result.Description = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Table:\s+(.*)$`).FindStringSubmatch(line); m != nil {
			result.Table = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Preference:\s+(\d+)$`).FindStringSubmatch(line); m != nil {
			result.Preference, _ = strconv.Atoi(m[1])
			continue
		}

		if m := regexp.MustCompile(`^\s+Input filter:\s+([^\s]+)$`).FindStringSubmatch(line); m != nil {
			result.InputFilter = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Output filter:\s+([^\s]+)$`).FindStringSubmatch(line); m != nil {
			result.OutputFilter = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Import limit:\s+(\d+)$`).FindStringSubmatch(line); m != nil {
			result.ImportLimit = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Action:\s+(\w+)$`).FindStringSubmatch(line); m != nil {
			result.LimitAction = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Routes:\s+(\d+)\s+imported,\s+(\d+)\s+exported`).FindStringSubmatch(line); m != nil {
			result.Routes = &BgpProtocolBgpRoutes{
				Imported: m[1],
				Exported: m[2],
			}
			continue
		}
		if m := regexp.MustCompile(`^\s+Routes:\s+(\d+)\s+imported,\s+(\d+)\s+filtered,\s+(\d+)\s+exported,\s+(\d+)\s+preferred`).FindStringSubmatch(line); m != nil {
			result.Routes = &BgpProtocolBgpRoutes{
				Imported:  m[1],
				Filtered:  m[2],
				Exported:  m[3],
				Preferred: m[4],
			}
			continue
		}

		if m := regexp.MustCompile(`^\s+BGP state:\s+(\w+)$`).FindStringSubmatch(line); m != nil {
			result.BgpState = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Neighbor address:\s+([^\s]+)$`).FindStringSubmatch(line); m != nil {
			result.NeighborAddress = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Neighbor AS:\s+(\d+)$`).FindStringSubmatch(line); m != nil {
			result.NeighborAS = atoi(m[1])
			continue
		}

		if m := regexp.MustCompile(`^\s+Neighbor ID:\s+([^\s]+)$`).FindStringSubmatch(line); m != nil {
			result.NeighborID = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Session:\s+(.*)$`).FindStringSubmatch(line); m != nil {
			result.BgpSession = strings.Fields(m[1])
			continue
		}

		if m := regexp.MustCompile(`^\s+Source address:\s+([^\s]+)$`).FindStringSubmatch(line); m != nil {
			result.SourceAddress = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Route limit:\s+(\d+)/(\d+)$`).FindStringSubmatch(line); m != nil {
			result.RouteLimitAt = m[1]
			continue
		}

		if m := regexp.MustCompile(`^\s+Hold timer:\s+([\d.]+)/([\d.]+)$`).FindStringSubmatch(line); m != nil {
			result.HoldTimerNow = atoi(m[1])
			result.HoldTimer = atoi(m[2])
			continue
		}

		if m := regexp.MustCompile(`^\s+Keepalive timer:\s+([\d.]+)/([\d.]+)$`).FindStringSubmatch(line); m != nil {
			result.KeepaliveNow = atoi(m[1])
			result.Keepalive = atoi(m[2])
			continue
		}

		if m := regexp.MustCompile(`^\s+Import updates:\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)$`).FindStringSubmatch(line); m != nil {
			if result.RouteChanges == nil {
				result.RouteChanges = &BgpProtocolRouteChanges{}
			}
			result.RouteChanges.ImportUpdates = &BgpProtocolRouteChangeDetail{
				Received: parseOptionalIntAsString(m[1]),
				Rejected: parseOptionalIntAsString(m[2]),
				Filtered: parseOptionalIntAsString(m[3]),
				Ignored:  parseOptionalIntAsString(m[4]),
				Accepted: parseOptionalIntAsString(m[5]),
			}
			continue
		}

		if m := regexp.MustCompile(`^\s+Import withdraws:\s+(\d+)\s+(\d+|--+)\s+(\d+|--+)\s+(\d+|--+)\s+(\d+|--+)$`).FindStringSubmatch(line); m != nil {
			if result.RouteChanges == nil {
				result.RouteChanges = &BgpProtocolRouteChanges{}
			}
			result.RouteChanges.ImportWithdraws = &BgpProtocolRouteChangeDetail{
				Received: parseOptionalIntAsString(m[1]),
				Rejected: parseOptionalIntAsString(m[2]),
				Filtered: parseOptionalIntAsString(m[3]),
				Ignored:  parseOptionalIntAsString(m[4]),
				Accepted: parseOptionalIntAsString(m[5]),
			}
			continue
		}

		if m := regexp.MustCompile(`^\s+Export updates:\s+(\d+)\s+(\d+|--+)\s+(\d+|--+)\s+(\d+|--+)\s+(\d+|--+)$`).FindStringSubmatch(line); m != nil {
			if result.RouteChanges == nil {
				result.RouteChanges = &BgpProtocolRouteChanges{}
			}
			result.RouteChanges.ExportUpdates = &BgpProtocolRouteChangeDetail{
				Received: parseOptionalIntAsString(m[1]),
				Rejected: parseOptionalIntAsString(m[2]),
				Filtered: parseOptionalIntAsString(m[3]),
				Ignored:  parseOptionalIntAsString(m[4]),
				Accepted: parseOptionalIntAsString(m[5]),
			}
			continue
		}

		if m := regexp.MustCompile(`^\s+Export withdraws:\s+(\d+)\s+(\d+|--+)\s+(\d+|--+)\s+(\d+|--+)\s+(\d+|--+)$`).FindStringSubmatch(line); m != nil {
			if result.RouteChanges == nil {
				result.RouteChanges = &BgpProtocolRouteChanges{}
			}
			result.RouteChanges.ExportWithdraws = &BgpProtocolRouteChangeDetail{
				Received: parseOptionalIntAsString(m[1]),
				Rejected: parseOptionalIntAsString(m[2]),
				Filtered: parseOptionalIntAsString(m[3]),
				Ignored:  parseOptionalIntAsString(m[4]),
				Accepted: parseOptionalIntAsString(m[5]),
			}
			continue
		}
	}

	if result.RouteLimitAt == "" && result.Routes != nil {
		result.RouteLimitAt = result.Routes.Imported
	}

	if seenChannels["ipv4"] && seenChannels["ipv6"] {
		// TODO: MP-BGP not supported yet
		result = BgpProtocol{}
	}

	return result
}

func ParseBGPProtocols(data string) []BgpProtocol {
	var (
		results      []BgpProtocol
		currentBlock []string
		inBlock      bool
	)

	lines := strings.Split(data, "\n")

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r")

		if strings.HasPrefix(line, "BIRD") ||
			strings.HasPrefix(line, "Access restricted") {
			continue
		}

		if isBGPHeaderLine(line) {
			if inBlock && len(currentBlock) > 0 {
				p := ParseBGPProtocol(strings.Join(currentBlock, "\n"))
				if p.IsValid() {
					results = append(results, p)
				}
			}

			currentBlock = []string{line}
			inBlock = true
			continue
		}

		if inBlock && strings.TrimSpace(line) != "" {
			currentBlock = append(currentBlock, line)
		}
	}

	if inBlock && len(currentBlock) > 0 {
		p := ParseBGPProtocol(strings.Join(currentBlock, "\n"))
		if p.IsValid() {
			results = append(results, p)
		}
	}

	return results
}

func isBGPHeaderLine(line string) bool {
	headerRE := regexp.MustCompile(`^(\S+)\s+BGP\s+([-\w]+|\.{3,}|-+)\s+(\w+)\s+([0-9]{4}-[0-9]{2}-[0-9]{2}|[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}(?:\.[0-9]{1,3})?)\s*(.*)$`)
	return headerRE.FindStringSubmatch(line) != nil
}
