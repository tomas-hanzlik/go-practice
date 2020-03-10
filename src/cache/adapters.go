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
}

type INoisyAdapter interface {
	IAdapter
	Stats() string
}

type CommandLineInputAdapter struct {
	queue  types.ItemsQueue
	reader *bufio.Reader
	sync.Mutex
}

func NewCommandLineInputAdapter(rd io.Reader, bufferSize int64) IAdapter {
	var adapter IAdapter = &CommandLineInputAdapter{
		queue:  types.ItemsQueue{Capacity: bufferSize},
		reader: bufio.NewReader(rd),
	}
	go adapter.(*CommandLineInputAdapter).readFromStdin()

	return adapter
}

func (adapter *CommandLineInputAdapter) readFromStdin() {
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
			fmt.Println("Skipping. Unexpected data on input: " + err.Error())
			continue
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

	adapter.Lock()
	defer adapter.Unlock()
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
	sync.Mutex
}

func NewRandomInputAdapter(frequency int32, amount int32, bufferSize int64) IAdapter {
	var adapter IAdapter = &RandomInputAdapter{
		queue:          types.ItemsQueue{Capacity: bufferSize},
		overallCounter: 0,
		lastlyReturned: 0,
		frequency:      frequency,
		amount:         amount,
	}
	if frequency > 0 {
		go executePeriodic(frequency, adapter.(*RandomInputAdapter).generateData)
	}

	return adapter
}

func (adapter *RandomInputAdapter) Stats() string {
	return "[RandomInputAdapter] Collecting items: " + strconv.Itoa(int(adapter.lastlyReturned)) + " / overall: " + strconv.Itoa(int(adapter.overallCounter))
}

func (adapter *RandomInputAdapter) GetData() []*types.CacheItem {
	buffer := []*types.CacheItem{}
	adapter.lastlyReturned = 0

	adapter.Lock()
	defer adapter.Unlock()

	for !adapter.queue.IsEmpty() {
		adapter.lastlyReturned++
		adapter.overallCounter++
		item := adapter.queue.Deq()
		buffer = append(buffer, &item)
	}

	return buffer
}

func (adapter *RandomInputAdapter) generateData() {
	adapter.Lock()
	defer adapter.Unlock()

	for i := int32(0); i < adapter.amount; i++ {
		adapter.queue.Enq(types.CacheItem{
			Key:   strconv.Itoa(rand.Int()),
			Value: strconv.Itoa(rand.Int()),
		})
	}
}
