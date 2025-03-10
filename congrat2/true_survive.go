package congrat2

import (
	"active/clock"
	"active/utils"
)

const (
	BaseDispersion = 3e-6
	FailFlag       = "FAILED"
)

func getPeers(servers []*KeKeyTimestamp) []*clock.Peer {
	var peers []*clock.Peer
	for _, server := range servers {
		if server.NTPv4Address == FailFlag {
			continue
		}
		t1 := utils.TimestampValue(server.T1)
		t2 := utils.TimestampValue(server.T2)
		t3 := utils.TimestampValue(server.T3)
		t4 := utils.TimestampValue(server.T4)
		rootDelay := utils.RootDelayToValue(server.RootDelay)
		rootDispersion := utils.RootDelayToValue(server.RootDispersion)
		now := utils.TimestampValue(utils.GetTimestamp(utils.GlobalNowTime()))
		sample := []*clock.OriginSample{clock.NewOriginSample(t1, t2, t3, t4, BaseDispersion)}
		peers = append(peers, clock.NewPeer(sample, server.IPAddress, rootDelay, rootDispersion, now))
	}
	return peers
}

func whatsoever(peers []*clock.Peer, minCandidates, maxSurvivors int) {
	truechimers := clock.SelectPeers(peers, minCandidates, false)
	survivors, selectionJitter := clock.ClusterAlgorithm(truechimers, maxSurvivors)
	GlobalSystemClock = clock.CombineAlgorithm(survivors, selectionJitter)
}
