package congrat2

import (
	"active/clock"
	"active/congrat1"
	"active/utils"
	"bufio"
	"fmt"
	"os"
)

const (
	BaseDispersion = 3e-6
	FailFlag       = "FAILED"
)

func getPeers(servers []*KeKeyTimestamp, survivorSamples map[string]*clock.OriginSample) []*clock.Peer {
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
		if former, ok := survivorSamples[server.IPAddress]; ok {
			sample = append(sample, former)
		}
		peers = append(peers, clock.NewPeer(sample, server.IPAddress, rootDelay, rootDispersion, now))
	}
	return peers
}

func whatsoever(peers []*clock.Peer, minCandidates, minSurvivors int) {
	fmt.Printf("len(peers): %d\n", len(peers))
	// 筛选 truechimers
	truechimers := clock.SelectPeers(peers, minCandidates, false)
	fmt.Printf("len(truechimers) = %d\n", len(truechimers))
	// 聚类
	survivors, selectionJitter := clock.ClusterAlgorithm(truechimers, minSurvivors)
	fmt.Printf("selectionJitter = %.10f\n", selectionJitter)
	fmt.Printf("len(survivors) = %d\n", len(survivors))
	// 组合时钟
	newSystemClock := clock.CombineAlgorithm(survivors, selectionJitter)
	// 打印全局变量
	fmt.Printf("Offset = %.10f -> %.10f\nJitter = %.10f -> %.10f\nRootDelay = %.10f -> %.10f\nRootDispersion = %.10f -> %.10f\n",
		clock.GlobalSystemClock.Offset, newSystemClock.Offset, clock.GlobalSystemClock.Jitter, newSystemClock.Jitter,
		clock.GlobalSystemClock.RootDelay, newSystemClock.RootDelay, clock.GlobalSystemClock.RootDispersion, newSystemClock.RootDispersion)
	clock.GlobalSystemClock = newSystemClock
	// 以批次号作为文件名，写入文件
	path := fmt.Sprintf("C:\\Corner\\TMP\\BisheData\\clock\\%d.txt", congrat1.CurrentBatchID)
	err := writeToFile(path, survivors, selectionJitter)
	if err != nil {
		fmt.Printf("error writing to file: %v", err)
	}
}

func writeToFile(path string, survivors []*clock.Peer, selectionJitter float64) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	writer := bufio.NewWriter(file)
	// 第一行是选择抖动
	_, err = writer.WriteString(fmt.Sprintf("%.10f\n", selectionJitter))
	if err != nil {
		return err
	}
	// 接下来几行分别是系统变量中的 Offset, Jitter, RootDelay, RootDispersion
	_, err = writer.WriteString(fmt.Sprintf("%.10f\n%.10f\n%.10f\n%.10f\n",
		clock.GlobalSystemClock.Offset, clock.GlobalSystemClock.Jitter, clock.GlobalSystemClock.RootDelay, clock.GlobalSystemClock.RootDispersion))
	if err != nil {
		return err
	}
	// 下一行是 survivors 数量
	_, err = writer.WriteString(fmt.Sprintf("%d\n", len(survivors)))
	if err != nil {
		return err
	}
	// 接下来每行是一个 survivor 的 IP 地址（其实可以顺便存一下 t1-t4，但是考虑到后面反正也要查数据库所以无所谓）
	for _, survivor := range survivors {
		_, err = writer.WriteString(fmt.Sprintf("%s\n", survivor.IP))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
