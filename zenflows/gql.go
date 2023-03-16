package zenflows

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	zenroom "github.com/dyne/Zenroom/bindings/golang/zenroom"
	"github.com/interfacerproject/zenflows-bank/config"
	"io"
	"net/http"
)

const SIGN = `
Scenario eddsa: sign a graph query
Given I have a 'base64' named 'gql'
Given I have a 'keyring'
# Fix Apollo's mingling with query string
When I remove spaces in 'gql'
and I compact ascii strings in 'gql'
When I create the eddsa signature of 'gql'
And I create the hash of 'gql'
Then print 'eddsa signature' as 'base64'
Then print 'gql' as 'base64'
Then print 'hash' as 'hex'
`

const GQL_PERSON string = "query($id: ID!) {person(id: $id) {id name note ethereumAddress}}"

type Agent struct {
	Sk          string
	ZenflowsUrl string
}

func (za *Agent) signRequest(jsonData []byte) (string, string) {
	data := fmt.Sprintf(`{"gql": "%s"}`, b64.StdEncoding.EncodeToString(jsonData))
	keys := fmt.Sprintf(`{"keyring": {"eddsa": "%s"}}`, za.Sk)
	result, success := zenroom.ZencodeExec(SIGN, "", data, keys)
	if !success {
		panic(result.Logs)
	}
	var resDecoded map[string]string
	if err := json.Unmarshal([]byte(result.Output), &resDecoded); err != nil {
		panic(err)
	}
	return "zenflows-sign", resDecoded["eddsa_signature"]
}

type Person struct {
	Id              string
	Name            string
	Note            string
	EthereumAddress string
}

func (za *Agent) GetPerson(id string) (*Person, error) {
	query, err := json.Marshal(map[string]interface{}{
		"query": GQL_PERSON,
		"variables": map[string]string{
			"id": id,
		},
	})

	body, err := za.makeRequest(query)
	if err != nil {
		return nil, err
	}
	var result map[string]map[string]map[string]string
	json.Unmarshal(body, &result)
	if result["data"]["person"] == nil {
		return nil, errors.New("Empty response from zenflows")
	}
	return &Person{
		Id:              result["data"]["person"]["id"],
		Name:            result["data"]["person"]["name"],
		Note:            result["data"]["person"]["note"],
		EthereumAddress: result["data"]["person"]["ethereumAddress"],
	}, nil
}

func (za *Agent) makeRequest(query []byte) ([]byte, error) {
	r, err := http.NewRequest("POST", za.ZenflowsUrl, bytes.NewReader(query))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(za.signRequest(query))
	r.Header.Add("zenflows-user", config.Config.ZenflowsUser)
	// TODO: move it outside
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
