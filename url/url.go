package url
import (
  "log"
  "strings"
  "github.com/hmkoba/checktool/util"
)

func CreateUrlList() []string {

  result := []string{}
  setting, err := readSetting("")
  if err != nil {
    return result
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
    result = append(result, url)
  }
  result = util.SliceUniq(result)

  return result
}
