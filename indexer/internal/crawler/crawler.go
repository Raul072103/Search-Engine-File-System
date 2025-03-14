package crawler

import (
	"MyFileExporer/indexer/internal/queue"
	"MyFileExporer/indexer/internal/repo/file"
	"context"
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
	IgnorePatterns []string
	RootDir        string
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
	c.logger.Info("Starting crawler", zap.String("root", c.config.RootDir))
	c.Crawl(ctx, c.config.RootDir)
	c.logger.Info("Crawler finished")
}

func (c *crawler) Crawl(ctx context.Context, path string) {
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
			return
		}

		insertEvent := queue.DBEvent{
			Type: queue.InsertEvent,
			File: *fileModel,
		}

		c.eventsQueue.Push(insertEvent)

		if fileModel.IsDir {
			entries, err := os.ReadDir(path)
			if err != nil {
				c.logger.Error(
					"Error reading directory for further traversing",
					zap.String("path", path),
					zap.Error(err))

				return
			}

			// Recur for each entry
			for _, entry := range entries {
				entryPath := filepath.Join(path, entry.Name())
				c.Crawl(ctx, entryPath)
			}
		}
	}
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
