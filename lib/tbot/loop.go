/*
 * Teleport
 * Copyright (C) 2024  Gravitational, Inc.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package tbot

import (
	"context"
	"log/slog"
	"math"
	"time"

	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gravitational/teleport/api/utils/retryutils"
	"github.com/gravitational/teleport/lib/tbot/readyz"
)

var (
	loopIterationsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tbot_task_iterations_total",
			Help: "Number of task iteration attempts, not counting retries",
		}, []string{"service", "name"},
	)
	loopIterationsSuccessCounter = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tbot_task_iterations_successful",
			Help:    "Histogram of task iterations that ultimately succeeded, bucketed by number of retries before success",
			Buckets: []float64{0, 1, 2, 3, 4, 5},
		}, []string{"service", "name"},
	)
	loopIterationsFailureCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tbot_task_iterations_failed",
			Help: "Number of task iterations that ultimately failed, not counting retries",
		}, []string{"service", "name"},
	)
	loopIterationTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tbot_task_iteration_duration_seconds",
			Help:    "Time between beginning and ultimate end of one task iteration regardless of outcome, including all retries",
			Buckets: prometheus.ExponentialBuckets(0.1, 1.75, 6),
		}, []string{"service", "name"},
	)
)

type runOnIntervalConfig struct {
	service string
	name    string
	f       func(ctx context.Context) error
	clock   clockwork.Clock
	// reloadCh allows the task to be triggered immediately, ideal for handling
	// CA rotations or a manual signal from a user.
	// reloadCh can be nil, in which case, the task will only run on the
	// interval.
	reloadCh             chan struct{}
	log                  *slog.Logger
	interval             time.Duration
	retryLimit           int
	exitOnRetryExhausted bool
	waitBeforeFirstRun   bool
	// identityReadyCh allows the service to wait until the internal bot identity
	// renewal has completed before running, to avoid spamming the logs if the
	// service doesn't support gracefully degrading when there is no API client
	// available.
	identityReadyCh <-chan struct{}
	statusReporter  readyz.Reporter
}

func (cfg *runOnIntervalConfig) checkAndSetDefaults() error {
	switch {
	case cfg.interval <= 0:
		return trace.BadParameter("interval must be greater than 0")
	case cfg.retryLimit < 0:
		return trace.BadParameter("retryLimit must be greater than or equal to 0")
	case cfg.log == nil:
		return trace.BadParameter("log is required")
	case cfg.f == nil:
		return trace.BadParameter("f is required")
	case cfg.name == "":
		return trace.BadParameter("name is required")
	}

	if cfg.clock == nil {
		cfg.clock = clockwork.NewRealClock()
	}
	if cfg.statusReporter == nil {
		cfg.statusReporter = readyz.NoopReporter()
	}

	return nil
}

// runOnInterval runs a function on a given interval, with retries and jitter.
//
// TODO(noah): Emit Prometheus metrics for:
// - Time of next attempt
func runOnInterval(ctx context.Context, cfg runOnIntervalConfig) error {
	if err := cfg.checkAndSetDefaults(); err != nil {
		return err
	}

	log := cfg.log.With("task", cfg.name)

	if cfg.identityReadyCh != nil {
		select {
		case <-cfg.identityReadyCh:
		default:
			log.InfoContext(ctx, "Waiting for internal bot identity to be renewed before running")
			select {
			case <-cfg.identityReadyCh:
			case <-ctx.Done():
				return nil
			}
		}
	}

	ticker := cfg.clock.NewTicker(cfg.interval)
	defer ticker.Stop()
	jitter := retryutils.DefaultJitter
	firstRun := true
	for {
		if !firstRun || (firstRun && cfg.waitBeforeFirstRun) {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.Chan():
			case <-cfg.reloadCh:
			}
		}
		firstRun = false

		loopIterationsCounter.WithLabelValues(cfg.service, cfg.name).Inc()
		startTime := time.Now()

		var err error
		for attempt := 1; attempt <= cfg.retryLimit; attempt++ {
			log.InfoContext(
				ctx,
				"Attempting task",
				"attempt", attempt,
				"retry_limit", cfg.retryLimit,
			)
			err = cfg.f(ctx)
			if err == nil {
				cfg.statusReporter.Report(readyz.Healthy)
				loopIterationsSuccessCounter.WithLabelValues(cfg.service, cfg.name).Observe(float64(attempt - 1))
				break
			}

			if attempt != cfg.retryLimit {
				// exponentially back off with jitter, starting at 1 second.
				backoffTime := time.Second * time.Duration(math.Pow(2, float64(attempt-1)))
				backoffTime = jitter(backoffTime)
				cfg.log.WarnContext(
					ctx,
					"Task failed. Backing off and retrying",
					"attempt", attempt,
					"retry_limit", cfg.retryLimit,
					"backoff", backoffTime,
					"error", err,
				)
				select {
				case <-ctx.Done():
					// Note: will discard metric update for this loop. It
					// probably won't be collected if we're shutting down,
					// anyway.
					return nil
				case <-cfg.clock.After(backoffTime):
				}
			}
		}

		loopIterationTime.WithLabelValues(cfg.service, cfg.name).Observe(time.Since(startTime).Seconds())

		if err != nil {
			cfg.statusReporter.ReportReason(readyz.Unhealthy, err.Error())
			loopIterationsFailureCounter.WithLabelValues(cfg.service, cfg.name).Inc()

			if cfg.exitOnRetryExhausted {
				log.ErrorContext(
					ctx,
					"All retry attempts exhausted. Exiting",
					"error", err,
					"retry_limit", cfg.retryLimit,
				)
				return trace.Wrap(err)
			}
			log.WarnContext(
				ctx,
				"All retry attempts exhausted. Will wait for next interval",
				"retry_limit", cfg.retryLimit,
				"interval", cfg.interval,
			)
		} else {
			log.InfoContext(
				ctx,
				"Task succeeded. Waiting interval",
				"interval", cfg.interval,
			)
		}
	}
}
