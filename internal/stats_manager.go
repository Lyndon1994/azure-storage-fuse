package internal

import (
	"blobfuse2/common/log"
	"sync"
)

type ChannelReader func()

type StatsCollector struct {
	componentName string
	channel       chan interface{}
	workerDone    sync.WaitGroup
	reader        ChannelReader
}

type Stats struct {
	ComponentName string
	Operation     string
	Value         map[string]int64
}

// func OpenMemoryMappedFiles(fileName string) ([]byte, error) {
// 	f, err := os.Create(fileName)
// 	if err != nil {
// 		return nil, err
// 	}
// }

func NewStatsCollector(componentName string, reader ChannelReader) (*StatsCollector, error) {
	sc := &StatsCollector{componentName: componentName}
	sc.channel = make(chan interface{}, 100000)
	sc.reader = reader

	return sc, nil
}

func (sc *StatsCollector) Init() {
	sc.workerDone.Add(1)
	go sc.statsDumper()
}

func (sc *StatsCollector) Destroy() error {
	close(sc.channel)
	sc.workerDone.Wait()
	return nil
}

func (sc *StatsCollector) AddStats(stats interface{}) {
	sc.channel <- stats
}

func (sc *StatsCollector) statsDumper() {
	defer sc.workerDone.Done()

	for st := range sc.channel {
		log.Debug("StatsManager::StatsDumper : %v stats: %v", sc.componentName, st)
	}
}
