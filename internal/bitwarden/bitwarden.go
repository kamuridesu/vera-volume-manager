package bitwarden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	c "github.com/kamuridesu/vera-volume-manager/internal/config"
	u "github.com/kamuridesu/vera-volume-manager/internal/utils"
)

type Bitwarden struct {
	Config  c.Bitwarden
	Locked  bool
	Session string
}

func NewBitwarden(config c.Bitwarden) *Bitwarden {
	return &Bitwarden{
		Config: config,
		Locked: true,
	}
}

func (b *Bitwarden) Unlock() error {
	password, err := u.DecodeBase64String(b.Config.PasswordB64)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader([]byte(fmt.Sprintf(`{"password": "%s"}`, password)))
	res, err := http.Post(b.Config.Url+"/unlock", "application/json", bodyReader)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("could not unlock session, response: %s", resBody)
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return err
	}
	val, ok := response["data"]
	if ok {
		data := val.(map[string]interface{})
		val, ok = data["raw"]
		if ok {
			b.Session = val.(string)
			b.Locked = false
			os.Setenv("BW_SESSION", b.Session)
			return nil
		}
		return fmt.Errorf("could not find raw data in response")
	}
	return fmt.Errorf("could not find data in response")
}

func (b *Bitwarden) Lock() error {
	if b.Locked {
		return nil
	}
	req, err := http.Get(b.Config.Url + "/lock")
	if err != nil {
		return err
	}
	defer req.Body.Close()
	if req.StatusCode != 200 {
		return fmt.Errorf("could not lock session, response: %s", req.Body)
	}
	os.Unsetenv("BW_SESSION")
	b.Locked = true
	return nil
}

func (b *Bitwarden) GetItem(id string) (string, error) {
	if b.Locked {
		return "", fmt.Errorf("session is locked")
	}
	res, err := http.Get(b.Config.Url + "/object/item/" + id)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("could not get item, response: %s", resBody)
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return "", err
	}
	password, err := b.getPassword(&response)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (b *Bitwarden) GetItemByName(name string) (string, error) {
	items, err := b.ListItems()
	if err != nil {
		return "", err
	}
	data, ok := (*items)["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not find data in response")
	}
	iter, ok := data["data"].([]interface{})
	if !ok {
		return "", fmt.Errorf("could not find data in response")
	}
	for _, item := range iter {
		itemData := item.(map[string]interface{})
		if itemData["name"] == name {
			password, err := b.getPassword(&itemData)
			if err != nil {
				return "", err
			}
			return password, nil
		}
	}
	return "", fmt.Errorf("could not find item with name %s", name)
}

func (b *Bitwarden) getPassword(item *map[string]interface{}) (string, error) {
	login, ok := (*item)["login"].(map[string]interface{})
	if !ok {
		password, ok := (*item)["notes"]
		if !ok {
			return "", fmt.Errorf("could not find notes in response")
		}
		return password.(string), nil
	}
	password, ok := login["password"]
	if !ok {
		return "", fmt.Errorf("could not find password in response")
	}
	return password.(string), nil
}

func (b *Bitwarden) ListItems() (*map[string]interface{}, error) {
	if b.Locked {
		return nil, fmt.Errorf("session is locked")
	}
	req, err := http.Get(b.Config.Url + "/list/object/items")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	resBody, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if req.StatusCode != 200 {
		return nil, fmt.Errorf("could not get item, response: %s", resBody)
	}
	response := map[string]interface{}{}
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (b *Bitwarden) GetPassword() (string, error) {
	return b.GetItemByName(b.Config.CredentialName)
}
