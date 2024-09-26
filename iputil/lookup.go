package iputil

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func PTRRecords(path, ptrPath string) error {
	closeFile := func(f *os.File) { _ = f.Close() }

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer closeFile(file)
	scanner := bufio.NewScanner(file)

	ptrFile, err := os.OpenFile(ptrPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer closeFile(ptrFile)
	ptrWriter := bufio.NewWriter(ptrFile)

	resolver := &net.Resolver{}

	// 创建 worker 的数量
	const numWorkers = 1000
	jobs := make(chan string, numWorkers)
	results := make(chan string, numWorkers)
	wg := sync.WaitGroup{}

	// 每 30 秒打印一次时间
	go func() {
		var elapsed int
		for range time.Tick(30 * time.Second) {
			elapsed++
			fmt.Printf("%d seconds\n", elapsed*30)
		}
	}()

	// 启动工作协程
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		ii := i
		go func() {
			defer wg.Done()
			var idx int
			for ip := range jobs {
				idx++
				fmt.Println(ii, idx)
				result, err := performLookup(resolver, ip)
				if err != nil {
					results <- ip + "\terror\n"
				} else {
					results <- ip + "\t" + result + "\n"
				}
			}
		}()
	}

	// 主线程扫描 IP 并分发到 jobs
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			ip := strings.Split(line, "\t")[0]
			jobs <- ip
		}
		close(jobs)
	}()

	// 等待所有工作协程完成并关闭结果通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 处理结果
	for result := range results {
		if _, err := ptrWriter.WriteString(result); err != nil {
			return err
		}
	}

	return ptrWriter.Flush()
}

func performLookup(resolver *net.Resolver, ip string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	names, err := resolver.LookupAddr(ctx, ip)
	if err != nil {
		return "", err
	}
	return strings.Join(names, "\t"), nil
}
