package kafka_connect

import(
  "net/http"
  "encoding/json"
  "crypto/tls"
  "time"
)

func ListVersions(baseUrl string, subject string, validateSsl bool) ([]string, error) {
  endpoint := "/subjects/" + subject + "/versions"
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: validateSsl},
  }
  var client = &http.Client{Transport: transport, Timeout: 10 * time.Second}
  response, err := client.Get(baseUrl + endpoint)

  if err != nil {
    return nil, err
  }
  defer response.Body.Close()

  var versions []string
  errors := json.NewDecoder(response.Body).Decode(&versions)

  return versions, errors
}

func ListSubjects(baseUrl string, validateSsl bool) ([]string, error) {
  endpoint := "/subjects"
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: validateSsl},
  }
  var client = &http.Client{Transport: transport, Timeout: 10 * time.Second}
  response, err := client.Get(baseUrl + endpoint)

  if err != nil {
    return nil, err
  }
  defer response.Body.Close()

  var subjects []string
  errors := json.NewDecoder(response.Body).Decode(&subjects)

  return subjects, errors
}
