package main

import (
  "github.com/PuerkitoBio/goquery"
  "fmt"
  "sync"
  "log"
  "os"
)

func main() {

  setting, err := read_setting()
  if err != nil {
    return
  }

  urls  := []string {
  }

  // 出力先
  fd := make(map[int]*os.File)
  for i, scrapingItem := range setting.ScrapingItems {
    if scrapingItem.OutputFile != "" {
      f, err := os.OpenFile(scrapingItem.OutputFile, os.O_WRONLY|os.O_CREATE, 0666)
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
  for _, url := range urls {
    ch <- true
    w.Add(1)
    go scraping_url(url, setting, ch, &w, fd)
  }
  w.Wait()

}

/*
  スクレイピングメイン処理
*/
func scraping_url(url string, setting scrapingSetting, ch chan bool, w *sync.WaitGroup, fd map[int]*os.File) {

  defer func() { <-ch }()
  defer w.Done()

  // 初期設定
  doc, err := goquery.NewDocument(url)
  if err != nil {
      log.Fatal(err)
      return
  }

  hasNext := true
  for hasNext {
    fmt.Println(url)
    fmt.Println("---------------------------------")

    for i, scrapingItem := range setting.ScrapingItems {
      fmt.Println(scrapingItem.Name)
      result_line := scraping_items(doc, scrapingItem.Items, scrapingItem.PrintUrl)
      // 出力
      if scrapingItem.OutputFile == "" || fd[i] == nil {
        fmt.Print(result_line)
      } else {
        fd[i].Write([]byte(result_line))
      }
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
func scraping_items(doc *goquery.Document, items []items, printUrl bool) string {

  line := ""
  if printUrl {
    line = format_line(line, doc.Url.String())
  }
  for _, item := range items {

    doc.Find(item.Selector).Each(func(_ int, s *goquery.Selection) {

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
    })
  }
  return line + "\n"
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
