package url
import (
  "log"
  "strings"
  "fmt"
  "github.com/hmkoba/checktool/util"
)

func CreateUrlListSp() (map[string]string, map[string][]string) {

  checkList := make(map[string][]string)

  result := make(map[string]string)
  setting, err := readSetting("")
  if err != nil {
    return result, checkList
  }

  // 入力元
  lines, err := util.ReadListFile(setting.InputFile)
  if err != nil {
    log.Fatal(err)
  }


  for i, line := range lines {
    if setting.LineHeader && i == 0 { continue }

    url := setting.Domain + "/"

    items := strings.Split(line, ",")

    // TODO
    key := items[0] + "_" + items[1] + "_" + items[2] + "_" + items[3] + "_" + items[4] + "_" + items[7]
    checkList[key] = append(checkList[key], items[5] + "_" + items[6])

    for _, urlItem := range setting.UrlItems {
      switch urlItem.Attr {
        case "input" :
          if len(urlItem.UrlItems) > 0 {
            if len(urlItem.Exclusion) > 0 && urlItem.Exclusion == items[urlItem.Column] {
              continue
            }
            for _, childItem := range urlItem.UrlItems {
              url += strings.Replace(childItem.Out, "{$}", items[childItem.Column], -1) + "/"
            }

          } else if len(items[urlItem.Column]) > 0 {
            url += strings.Replace(urlItem.Out, "{$}", items[urlItem.Column], -1) + "/"
          }

        case "fix" :
          url += urlItem.Out + "/"
      }

    }
    if _, ok := result[key]; ! ok {
      result[key] = url
    }
  }

  fmt.Println(checkList)

  return result, checkList
}
