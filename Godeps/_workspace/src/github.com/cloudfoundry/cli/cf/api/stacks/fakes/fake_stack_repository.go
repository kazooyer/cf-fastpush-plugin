// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/api/stacks"
	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/models"
)

type FakeStackRepository struct {
	FindByNameStub        func(name string) (stack models.Stack, apiErr error)
	findByNameMutex       sync.RWMutex
	findByNameArgsForCall []struct {
		name string
	}
	findByNameReturns struct {
		result1 models.Stack
		result2 error
	}
	FindAllStub        func() (stacks []models.Stack, apiErr error)
	findAllMutex       sync.RWMutex
	findAllArgsForCall []struct{}
	findAllReturns     struct {
		result1 []models.Stack
		result2 error
	}
}

func (fake *FakeStackRepository) FindByName(name string) (stack models.Stack, apiErr error) {
	fake.findByNameMutex.Lock()
	fake.findByNameArgsForCall = append(fake.findByNameArgsForCall, struct {
		name string
	}{name})
	fake.findByNameMutex.Unlock()
	if fake.FindByNameStub != nil {
		return fake.FindByNameStub(name)
	} else {
		return fake.findByNameReturns.result1, fake.findByNameReturns.result2
	}
}

func (fake *FakeStackRepository) FindByNameCallCount() int {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return len(fake.findByNameArgsForCall)
}

func (fake *FakeStackRepository) FindByNameArgsForCall(i int) string {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return fake.findByNameArgsForCall[i].name
}

func (fake *FakeStackRepository) FindByNameReturns(result1 models.Stack, result2 error) {
	fake.FindByNameStub = nil
	fake.findByNameReturns = struct {
		result1 models.Stack
		result2 error
	}{result1, result2}
}

func (fake *FakeStackRepository) FindAll() (stacks []models.Stack, apiErr error) {
	fake.findAllMutex.Lock()
	fake.findAllArgsForCall = append(fake.findAllArgsForCall, struct{}{})
	fake.findAllMutex.Unlock()
	if fake.FindAllStub != nil {
		return fake.FindAllStub()
	} else {
		return fake.findAllReturns.result1, fake.findAllReturns.result2
	}
}

func (fake *FakeStackRepository) FindAllCallCount() int {
	fake.findAllMutex.RLock()
	defer fake.findAllMutex.RUnlock()
	return len(fake.findAllArgsForCall)
}

func (fake *FakeStackRepository) FindAllReturns(result1 []models.Stack, result2 error) {
	fake.FindAllStub = nil
	fake.findAllReturns = struct {
		result1 []models.Stack
		result2 error
	}{result1, result2}
}

var _ stacks.StackRepository = new(FakeStackRepository)
