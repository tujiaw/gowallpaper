package util

import (
	"golang.org/x/image/bmp"
	"image/jpeg"
	"os"
	"time"
)

const (
	DateFormat  = "2006-01-02"
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
	return time.Now().Format(DateFormat)
}

func FormatDate(t time.Time)string {
	return t.Format(DateFormat)
}

func FromDate(strTime string)(time.Time, error)   {
	return time.ParseInLocation(DateFormat, strTime, time.Local)
}

func Jpg2Bmp(from string, to string) error {
	imageInput, err := os.Open(from)
	if err != nil {
		return err
	}
	defer imageInput.Close()
	src, err := jpeg.Decode(imageInput)
	if err != nil {
		return err
	}

	outfile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer outfile.Close()
	err = bmp.Encode(outfile, src)
	if err != nil {
		return err
	}
	return nil
}