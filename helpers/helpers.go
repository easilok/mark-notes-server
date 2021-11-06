package helpers

import str "strings"

type String string

func (text String) TitleFromMarkdown () string{
  baseText := string(text)
  splitStr := str.Split(baseText, "\n")

  if len(splitStr) > 0 {
    title := str.TrimSpace(str.ReplaceAll(splitStr[0], "#", ""))
    return title
  }
  return baseText;
}

func (text String) GetFilename (ext string) string {
  baseText := string(text)
  if str.HasSuffix(baseText, ext) {
    return baseText[0:len(baseText) - len(ext)]
  }
  return baseText
}
