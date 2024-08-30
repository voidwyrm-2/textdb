package main

import (
	"fmt"
	"strings"
)

/*
func compileEntries(entries []Entry) string {
	return ""
}
*/

// SDS stands for Simple Database Script

func interpretSDS(text string, isFolder bool, offset int) ([]Entry, int, error) {
	var lines = ForEach(strings.Split(text, "\n"), func(v string, acc []string) []string {
		acc = append(acc, strings.TrimSpace(v))
		return acc
	})

	var out []Entry

	errF := func(ln int, msg string) ([]Entry, int, error) {
		return []Entry{}, 0, fmt.Errorf("error on line %d: %s", ln+1+offset, msg)
	}

	nameAcc := ""
	acc := []string{}
	collecting := false

	//fmt.Println(lines)

	/*
		getIndex := func(a []string, index int) string {
			if index < len(a) && index > -1 {
				return a[index]
			}
			return ""
		}
	*/

	ln := 0
	for ln < len(lines) {
		//fmt.Println(ln, lines[ln], isFolder, collecting)
		l := lines[ln]
		if l == "#end" {
			if isFolder && !collecting {
				return out, ln + 1, nil
			} else if !collecting {
				return errF(ln, "'#end' cannot be used outside of a '#folder' or '#namedentry' block")
			}
			out = append(out, NewNamedEntry(nameAcc, strings.TrimSpace(strings.Join(acc, "\n"))))
			acc = []string{}
			collecting = false
			ln++
		} else if collecting {
			acc = append(acc, l)
			ln++
			continue
		} else if l == "" {
			ln++
			continue
		} else if strings.HasPrefix(l, "#folder ") {
			folderName := strings.TrimSpace(l[7:])
			if folderName == "" {
				return errF(ln, "folder names cannot be empty")
			}
			content, offset, err := interpretSDS(strings.Join(lines[ln+1:], "\n"), true, ln+1)
			if err != nil {
				return []Entry{}, 0, err
			}
			out = append(out, NewFolderEntry(folderName, content))
			ln += offset + 1
			fmt.Println(lines[ln])
		} else if strings.HasPrefix(l, "#rawentry ") {
			content := strings.TrimSpace(l[9:])
			if content == "" {
				return errF(ln, "raw entries cannot be empty")
			}
			out = append(out, NewRawEntry(content))
			ln++
		} else if strings.HasPrefix(l, "#namedentry ") {
			entryName := strings.TrimSpace(l[11:])
			if entryName == "" {
				return errF(ln, "named entry names cannot be empty")
			}
			nameAcc = entryName
			collecting = true
			ln++
		} else {
			return errF(ln, fmt.Sprintf("unknown instruction '%s'", strings.TrimSpace(strings.Split(l, " ")[0])))
		}
	}

	if isFolder {
		return errF(len(lines)-1, "folder entries must be ended with '#end'")
	} else if collecting {
		return errF(len(lines)-1, "named entries must be ended with '#end'")
	}

	return out, 0, nil
}
