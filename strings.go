package main

import "fmt"

type fStr string // fancy, fancy string

func (s fStr) s(style string) string {
    if isWindows() { return fmt.Sprintf("%s", s) }
    var fancyString string
    switch style {
    case "red":
        fancyString = fmt.Sprintf("\x1b[0;31;49m%s\x1b[0m", s)
    case "blue":
        fancyString = fmt.Sprintf("\x1b[0;34;49m%s\x1b[0m", s)
    case "green":
        fancyString = fmt.Sprintf("\x1b[0;32;49m%s\x1b[0m", s)
    case "yellow":
        fancyString = fmt.Sprintf("\x1b[0;33;49m%s\x1b[0m", s)
    case "bold":
        fancyString = fmt.Sprintf("\x1b[01;33;49m%s\x1b[0m", s)
    default:
        fancyString = fmt.Sprintf("%s", s)
    }
    return fancyString
}

var pWarn, pNotice, pSuccess = fStr("! pogo:").s("red"), fStr("* pogo:").s("yellow"), fStr("* pogo:").s("green")