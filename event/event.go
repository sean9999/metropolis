package event

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

type Payload map[string]any

func (p Payload) String() string {
	j, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func (p Payload) Serialize() (string, error) {
	mjson, err := json.Marshal(p)
	if err != nil {
		return "err", err
	}
	b64output := base64.StdEncoding.EncodeToString(mjson)
	return b64output, nil
}

type Event struct {
	Data Payload `json:"data"`
	Id   uint64  `json:"id"`
	Is   bool    `json:"is"`
}

// Serialize formats the Event for SSE
func (e Event) Serialize() string {
	b, _ := e.Data.Serialize()
	return fmt.Sprintf("data: %s\nid: %d\n\n", b, e.Id)
}

var errOutofBounds = errors.New("index out of bounds for slice")
var NoEvent = Event{}
