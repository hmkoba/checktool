package main

import (
  "github.com/PuerkitoBio/goquery"
)

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
