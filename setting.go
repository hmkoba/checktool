package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

var setting_path = "./settings/setting.json"

type items struct {
    Selector string   `json:"selector"`
    Attr string       `json:"attr"`
    Attr2 string      `json:"attr2"`

    Items []struct {
      Selector string `json:"selector"`
      Attr string     `json:"attr"`
    }                 `json:"items"`

}

type scrapingItems struct {
  Name string       `json:"name"`
  OutputFile string `json:"output_file"`
  Encode string     `json:"encode"`
  Enclose string    `json:"enclose"`
  Separator string  `json:"separator"`
  PrintUrl bool     `json:"print_url"`
  Items []items     `json:"items"`
}

type scrapingSetting struct {
  Parallel int      `json:"parallel"`
  LineHeader bool   `json:"lineheader"`
  UrlFile string    `json:"url_file"`
  UserAgent string  `json:"user_agent"`
  NextPage struct {
    Selector string `json:"selector"`
    Attr string     `json:"attr"`
  }                 `json:"next_page"`
  ScrapingItems []scrapingItems `json:"scraping_items"`
}

/*
  スクレイピングの定義をjsonファイルから取得する
*/
func readSetting(path string) (scrapingSetting, error) {
  var setting scrapingSetting

  if(path != "") {
    setting_path = path
  }
  // JSONファイル読み込み
  bytes, err := ioutil.ReadFile(setting_path)
  if err != nil {
    fmt.Print("JSON read error:")
    log.Fatal(err)
    return setting, err
  }
  // JSONデコード
  err = json.Unmarshal(bytes, &setting)
  if err != nil {
    fmt.Print("JSON decode error:")
    log.Fatal(err)
    return setting, err
  }

  if setting.Parallel <= 0 {
    setting.Parallel = 1
  }

  return setting, err
}
