package clock

import (
	"fmt"
	"math"
)

func ClusterAlgorithm(peers []*Peer, maxSurvivors int) []*Peer {
	n := len(peers)
	ip2peer := make(map[string]*Peer)
	for _, peer := range peers {
		ip2peer[peer.IP] = peer
	}
	for {
		var maxSelectionJitterIP string
		maxSelectionJitter, minJitter := math.Inf(-1), math.Inf(1)
		for _, peer := range ip2peer {
			ip2peer[peer.IP] = peer
			minJitter = math.Min(minJitter, peer.Jitter)
			offset := peer.Offset
			var v float64
			for _, other := range peers {
				if peer == other {
					continue
				}
				v += math.Pow(offset-other.Offset, 2)
			}
			sj := math.Sqrt(v / float64(n-1))
			if sj > maxSelectionJitter {
				maxSelectionJitter = sj
				maxSelectionJitterIP = peer.IP
			}
		}
		fmt.Printf("remnant: %d, maxSelectionJitter: %.5f, minJitter: %.5f\n", len(ip2peer), maxSelectionJitter, minJitter)
		// 存活数量过多，或选择抖动最大值大于最小抖动，则踢出
		if len(ip2peer) > maxSurvivors || maxSelectionJitter > minJitter {
			delete(ip2peer, maxSelectionJitterIP)
		} else {
			break
		}
	}
	var survivors []*Peer
	for _, peer := range peers {
		if _, ok := ip2peer[peer.IP]; ok {
			survivors = append(survivors, peer)
		}
	}
	return survivors
}
