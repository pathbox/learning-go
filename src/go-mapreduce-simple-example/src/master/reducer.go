package master

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	logger "log"
)

func reducer(nr int, fileDir string) {
	mp := make(map[string]int) // store the frequence of words

	// read file and do reduce
	for {
		val, ok := <-ReduceChanIn
		if !ok {
			break
		}
		logger.Print("reducer called: ", nr)
		file, err := os.Open(val)
		if err != nil {
			errMsg := fmt.Sprintf("Read file(%s) error in reducer", val)
			logger.Printf(errMsg)
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()
			arr := strings.Split(str, " ")
			if len(arr) != 2 {
				errMsg := fmt.Sprintf("Read file(%s) error that len of line(%s) != 2(%d) in reducer", val, str, len(arr))
				logger.Printf(errMsg)
				continue
			}
			v, err := strconv.Atoi(arr[1])
			if err != nil {
				errMsg := fmt.Sprintf("Read file(%s) error that line(%s) parse error in reduer", val, str)
				logger.Printf(errMsg)
				continue
			}
			mp[arr[0]] += v
		}
		if err := scanner.Err(); err != nil {
			logger.Print("reducer: reading standard input:", err)
		}
		file.Close()
	}

	outputFilename := path.Join(fileDir, "reduce-output-"+strconv.Itoa(nr))
	outputFileHandler, err := os.Create(outputFilename)
	if err != nil {
		errMsg := fmt.Sprintf("Write file(%s) error in reducer(%d)", outputFilename, nr)
		logger.Printf(errMsg)
	} else {
		for k, v := range mp {
			str := fmt.Sprintf("%s %d\n", k, v)
			outputFileHandler.WriteString(str)
		}
		outputFileHandler.Close()
	}

	ReduceChanOut <- outputFilename
}
