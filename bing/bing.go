package bing

import (
	"encoding/json"
	"errors"
	"fmt"
	"gowallpaper/util"
	"gowallpaper/winapi"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type Wallpaper struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Url    string `json:"url"`
	Date   string `json:"date"`
}

type Bing struct {
	Dir                  string
	BmpTempPath          string
	CurrentWallpaperTime time.Time
	timer                *util.IntervalTimer
}

func GetImageName(imageUrl string) (string, error) {
	u, err := url.Parse(imageUrl)
	if err != nil {
		return "", err
	}
	names, exist := u.Query()["id"]
	if exist && len(names) > 0 {
		return names[0], nil
	}
	return "", errors.New("name not found")
}

func NewBing() *Bing {
	p := new(Bing)
	p.Dir = path.Join(os.TempDir(), "go_wallpaper")
	_ = os.MkdirAll(p.Dir, os.ModePerm)
	p.BmpTempPath = path.Join(p.Dir, "tmp.bmp")
	p.timer = util.NewIntervalTimer(0, func() {})
	log.Println("================================")
	fmt.Println("壁纸缓存目录:", p.Dir)
	return p
}

func (p *Bing) GetUrl(date time.Time) string {
	host := "https://www.ningto.com/api/bingdate"
	if date.Unix() > time.Now().Unix() {
		return host
	}
	return fmt.Sprintf("%s?date=%s", host, util.FormatDate(date))
}

func (p *Bing) GetWallpaper(date time.Time) (Wallpaper, error) {
	url := p.GetUrl(date)
	resp, err := http.Get(url)
	if err != nil {
		return Wallpaper{}, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Wallpaper{}, err
	}

	var wallpaper Wallpaper
	err = json.Unmarshal(b, &wallpaper)
	return wallpaper, err
}

func (p *Bing) FetchAndWrite(url string, localFilePath string) error {
	log.Println("fetch", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(localFilePath, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (p *Bing) SetWallpaper(date time.Time) {
	wallpaper, err := p.GetWallpaper(date)
	if err != nil {
		log.Println("GetWallpaper", err)
		return
	}

	name, err := GetImageName(wallpaper.Url)
	if err != nil {
		log.Println("GetImageName", err)
		return
	}

	jpgPath := path.Join(p.Dir, name)
	bmpPath := p.BmpTempPath
	if !util.PathExist(jpgPath) {
		if err = p.FetchAndWrite(wallpaper.Url, jpgPath); err != nil {
			log.Println("FetchAndWrite", err)
			return
		}
	}
	if err = util.Jpg2Bmp(jpgPath, bmpPath); err != nil {
		log.Println("Jpg2Bmp", err)
		return
	}
	log.Println(wallpaper.Title, wallpaper.Date)
	winapi.SetWallpaper(bmpPath)
	p.CurrentWallpaperTime = date
}

func (p *Bing) RunTask(d time.Duration, f func()) {
	p.timer.Stop()
	p.timer = util.NewIntervalTimer(d, f)
	if d.Seconds() == 0 {
		f()
	} else {
		p.timer.Start()
	}
}

func (p *Bing) Day() {
	p.RunTask(1*time.Hour, func() {
		if time.Now().Day() != p.CurrentWallpaperTime.Day() {
			p.SetWallpaper(time.Now())
		}
	})
}

func (p *Bing) Now() {
	p.RunTask(0, func() {
		p.SetWallpaper(time.Now())
	})
}

func (p *Bing) Prev() {
	p.RunTask(0, func() {
		if p.CurrentWallpaperTime.IsZero() {
			p.CurrentWallpaperTime = time.Now()
		}
		prevTime := p.CurrentWallpaperTime.Add(-1 * 24 * time.Hour)
		p.SetWallpaper(prevTime)
	})
}

func (p *Bing) Next() {
	p.RunTask(0, func() {
		if p.CurrentWallpaperTime.IsZero() {
			p.CurrentWallpaperTime = time.Now()
		}
		nextTime := p.CurrentWallpaperTime.Add(24 * time.Hour)
		p.SetWallpaper(nextTime)
	})
}

func (p *Bing) Rand(d time.Duration) {
	p.RunTask(d, func() {
		n := rand.Intn(100)
		cur := time.Now()
		for i := 0; i < n; i++ {
			cur = cur.Add(-24 * time.Hour)
		}
		p.SetWallpaper(cur)
	})
}
