package main

import (
  "github.com/PuerkitoBio/goquery"
  "strings"
  "golang.org/x/text/encoding/japanese"
  "golang.org/x/text/transform"
  "io"
  "io/ioutil"
)

/*
  同一グループのスクレイピング結果を１行にまとめる
*/
func formatLine(line string, add string, en string) string {
  if line == "" {
    return en+add+en
  }
  return line + separator + en+add+en
}

/*
  属性取得
*/
func getAttr(s *goquery.Selection, attr string) string {
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

func encodeString(s string, e string) string {
  if e == "ShiftJIS" {
    return toShiftJIS(s)
  }
  return s
}

func toShiftJIS(s string) string {
    return transformEncoding(strings.NewReader(s), japanese.ShiftJIS.NewEncoder())
}

func transformEncoding( r io.Reader, t transform.Transformer) string {
    ret, err := ioutil.ReadAll(transform.NewReader(r, t))
    if err == nil {
        return string(ret)
    } else {
        return ""
    }
}
