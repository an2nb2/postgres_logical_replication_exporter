package collector

import (
	"postgres_logical_replication_exporter/pg"
	"strconv"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "pg_logical_replication"
)

var (
	subscriptionDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "subscription_status"),
		"Status of subscription.",
		[]string{
			"host",
			"database",
			"name",
			"relname",
		},
		nil,
	)

	subscriptionLagDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "subscription_lag_bytes"),
		"The amount of WAL records generated in the primary, but not yet applied in the standby.",
		[]string{
			"host",
			"database",
			"name",
			"relname",
		},
		nil,
	)

	publicationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "publication_status"),
		"Status of publication.",
		[]string{
			"host",
			"database",
			"application_name",
			"state",
			"tmp",
		},
		nil,
	)

	publicationLagDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "publication_lag"),
		"The amount of WAL records generated in the primary, but not yet sent to the standby.",
		[]string{
			"host",
			"database",
			"application_name",
		},
		nil,
	)

	replicationSlotStatusDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "replication_slot_status"),
		"Status of replication slot.",
		[]string{
			"host",
			"database",
			"name",
		},
		nil,
	)
)

type Collector struct {
	primary *pg.DB
	standby *pg.DB
	logger  log.Logger
}

func NewCollector(primary, standby *pg.DB, logger log.Logger) *Collector {
	return &Collector{primary, standby, logger}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- subscriptionDesc
	ch <- subscriptionLagDesc
	ch <- publicationDesc
	ch <- publicationLagDesc
	ch <- replicationSlotStatusDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var err error

	err = c.CollectSubscriptions(ch)
	if err != nil {
		_ = level.Error(c.logger).Log("msg", "failed to collect subscriptions", "err", err)
	}

	err = c.CollectPublications(ch)
	if err != nil {
		_ = level.Error(c.logger).Log("msg", "failed to collect publications", "err", err)
	}

	err = c.CollectReplicationSlots(ch)
	if err != nil {
		_ = level.Error(c.logger).Log("msg", "failed to collect replication slots", "err", err)
	}
}

func (c *Collector) CollectSubscriptions(ch chan<- prometheus.Metric) error {
	subs, err := c.standby.Subscriptions()
	if err != nil {
		return err
	}

	_ = level.Debug(c.logger).Log("msg", "subscriptions", "count", len(subs))
	for _, s := range subs {
		ch <- prometheus.MustNewConstMetric(
			subscriptionDesc,
			prometheus.GaugeValue,
			float64(getStatus(s.Pid.Valid)),
			c.standby.Hostname,
			c.standby.Database,
			s.Name,
			s.Relname.String,
		)
	}

	currntlsn, err := c.primary.CurrentWalLsn()
	if err != nil {
		return err
	}

	for _, s := range subs {
		lag, err := calculateLag(currntlsn, s.ReceivedLsn)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			subscriptionLagDesc,
			prometheus.GaugeValue,
			float64(lag),
			c.standby.Hostname,
			c.standby.Database,
			s.Name,
			s.Relname.String,
		)
	}

	return nil
}

func (c *Collector) CollectPublications(ch chan<- prometheus.Metric) error {
	pubs, err := c.primary.Publications()
	if err != nil {
		return err
	}

	_ = level.Debug(c.logger).Log("msg", "publications", "count", len(pubs))
	for _, p := range pubs {
		ch <- prometheus.MustNewConstMetric(
			publicationDesc,
			prometheus.GaugeValue,
			float64(getStatus(p.Active)),
			c.primary.Hostname,
			c.primary.Database,
			p.Name,
			p.State,
			strconv.FormatBool(p.Tmp),
		)
	}

	currntlsn, err := c.primary.CurrentWalLsn()
	if err != nil {
		return err
	}

	for _, p := range pubs {
		lag, err := calculateLag(currntlsn, p.SentLsn)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			publicationLagDesc,
			prometheus.GaugeValue,
			float64(lag),
			c.primary.Hostname,
			c.primary.Database,
			p.Name,
		)
	}

	return nil
}

func (c *Collector) CollectReplicationSlots(ch chan<- prometheus.Metric) error {
	slots, err := c.primary.ReplicationSlots()
	if err != nil {
		return err
	}

	_ = level.Debug(c.logger).Log("msg", "replication slots", "count", len(slots))
	for _, s := range slots {
		ch <- prometheus.MustNewConstMetric(
			replicationSlotStatusDesc,
			prometheus.GaugeValue,
			float64(getStatus(s.Active)),
			c.primary.Hostname,
			c.primary.Database,
			s.Name,
		)
	}

	return nil
}
