package clock

import "fmt"

var (
	// Q 过程噪声协方差，一般不变
	Q = 0.001
	// R0 基础测量噪声协方差，一般不变
	R0 = 0.01
	// Alpha 影响因子，一般不变
	Alpha = 1.0
)

func KalmanFilter(prevOffset, measuredOffset, pPrev, rttErr float64) (offset, pNow float64) {
	// 预测阶段的预测协方差
	pPredict := pPrev + Q
	fmt.Printf("PPrev = %.10f, Q = %.10f, PPredict = %.10f\n", pPrev, Q, pPredict)
	// 测量更新阶段动态调整测量噪声协方差
	rk := R0 + Alpha*rttErr
	fmt.Printf("R0 = %.10f, Alpha = %.10f, RTTError = %.10f, Rk = %.10f\n", R0, Alpha, rttErr, rk)
	// 计算 Kalman 增益
	kk := pPredict / (pPredict + rk)
	fmt.Printf("PPredict = %.10f, Rk = %.10f, KK = %.10f\n", pPredict, rk, kk)
	// 更新时间偏差估计
	offset = prevOffset + kk*(measuredOffset-prevOffset)
	fmt.Printf("PrevOffset = %.10f, MeasuredOffset = %.10f, Offset = %.10f\n", prevOffset, measuredOffset, offset)
	// 更新协方差矩阵
	pNow = (1 - kk) * pPredict
	fmt.Printf("Next PPrev = %.10f\n", pNow)
	return
}
