package dropfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type DropFS struct {
	Root   string
	apiKey string
}

func NewFS(key, root string) *DropFS {
	return &DropFS{Root: root, apiKey: key}
}

func (fs *DropFS) Open(name string) (http.File, error) {
	name = filepath.Join(fs.Root, name)

	// special case root element
	if name == "/" {
		return &dFolder{
			baseFile: baseFile{
				name:        "root",
				PathLower:   "",
				PathDisplay: "/",

				fs: fs,
			},
		}, nil
	}

	file, err := FileGetMetadata{Path: name}.Do(fs)

	var hFile http.File
	switch file.Tag {
	case "file":
		rc, err := fs.makeDownloadRequest(file.PathLower)
		if err != nil {
			return hFile, os.ErrNotExist
		}

		hFile = &dFile{
			baseFile: baseFile{
				name:        file.Name,
				PathLower:   file.PathLower,
				PathDisplay: file.PathDisplay,

				ClientModified: file.ClientModified,
				ServerModified: file.ServerModified,

				fs: fs,
			},
			size: file.Size,
			rc:   rc,
		}
	case "folder":
		hFile = &dFolder{
			baseFile: baseFile{
				name:        file.Name,
				PathLower:   file.PathLower,
				PathDisplay: file.PathDisplay,

				ClientModified: file.ClientModified,
				ServerModified: file.ServerModified,

				fs: fs,
			},
		}
	default:
		return nil, os.ErrInvalid
	}

	return hFile, err
}

func (fs *DropFS) makeJSONRequest(action, method string, data, result interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return err
	}

	log.Printf("[%s - %s] %s", method, action, &buf)

	req, err := http.NewRequest(method, "https://api.dropboxapi.com/2/"+action, &buf)
	req.Header.Set("Authorization", "Bearer "+fs.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[ERROR %s] %s\n", resp.Status, body)

		return errors.New("dropFS: unknown request failure")
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

func (fs *DropFS) makeDownloadRequest(name string) (body io.ReadCloser, err error) {
	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/download", nil)

	req.Header.Set("Dropbox-API-Arg", `{"path":"`+name+`"}`)
	req.Header.Set("Authorization", "Bearer "+fs.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("[ERROR %s] %s\n", resp.Status, body)

		return nil, errors.New("dropFS: unknown request failure")
	}

	return resp.Body, nil
}
