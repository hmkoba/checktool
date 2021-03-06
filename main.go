package main

import (
  "github.com/PuerkitoBio/goquery"
  "fmt"
  "sync"
  "log"
  "os"
  "flag"
  "net/http"
  "github.com/hmkoba/checktool/util"
//  "github.com/hmkoba/checktool/url"
)

func main() {
  flag.Parse()
  setting, err := readSetting(flag.Arg(0))
  if err != nil {
    return
  }

//urls := url.CreateUrlList();

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
  for _, url := range urls {
    ch <- true
    w.Add(1)
    go scrapingUrl(url, setting, ch, &w, fd)
  }
  w.Wait()

}

/*
  スクレイピングメイン処理
*/
func scrapingUrl(url string, setting scrapingSetting, ch chan bool, w *sync.WaitGroup, fd map[int]*os.File) {

  defer func() { <-ch }()
  defer w.Done()

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
      log.Fatal(err)
      return
  }
  if(setting.UserAgent != "") {
    req.Header.Add("User-Agent", setting.UserAgent)
  }
  cl := &http.Client{}
  res, err := cl.Do(req)

  // 初期設定
  doc, err := goquery.NewDocumentFromResponse(res)
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

      for _, line := range scrapingDocument(doc, scrapingItem) {
        // 出力
        if scrapingItem.OutputFile == "" || fd[i] == nil {
          fmt.Println(line)
        } else {
          fd[i].Write([]byte(util.EncodeString(line, scrapingItem.Encode)+"\n"))
        }
      }
    }

    if setting.NextPage.Selector != "" {
      doc, hasNext = getNextDocument(doc, setting)
    } else {
      hasNext = false
    }
  }
}

/*
  取得
*/
func scrapingDocument(doc *goquery.Document, sc scrapingItems) []string {

  result := []string{}
  line := initLine(doc.Url.String(), sc.PrintUrl, sc.Enclose)

  for _, item := range sc.Items {

    is := doc.Find(item.Selector)
    if(is.Length() <= 0) {
      line = util.FormatLine(line, "", sc.Enclose, sc.Separator)
      continue
    }
    is.Each(func(_ int, s *goquery.Selection) {
      if item.Attr != "" {
        line = util.FormatLine(line, util.GetAttr(s, item.Attr), sc.Enclose, sc.Separator)
      }
      if item.Attr2 != "" {
        line = util.FormatLine(line, util.GetAttr(s, item.Attr2), sc.Enclose, sc.Separator)
      }

      if len(item.Items) > 0 {
        cl := ""
        for _, child_item := range item.Items {
          cis := s.Find(child_item.Selector)
          if(cis.Length() <= 0) {
            cl = util.FormatLine(cl, "", sc.Enclose, sc.Separator)
            continue
          }
          cis.Each(func(_ int, cs *goquery.Selection) {
            cl = util.FormatLine(cl, util.GetAttr(cs, child_item.Attr), sc.Enclose, sc.Separator)
          })
        }
        result = append(result, line + sc.Separator + cl)
        line = initLine(doc.Url.String(), sc.PrintUrl, sc.Enclose)
      }
    })
  }

  if len(result) == 0 {
    result = append(result, line)
  }
  return result
}

func initLine(s string, p bool, e string) string{
  if p {
    return util.FormatLine("", s, e, "")
  }
  return ""

}

/*
  次ページ取得
*/
func getNextDocument(doc *goquery.Document, setting scrapingSetting) (*goquery.Document, bool) {
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
