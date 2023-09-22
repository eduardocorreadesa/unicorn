package service

import (
	"context"
	"math/rand"
	"sync"
	"time"
	"unicorn/internal/configuration"
	"unicorn/internal/domain"
	"unicorn/internal/repository"

	"github.com/google/uuid"
)

type UnicornService interface {
	CreateAsync(int) string
	GetResult(string) (domain.Result, bool)
}

func NewUnicornService(ctx context.Context, config configuration.Config, unicornRepo *repository.UnicornRepo) UnicornService {
	return &unicorn{
		ctx:         ctx,
		config:      config,
		unicornRepo: *unicornRepo,
		results:     make(map[string]domain.Result),
	}
}

type unicorn struct {
	ctx         context.Context
	config      configuration.Config
	unicornRepo repository.UnicornRepo
	mu          sync.Mutex
	results     map[string]domain.Result
}

func (u *unicorn) CreateAsync(amount int) string {
	u.mu.Lock()
	defer u.mu.Unlock()

	id := uuid.New().String()

	go func(id string) {
		sleep_time := time.Duration(rand.Intn(1000)) * time.Millisecond
		items := []domain.Unicorn{}
		for j := 0; j < amount; j++ {
			name := u.unicornRepo.UnicornAdj[rand.Intn(1345)] + "-" + u.unicornRepo.UnicornNames[rand.Intn(5800)]

			item := domain.Unicorn{
				Name:         name,
				Capabilities: u.capabilities(3),
			}

			items = append(items, item)
			time.Sleep(sleep_time)
		}

		result := domain.Result{
			ID:       id,
			Unicorns: items,
		}

		u.mu.Lock()
		defer u.mu.Unlock()

		u.results[id] = result
	}(id)

	return id
}

func (u *unicorn) GetResult(id string) (domain.Result, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()

	result, ok := u.results[id]
	return result, ok
}

func (u *unicorn) capabilities(qtd int) []string {
	capabilities := make([]string, 0, qtd)
	capacityMap := make(map[string]bool)
	capacityTotal := len(u.unicornRepo.UnicornCapabilities)

	for len(capabilities) < qtd && len(capacityMap) < capacityTotal {
		randomIndex := rand.Intn(capacityTotal)
		cap := u.unicornRepo.UnicornCapabilities[randomIndex]

		if !capacityMap[cap] {
			capabilities = append(capabilities, cap)
			capacityMap[cap] = true
		}
	}

	return capabilities
}
