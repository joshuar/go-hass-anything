// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package config

import (
	"sync"
)

// Ensure, that AgentMock does implement Agent.
// If this is not the case, regenerate this file with moq.
var _ Agent = &AgentMock{}

// AgentMock is a mock implementation of Agent.
//
//	func TestSomethingThatUsesAgent(t *testing.T) {
//
//		// make and configure a mocked Agent
//		mockedAgent := &AgentMock{
//			DeleteFunc: func(s string) error {
//				panic("mock out the Delete method")
//			},
//			GetFunc: func(s string, ifaceVal interface{}) error {
//				panic("mock out the Get method")
//			},
//			SetFunc: func(s string, ifaceVal interface{}) error {
//				panic("mock out the Set method")
//			},
//		}
//
//		// use mockedAgent in code that requires Agent
//		// and then make assertions.
//
//	}
type AgentMock struct {
	// DeleteFunc mocks the Delete method.
	DeleteFunc func(s string) error

	// GetFunc mocks the Get method.
	GetFunc func(s string, ifaceVal interface{}) error

	// SetFunc mocks the Set method.
	SetFunc func(s string, ifaceVal interface{}) error

	// calls tracks calls to the methods.
	calls struct {
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// S is the s argument value.
			S string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// S is the s argument value.
			S string
			// IfaceVal is the ifaceVal argument value.
			IfaceVal interface{}
		}
		// Set holds details about calls to the Set method.
		Set []struct {
			// S is the s argument value.
			S string
			// IfaceVal is the ifaceVal argument value.
			IfaceVal interface{}
		}
	}
	lockDelete sync.RWMutex
	lockGet    sync.RWMutex
	lockSet    sync.RWMutex
}

// Delete calls DeleteFunc.
func (mock *AgentMock) Delete(s string) error {
	if mock.DeleteFunc == nil {
		panic("AgentMock.DeleteFunc: method is nil but Agent.Delete was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(s)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedAgent.DeleteCalls())
func (mock *AgentMock) DeleteCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *AgentMock) Get(s string, ifaceVal interface{}) error {
	if mock.GetFunc == nil {
		panic("AgentMock.GetFunc: method is nil but Agent.Get was just called")
	}
	callInfo := struct {
		S        string
		IfaceVal interface{}
	}{
		S:        s,
		IfaceVal: ifaceVal,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(s, ifaceVal)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedAgent.GetCalls())
func (mock *AgentMock) GetCalls() []struct {
	S        string
	IfaceVal interface{}
} {
	var calls []struct {
		S        string
		IfaceVal interface{}
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Set calls SetFunc.
func (mock *AgentMock) Set(s string, ifaceVal interface{}) error {
	if mock.SetFunc == nil {
		panic("AgentMock.SetFunc: method is nil but Agent.Set was just called")
	}
	callInfo := struct {
		S        string
		IfaceVal interface{}
	}{
		S:        s,
		IfaceVal: ifaceVal,
	}
	mock.lockSet.Lock()
	mock.calls.Set = append(mock.calls.Set, callInfo)
	mock.lockSet.Unlock()
	return mock.SetFunc(s, ifaceVal)
}

// SetCalls gets all the calls that were made to Set.
// Check the length with:
//
//	len(mockedAgent.SetCalls())
func (mock *AgentMock) SetCalls() []struct {
	S        string
	IfaceVal interface{}
} {
	var calls []struct {
		S        string
		IfaceVal interface{}
	}
	mock.lockSet.RLock()
	calls = mock.calls.Set
	mock.lockSet.RUnlock()
	return calls
}
