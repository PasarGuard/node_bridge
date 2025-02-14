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
	port              = 62050
	nodeAddr          = "172.27.158.135"
	serverCA          = "certs/ssl_cert.pem"
	clientSslCertFile = "certs/ssl_client_cert.pem"
	clientSslKeyFile  = "certs/ssl_client_key.pem"
	configPath        = "config/xray.json"
)

var (
	serverCAFile   []byte
	clientCertFile []byte
	clientKeyFile  []byte
	configFile     string
)

func init() {
	var err error
	serverCAFile, err = os.ReadFile(serverCA)
	if err != nil {
		log.Fatal(err)
	}

	clientCertFile, err = os.ReadFile(clientSslCertFile)
	if err != nil {
		log.Fatal(err)
	}

	clientKeyFile, err = os.ReadFile(clientSslKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	configFile, err = tools.ReadFileAsString(configPath)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGrpcNode(t *testing.T) {
	node, err := NewNode(nodeAddr, port, clientCertFile, clientKeyFile, serverCAFile, nil, GRPC)
	if err != nil {
		t.Fatal(err)
	}

	if err = node.Start(configFile, common.BackendType_XRAY, nil); err != nil {
		t.Fatal(err)
	}

	defer node.Stop()

	info, err := node.Info()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", info)

	user := common.CreateUser(
		"test_user",
		common.CreateProxies(
			common.CreateVmess(uuid.New().String()), common.CreateVless(uuid.New().String(), ""),
			common.CreateTrojan("random data"), common.CreateShadowsocks("random", "aes-256-gcm"),
		),
		[]string{},
	)

	if err = node.UpdateUser(user); err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			logChan, err := node.GetLogs()
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
	node, err := NewNode(nodeAddr, port, clientCertFile, clientKeyFile, serverCAFile, nil, REST)
	if err != nil {
		t.Fatal(err)
	}

	if err = node.Start(configFile, common.BackendType_XRAY, nil); err != nil {
		t.Fatal(err)
	}

	defer node.Stop()

	info, err := node.Info()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", info)

	user := common.CreateUser(
		"test_user",
		common.CreateProxies(
			common.CreateVmess(uuid.New().String()), common.CreateVless(uuid.New().String(), ""),
			common.CreateTrojan("random data"), common.CreateShadowsocks("random", "aes-256-gcm"),
		),
		[]string{},
	)

	if err = node.UpdateUser(user); err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			logChan, err := node.GetLogs()
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

	stats, err := node.GetOutboundsStats(true)
	if err != nil {
		t.Fatal(err)
	}

	for _, stat := range stats.GetStats() {
		fmt.Printf("%+v\n", stat)
	}
}
