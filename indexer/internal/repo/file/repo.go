package file

import (
	"MyFileExporer/common/models"
	"MyFileExporer/common/utils"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Repo interface {
	Read(path string) (*models.File, error)
	ReadDirectoryFiles(directoryPath string) ([]models.File, error)
	Stats(path string) (*models.File, error)
}

type fileRepo struct {
	typeMap models.FileTypesConfig
}

func NewRepo(typeMap models.FileTypesConfig) Repo {
	return &fileRepo{typeMap: typeMap}
}

var (
	ErrFileNotFound   = errors.New("the given path doesn't exist")
	ErrFileIDOverflow = errors.New("file ID to big to store in database")
)

// ReadFile reads the contents of the file, if it exists, at the given path and returns the content of that file.
// TODO() add support .docx and .pdf files as well
func (fr *fileRepo) Read(path string) (*models.File, error) {
	file, err := fr.Stats(path)
	if err != nil {
		return file, err
	}

	if file.Extension != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		file.Content.Text = string(data)
		file.Content.UpdatedAt = file.UpdatedAt
	}

	// Calculate Hash
	if file.Extension != "" {
		file.Hash, err = utils.CalculateFileMD5(path)
		if err != nil {
			return nil, err
		}
	} else {
		file.Hash, err = utils.CalculateDirectoryMD5(path)
		if err != nil {
			return nil, err
		}
	}

	fileID, err := utils.GetFileID(path)
	if err != nil {
		return nil, err
	}

	if *fileID > math.MaxInt64 {
		return nil, ErrFileIDOverflow
	}

	file.WindowsFileID = int64(*fileID)

	return file, nil
}

// Stats if the file exists it returns an instance of models.File, else a nil.
func (fr *fileRepo) Stats(path string) (*models.File, error) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}

	if err != nil {
		return nil, err
	}

	file := &models.File{
		Path:      path,
		Name:      fileInfo.Name(),
		Size:      fileInfo.Size(),
		Mode:      uint32(fileInfo.Mode()),
		Extension: getFileExtension(path, fileInfo.IsDir()),
		UpdatedAt: fileInfo.ModTime(),
	}

	var typeId int32
	if fileInfo.IsDir() {
		typeId = 0
	} else {
		typeId = fr.typeMap.ExtensionMappings[file.Extension]
	}

	file.Type.TypeID = typeId
	file.Type.UpdatedAt = file.UpdatedAt
	file.Rank = calculateStaticFileScore(file)

	return file, err
}

// ReadDirectoryFiles reads all files in a given directory and returns them as a slice of models.File.
func (fr *fileRepo) ReadDirectoryFiles(directoryPath string) ([]models.File, error) {
	var files []models.File

	filesInfo, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range filesInfo {
		filePath := filepath.Join(directoryPath, fileInfo.Name())

		file, err := fr.Read(filePath)
		if err != nil {
			continue
		}

		files = append(files, *file)
	}

	return files, nil
}

func getFileExtension(path string, isDir bool) string {
	if isDir {
		return ""
	}
	return filepath.Ext(path)
}

// calculateStaticFileScore calculates a score at insertion time based on factors considered to affect the ranking
func calculateStaticFileScore(file *models.File) float64 {
	pathDepth := float64(pathDepth(file.Path))
	size := float64(file.Size) / 1_000_000
	updatedAt := file.UpdatedAt
	currentTime := time.Now()
	timeDifference := currentTime.Sub(updatedAt).Seconds()

	sizePenalty := 1 / (1 + size)
	timePenalty := 1 / (1 + timeDifference)
	depthWeight := 1 / (1 + pathDepth)

	return sizePenalty*0.2 + timePenalty*0.6 + depthWeight*0.2
}

func pathDepth(path string) int {
	return strings.Count(path, "\\")
}
