package url

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

var setting_path = "settings/url.json"

type urlItem struct {
  Attr string     `json:"attr"`
  Column int      `json:"column"`
  Out string      `json:"out"`
  Condition string `json:"condition"`
  Exclusion string `json:"exclusion"`
  UrlItems []struct {
    Attr string     `json:"attr"`
    Column int      `json:"column"`
    Out string      `json:"out"`
  }  `json:"url_items"`
}

type urlSetting struct {
    InputFile string  `json:"input_file"`
    LineHeader bool   `json:"lineheader"`
    Domain string     `json:"domain"`

    UrlItems []urlItem  `json:"url_items"`
}

/*
  スクレイピングの定義をjsonファイルから取得する
*/
func readSetting(path string) (urlSetting, error) {
  var setting urlSetting

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

  return setting, err
}
