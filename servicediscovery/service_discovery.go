package servicediscovery

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/Trendyol/go-dcp-client/membership/info"
	"github.com/Trendyol/go-dcp-client/servicediscovery/model"
)

type ServiceDiscovery interface {
	Add(service *model.Service)
	Remove(name string)
	RemoveAll()
	AssignLeader(leaderService *model.Service)
	RemoveLeader()
	ReassignLeader() error
	StartHealthCheck()
	StopHealthCheck()
	StartRebalance()
	StopRebalance()
	GetAll() []string
	SetInfo(memberNumber int, totalMembers int)
	BeLeader()
	DontBeLeader()
}

type serviceDiscovery struct {
	leaderService       *model.Service
	services            map[string]*model.Service
	healthCheckSchedule *time.Ticker
	rebalanceSchedule   *time.Ticker
	info                *info.Model
	infoHandler         info.Handler
	amILeader           bool
	servicesLock        *sync.RWMutex
}

func (s *serviceDiscovery) Add(service *model.Service) {
	s.services[service.Name] = service
}

func (s *serviceDiscovery) Remove(name string) {
	if _, ok := s.services[name]; ok {
		_ = s.services[name].Client.Close()

		delete(s.services, name)
	}
}

func (s *serviceDiscovery) RemoveAll() {
	for name := range s.services {
		s.Remove(name)
	}
}

func (s *serviceDiscovery) BeLeader() {
	s.amILeader = true
}

func (s *serviceDiscovery) DontBeLeader() {
	s.amILeader = false
}

func (s *serviceDiscovery) AssignLeader(leaderService *model.Service) {
	s.leaderService = leaderService
}

func (s *serviceDiscovery) RemoveLeader() {
	if s.leaderService == nil {
		return
	}

	_ = s.leaderService.Client.Close()

	s.leaderService = nil
}

func (s *serviceDiscovery) ReassignLeader() error {
	if s.leaderService == nil {
		return fmt.Errorf("leader is not assigned")
	}

	err := s.leaderService.Client.Reconnect()

	if err == nil {
		err = s.leaderService.Client.Register()
	}

	return err
}

func (s *serviceDiscovery) healthCheckToService(service *model.Service, errorCallback func()) {
	if service == nil {
		return
	}

	if err := service.Client.Ping(); err != nil {
		errorCallback()
	}
}

func (s *serviceDiscovery) StartHealthCheck() {
	s.healthCheckSchedule = time.NewTicker(10 * time.Second)

	go func() {
		for range s.healthCheckSchedule.C {
			s.servicesLock.Lock()

			s.healthCheckToService(s.leaderService, func() {
				log.Println("leader is down")

				tempLeaderService := s.leaderService

				if err := s.ReassignLeader(); err != nil {
					if tempLeaderService != s.leaderService {
						_ = tempLeaderService.Client.Close()
					} else {
						s.RemoveLeader()
					}
				}
			})

			for name, service := range s.services {
				s.healthCheckToService(service, func() {
					s.Remove(name)

					log.Printf("client %s disconnected", name)
				})
			}

			s.servicesLock.Unlock()
		}
	}()
}

func (s *serviceDiscovery) StopHealthCheck() {
	s.healthCheckSchedule.Stop()
}

func (s *serviceDiscovery) StartRebalance() {
	s.rebalanceSchedule = time.NewTicker(10 * time.Second)

	go func() {
		for range s.rebalanceSchedule.C {
			if !s.amILeader {
				continue
			}

			s.servicesLock.RLock()

			names := s.GetAll()
			totalMembers := len(names) + 1

			s.SetInfo(1, totalMembers)

			for index, name := range names {
				if service, ok := s.services[name]; ok {
					if err := service.Client.Rebalance(index+2, totalMembers); err != nil {
						log.Printf("rebalance failed for %s", name)
					}
				}
			}

			s.servicesLock.RUnlock()
		}
	}()
}

func (s *serviceDiscovery) StopRebalance() {
	if s.rebalanceSchedule != nil {
		s.rebalanceSchedule.Stop()
	}
}

func (s *serviceDiscovery) GetAll() []string {
	var names []string

	for name := range s.services {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func (s *serviceDiscovery) SetInfo(memberNumber int, totalMembers int) {
	newInfo := &info.Model{
		MemberNumber: memberNumber,
		TotalMembers: totalMembers,
	}

	if newInfo.IsChanged(s.info) {
		s.info = newInfo

		log.Printf("new info arrived, member number: %d, total members: %d", memberNumber, totalMembers)

		s.infoHandler.OnModelChange(newInfo)
	}
}

func NewServiceDiscovery(infoHandler info.Handler) ServiceDiscovery {
	return &serviceDiscovery{
		services:     make(map[string]*model.Service),
		infoHandler:  infoHandler,
		servicesLock: &sync.RWMutex{},
	}
}
