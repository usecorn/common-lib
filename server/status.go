package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Database interface {
	PingContext(ctx context.Context) error
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type StatusChecker struct {
	keys       map[string]struct{}
	databases  map[string]Database
	ethclients map[string]EthClient
	log        logrus.Ext1FieldLogger
}

func NewStatusChecker(log logrus.Ext1FieldLogger) *StatusChecker {
	return &StatusChecker{
		log:        log,
		keys:       make(map[string]struct{}),
		databases:  make(map[string]Database),
		ethclients: make(map[string]EthClient),
	}
}

func (sc *StatusChecker) AddDatabase(key string, db Database) error {
	if _, ok := sc.keys[key]; ok {
		return errors.Errorf("key %s already exists", key)
	}
	sc.keys[key] = struct{}{}
	sc.databases[key] = db
	return nil
}

func (sc *StatusChecker) AddEthClient(key string, ec EthClient) error {
	if _, ok := sc.keys[key]; ok {
		return errors.Errorf("key %s already exists", key)
	}
	sc.keys[key] = struct{}{}
	sc.ethclients[key] = ec
	return nil
}

type healthCheckResult struct {
	key string
	err error
}

func (sc *StatusChecker) Check(c *gin.Context) {
	resChan := make(chan healthCheckResult)

	for key, db := range sc.databases {
		go func(key string, db Database) {
			err := db.PingContext(c)
			resChan <- healthCheckResult{key, err}
		}(key, db)
	}

	for key, ec := range sc.ethclients {
		go func(key string, ec EthClient) {
			_, err := ec.BlockNumber(c)
			resChan <- healthCheckResult{key, err}
		}(key, ec)
	}
	healthy := true
	status := map[string]string{}
	for range sc.keys {
		res := <-resChan
		if res.err != nil {
			if sc.log != nil {
				sc.log.WithField("key", res.key).WithError(res.err).Error("health check failed")
			}
			healthy = false
			status[res.key] = res.err.Error()
		} else {
			status[res.key] = "ok"
		}
	}

	if healthy {
		c.JSON(http.StatusOK, status)
	} else {
		c.JSON(http.StatusInternalServerError, status)
	}
}
