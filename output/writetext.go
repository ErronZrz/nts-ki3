package output

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	configPath         = "../resources/"
	outputPathKey      = "output.dir_path"
	fileTimeFormat     = "/2006-01-02_15-04-05_"
	dividingLineFormat = "------------ 15:04:05.000 ------------\n"
	beforeParsed       = "--- parsed ---\n"
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

func WriteToFile(raw, parsed, info string, seq int, rcvTime, now time.Time) {
	dirPath := viper.GetString(outputPathKey)
	info = strings.Replace(info, "/", "_", 1)
	filePath := dirPath + now.Format(fileTimeFormat) + info + ".txt"

	seqLine := "#" + strconv.Itoa(seq) + "\n"
	dividingLine := rcvTime.Format(dividingLineFormat)

	commonWrite(filePath, []string{seqLine, dividingLine, raw, beforeParsed, parsed})
}

func commonWrite(filePath string, strs []string) {
	var file *os.File

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("error creating file %s: %v", filePath, err)
			return
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("error opening file %s: %v", filePath, err)
			return
		}
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("error closing file %s: %v", filePath, err)
		}
	}(file)

	writer := bufio.NewWriter(file)

	for _, s := range strs {
		_, err = writer.WriteString(s)
		if err != nil {
			fmt.Printf("error writing string %s: %v", s, err)
			return
		}
	}

	err = writer.Flush()
	if err != nil {
		fmt.Printf("error flushing writer: %v", err)
	}
}
