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

func GenerateUrl() string {
  var endpoint = "/connectors/" + args.Connector + "/status"
  var hostString = args.ProtocolString + "://" + args.Host + ":" + strconv.Itoa(args.Port)
  url := hostString + endpoint

  return url
}

func CheckKafkaConnect(url string, target interface{}) error {
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: args.DontValidateSsl},
  }
  var client = &http.Client{Transport: transport, Timeout: 10 * time.Second}
  response, err := client.Get(url)

  if err != nil {
    return err
  }
  defer response.Body.Close()

  return json.NewDecoder(response.Body).Decode(target)
}
