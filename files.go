package dropfs

import (
	"io"
	"os"
	"time"
)

// File
type dFile struct {
	fileInfo
	size int64

	rc io.ReadCloser
}

func (f dFile) Size() int64 { return f.size }
func (dFile) IsDir() bool   { return false }

func (f *dFile) Read(p []byte) (n int, err error)             { return f.rc.Read(p) }
func (f *dFile) Seek(offset int64, whence int) (int64, error) { return 0, nil }
func (f *dFile) Close() error                                 { return f.rc.Close() }

func (f *dFile) Stat() (os.FileInfo, error) { return f, nil }

// Folder
type dFolder struct {
	fileInfo
}

func (dFolder) Size() int64 { return 0 }
func (dFolder) IsDir() bool { return true }

func (f dFolder) Readdir(count int) ([]os.FileInfo, error) {
	resp, err := FilesListFolder{Path: f.PathLower}.Do(f.fs)
	if err != nil {
		return nil, err
	}

	ilist := make([]os.FileInfo, len(resp.Entries))
	for i, v := range resp.Entries {
		fi := fileInfo{
			name:           v.Name,
			ClientModified: v.ClientModified,
			ServerModified: v.ServerModified,
		}

		switch v.Tag {
		case "file":
			ilist[i] = &dFile{
				fileInfo: fi,
				size:     v.Size,
			}
		case "folder":
			ilist[i] = &dFolder{fileInfo: fi}
		default:
			panic("dropfs: unexpected tag type")
		}
	}

	return ilist, nil
}

func (f *dFolder) Stat() (os.FileInfo, error) { return f, nil }

// Dumby
type fileInfo struct {
	name string

	PathLower   string
	PathDisplay string

	ClientModified string
	ServerModified string

	fs *DropFS
}

// http.File
func (fileInfo) Close() error                                 { return nil }
func (fileInfo) Read(p []byte) (n int, err error)             { panic("OVERRIDE ME") }
func (fileInfo) Seek(offset int64, whence int) (int64, error) { panic("OVERRIDE ME") }
func (fileInfo) Readdir(count int) ([]os.FileInfo, error)     { panic("OVERRIDE ME") }

// os.FileInfo
func (fi fileInfo) Name() string   { return fi.name }
func (fileInfo) Mode() os.FileMode { return 0777 }

func (fi fileInfo) ModTime() time.Time {
	t, _ := time.Parse(time.RFC3339, fi.ClientModified)
	if sm, err := time.Parse(time.RFC3339, fi.ServerModified); err == nil && t.Before(sm) {
		t = sm
	}

	return t
}

func (fileInfo) Sys() interface{} { return nil }
