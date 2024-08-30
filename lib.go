package main

func ForEach[arrType, accType, retType any](arr []arrType, cb func(v arrType, acc accType) retType) {

}

type Entry interface {
	Name() string
	Content() string
	Format() string
	IsFolder() (bool, []Entry)
}

type RawEntry struct {
	content string
}

func NewRawEntry(content string) RawEntry {
	return RawEntry{content: content}
}

func (re RawEntry) Name() string {
	return ""
}

func (re RawEntry) Content() string {
	return re.content
}

func (re RawEntry) Format() string {
	return "`" + re.content + "`"
}

func (re RawEntry) IsFolder() (bool, []Entry) {
	return false, []Entry{}
}

type NamedEntry struct {
	name string
	RawEntry
}

func NewNamedEntry(name, content string) NamedEntry {
	return NamedEntry{name: name, RawEntry: RawEntry{content: content}}
}

func (ne NamedEntry) Name() string {
	return ne.name
}

func (ne NamedEntry) Format() string {
	return ne.name
}

type FolderEntry struct {
	name    string
	entries []Entry
}

func NewFolderEntry(name string, entries []Entry) FolderEntry {
	return FolderEntry{name: name, entries: entries}
}

func (fe FolderEntry) Name() string {
	return fe.name
}

func (fe FolderEntry) Content() string {
	return ""
}

func (fe FolderEntry) Format() string {
	return "(folder)" + fe.name
}

func (fe FolderEntry) IsFolder() (bool, []Entry) {
	return true, fe.entries
}
