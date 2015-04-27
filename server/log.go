package server

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	LogTrace   *log.Logger
	LogInfo    *log.Logger
	LogWarning *log.Logger
	LogError   *log.Logger
	LogToFile  *log.Logger
)

func LogInit(traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer, prefix string) {

	LogTrace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime)

	LogInfo = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime)

	LogWarning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime)

	LogError = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime)

	DateNow := time.Now().Local().Format("20060102")
	FileName := "log/Log-" + DateNow + ".txt"

	file, err := os.OpenFile(FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file")
	}

	LogToFile = log.New(file,
		prefix+": ",
		log.Ldate|log.Ltime)

}

func Logging(msg string, prefix string) {
	//log record of everything

	upperPrefix := strings.ToUpper(prefix)

	LogInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr, upperPrefix)

	if upperPrefix == "ERROR" {
		LogError.Println(msg)
	} else if upperPrefix == "INFO" {
		LogInfo.Println(msg)
	} else if upperPrefix == "WARNING" {
		LogWarning.Println(msg)
	} else {
		LogTrace.Println(msg)
	}
	LogToFile.Println(msg)
}

func GetLogFileData(dateStr string, timeStr string) (string, error) {
	dateRep := strings.Replace(dateStr, "/", "", -1)
	fileName := "log/log-" + dateRep + ".txt"
	//fileName := "log/Log-20150423.txt"
	file, err := os.Open(fileName)

	if err != nil {
		return ("No log available!"), nil
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	logInfo := "Log Data at " + dateStr + " with time > " + timeStr
	timeFilter, _ := time.Parse("2006/01/02 15:04:05", dateStr+" "+timeStr)
	for scanner.Scan() {
		logStr := scanner.Text()
		splitLog := strings.Split(logStr, " ")
		logTimeStr := splitLog[2]
		logTime, _ := time.Parse("2006/01/02 15:04:05", dateStr+" "+logTimeStr)
		if logTime.Sub(timeFilter).Seconds() > 0 {
			logInfo = logInfo + "\n" + logStr + " (" + strconv.FormatFloat(logTime.Sub(timeFilter).Seconds(), 'f', 6, 64) + ")"
		}
		//timeDiff := strconv.FormatFloat(logTime.Sub(timeFilter).Seconds(), 'f', 6, 64)
		//logInfo = logInfo + "\n" + logTimeStr + " -- " + timeStr + " (" + timeDiff + ")"
	}
	return logInfo, nil
}
