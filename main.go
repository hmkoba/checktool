package main

import (
  "github.com/PuerkitoBio/goquery"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

type items struct {
    Selector string `json:"selector"`
    Attr string `json:"attr"`
    Attr2 string `json:"attr2"`

    Items []struct {
      Selector string `json:"selector"`
      Attr string `json:"attr"`
    } `json:"items"`

}

type scrapingSetting struct {
  Separator string `json:"separator"`
  NextPage struct {
    Selector string `json:"selector"`
    Attr string `json:"attr"`
  } `json:"next_page"`

  ScrapingItems []struct {
    Name string `json:"name"`
    Items []items `json:"items"`
  } `json:"scraping_items"`
}

var separator string  // 出力項目セパレータ（設定jsonから取得）
var setting_path = "./setting.json"

func main() {

  setting, err := read_setting()
  if err != nil {
    return
  }
  fmt.Printf("%+v", setting)

  scraping_url("", setting)

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

  return setting, err
}

/*
  同一グループのスクレイピング結果を１行にまとめる
*/
func format_line(line string, add string) string {
  if line == "" {
    return add
  }
  return line + separator + add
}

/*
  属性取得
*/
func get_attr(s *goquery.Selection, attr string) string {
  if attr == "text" {
    return s.Text()
  } else {
    ret, exists := s.First().Attr(attr)
    if exists {
      return ret
    }
    return ""
  }
}

/*
  スクレイピングメイン処理
*/
func scraping_url(url string, setting scrapingSetting) {

  // 初期設定
  separator = setting.Separator
  doc, err := goquery.NewDocument(url)
  if err != nil {
      log.Fatal(err)
      return
  }

  hasNext := true
  for hasNext {
    fmt.Println(url)
    fmt.Println("---------------------------------")

    for _, scrapingItem := range setting.ScrapingItems {
      fmt.Println(scrapingItem.Name)
      result_line := scraping_items(doc, scrapingItem.Items)
      fmt.Print(result_line)
    }

    if setting.NextPage.Selector != "" {
      doc, hasNext = get_next_document(doc, setting)
    } else {
      hasNext = false
    }
  }
}

/*
  取得
*/
func scraping_items(doc *goquery.Document, items []items) string {

  lines := ""
  for _, item := range items {

    doc.Find(item.Selector).Each(func(_ int, s *goquery.Selection) {
      line := ""

      if item.Attr != "" {
        line = format_line(line, get_attr(s, item.Attr))
      }
      if item.Attr2 != "" {
        line = format_line(line, get_attr(s, item.Attr2))
      }

      for _, child_item := range item.Items {
        s.Find(child_item.Selector).Each(func(_ int, cs *goquery.Selection) {
          line = format_line(line, get_attr(cs, child_item.Attr))
        })
      }
      lines += line + "\n"
    })
  }
  return lines
}

/*
  次ページ取得
*/
func get_next_document(doc *goquery.Document, setting scrapingSetting) (*goquery.Document, bool) {
  url, exists := doc.Find(setting.NextPage.Selector).First().Attr(setting.NextPage.Attr)

  if ! exists || url == "" {
    return doc, false
  }

  n, err := goquery.NewDocument(url)
  if err != nil {
      log.Fatal(err)
      return doc, false
  }
  return n, true
}
