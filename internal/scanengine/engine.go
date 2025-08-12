package scanengine

import (
	findingnore "jetfind/internal/findignore"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Config struct {
	Root       string
	NumWorkers int
	FindIgnore *findingnore.FindIgnore
}

type Scanner struct {
	config       Config
	taskWg       sync.WaitGroup
	workerWg     sync.WaitGroup
	visited      sync.Map
	workQueue    chan string
	resultsQueue chan string
}

func New(config Config) *Scanner {
	s := &Scanner{
		config:       config,
		workQueue:    make(chan string, 512),
		resultsQueue: make(chan string, 1024),
	}
	return s
}

func (s *Scanner) Run() <-chan string {
	if s.config.NumWorkers <= 0 {
		s.config.NumWorkers = runtime.NumCPU()
	}

	s.workerWg.Add(s.config.NumWorkers)
	for i := 0; i < s.config.NumWorkers; i++ {
		go func() {
			s.worker()
		}()
	}

	s.taskWg.Add(1)
	s.workQueue <- s.config.Root

	go func() {
		s.workerWg.Wait()
		close(s.workQueue)
	}()

	go func() {
		s.taskWg.Wait()
		close(s.resultsQueue)
	}()
	return s.resultsQueue

}

func (s *Scanner) worker() {
	defer s.workerWg.Done()
	for path := range s.workQueue {
		s.scan(path)
	}
}

func (s *Scanner) scan(path string) {
	defer s.taskWg.Done()
	canonicalPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return
	}

	if _, loaded := s.visited.LoadOrStore(canonicalPath, true); loaded {
		return
	}

	entries, err := os.ReadDir(path)

	if err != nil {
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if s.config.FindIgnore != nil && s.config.FindIgnore.ShouldIgnore(fullPath) {
			continue
		}

		if entry.IsDir() {
			s.taskWg.Add(1)
			go func(p string) {
				s.workQueue <- p
			}(fullPath)
		} else {
			if fileInfo, err := os.Stat(fullPath); err == nil && !fileInfo.IsDir() {
				s.resultsQueue <- fullPath
			}
		}
	}
}
