package main

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"fmt"
)

//Everything in our filesystem is an 'Entry'
type Entry interface {
	Name() string
	Attr() *fuse.Attr
	GetInfo() fuse.DirEntry
}

//Embeddable struct to ease of coding
type mEntry struct {
	attr *fuse.Attr
	name string
}

//Directory
type Dir struct {
	Entries []Entry
	mEntry
}

func MakeDir(name string) *Dir {
	d := new(Dir)
	d.name = name
	d.attr = &fuse.Attr{Mode: fuse.S_IFDIR | 0755}
	return d
}

func (d *Dir) AddEntry(e Entry) {
	if e == nil {
		return
	}
	d.Entries = append(d.Entries, e)
	d.attr.Nlink = len(d.Entries)
}

func (d *Dir) RemoveChild(name string) {
	for i, e := range d.Entries {
		if name == e.Name() {
			d.Entries = append(d.Entries[:i], d.Entries[i+1:]...)
			return
		}
	}
}

func (d *Dir) GetEntry(toks []string) Entry {
	if len(toks) == 0 {
		fmt.Println("returning self")
		return d
	}

	for _,e := range(d.Entries) {
		if e.Name() == toks[0] {
			if len(toks) == 1 {
				//This is it!
				return e
			} else {
				//Need to search deeper
				sub, ok := e.(*Dir)
				if !ok {
					return nil
				}
				return sub.GetEntry(toks[1:])
			}
		}
	}
	return nil
}

func (d *Dir) GetInfo() fuse.DirEntry {
	return fuse.DirEntry{Name: d.name, Mode: fuse.S_IFDIR}
}

func (d *Dir) Name() string {
	return  d.name
}

func (d *Dir) Attr() *fuse.Attr {
	return d.attr
}

//Normal file
type File struct {
	Content string
	FileData nodefs.File
	Chunks int
	RealSize int //Full size of the file
	LocalSize int //Size of this instance on disk
	mEntry
}

func MakeFile(name string) *File {
	f := new(File)
	f.name = name
	f.attr = &fuse.Attr{ Mode: fuse.S_IFREG | 0644, Size: uint64(len(name))}
	f.FileData = nodefs.NewDefaultFile()
	return f
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Attr() *fuse.Attr {
	return f.attr
}

func (f *File) GetInfo() fuse.DirEntry {
	return fuse.DirEntry{Name: f.name, Mode: fuse.S_IFREG}
}

//Link to a file or directory (or another link... interesting...)
type Link struct {
	to Entry
}

func (l *Link) Name() string {
	return l.to.Name()
}

func (l *Link) Attr() *fuse.Attr {
	return l.to.Attr()
}
