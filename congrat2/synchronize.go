package congrat2

import (
	"active/clock"
	"active/congrat1"
	"active/utils"
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

func SynchronizeOnce(db *sql.DB, m, minCandidates, minSurvivors int) error {
	// 0. 更新批次
	maxBatchID, err := congrat1.MaxBatchID(db)
	if err != nil {
		return err
	}
	congrat1.CurrentBatchID = maxBatchID + 1
	// 1. 从数据库读取可用服务器
	serverList, err := getLatestRecordsByAEADID(db, 15)
	if err != nil {
		return err
	}
	fmt.Printf("server num = %d\n", len(serverList))
	// 2. 从本地文件读取上一批次的 survivors
	path := fmt.Sprintf("C:\\Corner\\TMP\\BisheData\\clock\\%d.txt", maxBatchID)
	err, survivorIPMap, survivors := readLastSurvivors(path, serverList)
	if err != nil {
		return err
	}
	fmt.Printf("last survivor num = %d\n", len(survivors))
	// 3. 以分数作为概率筛选 m 台服务器，并且将队列中的服务器加入到排除列表
	selected := selectRecordsByScoreProbability(serverList, m, survivorIPMap)
	fmt.Printf("selected num = %d\n", len(selected))
	// 4. 合并两个列表，记录所选服务器之前的 4 个时间戳，后面会使用
	selected = append(selected, survivors...)
	prevSamples := make(map[string]*clock.OriginSample, len(selected))
	for _, s := range selected {
		t1 := utils.TimestampValue(s.T1)
		t2 := utils.TimestampValue(s.T2)
		t3 := utils.TimestampValue(s.T3)
		t4 := utils.TimestampValue(s.T4)
		prevSamples[s.IPAddress] = clock.NewOriginSample(t1, t2, t3, t4, BaseDispersion)
	}
	// 5. 进行同步
	for _, server := range selected {
		err = queryPort(db, server)
		if err != nil {
			return err
		}
		err = insertServerInfoSimple(db, server, "SYNC")
		if err != nil {
			return err
		}
		err = executeNTP(server, 15)
		if err != nil {
			fmt.Println(err)
			// 如果失败则标记失败
			server.NTPv4Address = FailFlag
			continue
		}
		err = insertKeyTimestamps2(db, server, 15)
		if err != nil {
			return err
		}
	}
	// 6. 生成对等体信息
	peers := getPeers(selected, prevSamples)
	// 7. 选出 truechimers、聚类、合并（使用 Kalman 滤波）
	whatsoever(peers, minCandidates, minSurvivors, true)
	// 8. 更新可用性与分数
	return congrat1.UpdateAvailabilityAndScore(db)
}

func readLastSurvivors(path string, serverList []*KeKeyTimestamp) (error, map[string]bool, []*KeKeyTimestamp) {
	// 读取文件
	file, err := os.Open(path)
	if err != nil {
		return err, nil, nil
	}
	defer func() { _ = file.Close() }()
	scanner := bufio.NewScanner(file)
	// 现在貌似用不到选择抖动，所以就不读取了
	_ = scanner.Scan()
	// 读取 6 个系统变量
	floats := make([]float64, 6)
	for i := 0; i < 6; i++ {
		_ = scanner.Scan()
		floats[i], err = strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			return err, nil, nil
		}
	}
	clock.GlobalSystemClock = &clock.SystemClock{
		Offset:         floats[0],
		Cumsum:         floats[1],
		Jitter:         floats[2],
		RootDelay:      floats[3],
		RootDispersion: floats[4],
		PPrev:          floats[5],
	}
	// 读取 survivors 数量
	_ = scanner.Scan()
	num, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return err, nil, nil
	}
	// 读取上一轮的 survivors IP 地址
	survivors := make([]*KeKeyTimestamp, 0, num)
	ips := make(map[string]bool)
	for i := 0; i < num; i++ {
		_ = scanner.Scan()
		ip := scanner.Text()
		ips[ip] = true
	}
	// 从数据库中筛选出上一轮 survivors
	for _, server := range serverList {
		if ips[server.IPAddress] {
			survivors = append(survivors, server)
		}
	}
	return nil, ips, survivors
}
