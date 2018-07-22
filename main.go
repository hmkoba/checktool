package main

import (
  "github.com/PuerkitoBio/goquery"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
)

type scrapingSetting struct {
  NextPage struct {
    Selector string `json:"selector"`
    Attr string `json:"attr"`
  } `json:"next_page"`

  ScrapingItems []struct {
    Name string `json:"name"`
    Items []struct {
      Selector string `json:"selector"`
      Attr string `json:"attr"`
      IsParent bool `json:"is_parent"`

      Items []struct {
        Selector string `json:"selector"`
        Attr string `json:"attr"`
        IsParent bool `json:"is_parent"`
      } `json:"items"`
    } `json:"items"`
  } `json:"scraping_items"`
}

func main() {

  scrapingSetting, err := read_setting()
  if err != nil {
    return
  }
  fmt.Printf("%+v", scrapingSetting)

  scraping_url("", scrapingSetting)

}

func read_setting() (scrapingSetting, error) {
  file_path := "./setting.json"
  // JSONファイル読み込み
  bytes, err := ioutil.ReadFile(file_path)
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

func format_line(line string, add string) string {
  if line == "" {
    return add
  }
  return line + "\t" + add
}

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

func scraping_url(url string, setting scrapingSetting) {

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

      for _, item := range scrapingItem.Items {
        doc.Find(item.Selector).Each(func(_ int, s *goquery.Selection) {
          result_line := ""

          if item.Attr != "" {
            result_line = format_line(result_line, get_attr(s, item.Attr))
          }

          for _, child_item := range item.Items {
            s.Find(child_item.Selector).Each(func(_ int, cs *goquery.Selection) {
              result_line = format_line(result_line, get_attr(cs, child_item.Attr))
            })
          }
          result_line += "\n"
          fmt.Print(result_line)
        })
      }
    }

    if setting.NextPage.Selector != "" {

      url, exists := doc.Find(setting.NextPage.Selector).First().Attr(setting.NextPage.Attr)

      if ! exists || url == "" {
        hasNext = false
      } else {
        doc, err = goquery.NewDocument(url)
        if err != nil {
            log.Fatal(err)
            return
        }
      }
    }
  }
}
