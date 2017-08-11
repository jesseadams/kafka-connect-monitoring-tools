package main

import (
  "io"
  "net/http"
  "log"
  "os"
  "fmt"
  "github.com/jesseadams/kafka-connect-monitoring-tools/common"
)

func PublishPrometheusMetrics(w http.ResponseWriter, req *http.Request) {
  hostString := os.Getenv("KAFKA_CONNECT_URL")
  connectors, err := kafka_connect.ListConnectors(hostString, true)
  if err != nil {
    fmt.Println(err)
  }

  connectorCount := len(connectors)
  io.WriteString(w, "# TYPE kafka_connect_connectorcount gauge\n")
  line := fmt.Sprintf("kafka_connect_connectorcount %.1f\n", float64(connectorCount))
  io.WriteString(w, line)

  io.WriteString(w, "# TYPE kafka_connect_runningtaskscount gauge\n")
  for _, connector := range connectors {
      status := new(kafka_connect.KafkaConnectorStatus)
      err = kafka_connect.CheckStatus(hostString, connector, status, true)

      if err != nil {
        fmt.Println(err)
      }

      runningTasksCount := 0.0
      for _, task := range status.Tasks {
        if task.State == "RUNNING" {
          runningTasksCount += 1.0
        }
      }
      line := fmt.Sprintf("kafka_connect_runningtaskscount{connector=\"%s\"} %.1f\n", connector, runningTasksCount)
      io.WriteString(w, line)
  }
}

func main() {
  http.HandleFunc("/metrics", PublishPrometheusMetrics)
  log.Fatal(http.ListenAndServe(":9045", nil))
}
