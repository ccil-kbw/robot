package v1

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Get() (*Resp, error) {

	body := GetRAW()

	iqamaResp, err := UnmarshalResp(body)
	if err != nil {
		return nil, err
	}
	return &iqamaResp, err
}

func GetRAW() []byte {
	resp, err := http.Get("https://iqama.ccil-kbw.com/iqamatimes.php")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}
