package cache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Cache ...
type Cache struct {
}

const (
	setURL = "http://202.182.118.203:5645/set?key=%s&value=%s"
	getURL = "http://202.182.118.203:5645/get?key=%s"
)

// Set 设置缓存
func (ca *Cache) Set(key, value string) error {
	url := fmt.Sprintf(setURL, key, value)
	str, err := httpGet(url)
	if err != nil {
		return err
	}
	if string(str) == "error" {
		return errors.New(" server res error ")
	}
	return nil
}

// Get 获取缓存
func (ca *Cache) Get(key string) (string, error) {
	url := fmt.Sprintf(getURL, key)
	str, err := httpGet(url)
	if err != nil {
		return "", err
	}
	if string(str) == "error" {
		return "", errors.New(" server res error ")
	}
	return string(str), nil
}

// httpGet get 请求
func httpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("response.StatusCode error")
	}
	return ioutil.ReadAll(response.Body)
}
