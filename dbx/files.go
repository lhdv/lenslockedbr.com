package dbx

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	dbxfiles "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

type Folder struct {
	Name string
	Path string
}

type File struct {
	Name string
	Path string
}

func List(accessToken string, path string) ([]Folder, []File, error) {

	var folders []Folder
	var files []File

	config := dropbox.Config {
		Token: accessToken,
	}
	
	client := dbxfiles.New(config)
	args := &dbxfiles.ListFolderArg {
		Path: path,
	}

	res, err := client.ListFolder(args)
	if err != nil {
		return nil, nil, err
	} 

	for _, entry := range res.Entries {
		switch meta := entry.(type) {
		case *dbxfiles.FolderMetadata:
			folders = append(folders, Folder {
				Name: meta.Name,
				Path: meta.PathLower,
			})
		case *dbxfiles.FileMetadata:
			files = append(files, File {
				Name: meta.Name,
				Path: meta.PathLower,
			})
		}
	} 

	return folders, files, nil
}
