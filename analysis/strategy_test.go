package analysis

import "testing"

func TestDrawBarChart(t *testing.T) {
	data1 := [][]float64{
		{2202, 3486, 3701},
		{3132, 4742, 6829},
		{2312, 3670, 3854},
		{3085, 4675, 6918},
		{2177, 3519, 3685},
		{3104, 4650, 6960},
	}
	//data2 := [][]float64{
	//	{2179, 3303, 3449},
	//	{3196, 4377, 5415},
	//	{2288, 3344, 3516},
	//	{3179, 4394, 5423},
	//	{2266, 3230, 3433},
	//	{3164, 4401, 5496},
	//}
	data2 := [][]float64{
		{3269, 3303},
		{4298, 4377},
		{3284, 3344},
		{4313, 4394},
		{3198, 3230},
		{4305, 4401},
	}
	path1 := "D:/Desktop/1.png"
	path2 := "D:/Desktop/3.png"

	err := DrawBarChart(data1, data2, path1, path2)
	if err != nil {
		t.Error(err)
	}
}

func TestDrawNTS(t *testing.T) {
	err := DrawNTS("D:/Desktop")
	if err != nil {
		t.Error(err)
	}
}
