// This file was generated by counterfeiter
package grootfakes

import (
	"sync"

	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/lager"
)

type FakeGarbageCollector struct {
	MarkUnusedStub        func(logger lager.Logger, keepBaseImages []string) error
	markUnusedMutex       sync.RWMutex
	markUnusedArgsForCall []struct {
		logger         lager.Logger
		keepBaseImages []string
	}
	markUnusedReturns struct {
		result1 error
	}
	markUnusedReturnsOnCall map[int]struct {
		result1 error
	}
	CollectStub        func(logger lager.Logger) error
	collectMutex       sync.RWMutex
	collectArgsForCall []struct {
		logger lager.Logger
	}
	collectReturns struct {
		result1 error
	}
	collectReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGarbageCollector) MarkUnused(logger lager.Logger, keepBaseImages []string) error {
	var keepBaseImagesCopy []string
	if keepBaseImages != nil {
		keepBaseImagesCopy = make([]string, len(keepBaseImages))
		copy(keepBaseImagesCopy, keepBaseImages)
	}
	fake.markUnusedMutex.Lock()
	ret, specificReturn := fake.markUnusedReturnsOnCall[len(fake.markUnusedArgsForCall)]
	fake.markUnusedArgsForCall = append(fake.markUnusedArgsForCall, struct {
		logger         lager.Logger
		keepBaseImages []string
	}{logger, keepBaseImagesCopy})
	fake.recordInvocation("MarkUnused", []interface{}{logger, keepBaseImagesCopy})
	fake.markUnusedMutex.Unlock()
	if fake.MarkUnusedStub != nil {
		return fake.MarkUnusedStub(logger, keepBaseImages)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.markUnusedReturns.result1
}

func (fake *FakeGarbageCollector) MarkUnusedCallCount() int {
	fake.markUnusedMutex.RLock()
	defer fake.markUnusedMutex.RUnlock()
	return len(fake.markUnusedArgsForCall)
}

func (fake *FakeGarbageCollector) MarkUnusedArgsForCall(i int) (lager.Logger, []string) {
	fake.markUnusedMutex.RLock()
	defer fake.markUnusedMutex.RUnlock()
	return fake.markUnusedArgsForCall[i].logger, fake.markUnusedArgsForCall[i].keepBaseImages
}

func (fake *FakeGarbageCollector) MarkUnusedReturns(result1 error) {
	fake.MarkUnusedStub = nil
	fake.markUnusedReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGarbageCollector) MarkUnusedReturnsOnCall(i int, result1 error) {
	fake.MarkUnusedStub = nil
	if fake.markUnusedReturnsOnCall == nil {
		fake.markUnusedReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.markUnusedReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGarbageCollector) Collect(logger lager.Logger) error {
	fake.collectMutex.Lock()
	ret, specificReturn := fake.collectReturnsOnCall[len(fake.collectArgsForCall)]
	fake.collectArgsForCall = append(fake.collectArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.recordInvocation("Collect", []interface{}{logger})
	fake.collectMutex.Unlock()
	if fake.CollectStub != nil {
		return fake.CollectStub(logger)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.collectReturns.result1
}

func (fake *FakeGarbageCollector) CollectCallCount() int {
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	return len(fake.collectArgsForCall)
}

func (fake *FakeGarbageCollector) CollectArgsForCall(i int) lager.Logger {
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	return fake.collectArgsForCall[i].logger
}

func (fake *FakeGarbageCollector) CollectReturns(result1 error) {
	fake.CollectStub = nil
	fake.collectReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGarbageCollector) CollectReturnsOnCall(i int, result1 error) {
	fake.CollectStub = nil
	if fake.collectReturnsOnCall == nil {
		fake.collectReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.collectReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGarbageCollector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.markUnusedMutex.RLock()
	defer fake.markUnusedMutex.RUnlock()
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeGarbageCollector) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ groot.GarbageCollector = new(FakeGarbageCollector)
