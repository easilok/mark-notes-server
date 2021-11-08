package helpers

import (
  "os"
  str "strings"
)

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

func TrashFile(rootFolder string, filepath string) error{

  fullPath := rootFolder + string(os.PathSeparator) + filepath
  trashFolderPath := rootFolder + string(os.PathSeparator) + "trash"
  trashedFilePath := trashFolderPath + string(os.PathSeparator) + filepath

  if _, err := os.Open(trashFolderPath); err != nil {
    os.Mkdir(trashFolderPath, 0755)
  }

  err := os.Rename(fullPath, trashedFilePath)

  return err
}
