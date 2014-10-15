package po

import (
    "strings"
    spec "github.com/Sam-Izdat/pogo/gtspec"
)

func pop(a *[]string) (x string) {
    lna := len(*a)
    if lna < 1 { return }
    x, *a = (*a)[lna-1], (*a)[:lna-1]   
    return
}

func shift(a *[]string) (x string) {
    lna := len(*a)
    if lna < 1 { return }
    x, *a = (*a)[0], (*a)[1:lna]   
    return
}

// Very naive line folding but it should do the trick. Lines will be folded 
// at escaped newlines or along spaces after line exceeds maxLen runes.
func foldLine(str string) (res []string) {
    maxLen := 76

    str = strings.Join(strings.Split(str, `\n`), `\n`+"\n")
    lines := strings.Split(str, "\n")
    for _, line := range lines {
        a := []rune(line)
        tline := ""
        for i, r := range a {
            tline = tline + string(r)
            if i == len(a) - 1 || (len(tline) > maxLen && string(r) == " ") {
                res = append(res, tline)
                tline = ""
            }
        }        
    }
    return
}

// Removes duplicate messages (identical ctxt and id); consolidates references
func RemoveDuplicates(msgs *[]spec.Msg) {
    found := make(map[string]int)
    j := 0
    for i, x := range *msgs {
        if found[x.Ctxt+"\x04"+x.Id] == 0 {
            found[x.Ctxt+"\x04"+x.Id] = j+1
            (*msgs)[j] = (*msgs)[i]
            j++
        } else {
            (*msgs)[found[x.Ctxt+"\x04"+x.Id]-1].Comments["reference"] = append(
                (*msgs)[found[x.Ctxt+"\x04"+x.Id]-1].Comments["reference"], 
                x.Comments["reference"]...)
        }
    }
    *msgs = (*msgs)[:j]
}