package utils

import (
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestCSVConvert(t *testing.T) {
	path := "C:\\Corner\\TMP\\BisheData\\data.csv"
	outPath := "C:\\Corner\\TMP\\BisheData\\data1.csv"
	err := ProcessCSV(path, outPath)
	if err != nil {
		t.Errorf("处理CSV文件出错: %v", err)
	}
}

// 将十六进制字符串（如 0x000002A4）转为 []byte
func HexStringToBytes(hexStr string) ([]byte, error) {
	trimmed := strings.TrimPrefix(hexStr, "0x")
	if len(trimmed) < 8 {
		trimmed = strings.Repeat("0", 8-len(trimmed)) + trimmed
	}
	return hex.DecodeString(trimmed)
}

// 主处理函数：读取输入文件并写入输出文件
func ProcessCSV(inputPath, outputPath string) error {
	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开输入文件失败: %v", err)
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	reader.FieldsPerRecord = 3

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer func() { _ = outFile.Close() }()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	for {
		record, err := reader.Read()
		if err != nil {
			break // EOF 或其他错误
		}

		score := record[0]

		rootDelayBytes, err := HexStringToBytes(record[1])
		if err != nil {
			return fmt.Errorf("解析 root_delay 失败: %v", err)
		}
		rootDispersionBytes, err := HexStringToBytes(record[2])
		if err != nil {
			return fmt.Errorf("解析 root_dispersion 失败: %v", err)
		}

		rootDelay := RootDelayToValue(rootDelayBytes)
		rootDispersion := RootDelayToValue(rootDispersionBytes)

		newRecord := []string{
			score,
			fmt.Sprintf("%.10f", rootDelay),
			fmt.Sprintf("%.10f", rootDispersion),
		}
		_ = writer.Write(newRecord)
	}

	return nil
}
