package exporter

import (
	"errors"
	"github.com/ktsstudio/selectel-exporter/pkg/config"
	"github.com/ktsstudio/selectel-exporter/pkg/selapi"
	"log"
	"sync"
	"time"
)

type selectelCollector interface {
	GetInfo() string
	Collect(e *exporter) error
}

type exporter struct {
	token string
	region string
	project selapi.Project
	openstackAccountToken string
	lastTokenUpdate time.Time

	refreshPeriod time.Duration
	stopCh chan bool
	wg sync.WaitGroup
	collectors []selectelCollector

	datastores []selapi.Datastore
}

func Init(config *config.ExporterConfig, refreshPeriod time.Duration) (*exporter, error) {
	resp, err := selapi.FetchProjects(config.Token)
	if err != nil {
		return nil, err
	}

	// todo filter by project from config
	var project selapi.Project
	if len(resp.Projects) > 1 {
		project = resp.Projects[0]
		log.Println("%d projects have been found, %s was selected", len(resp.Projects), project.Name)
	} else if len(resp.Projects) == 1 {
		project = resp.Projects[0]
	} else {
		return nil, errors.New("there are no projects")
	}

	e := &exporter{
		token: config.Token,
		region: config.Region,
		project: project,
		refreshPeriod: refreshPeriod,
		stopCh: make(chan bool),
	}
	err = e.obtainToken()
	if err != nil {
		return nil, err
	}
	err = e.fetchDatastores()
	if err != nil {
		return nil, err
	}
	e.loadCollectors()
	go e.loop()

	return e, nil
}

func (e *exporter) obtainToken() error {
	res, err := selapi.ObtainToken(e.token, e.project.Id)
	if err != nil {
		return err
	}
	log.Println("openstack account token has been obtained successfully")
	e.openstackAccountToken = res.Token.Id
	e.lastTokenUpdate = time.Now()
	return nil
}

func (e *exporter) checkToken() error {
	 if e.lastTokenUpdate.Sub(time.Now()) > 24 * time.Hour {
	 	return e.obtainToken()
	 }
	 return nil
}

func (e *exporter) fetchDatastores() error {
	res, err := selapi.FetchDatastores(e.openstackAccountToken, e.region)
	if err != nil {
		return err
	}
	e.datastores = res.Datastores
	return nil
}

func (e *exporter) loadCollectors() {
	for _, ds := range e.datastores {
		e.collectors = append(e.collectors, NewDatastoreCollector(e.project, ds))
		e.collectors = append(e.collectors, NewDatabaseCollector(e.project, ds))
	}
	e.collectors = append(e.collectors, NewBalanceCollector(e.project))
}

func (e *exporter) runCollectors() {
	var wg sync.WaitGroup
	for _, col := range e.collectors {
		wg.Add(1)
		col := col
		go func() {
			defer wg.Done()
			log.Println(col.GetInfo())
			err := col.Collect(e)
			if err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()
}

func (e *exporter) loop()  {
	log.Println("exporter loop has started")
	e.wg.Add(1)
	for {
		err := e.checkToken()
		if err != nil {
			log.Println(err)
		} else {
			e.runCollectors()
		}
		select {
		case <-e.stopCh:
			log.Println("exporter loop has stopped")
			e.wg.Done()
			return
		case <-time.After(e.refreshPeriod * time.Second):
			continue
		}
	}
}

func (e *exporter) Stop() {
	e.stopCh <- true
	e.wg.Wait()
}
