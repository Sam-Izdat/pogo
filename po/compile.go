package po

import (
	spec "github.com/Sam-Izdat/pogo/gtspec"
	"strconv"
	"strings"
	"time"
)

func drawComments(comments spec.CommentPack) string {
	var lines []string
	ref := []spec.CommentSpec{
		{"translator", "# "},
		{"reference", "#: "},
		{"extracted", "#. "},
		{"flag", "#, "},
		{"previous", "#| "},
	}

	for _, t := range ref {
		if len(comments[t.Key]) == 0 {
			continue
		}
		for _, cv := range comments[t.Key] {
			lines = append(lines, t.Prefix+cv)
		}
	}
	return strings.Join(lines, "\n")
}

func drawBlock(msg spec.Msg, target string) string {
	var (
		response    []string
		comments    spec.CommentPack = msg.Comments
		msgctxt     string           = msg.Ctxt
		msgid       string           = msg.Id
		msgidPlural string           = msg.IdPlural
		msgstr      string           = msg.Str
	)

	if len(comments) > 0 {
		commentsStr := drawComments(comments)
		if len(commentsStr) > 0 {
			response = append(response, commentsStr)
		}
	}

	if len(msgctxt) > 0 {
		response = append(response, addPOString("msgctxt", msgctxt))
	}

	response = append(response, addPOString("msgid", msgid))

	if len(msgidPlural) > 0 {
		response = append(response, addPOString("msgid_plural", msgidPlural))
		if target == "" {
			response = append(response, addPOString("msgstr[0]", ""))
			response = append(response, addPOString("msgstr[1]", ""))
		} else {
			// Get number of plurals
			loc := strings.Split(target, "_")[0]
			tmp := spec.Plurals[loc].Header()
			tmp = strings.Split(tmp, ";")[0]
			tmp = strings.Split(tmp, "=")[1]
			nplurals, err := strconv.Atoi(tmp)
			if err != nil {
				response = append(response, addPOString("msgstr[0]", ""))
			} else {
				for i := 0; i < nplurals; i++ {
					response = append(response, addPOString("msgstr["+strconv.Itoa(i)+"]", ""))
				}
			}
		}
	} else {
		response = append(response, addPOString("msgstr", msgstr))
	}

	return strings.Join(response, "\n")
}

func addPOString(key, value string) string {

	var lines []string = foldLine(value)

	switch len(lines) {
	case 0:
		return key + ` ""`
	case 1:
		return key + ` "` + lines[0] + `"`
	default:
		return key + ` ""` + "\n" + `"` + strings.Join(lines, "\"\n\"") + `"`
	}
}

func Compile(msgs []spec.Msg, target, name, pf string) string {
	cdate, rdate := time.Now().Local().String(), "YEAR-MO-DA HO:MI +ZONE"
	ltrans, lteam := "FULL NAME <EMAIL@ADDRESS>", "TEAM NAME <EMAIL@ADDRESS>"
	mimever, cttype, ctenc := "1.0", "text/plain; charset=UTF-8", "8bit"

	var tmp []string
	for _, s := range strings.Split(o.Po.Comment, "\n") {
		tmp = append(tmp, "# "+s)
	}
	var comments string
	if name != "" {
		comments += "# Translation of " + o.General.ProjectName +
			" into " + name + " (" + target + ")\n# \n"
	}
	comments += strings.Join(tmp, "\n") + "\n"

	header := (`` + comments +
		`msgid ""` + "\n" +
		`msgstr ""` + "\n" +
		`"Project-Id-Version: ` + o.General.ProjectName + `\n"` + "\n" +
		`"Report-Msgid-Bugs-To: ` + o.Po.ReportBugs + `\n"` + "\n" +
		`"POT-Creation-Date: ` + cdate + `\n"` + "\n" +
		`"PO-Revision-Date: ` + rdate + `\n"` + "\n" +
		`"Last-Translator: ` + ltrans + `\n"` + "\n" +
		`"Language-Team: ` + lteam + `\n"` + "\n" +
		`"Language: ` + target + `\n"` + "\n" +
		`"MIME-Version: ` + mimever + `\n"` + "\n" +
		`"Content-Type: ` + cttype + `\n"` + "\n" +
		`"Content-Transfer-Encoding: ` + ctenc + `\n"` + "\n" +
		`"Plural-Forms: ` + pf + `\n"` +
		``)

	var response []string
	response = append(response, header)

	for _, v := range msgs {
		response = append(response, drawBlock(v, target))
	}

	return strings.Join(response, "\n\n")
}
