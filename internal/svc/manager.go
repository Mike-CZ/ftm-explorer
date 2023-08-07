package svc

import (
	"ftm-explorer/internal/config"
	"ftm-explorer/internal/logger"
	"ftm-explorer/internal/repository"
	"sync"
)

// Manager represents the manager controlling services lifetime.
type Manager struct {
	cfg  *config.Config
	repo repository.IRepository
	wg   sync.WaitGroup
	svc  []iService
	log  logger.ILogger

	// managed services
	blkScanner   *blockScanner
	blkObserver  *blockObserver
	metaObserver *metadataObserver
}

// NewServiceManager returns a new service manager.
func NewServiceManager(cfg *config.Config, repo *repository.Repository, log logger.ILogger) *Manager {
	// prepare the manager
	mgr := Manager{
		cfg:  cfg,
		repo: repo,
		svc:  make([]iService, 0),
		log:  log.ModuleLogger("svc_manager"),
	}
	mgr.init()
	return &mgr
}

// Start starts all the services prepared to be run.
func (mgr *Manager) Start() {
	// start services
	for _, s := range mgr.svc {
		s.start()
	}
}

// Close terminates the service manager
// and all the managed services along with it.
func (mgr *Manager) Close() {
	mgr.log.Notice("services are being terminated")

	for _, s := range mgr.svc {
		mgr.log.Noticef("closing %s", s.name())
		s.close()
	}

	mgr.wg.Wait()
	mgr.log.Notice("services closed")
}

// init initializes the services in the correct order.
func (mgr *Manager) init() {
	// make services
	mgr.blkScanner = newBlockScanner(mgr)
	mgr.svc = append(mgr.svc, mgr.blkScanner)

	mgr.blkObserver = newBlockObserver(mgr, mgr.blkScanner.scannedBlocks())
	mgr.svc = append(mgr.svc, mgr.blkObserver)

	mgr.metaObserver = newMetadataObserver(mgr)
	mgr.svc = append(mgr.svc, mgr.metaObserver)
}

// started signals to the manager that the calling service
// has been started and is functioning.
func (mgr *Manager) started(svc iService) {
	mgr.wg.Add(1)
	mgr.log.Noticef("%s is running", svc.name())
}

// finished signals to the manager that the calling service
// has been terminated and is no longer running.
func (mgr *Manager) finished(svc iService) {
	mgr.wg.Done()
	mgr.log.Noticef("%s terminated", svc.name())
}
