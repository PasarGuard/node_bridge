package rest

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"

	"github.com/pasarguard/node_bridge/common"
	"github.com/pasarguard/node_bridge/controller"
)

func (n *Node) SyncUsers(users []*common.User) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()

		if len(users) == 0 {
			chunk := &common.UsersChunk{
				Users: nil,
				Index: 0,
				Last:  true,
			}
			if err := sendChunk(pw, chunk); err != nil {
				pw.CloseWithError(err)
			}
			return
		}

		for i := 0; i < len(users); i += controller.MaxChunkSize {
			end := i + controller.MaxChunkSize
			if end > len(users) {
				end = len(users)
			}
			chunk := &common.UsersChunk{
				Users: users[i:end],
				Index: uint64(i / controller.MaxChunkSize),
				Last:  end == len(users),
			}
			if err := sendChunk(pw, chunk); err != nil {
				pw.CloseWithError(err)
				return
			}
		}
	}()

	req, err := http.NewRequest("PUT", n.baseUrl+"/users/sync/chunked", pr)
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", n.ApiKey())
	req.Header.Set("Content-Type", "application/x-protobuf")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func sendChunk(w io.Writer, chunk *common.UsersChunk) error {
	data, err := proto.Marshal(chunk)
	if err != nil {
		return err
	}

	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(data)))
	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}
