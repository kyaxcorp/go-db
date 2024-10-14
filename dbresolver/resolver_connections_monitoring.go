package dbresolver

import (
	"time"

	"github.com/gookit/color"
	"github.com/kyaxcorp/go-db/codes"
	"github.com/kyaxcorp/go-helper/gor"
	"github.com/rs/zerolog"
)

func (r *resolver) startConnectionsMonitoring() error {

	info := func() *zerolog.Event {
		return r.LInfoF("connections_monitoring")
	}
	_error := func() *zerolog.Event {
		return r.LErrorF("connections_monitoring")
	}

	info().Msg("starting connections monitoring...")

	connMonitoring, _err := gor.GO(r.ctx, gor.Config{
		Logger: r.dbResolver.Logger,
		OnRun: func(instance *gor.GInstance) {
			info().Msg(color.LightGreen.Render("connections monitoring started..."))
			for {
				info().Msg(color.LightYellow.Render("monitoring..."))

				if instance.IsTerminating() {
					break
				}

				// Set inactive & active replicas!

				// TODO: we should set some settings which define:
				// how often it should ping the host
				// is monitoring enabled...usually it should be always available...

				var activeSources []detailedConnPool
				var inactiveSources []detailedConnPool
				var activeReplicas []detailedConnPool
				var inactiveReplicas []detailedConnPool

				for _, source := range r.sources {

					startTime := time.Now().Nanosecond()
					_err := checkConnection(
						source.pool,
						instance.GetContext(),
					)
					endTime := time.Now().Nanosecond()
					latency := endTime - startTime

					source.latencyNano = latency

					if _err != nil {
						// error...
						// Set as inactive...
						inactiveSources = append(inactiveSources, source)
					} else {
						// set as active
						activeSources = append(activeSources, source)
					}

					if instance.IsTerminating() {
						break
					}
				}

				for _, replica := range r.replicas {

					startTime := time.Now().Nanosecond()
					_err := checkConnection(
						replica.pool,
						instance.GetContext(),
					)
					endTime := time.Now().Nanosecond()
					latency := endTime - startTime

					replica.latencyNano = latency

					if _err != nil {
						// error...
						// Set as inactive...
						inactiveReplicas = append(inactiveReplicas, replica)
					} else {
						// set as active
						activeReplicas = append(activeReplicas, replica)
					}

					if instance.IsTerminating() {
						break
					}
				}

				// Active Sources
				r.activeSourcesLock.Lock()
				r.activeSources = activeSources
				r.nrOfActiveSources.Set(len(activeSources))
				r.activeSourcesLock.Unlock()

				// Inactive Sources
				r.inactiveSourcesLock.Lock()
				r.inactiveSources = inactiveSources
				r.nrOfInActiveSources.Set(len(inactiveSources))
				r.inactiveSourcesLock.Unlock()

				// Active Replicas
				r.activeReplicasLock.Lock()
				r.activeReplicas = activeReplicas
				r.nrOfActiveReplicas.Set(len(activeReplicas))
				r.activeReplicasLock.Unlock()
				// Inactive Replicas
				r.inactiveReplicasLock.Lock()
				r.inactiveReplicas = inactiveReplicas
				r.nrOfInActiveReplicas.Set(len(inactiveReplicas))
				r.inactiveReplicasLock.Unlock()
				// Set that has started!
				r.isMonitoringActive.True()

				// Set the status of the resolver...
				nrOfActiveSources := r.nrOfActiveSources.Get()
				nrOfActiveReplicas := r.nrOfActiveReplicas.Get()
				if nrOfActiveSources > 0 {
					r.resolverStatus.Set(ResolverReadyToProcess)
				} else if nrOfActiveReplicas > 0 && nrOfActiveSources == 0 {
					// it can read from replicas, but it can't do any write operations
					r.resolverStatus.Set(ResolverReadOnly)
				} else {
					// it can't do anything...
					r.resolverStatus.Set(ResolverAllNodesOffline)
				}

				// Let's add some color to status
				statusColor := color.LightGreen
				var statusMessage string
				switch r.resolverStatus.Get() {
				case ResolverReadyToProcess:
					statusColor = color.LightGreen
					statusMessage = "READY TO PROCESS"
				case ResolverReadOnly:
					statusColor = color.LightYellow
					statusMessage = "READ-ONLY"
				case ResolverAllNodesOffline:
					statusColor = color.LightRed
					statusMessage = "NO-ACTIVE-CONNECTIONS"
				}

				info().
					Uint16("resolver_status", r.resolverStatus.Get()).
					Int("nr_of_active_sources", r.nrOfActiveSources.Get()).
					Int("nr_of_active_replicas", r.nrOfActiveReplicas.Get()).
					Int("nr_of_inactive_sources", r.nrOfInActiveSources.Get()).
					Int("nr_of_inactive_replicas", r.nrOfInActiveReplicas.Get()).
					Msg(statusColor.Render(statusMessage))

				//if instance.IsTerminating() {
				//	break
				//}
				select {
				case <-instance.GetContext().Done():
					break
				case <-time.After(time.Second):
				}
				//time.Sleep(time.Second)
			}
		},
	})

	if _err != nil {
		// Error... return this error!
		_error().Err(codes.ErrFailedToGenerateResolversManagerMonitoringGoroutine).Msg("")
		return codes.ErrFailedToGenerateResolversManagerMonitoringGoroutine
	}
	// Set the connection revival here...
	r.connMonitoring = connMonitoring
	_err = r.connMonitoring.Run()
	if _err != nil {
		_error().Err(codes.ErrFailedToRunResolversManagerMonitoringGoroutine).Msg("")
		return codes.ErrFailedToRunResolversManagerMonitoringGoroutine
	}

	return nil
}
