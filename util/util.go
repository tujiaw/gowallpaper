package util

import (
	"os"
	"time"
)

const (
	DATE_FORMAT = "2006-01-02"
	TIME_FORMAT = "2006-01-02 15:04:05"
)

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func CurrentDate() string {
	return time.Now().Format(DATE_FORMAT)
}

func FormatDate(t time.Time)string {
	return t.Format(DATE_FORMAT)
}

func FromDate(strTime string)(time.Time, error)   {
	return time.ParseInLocation(DATE_FORMAT, strTime, time.Local)
}
