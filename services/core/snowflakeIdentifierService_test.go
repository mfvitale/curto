package core

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/syncmap"
)

var wg sync.WaitGroup

func TestSimpleGeneration(t *testing.T) {

	assert := assert.New(t)

	snowflakeService := NewDefaultSnowflakeGenerator(1, 1)

	id, err := snowflakeService.NextID()

	assert.Nil(err)
	assert.NotNil(id)
}

func TestTimeInPast(t *testing.T) {

	assert := assert.New(t)

	timeInFuture := time.Now().Add(time.Hour).UnixMilli()
	snowflakeService := NewSnowflakeGenerator(1, 1, timeInFuture, 0)

	_, err := snowflakeService.NextID()

	assert.ErrorContains(err, "time is moving backwards, waiting until")
}

func TestSequenceOverflow(t *testing.T) {

	assert := assert.New(t)

	snowflakeService := NewSnowflakeGenerator(1, 1, time.Now().UnixMilli(), 4095)

	id, err := snowflakeService.NextID()

	assert.Nil(err)
	assert.NotNil(id)
}

func TestConcurrentGeneration(t *testing.T) {

	assert := assert.New(t)

	snowflakeService := NewDefaultSnowflakeGenerator(1, 1)

	count := 10000
	wg.Add(count)
	generatedIds := syncmap.Map{}

	// Concurrently count goroutines for snowFlake ID generation
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id, _ := snowflakeService.NextID()
			generatedIds.Store(id, 1)
		}()
	}
	wg.Wait()

	assert.Equal(lenSyncMap(&generatedIds), count)
}

func lenSyncMap(m *sync.Map) int {
	var i int
	m.Range(func(k, v interface{}) bool {
		i++
		return true
	})
	return i
}
