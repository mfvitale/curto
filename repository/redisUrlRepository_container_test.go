package repository

import (
    "context"
    "fmt"
    "os"
    "testing"

    "github.com/go-redis/redis/v8"
    "github.com/stretchr/testify/assert"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

var rdb *redis.Client
var redisContainerInstance *redisContainer
var repository UrlRepository

func setup() {

    var err error
    redisContainerInstance, err = setupRedis(ctx)
    if err != nil {
        panic(err)
    }
    rdb = redis.NewClient(&redis.Options{
        Addr:    redisContainerInstance.URI,
        Password: "password",
        DB:       0,  // use default DB
    })

    repository = NewRedisUrlRepository(rdb)
}

func shutdown() {
    fmt.Println("Shutting down container")
    err := redisContainerInstance.Terminate(ctx)
	if err != nil {
		fmt.Println("Error during container shutdown")
	}
}

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    shutdown()
    os.Exit(code)
}

func TestStore(t *testing.T) {

    assert := assert.New(t)

    err := repository.Store("du45g", "http://www.google.it")
	if err != nil {
		t.Fail()
	}
    url, err := repository.Get("du45g")

    assert.NoError(err)
    assert.Equal("http://www.google.it", url)

    flushRedis(ctx, *rdb)
}

func TestGetNotExistingKey(t *testing.T) {

    assert := assert.New(t)

    _, err := repository.Get("not_exist")

    assert.Equal("Short URL not_exist not found", err.Error())

    flushRedis(ctx, *rdb)
}


type redisContainer struct {
    testcontainers.Container
    URI string
}

func setupRedis(ctx context.Context) (*redisContainer, error) {

    req := testcontainers.ContainerRequest{
        Image:        "redis:6.2.3",
        ExposedPorts: []string{"6379/tcp"},
        Cmd: []string{"--requirepass 'password'"},
        WaitingFor:   wait.ForLog("* Ready to accept connections"),
    }

    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        return nil, err
    }

    endpoint, err := container.Endpoint(ctx, "")
    if err != nil {
        return nil, err
    }

    fmt.Printf("Redis is running on %s\n", endpoint)

    return &redisContainer{Container: container, URI: endpoint}, nil
}

func flushRedis(ctx context.Context, client redis.Client) {

	client.FlushAll(ctx)
}