package vk

import (
	"github.com/thedoctorde/vk-api"
	"net/http"
	"sync"
)

type Manager struct {
	Client      *vkapi.Client
	Groups      map[int64]string // id - name
	MutexGroups sync.Mutex
}

func NewManager(accessToken string) (*Manager, error) {
	apiClient := vkapi.NewApiClient()
	myClient := &http.Client{}
	apiClient.SetHTTPClient(myClient)
	apiClient.SetAccessToken(accessToken)
	client, err := vkapi.NewClientFromAPIClient(apiClient)
	return &Manager{
		Client: client,
		Groups: make(map[int64]string, 0),
	}, err
}

func (m *Manager) FillGroups(ids []int64) {
	for _, id := range ids {
		m.Groups[id] = ""
	}
}
