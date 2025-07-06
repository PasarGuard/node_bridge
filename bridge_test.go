package gozargah_node_bridge

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/tools"
)

var (
	nodeAddr   = "172.27.158.135"
	serverCA   = "certs/ssl_cert.pem"
	apiKey     = "d04d8680-942d-4365-992f-9f482275691d"
	configPath = "config/xray.json"
	keepAlive  = uint64(60)
)

var (
	serverCAFile []byte
	uuidKey      uuid.UUID
	configFile   string
	opts         []NodeOption
)

func init() {
	var err error
	serverCAFile, err = os.ReadFile(serverCA)
	if err != nil {
		log.Fatal(err)
	}

	uuidKey, err = uuid.Parse(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	configFile, err = tools.ReadFileAsString(configPath)
	if err != nil {
		log.Fatal(err)
	}

	opts = []NodeOption{
		WithPort(2096),
		WithAPIKey(uuidKey),
		WithServerCA(serverCAFile),
		WithLogChannelSize(100),
	}
}

var user = common.CreateUser(
	"test_user",
	common.CreateProxies(
		common.CreateVmess(uuid.New().String()),
		common.CreateVless(uuid.New().String(), ""),
		common.CreateTrojan("random data"),
		common.CreateShadowsocks("random", "aes-256-gcm"),
	),
	[]string{},
)

func TestGrpcNode(t *testing.T) {
	node, err := New(nodeAddr, GRPC, opts...)
	if err != nil {
		t.Fatal(err)
	}

	if err = node.Start(configFile, common.BackendType_XRAY, nil, keepAlive); err != nil {
		t.Fatal(err)
	}

	defer node.Stop()

	info, err := node.Info()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", info)

	node.UpdateUser(user)

	go func() {
		for {
			logChan, err := node.Logs()
			if err != nil {
				t.Error(err)
			}
			newLog, ok := <-logChan
			if !ok {
				return
			}
			fmt.Println(newLog)
		}
	}()

	time.Sleep(2 * time.Second)
}

func TestRestNode(t *testing.T) {
	node, err := New(nodeAddr, REST, opts...)
	if err != nil {
		t.Fatal(err)
	}

	if err = node.Start(configFile, common.BackendType_XRAY, nil, keepAlive); err != nil {
		t.Fatal(err)
	}

	defer node.Stop()

	info, err := node.Info()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", info)

	node.UpdateUser(user)

	go func() {
		for {
			logChan, err := node.Logs()
			if err != nil {
				t.Error(err)
			}
			newLog, ok := <-logChan
			if !ok {
				return
			}
			fmt.Println(newLog)
		}
	}()

	time.Sleep(3 * time.Second)

	stats, err := node.GetStats(true, "", common.StatType_Outbounds)
	if err != nil {
		t.Fatal(err)
	}

	for _, stat := range stats.GetStats() {
		fmt.Printf("%+v\n", stat)
	}
}
