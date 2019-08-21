package datastore

import (
	"sync"
)

type SystemInfo struct {
	System      string
	Status      string
	Load5       string
	Load10      string
	Load15      string
	CpuUs       string
	CpuSy       string
	CpuWa       string
	MemoryPerc  string
	MemoryBytes string
	SwapPerc    string
	SwapBytes   string
}

type HostInfo struct {
	Host      string
	Status    string
	Protocols string
}

type DataStore struct {
	systemInfo *SystemInfo
	hostsInfo  []*HostInfo

	systemInfoMutex *sync.Mutex
	hostsInfoMutex  *sync.Mutex
}

func New() *DataStore {
	ds := DataStore{}

	ds.systemInfoMutex = &sync.Mutex{}
	ds.hostsInfoMutex = &sync.Mutex{}

	return &ds
}

func (ds *DataStore) ReadSystemInfo() *SystemInfo {
	ds.systemInfoMutex.Lock()
	defer ds.systemInfoMutex.Unlock()

	return ds.systemInfo
}

func (ds *DataStore) WriteSystemInfo(si *SystemInfo) {
	ds.systemInfoMutex.Lock()
	defer ds.systemInfoMutex.Unlock()

	ds.systemInfo = si
}

func (ds *DataStore) ReadHostInfos() []*HostInfo {
	ds.hostsInfoMutex.Lock()
	defer ds.hostsInfoMutex.Unlock()

	return ds.hostsInfo
}

func (ds *DataStore) WriteHostInfos(hi []*HostInfo) {
	ds.hostsInfoMutex.Lock()
	defer ds.hostsInfoMutex.Unlock()

	ds.hostsInfo = hi
}
