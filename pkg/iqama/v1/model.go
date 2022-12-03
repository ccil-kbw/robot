// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    resp, err := UnmarshalResp(bytes)
//    bytes, err = resp.Marshal()

package v1

import "encoding/json"

func UnmarshalResp(data []byte) (Resp, error) {
	var r Resp
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Resp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Resp struct {
	Dates                  Dates         `json:"dates"`
	Fajr                   Prayer        `json:"fajr"`
	Sunrise                string        `json:"sunrise"`
	Dhuhr                  Prayer        `json:"dhuhr"`
	Asr                    Prayer        `json:"asr"`
	Maghrib                Prayer        `json:"maghrib"`
	Isha                   Prayer        `json:"isha"`
	Jumua                  Jumua         `json:"jumua"`
	SalatsToChangeTomorrow []interface{} `json:"salatsToChangeTomorrow"`
	SalatsThatChangeToday  interface{}   `json:"salatsThatChangeToday"`
}

type Prayer struct {
	Adhan string `json:"adhan"`
	Iqama string `json:"iqama"`
}

type Dates struct {
	Hijri  string `json:"hijri"`
	Miladi string `json:"miladi"`
}

type Jumua struct {
	Fr string `json:"fr"`
	Ar string `json:"ar"`
}
