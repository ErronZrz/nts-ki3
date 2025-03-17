package clock

import (
	"fmt"
)

var (
	// InitialP 初始化误差协方差，一般不变（单位平方秒）
	InitialP = 1.0
	// Q 过程噪声协方差，一般不变（单位平方秒）
	Q = 0.00008
	// R0 基础测量噪声协方差，一般不变（单位平方秒）
	R0 = 0.0001
	// Alpha 影响因子，一般不变（单位秒）
	Alpha = 0.01
	// KalmanGain Kalman 增益，用于给别的包读取和修改
	KalmanGain = 1.0
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
	KalmanGain = kk
	fmt.Printf("PPredict = %.10f, Rk = %.10f, Kk = %.10f\n", pPredict, rk, kk)
	// 更新时间偏差估计，这里不应该加上 prevOffset，因为认为之前已经调整过时间
	// offset = prevOffset + kk*(measuredOffset-prevOffset)
	offset = kk * measuredOffset
	fmt.Printf("PrevOffset = %.10f, MeasuredOffset = %.10f, Offset = %.10f\n", prevOffset, measuredOffset, offset)
	// 更新协方差矩阵
	pNow = (1 - kk) * pPredict
	fmt.Printf("Next PPrev = %.10f\n", pNow)

	fmt.Printf("Result: %.6f %.6f %.6f %.6f %.6f %.6f %.6f %.6f",
		pPrev, pPredict, rttErr, rk, kk, pNow, measuredOffset, offset)
	return
}
