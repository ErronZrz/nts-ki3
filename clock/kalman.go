package clock

import (
	"fmt"
)

var (
	// InitialSkew 初始化漂移频率，一般不变（单位秒/秒）
	InitialSkew = 0.0001
	// InitialP 初始化误差协方差，一般不变（单位平方秒）
	InitialP = [2][2]float64{
		{1.0, 0.0},
		{0.0, 1.0},
	}
	// Q 过程噪声协方差，一般不变（单位平方秒）
	Q = [2][2]float64{
		{0.00008, 0.0},
		{0.0, 0.000001},
	}
	// R0 基础测量噪声协方差，一般不变（单位平方秒）
	R0 = 0.0001
	// Alpha 影响因子，一般不变（单位秒）
	Alpha = 0.01
	// KalmanGain Kalman 增益，用于给别的包读取和修改
	KalmanGain [2]float64
)

// KalmanState 表示当前的估计状态
type KalmanState struct {
	// x_k
	Offset float64
	// θ_k
	Skew float64
	// 协方差矩阵 P_k
	P [2][2]float64
}

// KalmanFilterSkew 使用 2 维状态的 Kalman 滤波器：包含偏差和漂移率
func KalmanFilterSkew(prev KalmanState, measuredOffset, rttErr, dt float64) (next KalmanState) {
	// 状态预测
	xPrior := prev.Offset + prev.Skew*dt
	skewPrior := prev.Skew

	// 打印状态预测
	fmt.Printf("Predicted Offset: %.9f\n", xPrior)
	fmt.Printf("Predicted Skew: %.9f\n", skewPrior)

	// 协方差预测：P_k|k-1 = F * P_k-1 * F^T + Q
	F := [2][2]float64{
		{1, dt},
		{0, 1},
	}

	// 打印 F 矩阵
	fmt.Printf("F matrix: \n")
	fmt.Printf("%.9f %.9f\n", F[0][0], F[0][1])
	fmt.Printf("%.9f %.9f\n", F[1][0], F[1][1])

	// F * P
	var FP [2][2]float64
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			FP[i][j] = F[i][0]*prev.P[0][j] + F[i][1]*prev.P[1][j]
		}
	}

	// 打印 F * P
	fmt.Printf("F * P: \n")
	fmt.Printf("%.9f %.9f\n", FP[0][0], FP[0][1])
	fmt.Printf("%.9f %.9f\n", FP[1][0], FP[1][1])

	// P_predict = FP * F^T + Q
	var P [2][2]float64
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			P[i][j] = FP[i][0]*F[j][0] + FP[i][1]*F[j][1] + Q[i][j]
		}
	}

	// 打印协方差预测 P
	fmt.Printf("Predicted P: \n")
	fmt.Printf("%.9f %.9f\n", P[0][0], P[0][1])
	fmt.Printf("%.9f %.9f\n", P[1][0], P[1][1])

	// 计算测量噪声调整：rk = R0 + Alpha * rttErr
	rk := R0 + Alpha*rttErr
	// 打印测量噪声 rk 和 RTT 错误
	fmt.Printf("RTT Error: %.9f\n", rttErr)
	fmt.Printf("R0: %.9f, Alpha: %.9f, rk: %.9f\n", R0, Alpha, rk)

	// 计算 S = H * P * H^T + rk (H = [1 0])
	S := P[0][0] + rk
	fmt.Printf("S = P[0][0] + rk = %.9f + %.9f = %.9f\n", P[0][0], rk, S)

	// 计算 Kalman 增益 K = P * H^T / S
	K := [2]float64{
		P[0][0] / S,
		P[1][0] / S,
	}
	// 打印 Kalman 增益
	fmt.Printf("Kalman Gain K: %.9f %.9f\n", K[0], K[1])

	// 更新状态估计：x_k = x_k|k-1 + K * (z_k - H * x_k|k-1)
	// y := measuredOffset - xPrior
	xNew := K[0] * measuredOffset
	skewNew := skewPrior + K[1]*measuredOffset

	// 打印更新后的偏差和漂移率
	fmt.Printf("Updated Offset: %.9f\n", xNew)
	fmt.Printf("Updated Skew: %.9f\n", skewNew)

	// 更新协方差 P = (I - K * H) * P_predict
	var pNew [2][2]float64
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			pNew[i][j] = P[i][j] - K[i]*P[0][j]
		}
	}

	// 打印更新后的协方差 P
	fmt.Printf("Updated P: \n")
	fmt.Printf("%.9f %.9f\n", pNew[0][0], pNew[0][1])
	fmt.Printf("%.9f %.9f\n", pNew[1][0], pNew[1][1])

	// 返回更新后的状态
	next = KalmanState{
		Offset: xNew,
		Skew:   skewNew,
		P:      pNew,
	}

	// 打印最终的结果
	fmt.Printf("Result: Offset=%.9f Skew=%.9f\n", next.Offset, next.Skew)
	fmt.Printf("Final P: \n")
	fmt.Printf("%.9f %.9f\n", next.P[0][0], next.P[0][1])
	fmt.Printf("%.9f %.9f\n", next.P[1][0], next.P[1][1])

	return
}
