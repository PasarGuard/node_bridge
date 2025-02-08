package gozargah_node_bridge

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/m03ed/gozargah_node_bridge/common"
	"github.com/m03ed/gozargah_node_bridge/tools"
	"testing"
	"time"
)

var (
	port              = 62050
	nodeAddr          = "172.27.158.135"
	serverCA          = "certs/ssl_cert.pem"
	clientSslCertFile = "certs/ssl_client_cert.pem"
	clientSslKeyFile  = "certs/ssl_client_key.pem"
	configPath        = "config/xray.json"
)

func TestGrpcNode(t *testing.T) {
	serverCAFile, err := tools.ReadFileAsString(serverCA)
	if err != nil {
		t.Fatal(err)
	}

	clientCertFile, err := tools.ReadFileAsString(clientSslCertFile)
	if err != nil {
		t.Fatal(err)
	}

	clientKeyFile, err := tools.ReadFileAsString(clientSslKeyFile)
	if err != nil {
		t.Fatal(err)
	}

	node, err := NewNode(nodeAddr, port, clientCertFile, clientKeyFile, serverCAFile, GRPC)
	if err != nil {
		t.Fatal(err)
	}

	configFile, err := tools.ReadFileAsString(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if err = node.Start(configFile, common.BackendType_XRAY, nil); err != nil {
		t.Fatal(err)
	}

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

	defer node.Stop()
}
