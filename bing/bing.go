package bing

import (
	"encoding/json"
	"errors"
	"gowallpaper/util"
	"gowallpaper/winapi"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"time"
)

type Wallpaper struct {
	Date string `json:"date"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

type Bing struct {
	Url      string
	Dir      string
	JsonPath string
	BmpTempPath string
	WallpaperList []Wallpaper
	CurrentWallpaperTime time.Time
	ticker *time.Ticker
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

func FindWallpaper(list []Wallpaper, date string)(Wallpaper, error) {
	for _, v := range list {
		if v.Date == date {
			return v, nil
		}
	}
	return Wallpaper{}, errors.New("not found")
}

func NewBing() *Bing {
	p := new(Bing)
	p.Url = "https://www.ningto.com/public/bing/his.json"
	p.Dir = path.Join(os.TempDir(), "go_wallpaper")
	_ = os.MkdirAll(p.Dir, os.ModePerm)
	p.JsonPath = path.Join(p.Dir, "his.json")
	p.BmpTempPath = path.Join(p.Dir, "tmp.bmp")
	log.Println("tmp dir:", p.Dir)
	if err := p.Init(); err != nil {
		panic(err)
	}
	return p
}

func(p *Bing)Init() error {
	if !util.PathExist(p.JsonPath) {
		if err := p.FetchAndWrite(p.Url, p.JsonPath); err != nil {
			return err
		}
	}
	b, err := ioutil.ReadFile(p.JsonPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &p.WallpaperList); err != nil {
		return err
	}
	sort.Slice(p.WallpaperList, func(i, j int)bool {
		return p.WallpaperList[i].Date > p.WallpaperList[j].Date
	})
	return nil
}

func(p *Bing)FetchAndWrite(url string, localFilePath string)error {
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

func(p *Bing)SetWallpaper(date string) {
	if len(p.WallpaperList) == 0 {
		log.Println("wallpaper list is empty")
		return
	}

	if p.WallpaperList[0].Date != util.CurrentDate() {
		if err := os.Remove(p.JsonPath); err != nil {
			log.Println(err)
		}
		if err := p.Init(); err != nil {
			log.Println(err)
		}
	}

	v, err := FindWallpaper(p.WallpaperList, date)
	if err != nil {
		log.Println(err)
		return
	}

	name, err := GetImageName(v.Url)
	if err != nil {
		log.Println(err)
		return
	}

	jpgPath := path.Join(p.Dir, name)
	bmpPath := p.BmpTempPath
	if !util.PathExist(jpgPath) {
		if err = p.FetchAndWrite(v.Url, jpgPath); err != nil {
			log.Println(err)
			return
		}
	}
	if err = util.Jpg2Bmp(jpgPath, bmpPath); err != nil {
		log.Println(err)
		return
	}
	log.Println(v.Title, v.Date)
	winapi.SetWallpaper(bmpPath)
	p.CurrentWallpaperTime, _ = util.FromDate(date)
}

func(p *Bing)RunTask(d time.Duration, f func()) {
	if p.ticker != nil {
		p.ticker.Stop()
		p.ticker = nil
	}

	if d.Seconds() == 0 {
		f()
		return
	}

	p.ticker = time.NewTicker(d)
	go func() {
		for range p.ticker.C {
			f()
		}
	}()
}

func(p *Bing)Day() {
	p.RunTask(1 * time.Hour, func() {
		if time.Now().Day() != p.CurrentWallpaperTime.Day() {
			p.SetWallpaper(util.FormatDate(time.Now()))
		}
	})
}

func(p *Bing)Now() {
	p.RunTask(0, func() {
		p.SetWallpaper(util.FormatDate(time.Now()))
	})
}

func(p *Bing)Prev() {
	p.RunTask(0, func() {
		if p.CurrentWallpaperTime.IsZero() {
			p.CurrentWallpaperTime = time.Now()
		}
		prevTime := p.CurrentWallpaperTime.Add(-1 * 24 * time.Hour)
		p.SetWallpaper(util.FormatDate(prevTime))
	})
}

func(p *Bing)Next() {
	p.RunTask(0, func() {
		if p.CurrentWallpaperTime.IsZero() {
			p.CurrentWallpaperTime = time.Now()
		}
		nextTime := p.CurrentWallpaperTime.Add(24 * time.Hour)
		p.SetWallpaper(util.FormatDate(nextTime))
	})
}

func(p *Bing)Rand(d time.Duration) {
	p.RunTask(d, func() {
		n := rand.Intn(len(p.WallpaperList))
		if n >= 0 && n < len(p.WallpaperList) {
			p.SetWallpaper(p.WallpaperList[n].Date)
		}
	})
}
