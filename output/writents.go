package output

import (
	"github.com/spf13/viper"
	"time"
)

const (
	ntsOutputPathKey      = "output.nts_path"
	ntsFileTimeFormat     = "/2006-01-02_"
	ntsDividingLineFormat = "------------ 15:04:05 ------------\n"
)

func init() {
	viper.AddConfigPath(configPath)
	viper.SetConfigType("yaml")
	viper.SetConfigName("properties")
	err := viper.ReadInConfig()
	if err != nil {
		// fmt.Printf("error reading resource file: %v", err)
	}
}

func WriteNTSToFile(raw, parsed, host string) {
	now := time.Now()
	dirPath := viper.GetString(ntsOutputPathKey)
	filePath := dirPath + now.Format(ntsFileTimeFormat) + host + ".txt"
	dividingLine := now.Format(ntsDividingLineFormat)

	commonWrite(filePath, []string{dividingLine, raw, beforeParsed, parsed, "\n\n"})
}

func WriteNTSDetectToFile(content, host string) {
	now := time.Now()
	dirPath := viper.GetString(ntsOutputPathKey)
	filePath := dirPath + now.Format(ntsFileTimeFormat) + host + "_detect.txt"
	dividingLine := now.Format(ntsDividingLineFormat)

	commonWrite(filePath, []string{dividingLine, content, "\n\n"})
}
