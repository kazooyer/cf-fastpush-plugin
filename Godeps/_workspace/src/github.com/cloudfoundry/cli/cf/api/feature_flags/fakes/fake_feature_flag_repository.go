// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/api/feature_flags"
	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/models"
)

type FakeFeatureFlagRepository struct {
	ListStub        func() ([]models.FeatureFlag, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct{}
	listReturns     struct {
		result1 []models.FeatureFlag
		result2 error
	}
	FindByNameStub        func(string) (models.FeatureFlag, error)
	findByNameMutex       sync.RWMutex
	findByNameArgsForCall []struct {
		arg1 string
	}
	findByNameReturns struct {
		result1 models.FeatureFlag
		result2 error
	}
	UpdateStub        func(string, bool) error
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 string
		arg2 bool
	}
	updateReturns struct {
		result1 error
	}
}

func (fake *FakeFeatureFlagRepository) List() ([]models.FeatureFlag, error) {
	fake.listMutex.Lock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct{}{})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub()
	} else {
		return fake.listReturns.result1, fake.listReturns.result2
	}
}

func (fake *FakeFeatureFlagRepository) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeFeatureFlagRepository) ListReturns(result1 []models.FeatureFlag, result2 error) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 []models.FeatureFlag
		result2 error
	}{result1, result2}
}

func (fake *FakeFeatureFlagRepository) FindByName(arg1 string) (models.FeatureFlag, error) {
	fake.findByNameMutex.Lock()
	fake.findByNameArgsForCall = append(fake.findByNameArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.findByNameMutex.Unlock()
	if fake.FindByNameStub != nil {
		return fake.FindByNameStub(arg1)
	} else {
		return fake.findByNameReturns.result1, fake.findByNameReturns.result2
	}
}

func (fake *FakeFeatureFlagRepository) FindByNameCallCount() int {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return len(fake.findByNameArgsForCall)
}

func (fake *FakeFeatureFlagRepository) FindByNameArgsForCall(i int) string {
	fake.findByNameMutex.RLock()
	defer fake.findByNameMutex.RUnlock()
	return fake.findByNameArgsForCall[i].arg1
}

func (fake *FakeFeatureFlagRepository) FindByNameReturns(result1 models.FeatureFlag, result2 error) {
	fake.FindByNameStub = nil
	fake.findByNameReturns = struct {
		result1 models.FeatureFlag
		result2 error
	}{result1, result2}
}

func (fake *FakeFeatureFlagRepository) Update(arg1 string, arg2 bool) error {
	fake.updateMutex.Lock()
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 string
		arg2 bool
	}{arg1, arg2})
	fake.updateMutex.Unlock()
	if fake.UpdateStub != nil {
		return fake.UpdateStub(arg1, arg2)
	} else {
		return fake.updateReturns.result1
	}
}

func (fake *FakeFeatureFlagRepository) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeFeatureFlagRepository) UpdateArgsForCall(i int) (string, bool) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return fake.updateArgsForCall[i].arg1, fake.updateArgsForCall[i].arg2
}

func (fake *FakeFeatureFlagRepository) UpdateReturns(result1 error) {
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 error
	}{result1}
}

var _ feature_flags.FeatureFlagRepository = new(FakeFeatureFlagRepository)
