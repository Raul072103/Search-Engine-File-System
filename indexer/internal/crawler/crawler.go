package crawler

import (
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/file"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

// Crawler is the basis for the component which crawls every file and directory starting from the root path.
// Logs any errors happening throughout the process and jumps over files specified in the configuration.
type Crawler interface {
	Run(ctx context.Context)
	Crawl(ctx context.Context, path string)
}

type crawler struct {
	fileRepo    file.Repo
	eventsQueue *queue.InMemoryQueue
	logger      *zap.Logger
	config      Config
}

type Config struct {
	IgnorePatterns   []string `json:"ignore_patterns"`
	LastTraversedDir string   `json:"last_traversed_dir"`
	RootDir          string   `json:"root_dir"`
	CrawlerDone      bool     `json:"crawler_done"`
}

func New(fileRepo file.Repo, eventsQueue *queue.InMemoryQueue, logger *zap.Logger, config Config) Crawler {
	return &crawler{
		fileRepo:    fileRepo,
		eventsQueue: eventsQueue,
		logger:      logger,
		config:      config,
	}
}

func (c *crawler) Run(ctx context.Context) {
	if !c.config.CrawlerDone {
		var startDir string
		if c.config.LastTraversedDir == "" {
			startDir = c.config.RootDir
		} else {
			startDir = c.config.LastTraversedDir
		}

		c.logger.Info("Starting crawler", zap.String("dir", startDir))
		c.Crawl(ctx, startDir)

		err := c.saveCrawlerProgress()
		if err != nil {
			c.logger.Error("failed to update crawler status to done", zap.Error(err))
			return
		}

		c.logger.Info("Crawler finished")
	}
}

func (c *crawler) Crawl(ctx context.Context, startPath string) {
	var explorePaths = make([]string, 0)
	var parentID int64 = 0

	explorePaths = append(explorePaths, startPath)

	for {
		explorePathsLength := len(explorePaths)
		if explorePathsLength == 0 {
			c.logger.Info("Crawler finished traversing process")
			break
		}

		path := explorePaths[explorePathsLength-1]
		explorePaths = explorePaths[:explorePathsLength-1]

		select {
		case <-ctx.Done():
			c.logger.Info("Crawler stopped before going further", zap.String("path", path))
			return
		default:
			if c.matchesPattern(path) {
				return
			}

			fileModel, err := c.fileRepo.Read(path)
			if err != nil {
				c.logger.Info("Error reading file or dir", zap.String("path", path), zap.Error(err))
				explorePaths = explorePaths[:explorePathsLength-1]
				continue
			}

			fileModel.ParentFileID = parentID

			insertEvent := queue.DBEvent{
				Type: queue.InsertEvent,
				File: *fileModel,
			}

			c.eventsQueue.Push(insertEvent)

			if fileModel.Extension == "" {
				// Save the fileID as the next parentID
				parentID = fileModel.WindowsFileID

				entries, err := os.ReadDir(path)
				if err != nil {
					c.logger.Error(
						"Error reading directory for further traversing",
						zap.String("path", path),
						zap.Error(err))

					explorePaths = explorePaths[:explorePathsLength-1]
					continue
				}

				// Recur for each entry
				for _, entry := range entries {
					entryPath := filepath.Join(path, entry.Name())
					explorePaths = append(explorePaths, entryPath)
				}

				// Save as last traversed directory
				c.config.LastTraversedDir = path
				err = c.saveCrawlerProgress()
				if err != nil {
					c.logger.Error("error saving crawler progress to JSON file", zap.Error(err))
				}
			}
		}
	}

	c.config.CrawlerDone = true
}

func (c *crawler) matchesPattern(filePath string) bool {
	for _, pattern := range c.config.IgnorePatterns {
		match, err := filepath.Match(pattern, filepath.Base(filePath))
		if err != nil {
			c.logger.Error("Error matching pattern", zap.String("filePath", filePath), zap.String("pattern", pattern), zap.Error(err))
			return true
		}
		if match {
			c.logger.Info("File ignored", zap.String("filePath", filePath), zap.String("pattern", pattern))
			return true
		}
	}
	return false
}

func (c *crawler) saveCrawlerProgress() error {
	crawlerConfigFile, err := os.OpenFile("./config.json", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("error opening crawlerConfigFile: %w", err)
	}
	defer crawlerConfigFile.Close()

	encoder := json.NewEncoder(crawlerConfigFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(c.config)
	if err != nil {
		return fmt.Errorf("error encoding crawler config to crawlerConfigFile: %w", err)
	}

	return nil
}
