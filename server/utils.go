package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SetupMetrics sets up a ginmetrics.Monitor with the given parameters, will expose the metrics on the given path and listen on the given address.
func SetupMetrics(otherRouter *gin.Engine, path, listen string, durations []float64) *ginmetrics.Monitor {
	router := gin.New()

	metrics := ginmetrics.GetMonitor()

	metrics.SetMetricPath(path)
	metrics.SetSlowTime(3)
	metrics.SetDuration(durations)
	if otherRouter != nil {
		metrics.UseWithoutExposingEndpoint(otherRouter)
		metrics.Expose(router)
	} else {
		metrics.Use(router)
	}

	go func() {
		err := router.Run(listen)
		if err != nil {
			panic(err)
		}
	}()
	return metrics
}

type HealthCheck struct {
	IsUnhealthy *atomic.Bool
}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{
		IsUnhealthy: &atomic.Bool{},
	}
}

func (h *HealthCheck) Health(c *gin.Context) {
	if !h.IsUnhealthy.Load() {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"status": "unhealthy"})
}

// ListenWithGracefulShutdown listens on the given address with the given gin router, and will shutdown gracefully when a SIGINT or SIGTERM signal is received.
// Also setups up a health check endpoint on the given path, which will be unhealthy when the server is shutting down.
func ListenWithGracefulShutdown(ctx context.Context, log logrus.Ext1FieldLogger, router *gin.Engine, conf ServerConfig) error {
	// Wrap the gin router in http.Server so we can call Shutdown
	hc := NewHealthCheck()
	router.GET(conf.HealthCheckPath, hc.Health)
	srv := &http.Server{
		Addr:              conf.Listen,
		Handler:           router.Handler(),
		ReadTimeout:       conf.ReadTimeout,
		ReadHeaderTimeout: conf.ReadTimeout,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to listen and serve")
		}
	}()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)

	signal := <-quitCh
	log.WithField("signal", signal).Warn("shutting down server")
	ctx, cancel := context.WithTimeout(ctx, conf.ShutdownTimeout)
	defer cancel()
	// Set the health check to unhealthy, so we can stop accepting new requests
	hc.IsUnhealthy.Store(true)
	if err := srv.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown server")
	}
	<-ctx.Done()
	return nil

}

// GetInt64Param gets the int64 parameter from the gin context, and returns the value or the default value if the parameter is not set.
func GetInt64Param(c *gin.Context, key string, maxVal, defaultVal int64) (int64, error) {
	val := c.Param(key)
	if len(val) == 0 {
		return defaultVal, nil
	}
	out, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to parse int64 param \"%s\" value \"%s\"", key, val)
	}
	if maxVal != -1 && out > maxVal {
		return 0, errors.Errorf("param \"%s\" too large, must be less than %d", key, maxVal)
	}
	return out, nil
}

// SafeMetricsInc increments the given metric with the given label values, and logs an error if the increment fails.
// All parameters can safely be nil, if metric is nil, this function does nothing.
func SafeMetricsInc(log logrus.Ext1FieldLogger, metric *ginmetrics.Metric, labelValues []string) {
	if metric == nil {
		return
	}
	if err := metric.Inc(labelValues); err != nil {
		if log != nil {
			log.WithError(err).Error("failed to increment metric")
		}
	}
}

// SafeMetricsAdd increased the given metric by value with the given label values, and logs an error if the increase fails.
// All parameters can safely be nil, if metric is nil, this function does nothing.
func SafeMetricsAdd(log logrus.Ext1FieldLogger, metric *ginmetrics.Metric, labelValues []string, value float64) {
	if metric == nil {
		return
	}
	if err := metric.Add(labelValues, value); err != nil {
		if log != nil {
			log.WithError(err).Error("failed to increment metric")
		}
	}
}

// SafeMetricsGauge sets the given metric to the given value with the given label values, and logs an error if the set fails.
func SafeMetricsGauge(log logrus.Ext1FieldLogger, metric *ginmetrics.Metric, labelValues []string, value float64) {
	if metric == nil {
		return
	}
	if err := metric.SetGaugeValue(labelValues, value); err != nil {
		if log != nil {
			log.WithError(err).Error("failed to set gauge metric")
		}
	}
}
