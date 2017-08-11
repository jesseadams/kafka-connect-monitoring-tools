package kafka_connect

import(
  "net/http"
  "encoding/json"
  "crypto/tls"
  "time"
)

type KafkaConnector struct {
  State string
  WorkerId string
}

type KafkaConnectorTask struct {
  State string
  Id int
  WorkerId string
}

type KafkaConnectorStatus struct {
  Name string
  Connector KafkaConnector
  Tasks []KafkaConnectorTask
}

func CheckStatus(baseUrl string, connector string, target interface{}, validateSsl bool) error {
  endpoint := "/connectors/" + connector + "/status"
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: validateSsl},
  }
  var client = &http.Client{Transport: transport, Timeout: 10 * time.Second}
  response, err := client.Get(baseUrl + endpoint)

  if err != nil {
    return err
  }
  defer response.Body.Close()

  return json.NewDecoder(response.Body).Decode(target)
}

func ListConnectors(baseUrl string, validateSsl bool) ([]string, error) {
  endpoint := "/connectors"
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: validateSsl},
  }
  var client = &http.Client{Transport: transport, Timeout: 10 * time.Second}
  response, err := client.Get(baseUrl + endpoint)

  if err != nil {
    return nil, err
  }
  defer response.Body.Close()

  var connectors []string
  errors := json.NewDecoder(response.Body).Decode(&connectors)

  return connectors, errors
}
