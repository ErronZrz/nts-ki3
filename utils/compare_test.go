package utils

import "testing"

func TestCompare(t *testing.T) {
	path1 := "D:\\Desktop\\TMP\\ntpdata\\2024-04-10_ntske_1.txt"
	path2 := "D:\\Desktop\\TMP\\ntpdata\\2024-04-11_ntske_1.txt"
	CompareIps(path1, path2)
}
