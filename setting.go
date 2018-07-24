package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

var setting_path = "./setting.json"
var separator string  // 出力項目セパレータ（設定jsonから取得）

type items struct {
    Selector string `json:"selector"`
    Attr string `json:"attr"`
    Attr2 string `json:"attr2"`

    Items []struct {
      Selector string `json:"selector"`
      Attr string `json:"attr"`
    } `json:"items"`

}

type scrapingItems struct {
  Name string `json:"name"`
  OutputFile string `json:"output_file"`
  PrintUrl bool `json:"print_url"`
  Items []items `json:"items"`
}

type scrapingSetting struct {
  Parallel int `json:"parallel"`
  LineHeader bool `json:"lineheader"`
  Separator string `json:"separator"`
  NextPage struct {
    Selector string `json:"selector"`
    Attr string `json:"attr"`
  } `json:"next_page"`
  ScrapingItems []scrapingItems `json:"scraping_items"`
}

/*
  スクレイピングの定義をjsonファイルから取得する
*/
func read_setting() (scrapingSetting, error) {
  // JSONファイル読み込み
  bytes, err := ioutil.ReadFile(setting_path)
  if err != nil {
    fmt.Print("JSON read error:")
    log.Fatal(err)
  }
  // JSONデコード
  var setting scrapingSetting
  err = json.Unmarshal(bytes, &setting)
  if err != nil {
    fmt.Print("JSON decode error:")
    log.Fatal(err)
  }

  if setting.Parallel <= 0 {
    setting.Parallel = 1
  }

  separator = setting.Separator
  if separator == "" { separator = "," }

  return setting, err
}
