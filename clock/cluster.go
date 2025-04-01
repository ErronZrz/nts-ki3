package clock

import (
	"math"
)

func ClusterAlgorithm(peers []*Peer, minSurvivors int) (survivors []*Peer, selectionJitter float64) {
	n := len(peers)
	ip2peer := make(map[string]*Peer)
	for _, peer := range peers {
		ip2peer[peer.IP] = peer
	}
	for {
		var maxSelectionJitterIP string
		selectionJitter = math.Inf(-1)
		minJitter := math.Inf(1)
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
			if sj > selectionJitter {
				selectionJitter = sj
				maxSelectionJitterIP = peer.IP
			}
		}
		// fmt.Printf("remnant: %d, selectionJitter = %.10f, minJitter = %.10f\n", len(ip2peer), selectionJitter, minJitter)
		// 存活数量已达到下限，或者选择抖动最大值小于最小抖动，则完成聚类，否则踢出
		if len(ip2peer) <= minSurvivors || selectionJitter <= minJitter {
			break
		}
		delete(ip2peer, maxSelectionJitterIP)
	}
	for _, peer := range peers {
		if _, ok := ip2peer[peer.IP]; ok {
			survivors = append(survivors, peer)
		}
	}
	return
}
