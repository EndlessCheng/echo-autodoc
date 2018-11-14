package example

import (
	"time"
	"encoding/json"
	"github.com/speps/go-hashids"
)

type HashID uint

func (id HashID) String() string {
	hData := hashids.NewData()
	hData.MinLength = 8
	hData.Salt = "qwert$yuiop"
	hID, _ := hashids.NewWithData(hData)
	result, _ := hID.Encode([]int{int(id)})
	return result
}

func (id HashID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

//

var local, _ = time.LoadLocation("Asia/Shanghai")

type JSONTime time.Time

func (t JSONTime) String() string {
	return time.Time(t).In(local).Format("2006-01-02 15:04:05")
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
