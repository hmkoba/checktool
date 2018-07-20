package main

import (
  "github.com/PuerkitoBio/goquery"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

type scrapingItem struct {
  selector string `json:"selector"`
  attr string `json:"attr"`
}

func main() {

  items, err := read_setting()
  if err != nil {
    return
  }
  for _, item := range items {
    fmt.Printf("%s : %s\n", item.selector, item.attr)
}
//  scraping_url("")

}
func read_setting() ([]scrapingItem, error) {
  file_path := "./setting.json"
  // JSONファイル読み込み
  bytes, err := ioutil.ReadFile(file_path)
  if err != nil {
    log.Fatal(err)
  }
  // JSONデコード
  var items []scrapingItem
  err = json.Unmarshal(bytes, &items)
  if err != nil {
    log.Fatal(err)
  }

  return items, err
}

func scraping_url(url string) {

  doc, err := goquery.NewDocument(url)
  if err != nil {
      fmt.Print("url scarapping failed:"+url)
      return
  }

  hasNext, exists := true, true
  for hasNext {
    fmt.Println(url)
    fmt.Println("---------------------------------")
    doc.Find("h2.p-result__name").Each(func(_ int, s *goquery.Selection) {
          fmt.Println(s.Text())
    })
    url, exists = doc.Find("link[rel^='next']").First().Attr("href")
    if exists {
      doc, err = goquery.NewDocument(url)
      if err != nil {
          fmt.Print("url scarapping failed")
          return
      }
    }
  }
}