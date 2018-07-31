package url

import (
    "fmt"
    "testing"
    "github.com/hmkoba/checktool/url"
)

func TestCreateUrlListNormal(t *testing.T) {
  urls := url.CreateUrlList()

  for _, url := range urls {
    fmt.Println(url)
  }
  if len(urls) != 13 {
    t.Errorf("エラー")
  }
}
