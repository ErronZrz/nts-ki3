package congrat

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeExt(t *testing.T) {
	start := time.Now()
	// 每隔 5 秒执行一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	left := 24
	for left > 0 {
		<-ticker.C
		left--
		now := time.Now()
		fmt.Println(now.UnixNano())
		fmt.Println(start.Add(now.Sub(start)).UnixNano())
		fmt.Println()
	}
}
