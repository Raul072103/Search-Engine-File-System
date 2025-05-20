package cache

import (
	"MyFileExporer/backend/internal/repo/database"
	"MyFileExporer/common/models"
	"slices"
	"sync"
	"time"
)

// DefaultCacheSize is the default size for the slice containing files
const (
	DefaultCacheSize  = 2048
	DefaultTTLMapSize = 2
	JanitorWakeupTime = time.Minute * 10
	IntervalSeparator = time.Minute * 5
)

type Cache struct {
	cache *cache
}

type cache struct {
	// requestMap is the in-memory cache, keeps a
	requestsMap map[string]*cacheEntry
	// janitor handles the clearing of cache
	janitor *janitor
	// janitorMap is a map that gets updated each time a new request is cached
	janitorMap map[time.Time]map[string]struct{}
	mutex      sync.RWMutex
}

type cacheEntry struct {
	lastHit time.Time
	files   []models.File
}

type janitor struct {
	lastPing time.Time
	nextPing time.Time
}

// TODO() instantiate janitor and run it
func newCache() *cache {
	requestMap := make(map[string]*cacheEntry, DefaultCacheSize)
	janitorMap := make(map[time.Time]map[string]struct{}, DefaultTTLMapSize)

	return &cache{
		requestsMap: requestMap,
		janitor:     nil,
		janitorMap:  janitorMap,
	}
}

func (j *janitor) clean() {}

func (cache *cache) Find(fs *database.FileSearchRequest) []models.File {
	key := buildKeyFromRequest(fs)

	cache.mutex.RLock()
	defer cache.mutex.RUnlock()

	entry, cacheHit := cache.requestsMap[key]
	if cacheHit {
		// decide which interval to put the current entry
		if entry.lastHit.UnixMicro()-cache.janitor.lastPing.UnixMicro() < IntervalSeparator.Microseconds() {
			// first interval
			// do nothing
		} else {
			// second interval

			// this prevents from a deadlock
			cache.mutex.RUnlock()

			cache.mutex.Lock()

			delete(cache.janitorMap[cache.janitor.lastPing], key)
			// TODO() janitor always needs instantiate the next map
			cache.janitorMap[cache.janitor.nextPing][key] = struct{}{}

			cache.mutex.Unlock()

			cache.mutex.RLock() // this reaquires the lock
		}

		// update entry lastHit
		entry.lastHit = time.Now()

		return entry.files
	}

	return nil
}

func (cache *cache) Add(fs *database.FileSearchRequest, files []models.File) {
	key := buildKeyFromRequest(fs)

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry := &cacheEntry{
		lastHit: time.Now(),
		files:   files,
	}

	cache.requestsMap[key] = entry

	if entry.lastHit.UnixMicro()-cache.janitor.lastPing.UnixMicro() < IntervalSeparator.Microseconds() {
		// first interval
		cache.janitorMap[cache.janitor.lastPing][key] = struct{}{}
	} else {
		// second interval
		cache.janitorMap[cache.janitor.nextPing][key] = struct{}{}
	}
}

func buildKeyFromRequest(fs *database.FileSearchRequest) string {
	result := ""

	if fs.Words != nil {
		slices.Sort(*fs.Words)

		for _, word := range *fs.Words {
			result += word
		}
	}

	if fs.Name != nil {
		result += *fs.Name
	}

	if fs.Extension != nil {
		slices.Sort(*fs.Extension)

		for _, extension := range *fs.Extension {
			result += extension
		}
	}

	return result
}
