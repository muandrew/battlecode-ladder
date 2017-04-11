package utils

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func ReadBody(r *http.Response, t interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(t)
}

func GetBody(r *http.Response) string {
	defer r.Body.Close()
	contents, _ := ioutil.ReadAll(r.Body)
	return fmt.Sprintf("%s", contents)
}
