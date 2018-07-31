package url
import (
  "log"
  "fmt"
  "strings"
  "strconv"
  "github.com/hmkoba/checktool/util"
)

func CreateUrlList() []string {
  origin := "https://test.jp/{0}/{1}/z-{4}/c-{2}/t-{3}/list"

  p := [][]string {{"a","b"}}
  fmt.Println(p[0][0])

  // 入力元
  urls, err := util.ReadListFile("../j.csv")
  if err != nil {
    log.Fatal(err)
  }

  result := []string{}

  for i, s := range urls {
    if i == 0 { continue }

    items := strings.Split(s, ",")

    tmp := origin
    for i, item := range items {
      tmp = strings.Replace(tmp, "{" + strconv.Itoa(i) + "}", item, -1)
    }
    result = append(result, tmp)
  }
  result = util.SliceUniq(result)

  return result
}
