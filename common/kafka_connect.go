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

func CheckStatus(url string, target interface{}, validateSsl bool) error {
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: validateSsl},
  }
  var client = &http.Client{Transport: transport, Timeout: 10 * time.Second}
  response, err := client.Get(url)

  if err != nil {
    return err
  }
  defer response.Body.Close()

  return json.NewDecoder(response.Body).Decode(target)
}
