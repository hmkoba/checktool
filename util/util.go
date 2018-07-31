package util

import (
  "github.com/PuerkitoBio/goquery"
  "strings"
  "golang.org/x/text/encoding/japanese"
  "golang.org/x/text/transform"
  "io"
  "io/ioutil"
  "bufio"
  "os"
)

/*
  同一グループのスクレイピング結果を１行にまとめる
*/
func FormatLine(l string, add string, en string, sp string) string {
  if l == "" {
    return en + add + en
  }
  if sp == "" {
    return l + "," + en + add + en
  }
  return l + sp + en + add + en
}

/*
  属性取得
*/
func GetAttr(s *goquery.Selection, attr string) string {
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

func ReadListFile(path string) ([]string, error) {
  l := []string{}

  fp, err := os.Open(path)
  if err != nil {
    return l, err
  }
  defer fp.Close()

  s := bufio.NewScanner(fp)
  for s.Scan() {
    if len(s.Text()) > 0 {
      l = append(l, s.Text())
    }
  }
  return l, err
}

func EncodeString(s string, e string) string {
  if e == "ShiftJIS" {
    return toShiftJIS(s)
  }
  return s
}

func toShiftJIS(s string) string {
    return transformEncoding(strings.NewReader(s), japanese.ShiftJIS.NewEncoder())
}

func transformEncoding(r io.Reader, t transform.Transformer) string {
    ret, err := ioutil.ReadAll(transform.NewReader(r, t))
    if err == nil {
        return string(ret)
    } else {
        return ""
    }
}

func SliceUniq(ar []string) []string {
    r := []string{}
    m := map[string]bool{}
    for _, e := range ar {
        if !m[e] {
            m[e] = true
            r = append(r, e)
        }
    }
    return r
}
