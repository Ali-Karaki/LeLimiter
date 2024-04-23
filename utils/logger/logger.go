package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type logWriter struct {
	logFile *os.File
}

var LoggerPath = "utils/logger/"

var NoLuaPath = LoggerPath + "logs/nolua.log"
var LuaPath = LoggerPath + "logs/lua.log"
var LuaShaPath = LoggerPath + "logs/luasha.log"

var StatsPath = LoggerPath + "logs/stats.log"

var AwkPath = LoggerPath + "stats.awk"


func (writer *logWriter) Write(bytes []byte) (int, error) {
	timestamp := time.Now().UTC().Format("15:04:05.000000")
	formattedMessage := fmt.Sprintf("%s | %s", timestamp, string(bytes))
	return writer.logFile.Write([]byte(formattedMessage))
}

func InitLogger() {
	os.Remove(LuaPath) 
	os.Remove(NoLuaPath)
	os.Remove(LuaShaPath)
	os.Remove(StatsPath)

	logFile, err := os.OpenFile(LuaPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	customWriter := &logWriter{
		logFile: logFile,
	}

	log.SetFlags(0)
	log.SetOutput(customWriter)
}

func SwitchLogFile(path string) {
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(&logWriter{
		logFile: logFile,
	})
}

func RunAwk() {
    statsFile, err := os.Create(StatsPath)
    if err != nil {
        fmt.Println("Error creating output file:", err)
        return
    }
    defer statsFile.Close() 

    cmd := exec.Command("awk", "-f", AwkPath, NoLuaPath, LuaPath, LuaShaPath)
    cmd.Stdout = statsFile

    stderr := new(bytes.Buffer)
    cmd.Stderr = stderr

    err = cmd.Run()
    statsFile.Close() 
	
    if err != nil {
        fmt.Printf("Command execution failed: %s\n", err)
    }
    if stderr.Len() > 0 {
        fmt.Printf("Stderr: %s\n", stderr.String())
    } else {
        fmt.Println("No errors captured in stderr.")
    }

    if err == nil {
        fmt.Println("Command executed successfully, output saved to", StatsPath)
    }
}
