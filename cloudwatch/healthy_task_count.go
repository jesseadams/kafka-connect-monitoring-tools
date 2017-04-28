package main

import(
  "fmt"
  "os"
  "strconv"
  "time"
  "github.com/alexflint/go-arg"
  "github.com/uscis/kafka-connect-monitoring-tools/common"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/cloudwatch"
)

var args struct {
  Host string `arg:"required"`
  Connector string `arg:"required"`
  DontValidateSsl bool `arg:"--dont-validate-ssl"`
  Port int
  Insecure bool
  ProtocolString string
}

func main() {
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
    MetricData: []*cloudwatch.MetricDatum{ // Required
        { // Required
            MetricName: aws.String("Test"), // Required
            Dimensions: []*cloudwatch.Dimension{
                { // Required
                    Name:  aws.String("DimensionName"),  // Required
                    Value: aws.String("DimensionValue"), // Required
                },
                // More values...
            },
            StatisticValues: &cloudwatch.StatisticSet{
                Maximum:     aws.Float64(1.0), // Required
                Minimum:     aws.Float64(1.0), // Required
                SampleCount: aws.Float64(1.0), // Required
                Sum:         aws.Float64(1.0), // Required
            },
            Timestamp: aws.Time(time.Now()),
            Unit:      aws.String("StandardUnit"),
            Value:     aws.Float64(1.0),
        },
        // More values...
    },
    Namespace: aws.String("Jesse"), // Required
  }
  resp, err := svc.PutMetricData(params)

  if err != nil {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
      return
  }

  // Pretty-print the response data.
  fmt.Println(resp)
}
