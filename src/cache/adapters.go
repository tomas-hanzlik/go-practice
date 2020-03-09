package cache

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"

	types "tohan.net/go-practice/src/cache/types"
)

type IAdapter interface {
	GetData() []*types.CacheItem
	Run(wg *sync.WaitGroup)
}

type INoisyAdapter interface {
	IAdapter
	Stats() string
}

type CommandLineInputAdapter struct {
	queue  types.ItemsQueue
	reader *bufio.Reader
}

func NewCommandLineInputAdapter(rd io.Reader, bufferSize int64) IAdapter {
	var adapter IAdapter = &CommandLineInputAdapter{
		queue:  types.ItemsQueue{Capacity: bufferSize},
		reader: bufio.NewReader(rd),
	}

	return adapter
}

func (adapter *CommandLineInputAdapter) Run(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	adapter.ReadFromStdin()
}

func (adapter *CommandLineInputAdapter) ReadFromStdin() {
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) != 0 { // do not show if streamed via PIPE
		fmt.Println("Enter items in format `KEY:VALUE` separated by `\n`. Stop reading with cmd `STOP`:")
	}

	savedItemsCnt := int64(0)
	for {
		text, err := adapter.reader.ReadString('\n')
		// Check if we should stop reading
		if err == io.EOF {
			break
		} else if err != nil {
			panic("Unexpected data on input: " + err.Error())
		}

		// Normalize string
		text = strings.TrimSuffix(text, "\n")
		text = strings.TrimSpace(text)
		if text == "STOP" {
			break
		}

		// Parse text and check for correct data format
		data := strings.Split(text, ":")
		if err == nil && len(data) != 2 {
			fmt.Println("Skipping. Key:Value pair in wrong format:", text)
			continue
		} else {
			savedItemsCnt++
		}

		// save item
		adapter.queue.Enq(types.CacheItem{
			Key:   data[0],
			Value: data[1],
		})
	}
	fmt.Println("Number of collected items:", savedItemsCnt)
}

func (adapter *CommandLineInputAdapter) GetData() []*types.CacheItem {
	buffer := []*types.CacheItem{}

	adapter.queue.Lock()
	defer adapter.queue.Unlock()

	for !adapter.queue.IsEmpty() {
		item := adapter.queue.Deq()
		buffer = append(buffer, &item)
	}

	return buffer
}

type RandomInputAdapter struct {
	queue          types.ItemsQueue
	overallCounter int64
	lastlyReturned int64
	frequency      int32
	amount         int32
}

func NewRandomInputAdapter(frequency int32, amount int32, bufferSize int64) IAdapter {
	var adapter IAdapter = &RandomInputAdapter{
		queue:          types.ItemsQueue{Capacity: bufferSize},
		overallCounter: 0,
		lastlyReturned: 0,
		frequency:      frequency,
		amount:         amount,
	}
	return adapter
}

func (adapter *RandomInputAdapter) Run(wg *sync.WaitGroup) {
	executePeriodic(wg, adapter.frequency, adapter.GenerateData)
}

func (adapter *RandomInputAdapter) Stats() string {
	return "Taken items from current batch: " + strconv.Itoa(int(adapter.lastlyReturned)) + " / overall: " + strconv.Itoa(int(adapter.overallCounter))
}

func (adapter *RandomInputAdapter) GetData() []*types.CacheItem {
	buffer := []*types.CacheItem{}

	adapter.lastlyReturned = 0

	adapter.queue.Lock()
	defer adapter.queue.Unlock()

	for !adapter.queue.IsEmpty() {
		adapter.lastlyReturned++
		adapter.overallCounter++
		item := adapter.queue.Deq()
		buffer = append(buffer, &item)
	}

	return buffer
}

func (adapter *RandomInputAdapter) GenerateData() {
	adapter.queue.Lock()
	defer adapter.queue.Unlock()

	for i := int32(0); i < adapter.amount; i++ {
		adapter.queue.Enq(types.CacheItem{
			Key:   strconv.Itoa(rand.Int()),
			Value: strconv.Itoa(rand.Int()),
		})
	}
}
