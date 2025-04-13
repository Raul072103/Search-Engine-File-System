package comparator

import (
	"MyFileExporer/common/models"
	"MyFileExporer/common/utils"
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/file"
	"go.uber.org/zap"
)

type Directory interface {
	Run(directory models.File, directoryFiles []models.File) error
}

type directory struct {
	fileRepo    file.Repo
	eventsQueue *queue.InMemoryQueue
	logger      *zap.Logger
}

func New(fileRepo file.Repo, eventsQueue *queue.InMemoryQueue, logger *zap.Logger) Directory {
	return &directory{
		fileRepo:    fileRepo,
		eventsQueue: eventsQueue,
		logger:      logger,
	}
}

func (d *directory) Run(directory models.File, directoryFiles []models.File) error {
	directoryStillExists, err := utils.FileExists(directory.Path)
	if err != nil {
		return err
	}

	// check if directory exists, if not delete recursively everything under that directory
	if !directoryStillExists {
		recursiveDeleteEvent := queue.DBEvent{
			Type: queue.RecursiveDelete,
			File: directory,
		}

		d.eventsQueue.Push(recursiveDeleteEvent)
	} else {
		directoryHash, err := utils.CalculateDirectoryMD5(directory.Path)
		if err != nil {
			return err
		}

		// nothing changed, false alarm
		if directory.Hash == directoryHash {
			d.logger.Info("false alarm", zap.String("parent_dir", directory.Path))
			return nil
		}

		dbFilesMap := make(map[int64]models.File)

		for _, dbFile := range directoryFiles {
			dbFilesMap[dbFile.WindowsFileID] = dbFile
		}

		// something change, find all changes
		directoryReadFiles, err := d.fileRepo.ReadDirectoryFiles(directory.Path)
		if err != nil {
			return err
		}

		var newReadFiles = make([]models.File, 0)

		// iterate over currently read directories from the file system
		for _, fileSystemReadFile := range directoryReadFiles {
			dbFile, exists := dbFilesMap[fileSystemReadFile.WindowsFileID]

			if exists {
				// delete the entry from the map
				delete(dbFilesMap, dbFile.WindowsFileID)

				if dbFile.Hash == fileSystemReadFile.Hash {
					// same hash, nothing changed
					continue
				}

				deleteEvent := queue.DBEvent{
					Type: queue.DeleteEvent,
					File: dbFile,
				}
				insertEvent := queue.DBEvent{
					Type: queue.InsertEvent,
					File: fileSystemReadFile,
				}

				// ORDER MATTERS!
				d.eventsQueue.Push(insertEvent)
				d.eventsQueue.Push(deleteEvent)
			} else {
				newReadFiles = append(newReadFiles, fileSystemReadFile)
			}

		}

		// don't forget for newReadFiles
		for _, newFile := range newReadFiles {
			insertEvent := queue.DBEvent{
				Type: queue.InsertEvent,
				File: newFile,
			}

			d.eventsQueue.Push(insertEvent)
		}

		// don't forget for undeleted map entries
		for _, deletedFile := range dbFilesMap {
			deleteEvent := queue.DBEvent{
				Type: queue.DeleteEvent,
				File: deletedFile,
			}

			d.eventsQueue.Push(deleteEvent)
		}
	}

	return nil
}
