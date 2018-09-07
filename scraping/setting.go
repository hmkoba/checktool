package scraping

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

var setting_path = "./settings/setting.json"

type Item struct {
    Selector string   `json:"selector"`
    Attr string       `json:"attr"`
    Attr2 string      `json:"attr2"`

    Items []struct {
      Selector string `json:"selector"`
      Attr string     `json:"attr"`
    }                 `json:"items"`

}

type ScrapingItem struct {
  Name string       `json:"name"`
  OutputFile string `json:"output_file"`
  Encode string     `json:"encode"`
  Enclose string    `json:"enclose"`
  Separator string  `json:"separator"`
  InnerSeparator string  `json:"inner_separator"`
  PrintUrl bool     `json:"print_url"`
  Items []Item     `json:"items"`
}

type ScrapingSetting struct {
  Parallel int      `json:"parallel"`
  LineHeader bool   `json:"lineheader"`
  UrlFile string    `json:"url_file"`
  UserAgent string  `json:"user_agent"`
  NextPage struct {
    Selector string `json:"selector"`
    Attr string     `json:"attr"`
  }                 `json:"next_page"`
  ScrapingItems []ScrapingItem `json:"scraping_items"`
}

/*
  スクレイピングの定義をjsonファイルから取得する
*/
func ReadSetting(path string) (ScrapingSetting, error) {
  var setting ScrapingSetting

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
