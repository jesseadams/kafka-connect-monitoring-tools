### Kafka Connect Monitoring Tools

These tools utilize the `/connectors/$connector_name/status` endpoint of http://docs.confluent.io/current/connect/ to verify there are healthy tasks running.

#### Installation

Download the appropriate binary from the release page: https://github.com/jesseadams/kafka-connect-monitoring-tools/releases

#### Usage

##### checks/check_kafka_connect

This was designed to be a Nagios/Sensu like check. You can use exit codes as you would expect to determine success.

###### Parameters

```go
Host string `arg:"required"`
Connector string `arg:"required"`
DontValidateSsl bool `arg:"--dont-validate-ssl"`
TaskCount int
Port int
Insecure bool
ProtocolString string
```

Example:

`./check_kafka_connect --host foo.example.com --connector my-connector-name --taskcount 1 --dont-validate-ssl`

##### cloudwatch/healthy_task_count

This is used to feed a HealthyTaskCount metric into AWS CloudWatch. By default, it will call a PutMetricData for Namespace: KafkaConnect, Metric: HealthyTaskCount, Unit: Count. The dimension defaults to the hostname of the server. You'll likely want to setup a cron job that runs every 5 minutes on the servers.

```go
Host string `arg:"required"`
Connector string `arg:"required"`
DontValidateSsl bool `arg:"--dont-validate-ssl"`
DimensionName string
DimensionValue string
Namespace string
Port int
Insecure bool
ProtocolString string
```

Example:

`./healthy_task_count --host foo.example.com --connector my-connector-name --dont-validate-ssl`
