package main

import(
  "fmt"
  "os"
  "strconv"
  "net/http"
  "encoding/json"
  "crypto/tls"
  "time"
  "github.com/alexflint/go-arg"
  "github.com/uscis/kafka-connect-monitoring-tools/common/kafka_connect"
)

var args struct {
  Host string `arg:"required"`
  Connector string `arg:"required"`
  DontValidateSsl bool `arg:"--dont-validate-ssl"`
  TaskCount int
  Port int
  Insecure bool
  ProtocolString string
}

func main() {
  args.TaskCount = 1
  arg.MustParse(&args)

  if args.Port == 0 {
    if args.Insecure {
      args.Port = 80
      args.ProtocolString = "http"
    } else {
      args.Port = 443
      args.ProtocolString = "https"
    }
  }

  url := GenerateUrl()
  status := new(KafkaConnectorStatus)
  err := CheckKafkaConnect(url, status)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  tasksCount := len(status.Tasks)
  if tasksCount != args.TaskCount {
    fmt.Printf("Task count is off! Wanted: %d, Actual: %d\n", args.TaskCount, tasksCount)
    os.Exit(1)
  } else {
    fmt.Println("Task count OK")
  }

  failure := false
  for _, task := range status.Tasks {
    fmt.Printf("Task ID %d is %s\n", task.Id, task.State)
    if task.State != "RUNNING" {
      failure = true
    }
  }

  if failure {
    fmt.Println("One more more tasks are not running!")
    os.Exit(1)
  }
}
