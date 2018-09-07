package main

import (
  "sync"
  "log"
  "os"
  "flag"
  "fmt"
  "strings"
  "github.com/hmkoba/checktool/util"
  "github.com/hmkoba/checktool/scraping"
  "github.com/hmkoba/checktool/url"
)

func sp() {
  setting, err := scraping.ReadSetting("settings/setting_sp.json")
  if err != nil {
    return
  }
  urls, checkList := url.CreateUrlListSp();

  var w sync.WaitGroup
  ch := make(chan bool, setting.Parallel)
  rch := make(chan []string, setting.Parallel)

  for key, url := range urls {
    ch <- true
    w.Add(1)
    go scraping.ScrapingUrl(url, setting, ch, rch, &w, nil)
    scRet := <- rch
    fmt.Println(key)
    list := checkList[key]
    for _, l := range list {

      r := func () bool {
        for _, k := range scRet {
          items := strings.Split(k, "\t")
          w := strings.Split(strings.Trim(items[1], "[]"),",")
          tmp_key := items[0] + "_" + strings.Trim(w[5], "'")
          if l == tmp_key {
            return true
          }
        }
        return false
      }

      if false == r() {
        fmt.Println("NG:" + l)
      } else {
        fmt.Println("OK:" + l)

      }
    }
  }
  w.Wait()
}

func main() {
  flag.Parse()

  if flag.Arg(0) == "sp" {
    sp()
    return;
  }

  setting, err := scraping.ReadSetting(flag.Arg(0))
  if err != nil {
    return
  }

  // 入力元
  urls, err := util.ReadListFile(setting.UrlFile)
  if err != nil {
    log.Fatal(err)
    return
  }

  // 出力先
  fd := make(map[int]*os.File)
  for i, scrapingItem := range setting.ScrapingItems {
    if scrapingItem.OutputFile != "" {
      f, err := os.OpenFile(scrapingItem.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
      if err != nil {
          //エラー処理
          log.Fatal(err)
          return
      }
      fd[i] = f
    }
  }

  defer func() {
          for _, f := range fd {
            f.Close()
          }
        }()

  var w sync.WaitGroup
  ch := make(chan bool, setting.Parallel)
  rch := make(chan []string, setting.Parallel)
  for _, url := range urls {
    ch <- true
    w.Add(1)
    go scraping.ScrapingUrl(url, setting, ch, rch, &w, fd)
    <- rch
  }
  w.Wait()

}
