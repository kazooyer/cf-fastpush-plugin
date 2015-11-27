// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry/loggregatorlib/logmessage"
	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/api"
)

type FakeOldLogsRepository struct {
	RecentLogsForStub        func(appGuid string) ([]*logmessage.LogMessage, error)
	recentLogsForMutex       sync.RWMutex
	recentLogsForArgsForCall []struct {
		appGuid string
	}
	recentLogsForReturns struct {
		result1 []*logmessage.LogMessage
		result2 error
	}
	TailLogsForStub        func(appGuid string, onConnect func(), onMessage func(*logmessage.LogMessage)) error
	tailLogsForMutex       sync.RWMutex
	tailLogsForArgsForCall []struct {
		appGuid   string
		onConnect func()
		onMessage func(*logmessage.LogMessage)
	}
	tailLogsForReturns struct {
		result1 error
	}
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
}

func (fake *FakeOldLogsRepository) RecentLogsFor(appGuid string) ([]*logmessage.LogMessage, error) {
	fake.recentLogsForMutex.Lock()
	fake.recentLogsForArgsForCall = append(fake.recentLogsForArgsForCall, struct {
		appGuid string
	}{appGuid})
	fake.recentLogsForMutex.Unlock()
	if fake.RecentLogsForStub != nil {
		return fake.RecentLogsForStub(appGuid)
	} else {
		return fake.recentLogsForReturns.result1, fake.recentLogsForReturns.result2
	}
}

func (fake *FakeOldLogsRepository) RecentLogsForCallCount() int {
	fake.recentLogsForMutex.RLock()
	defer fake.recentLogsForMutex.RUnlock()
	return len(fake.recentLogsForArgsForCall)
}

func (fake *FakeOldLogsRepository) RecentLogsForArgsForCall(i int) string {
	fake.recentLogsForMutex.RLock()
	defer fake.recentLogsForMutex.RUnlock()
	return fake.recentLogsForArgsForCall[i].appGuid
}

func (fake *FakeOldLogsRepository) RecentLogsForReturns(result1 []*logmessage.LogMessage, result2 error) {
	fake.RecentLogsForStub = nil
	fake.recentLogsForReturns = struct {
		result1 []*logmessage.LogMessage
		result2 error
	}{result1, result2}
}

func (fake *FakeOldLogsRepository) TailLogsFor(appGuid string, onConnect func(), onMessage func(*logmessage.LogMessage)) error {
	fake.tailLogsForMutex.Lock()
	fake.tailLogsForArgsForCall = append(fake.tailLogsForArgsForCall, struct {
		appGuid   string
		onConnect func()
		onMessage func(*logmessage.LogMessage)
	}{appGuid, onConnect, onMessage})
	fake.tailLogsForMutex.Unlock()
	if fake.TailLogsForStub != nil {
		return fake.TailLogsForStub(appGuid, onConnect, onMessage)
	} else {
		return fake.tailLogsForReturns.result1
	}
}

func (fake *FakeOldLogsRepository) TailLogsForCallCount() int {
	fake.tailLogsForMutex.RLock()
	defer fake.tailLogsForMutex.RUnlock()
	return len(fake.tailLogsForArgsForCall)
}

func (fake *FakeOldLogsRepository) TailLogsForArgsForCall(i int) (string, func(), func(*logmessage.LogMessage)) {
	fake.tailLogsForMutex.RLock()
	defer fake.tailLogsForMutex.RUnlock()
	return fake.tailLogsForArgsForCall[i].appGuid, fake.tailLogsForArgsForCall[i].onConnect, fake.tailLogsForArgsForCall[i].onMessage
}

func (fake *FakeOldLogsRepository) TailLogsForReturns(result1 error) {
	fake.TailLogsForStub = nil
	fake.tailLogsForReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeOldLogsRepository) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *FakeOldLogsRepository) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

var _ api.OldLogsRepository = new(FakeOldLogsRepository)
