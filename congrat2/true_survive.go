package congrat2

import (
	"active/clock"
	"active/congrat1"
	"active/utils"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	BaseDispersion          = 3e-6
	FailFlag                = "FAILED"
	MinPanicEliminationRate = 0.05
	MaxPanicKalmanGain      = 0.05
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
		samples := []*clock.OriginSample{clock.NewOriginSample(t1, t2, t3, t4, BaseDispersion)}
		if former, ok := survivorSamples[server.IPAddress]; ok {
			samples = append(samples, former)
		}
		peers = append(peers, clock.NewPeer(samples, server.IPAddress, rootDelay, rootDispersion, now))
	}
	return peers
}

func whatsoever(peers []*clock.Peer, minCandidates, minSurvivors int, useKalman bool) {
	fmt.Printf("len(peers): %d\n", len(peers))
	// 筛选 truechimers
	truechimers := clock.SelectPeers(peers, minCandidates, true)
	fmt.Printf("len(truechimers) = %d\n", len(truechimers))
	// 计算淘汰率
	eliminationRate := 1.0 - float64(len(truechimers))/float64(len(peers))
	// 聚类
	survivors, selectionJitter := clock.ClusterAlgorithm(truechimers, minSurvivors)
	fmt.Printf("selectionJitter = %.10f\n", selectionJitter)
	fmt.Printf("len(survivors) = %d\n", len(survivors))
	// 组合时钟
	newSystemClock := clock.CombineAlgorithm(survivors, selectionJitter, useKalman)
	// 如果满足两个条件：淘汰率超过阈值、卡尔曼增益低于阈值，则触发恐慌模式
	if eliminationRate > MinPanicEliminationRate && clock.KalmanGain[0] < MaxPanicKalmanGain {
		fmt.Printf("PANIC MODE, elimination rate = %.1f%%, Kk = %.10f\n", eliminationRate*100, clock.KalmanGain[0])
	}
	// 打印全局变量
	p0 := clock.GlobalSystemClock.PPrev
	p1 := newSystemClock.PPrev
	fmt.Printf("Offset = %.10f -> %.10f\nCumsum = %.10f -> %.10f\nJitter = %.10f -> %.10f\n"+
		"RootDelay = %.10f -> %.10f\nRootDispersion = %.10f -> %.10f\nSkew =  %.10f -> %.10f\n"+
		"PPrev = %.10f %.10f -> %.10f %.10f\n        %.10f %.10f    %.10f %.10f\n",
		clock.GlobalSystemClock.Offset, newSystemClock.Offset,
		clock.GlobalSystemClock.Cumsum, newSystemClock.Cumsum,
		clock.GlobalSystemClock.Jitter, newSystemClock.Jitter,
		clock.GlobalSystemClock.RootDelay, newSystemClock.RootDelay,
		clock.GlobalSystemClock.RootDispersion, newSystemClock.RootDispersion,
		clock.GlobalSystemClock.Skew, newSystemClock.Skew,
		p0[0][0], p0[0][1], p1[0][0], p1[0][1], p0[1][0], p0[1][1], p1[1][0], p1[1][1])
	// 替换全局变量
	clock.GlobalSystemClock = newSystemClock
	// 以批次号作为文件名，写入文件
	path := filepath.Join(congrat1.BaseDir, "clock", fmt.Sprintf("%d.txt", congrat1.CurrentBatchID))
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
	// 接下来几行分别是系统变量中的 Offset, Jitter, RootDelay, RootDispersion, Skew, PPrev（矩阵占四行）
	gc := clock.GlobalSystemClock
	p := gc.PPrev
	_, err = writer.WriteString(fmt.Sprintf("%.10f\n%.10f\n%.10f\n%.10f\n%.10f\n%.10f\n%.10f\n%.10f\n%.10f\n%.10f\n",
		gc.Offset, gc.Cumsum, gc.Jitter,
		gc.RootDelay, gc.RootDispersion, gc.Skew,
		p[0][0], p[0][1], p[1][0], p[1][1]))
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
