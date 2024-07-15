package chfs

import (
	"testing"
)

func newPathObject(directory Path, filename Name, contents string, executable bool) PathObject {
	return PathObject{
		directory,
		filename,
		*EncodeChecksum(contents),
		executable,
	}
}

type Dataset struct {
	objects []PathObject
}

func TestRandomShit(t *testing.T) {
	fs := NewSimpleChFS()
	Init(fs.ref, fs.tree)

	files := []PathObject{
		newPathObject(NewPath("/root/"),*NewName("file1"),"Test File 1",false,),
		newPathObject(NewPath("/root/"),*NewName("file2"),"Test File 2",false,),
		newPathObject(NewPath("/root/"),*NewName("file3"),"Test File 3",false,),
	}

	fs.Write("HEAD", files)
}
