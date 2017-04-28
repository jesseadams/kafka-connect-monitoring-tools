package main

import(
  "fmt"
  "os"
  "strconv"
  "time"
  "github.com/alexflint/go-arg"
  "github.com/jesseadams/kafka-connect-monitoring-tools/common"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/cloudwatch"
)

var args struct {
  Host string `arg:"required"`
  Connector string `arg:"required"`
  DontValidateSsl bool `arg:"--dont-validate-ssl"`
  DimensionName string
  DimensionValue string
  Namespace string
  Port int
  Insecure bool
  ProtocolString string
}

func main() {
  args.Namespace = "KafkaConnect"
  args.DimensionName = "Host"
  args.DimensionValue, err = os.Hostname()

  if err != nil {
    args.DimensionValue = "Unknown"
  }

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

  endpoint := "/connectors/" + args.Connector + "/status"
  hostString := args.ProtocolString + "://" + args.Host + ":" + strconv.Itoa(args.Port)
  url := hostString + endpoint
  status := new(kafka_connect.KafkaConnectorStatus)
  err := kafka_connect.CheckStatus(url, status, args.DontValidateSsl)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  tasksCount := len(status.Tasks)
  fmt.Println(tasksCount)

  sess := session.Must(session.NewSession())
  svc := cloudwatch.New(sess)
  params := &cloudwatch.PutMetricDataInput{
    MetricData: []*cloudwatch.MetricDatum{
        {
            MetricName: aws.String("HealthyTaskCount"),
            Dimensions: []*cloudwatch.Dimension{
                {
                    Name:  aws.String(args.DimensionName),
                    Value: aws.String(args.DimensionValue),
                },
            },
            Timestamp: aws.Time(time.Now()),
            Unit:      aws.String("Count"),
            Value:     aws.Float64(float64(tasksCount)),
        },
    },
    Namespace: aws.String(args.Namespace),
  }
  resp, err := svc.PutMetricData(params)

  if err != nil {
      fmt.Println(err.Error())
      return
  }
}
