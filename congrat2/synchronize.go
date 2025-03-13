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
	// 3. 记录这些 survivors 之前的四个时间戳，后面会使用
	survivorSamples := make(map[string]*clock.OriginSample, len(survivors))
	for _, survivor := range survivors {
		t1 := utils.TimestampValue(survivor.T1)
		t2 := utils.TimestampValue(survivor.T2)
		t3 := utils.TimestampValue(survivor.T3)
		t4 := utils.TimestampValue(survivor.T4)
		survivorSamples[survivor.IPAddress] = clock.NewOriginSample(t1, t2, t3, t4, BaseDispersion)
	}
	// 4. 以分数作为概率筛选 m 台服务器，并且将队列中的服务器加入到排除列表
	selected := selectRecordsByScoreProbability(serverList, m, survivorIPMap)
	fmt.Printf("selected num = %d\n", len(selected))
	// 5. 合并两个列表，然后进行同步
	selected = append(selected, survivors...)
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
	peers := getPeers(selected, survivorSamples)
	// 7. 选出 truechimers、聚类、合并
	whatsoever(peers, minCandidates, minSurvivors)
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
	// 读取四个系统变量
	floats := make([]float64, 4)
	for i := 0; i < 4; i++ {
		_ = scanner.Scan()
		floats[i], err = strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			return err, nil, nil
		}
	}
	clock.GlobalSystemClock = &clock.SystemClock{
		Offset:         floats[0],
		Jitter:         floats[1],
		RootDelay:      floats[2],
		RootDispersion: floats[3],
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
