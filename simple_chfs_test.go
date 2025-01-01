package chfs

import (
	// "fmt"
	"testing"
	"github.com/stretchr/testify/assert"
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

func TestFileSystemIntegrity(t *testing.T) {
	testCases := []struct {
		name  string
		files []PathObject
		expectedPaths []string
	}{
		{
			"test write single file",
			[]PathObject{
				newPathObject(NewPath("/root"),*NewName("file1"),"Test File 1",false,),
			},
			[]string{"/root", "/root/file1"},
		},
		{
			"test write three files",
			[]PathObject{
				newPathObject(NewPath("/root"),*NewName("file1"),"Test File 1",false,),
				newPathObject(NewPath("/root"),*NewName("file2"),"Test File 2",false,),
				newPathObject(NewPath("/root"),*NewName("file3"),"Test File 3",false,),
			},
			[]string{"/root", "/root/file1", "/root/file2", "/root/file3"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := NewSimpleChFS()
			Init(fs.ref, fs.tree)

			// files := []PathObject{
			// 	newPathObject(NewPath("/root/"),*NewName("file1"),"Test File 1",false,),
			// 	// newPathObject(NewPath("/root/"),*NewName("file3"),"Test File 3",false,),
			// }

			fs.Write("HEAD", tc.files)

			search, err := fs.ListDir("HEAD", NewPath("/root/"))
			if err != nil {
				t.Error(err)
			}

			actual := make([]string, len(search))
			for i, file := range search {
				actual[i] = file.Path().String()
			}
			assert.ElementsMatch(
				t,
				actual,
				tc.expectedPaths,
				"Unexpected filesystem structure.",
			)
		})
	}

}
