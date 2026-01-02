package proxy_group

import (
	"github.com/coder-lulu/newbee-ops-rpc/ent/proxygroup"
)

// stringToSelectionStrategy converts string to proxygroup.SelectionStrategy enum
func stringToSelectionStrategy(s *string) proxygroup.SelectionStrategy {
	if s == nil {
		return proxygroup.SelectionStrategyRoundRobin
	}
	switch *s {
	case "round_robin":
		return proxygroup.SelectionStrategyRoundRobin
	case "random":
		return proxygroup.SelectionStrategyRandom
	case "least_connections":
		return proxygroup.SelectionStrategyLeastConnections
	case "weighted":
		return proxygroup.SelectionStrategyWeighted
	default:
		return proxygroup.SelectionStrategyRoundRobin
	}
}

// selectionStrategyToString converts proxygroup.SelectionStrategy enum to string
func selectionStrategyToString(strategy proxygroup.SelectionStrategy) string {
	return string(strategy)
}
