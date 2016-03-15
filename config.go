package producer

import (
	"time"

	"github.com/Sirupsen/logrus"
	k "github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/jpillora/backoff"
)

// Constants and default configuration take from:
// github.com/awslabs/amazon-kinesis-producer/.../KinesisProducerConfiguration.java
const (
	maxRecordSize        = 1 << 20 // 1MiB
	maxRequestSize       = 5 << 20 // 5MiB
	maxRecordsPerRequest = 500
	maxAggregationSize   = 51200 // 50KB
	// The KinesisProducerConfiguration set the default to 4294967295L;
	// it's kinda odd, because the maxAggregationSize is limit to 51200L;
	maxAggregationCount  = 4294967295
	defaultFlushInterval = time.Second
)

// Putter is the interface that wraps the KinesisAPI.PutRecords method.
type Putter interface {
	PutRecords(*k.PutRecordsInput) (*k.PutRecordsOutput, error)
}

type Config struct {
	// StreamName is the Kinesis stream.
	StreamName string

	// FlushInterval is a regular interval for flushing the buffer. Defaults to 1s.
	FlushInterval time.Duration

	// BatchCount determine the maximum number of items to pack into an Put
	// Must not exceed length. Defaults to 500.
	BatchCount int

	// BatchSize determine the maximum number of bytes to send with a PutRecords request.
	// Must not exceed 5MiB; Default to 5MiB.
	BatchSize int

	// AggregateBatchCount determine the maximum number of items to pack into an aggregated record.
	AggregateBatchCount int

	// AggregationBatchSize determine the maximum number of bytes to pack into an aggregated record.
	AggregateBatchSize int

	// BacklogCount determines the channel capacity before Put() will begin blocking. Default to `BatchLen`.
	BacklogCount int

	// Backoff determines the backoff strategy for record failures.
	Backoff backoff.Backoff

	// Logger is the logger used. Defaults to logrus.Log.
	Logger *logrus.Logger

	// Client is the Putter interface implementation.
	Client Putter

	// - Maximum number of connections to open to the backend.
	//   HTTP requests are sent in parallel over multiple connections.
}

// defaults for configuration
func (c *Config) defaults() {
	if c.Logger == nil {
		c.Logger = logrus.New()
	}
	if c.BatchCount == 0 {
		c.BatchCount = maxRecordsPerRequest
	}
	falseOrPanic(c.BatchCount > maxRecordsPerRequest, "kinesis: BatchCount exceeds 500")
	if c.BatchSize == 0 {
		c.BatchSize = maxRequestSize
	}
	falseOrPanic(c.BatchSize > maxRequestSize, "kinesis: BatchSize exceeds 5MiB")
	if c.BacklogCount == 0 {
		c.BacklogCount = maxRecordsPerRequest
	}
	if c.AggregateBatchCount == 0 {
		c.AggregateBatchCount = maxAggregationCount
	}
	falseOrPanic(c.AggregateBatchCount > maxAggregationCount, "kinesis: AggregateBatchCount exceeds 4294967295")
	if c.AggregateBatchSize == 0 {
		c.AggregateBatchSize = maxAggregationSize
	}
	falseOrPanic(c.AggregateBatchSize > maxAggregationSize, "kinesis: AggregateBatchSize exceeds 50KB")
	if c.FlushInterval == 0 {
		c.FlushInterval = time.Second
	}
}

func falseOrPanic(p bool, msg string) {
	if p {
		panic(msg)
	}
}