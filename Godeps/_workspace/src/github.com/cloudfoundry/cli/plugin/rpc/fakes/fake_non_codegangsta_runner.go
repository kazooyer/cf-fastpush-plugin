// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/cf/command_registry"
	"github.com/emirozer/cf-fastpush-plugin/Godeps/_workspace/src/github.com/cloudfoundry/cli/plugin/rpc"
)

type FakeNonCodegangstaRunner struct {
	CommandStub        func([]string, command_registry.Dependency, bool) error
	commandMutex       sync.RWMutex
	commandArgsForCall []struct {
		arg1 []string
		arg2 command_registry.Dependency
		arg3 bool
	}
	commandReturns struct {
		result1 error
	}
}

func (fake *FakeNonCodegangstaRunner) Command(arg1 []string, arg2 command_registry.Dependency, arg3 bool) error {
	fake.commandMutex.Lock()
	fake.commandArgsForCall = append(fake.commandArgsForCall, struct {
		arg1 []string
		arg2 command_registry.Dependency
		arg3 bool
	}{arg1, arg2, arg3})
	fake.commandMutex.Unlock()
	if fake.CommandStub != nil {
		return fake.CommandStub(arg1, arg2, arg3)
	} else {
		return fake.commandReturns.result1
	}
}

func (fake *FakeNonCodegangstaRunner) CommandCallCount() int {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	return len(fake.commandArgsForCall)
}

func (fake *FakeNonCodegangstaRunner) CommandArgsForCall(i int) ([]string, command_registry.Dependency, bool) {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	return fake.commandArgsForCall[i].arg1, fake.commandArgsForCall[i].arg2, fake.commandArgsForCall[i].arg3
}

func (fake *FakeNonCodegangstaRunner) CommandReturns(result1 error) {
	fake.CommandStub = nil
	fake.commandReturns = struct {
		result1 error
	}{result1}
}

var _ rpc.NonCodegangstaRunner = new(FakeNonCodegangstaRunner)
