package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	iotextypes "github.com/iotexproject/iotex-proto/golang/iotextypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

const (
	LabelEndpoint = "endpoint"

	Height                       = "height"
	NumActions                   = "num_actions"
	Tps                          = "tps"
	EpochNum                     = "epoch_num"
	EpochHeight                  = "epoch_height"
	EpochGravityChainStartHeight = "epoch_gravity_chain_start_height"
	TpsFloat                     = "tps_float"

	// Namespace is the metrics prefix
	Namespace = "iotex_chainmeta"
)

var (
	// BuildTime represents the time of the build
	BuildTime = "N/A"
	// Version represents the Build SHA-1 of the binary
	Version = "N/A"

	// labels are the static labels that come with every metric
	labels = []string{LabelEndpoint}

	// metatOpts is the number of Iotex Chain meta reported
	metaHeightOpts = prometheus.GaugeOpts{
		Name:      Height,
		Namespace: Namespace,
		Help:      "Gauge for Iotex Chain metadata Height",
	}
	metaNumActionsOpts = prometheus.GaugeOpts{
		Name:      NumActions,
		Namespace: Namespace,
		Help:      "Gauge for Iotex Chain metadata NumActions",
	}
	metaTpsOpts = prometheus.GaugeOpts{
		Name:      Tps,
		Namespace: Namespace,
		Help:      "Gauge for Iotex Chain metadata Tps",
	}
	metaEpochNumOpts = prometheus.GaugeOpts{
		Name:      EpochNum,
		Namespace: Namespace,
		Help:      "Gauge for Iotex Chain metadata EpochNum",
	}
	metaEpochHeightOpts = prometheus.GaugeOpts{
		Name:      EpochHeight,
		Namespace: Namespace,
		Help:      "Gauge for Iotex Chain metadata EpochHeight",
	}
	metaEpochGravityChainStartHeightOpts = prometheus.GaugeOpts{
		Name:      EpochGravityChainStartHeight,
		Namespace: Namespace,
		Help:      "Gauge for Iotex Chain metadata EpochGravityChainStartHeight",
	}
	metaTpsFloatOpts = prometheus.GaugeOpts{
		Name:      TpsFloat,
		Namespace: Namespace,
		Help:      "Gauge for Iotex Chain metadata TpsFloat",
	}
)

type exporter struct {
	endpoint string
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName(metaHeightOpts.Namespace, metaHeightOpts.Subsystem, metaHeightOpts.Name),
		metaHeightOpts.Help,
		labels,
		nil,
	)

	ch <- prometheus.NewDesc(
		prometheus.BuildFQName(metaNumActionsOpts.Namespace, metaNumActionsOpts.Subsystem, metaNumActionsOpts.Name),
		metaNumActionsOpts.Help,
		labels,
		nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName(metaTpsOpts.Namespace, metaTpsOpts.Subsystem, metaTpsOpts.Name),
		metaTpsOpts.Help,
		labels,
		nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName(metaEpochNumOpts.Namespace, metaEpochNumOpts.Subsystem, metaEpochNumOpts.Name),
		metaEpochNumOpts.Help,
		labels,
		nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName(metaEpochHeightOpts.Namespace, metaEpochHeightOpts.Subsystem, metaEpochHeightOpts.Name),
		metaEpochHeightOpts.Help,
		labels,
		nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName(metaEpochGravityChainStartHeightOpts.Namespace, metaEpochGravityChainStartHeightOpts.Subsystem, metaEpochGravityChainStartHeightOpts.Name),
		metaEpochGravityChainStartHeightOpts.Help,
		labels,
		nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName(metaTpsFloatOpts.Namespace, metaTpsFloatOpts.Subsystem, metaTpsFloatOpts.Name),
		metaTpsFloatOpts.Help,
		labels,
		nil,
	)
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	metadata, err := getMetadata(e.endpoint)
	if err != nil {
		log.Fatal(err)
		return
	}

	heightGv := prometheus.NewGaugeVec(metaHeightOpts, labels)
	heightGv.WithLabelValues(e.endpoint).Set(float64(metadata.Height))
	heightGv.Collect(ch)

	numActionsGv := prometheus.NewGaugeVec(metaNumActionsOpts, labels)
	numActionsGv.WithLabelValues(e.endpoint).Set(float64(metadata.NumActions))
	numActionsGv.Collect(ch)

	tpsGv := prometheus.NewGaugeVec(metaTpsOpts, labels)
	tpsGv.WithLabelValues(e.endpoint).Set(float64(metadata.Tps))
	tpsGv.Collect(ch)

	epochNumGv := prometheus.NewGaugeVec(metaEpochNumOpts, labels)
	epochNumGv.WithLabelValues(e.endpoint).Set(float64(metadata.Epoch.Num))
	epochNumGv.Collect(ch)

	epochHeightGv := prometheus.NewGaugeVec(metaEpochHeightOpts, labels)
	epochHeightGv.WithLabelValues(e.endpoint).Set(float64(metadata.Epoch.Height))
	epochHeightGv.Collect(ch)

	epochGravityChainStartHeightGv := prometheus.NewGaugeVec(metaEpochGravityChainStartHeightOpts, labels)
	epochGravityChainStartHeightGv.WithLabelValues(e.endpoint).Set(float64(metadata.Epoch.GravityChainStartHeight))
	epochGravityChainStartHeightGv.Collect(ch)

	tpsFloatGv := prometheus.NewGaugeVec(metaTpsFloatOpts, labels)
	tpsFloatGv.WithLabelValues(e.endpoint).Set(float64(metadata.TpsFloat))
	tpsFloatGv.Collect(ch)
}

func init() {
	prometheus.MustRegister(version.NewCollector("iotex_chainmeta_exporter"))
}

func main() {
	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" {
		endpoint = "api.mainnet.iotex.one:80"
	}

	exporter := &exporter{
		endpoint: endpoint,
	}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9961", nil))
}

func getMetadata(endpoint string) (*iotextypes.ChainMeta, error) {
	grpcCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(grpcCtx, endpoint, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	resp, err := iotexapi.NewAPIServiceClient(conn).
		GetChainMeta(context.Background(), &iotexapi.GetChainMetaRequest{})
	if err != nil {
		return nil, err
	}

	return resp.ChainMeta, nil
}
