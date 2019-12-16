// Code generated by counterfeiter. DO NOT EDIT.
package httpfakes

import (
	"sync"

	"github.com/storyscript/auth"
	"github.com/storyscript/auth/http"
)

type FakeUserRepository struct {
	SaveStub        func(auth.User) (string, error)
	saveMutex       sync.RWMutex
	saveArgsForCall []struct {
		arg1 auth.User
	}
	saveReturns struct {
		result1 string
		result2 error
	}
	saveReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeUserRepository) Save(arg1 auth.User) (string, error) {
	fake.saveMutex.Lock()
	ret, specificReturn := fake.saveReturnsOnCall[len(fake.saveArgsForCall)]
	fake.saveArgsForCall = append(fake.saveArgsForCall, struct {
		arg1 auth.User
	}{arg1})
	fake.recordInvocation("Save", []interface{}{arg1})
	fake.saveMutex.Unlock()
	if fake.SaveStub != nil {
		return fake.SaveStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.saveReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUserRepository) SaveCallCount() int {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	return len(fake.saveArgsForCall)
}

func (fake *FakeUserRepository) SaveCalls(stub func(auth.User) (string, error)) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = stub
}

func (fake *FakeUserRepository) SaveArgsForCall(i int) auth.User {
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	argsForCall := fake.saveArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUserRepository) SaveReturns(result1 string, result2 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	fake.saveReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeUserRepository) SaveReturnsOnCall(i int, result1 string, result2 error) {
	fake.saveMutex.Lock()
	defer fake.saveMutex.Unlock()
	fake.SaveStub = nil
	if fake.saveReturnsOnCall == nil {
		fake.saveReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.saveReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeUserRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.saveMutex.RLock()
	defer fake.saveMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeUserRepository) recordInvocation(key string, args []interface{}) {
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

var _ http.UserRepository = new(FakeUserRepository)
