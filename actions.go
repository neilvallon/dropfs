package dropfs

type FileGetMetadata struct {
	Path                            string `json:"path"`
	IncludeMediaInfo                bool   `json:"include_media_info,omitempty"`
	IncludeDeleted                  bool   `json:"include_deleted,omitempty"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members,omitempty"`
}

type FileGetMetadataResponse struct {
	FileResponse
}

func (a FileGetMetadata) Do(fs *DropFS) (FileGetMetadataResponse, error) {
	var resp FileGetMetadataResponse
	return resp, fs.makeJSONRequest("files/get_metadata", "POST", a, &resp)
}

type FilesListFolder struct {
	Path                            string `json:"path"`
	Recursive                       bool   `json:"recursive,omitempty"`
	IncludeMediaInfo                bool   `json:"include_media_info,omitempty"`
	IncludeDeleted                  bool   `json:"include_deleted,omitempty"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members,omitempty"`
}

type FilesListFolderResponse struct {
	Entries []FileResponse `json:"entries"`
	Cursor  string         `json:"cursor"`
	HasMore bool           `json:"has_more"`
}

func (a FilesListFolder) Do(fs *DropFS) (FilesListFolderResponse, error) {
	var resp FilesListFolderResponse
	return resp, fs.makeJSONRequest("files/list_folder", "POST", a, &resp)
}

type FileResponse struct {
	Tag         string `json:".tag"`
	Name        string `json:"name"`
	PathLower   string `json:"path_lower"`
	PathDisplay string `json:"path_display"`
	ID          string `json:"id"`

	ClientModified string `json:"client_modified"`
	ServerModified string `json:"server_modified"`
	Rev            string `json:"rev"`
	Size           int64  `json:"size"`
}
