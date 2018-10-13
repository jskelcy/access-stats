package collector

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/jskelcy/access-stats/pkg/types"
)

// FileWatcherConfig is configuration
type FileWatcherConfig struct {
	FileName string
}

// FileWatcher implements the collector by watching a configured file.
type FileWatcher struct {
	file     *os.File
	doneChan chan chan struct{}
	watcher  *fsnotify.Watcher
	watching bool
}

// NewFileWatcher returns a new FileWatcher from config.
func NewFileWatcher(cfg FileWatcherConfig) (*FileWatcher, error) {
	if _, err := os.Stat(cfg.FileName); err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	watcher.Add(cfg.FileName)
	file, err := os.Open(cfg.FileName)
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		doneChan: make(chan chan struct{}),
		watcher:  watcher,
		file:     file,
	}, nil
}

// Watch watches a configured file and reaturns a channel which
// will produce new lines written to the file.
func (f *FileWatcher) Watch() (<-chan types.Event, error) {
	if f.watching {
		return nil, fmt.Errorf("Collector already watching")
	}
	eventChan := make(chan types.Event)
	go f.asyncWatch(eventChan)
	f.watching = true

	return eventChan, nil
}

// Stop stops a file watching which is currently watching.
func (f *FileWatcher) Stop() (<-chan struct{}, error) {
	if !f.watching {
		return nil, fmt.Errorf("Collector is not watching")
	}

	cleanUpChan := make(chan struct{})
	f.doneChan <- cleanUpChan
	return cleanUpChan, nil
}

func (f *FileWatcher) asyncWatch(eventChan chan<- types.Event) {
	var err error
	fileReader := bufio.NewReader(f.file)

	// advance reader to the end of the file so we only
	// capture new writes
	for err != io.EOF {
		_, err = fileReader.ReadBytes('\n')
	}

	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				data, err := fileReader.ReadBytes('\n')
				if err != nil {
					eventChan <- types.Event{
						Err: err,
					}
					f.cleanup(eventChan)
					return
				}
				eventChan <- types.Event{
					Data: data,
				}
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				f.cleanup(eventChan)
				return
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				f.cleanup(eventChan)
				return
			}
			eventChan <- types.Event{
				Err: err,
			}
			f.cleanup(eventChan)
			return
		case cleanUpChan := <-f.doneChan:
			f.asyncCleanup(eventChan, cleanUpChan)
		}
	}
}

func (f *FileWatcher) cleanup(eventChan chan<- types.Event) {
	f.watching = false
	if err := f.watcher.Close(); err != nil {
		log.Printf("error cleaning up: %v", err)
	}
	if err := f.file.Close(); err != nil {
		log.Printf("error cleaning up: %v", err)
	}
	close(eventChan)
}

// asyncCleanup takes a channel which gets an empty struct when
// clean up is completed.
func (f *FileWatcher) asyncCleanup(
	eventChan chan<- types.Event,
	cleanUpChan chan struct{},
) {
	f.cleanup(eventChan)
	cleanUpChan <- struct{}{}
}
