// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package ui

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
//			AppIDFunc: func() string {
//				panic("mock out the AppID method")
//			},
//			AppNameFunc: func() string {
//				panic("mock out the AppName method")
//			},
//			AppVersionFunc: func() string {
//				panic("mock out the AppVersion method")
//			},
//			StopFunc: func()  {
//				panic("mock out the Stop method")
//			},
//		}
//
//		// use mockedAgent in code that requires Agent
//		// and then make assertions.
//
//	}
type AgentMock struct {
	// AppIDFunc mocks the AppID method.
	AppIDFunc func() string

	// AppNameFunc mocks the AppName method.
	AppNameFunc func() string

	// AppVersionFunc mocks the AppVersion method.
	AppVersionFunc func() string

	// StopFunc mocks the Stop method.
	StopFunc func()

	// calls tracks calls to the methods.
	calls struct {
		// AppID holds details about calls to the AppID method.
		AppID []struct {
		}
		// AppName holds details about calls to the AppName method.
		AppName []struct {
		}
		// AppVersion holds details about calls to the AppVersion method.
		AppVersion []struct {
		}
		// Stop holds details about calls to the Stop method.
		Stop []struct {
		}
	}
	lockAppID      sync.RWMutex
	lockAppName    sync.RWMutex
	lockAppVersion sync.RWMutex
	lockStop       sync.RWMutex
}

// AppID calls AppIDFunc.
func (mock *AgentMock) AppID() string {
	if mock.AppIDFunc == nil {
		panic("AgentMock.AppIDFunc: method is nil but Agent.AppID was just called")
	}
	callInfo := struct {
	}{}
	mock.lockAppID.Lock()
	mock.calls.AppID = append(mock.calls.AppID, callInfo)
	mock.lockAppID.Unlock()
	return mock.AppIDFunc()
}

// AppIDCalls gets all the calls that were made to AppID.
// Check the length with:
//
//	len(mockedAgent.AppIDCalls())
func (mock *AgentMock) AppIDCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAppID.RLock()
	calls = mock.calls.AppID
	mock.lockAppID.RUnlock()
	return calls
}

// AppName calls AppNameFunc.
func (mock *AgentMock) AppName() string {
	if mock.AppNameFunc == nil {
		panic("AgentMock.AppNameFunc: method is nil but Agent.AppName was just called")
	}
	callInfo := struct {
	}{}
	mock.lockAppName.Lock()
	mock.calls.AppName = append(mock.calls.AppName, callInfo)
	mock.lockAppName.Unlock()
	return mock.AppNameFunc()
}

// AppNameCalls gets all the calls that were made to AppName.
// Check the length with:
//
//	len(mockedAgent.AppNameCalls())
func (mock *AgentMock) AppNameCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAppName.RLock()
	calls = mock.calls.AppName
	mock.lockAppName.RUnlock()
	return calls
}

// AppVersion calls AppVersionFunc.
func (mock *AgentMock) AppVersion() string {
	if mock.AppVersionFunc == nil {
		panic("AgentMock.AppVersionFunc: method is nil but Agent.AppVersion was just called")
	}
	callInfo := struct {
	}{}
	mock.lockAppVersion.Lock()
	mock.calls.AppVersion = append(mock.calls.AppVersion, callInfo)
	mock.lockAppVersion.Unlock()
	return mock.AppVersionFunc()
}

// AppVersionCalls gets all the calls that were made to AppVersion.
// Check the length with:
//
//	len(mockedAgent.AppVersionCalls())
func (mock *AgentMock) AppVersionCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAppVersion.RLock()
	calls = mock.calls.AppVersion
	mock.lockAppVersion.RUnlock()
	return calls
}

// Stop calls StopFunc.
func (mock *AgentMock) Stop() {
	if mock.StopFunc == nil {
		panic("AgentMock.StopFunc: method is nil but Agent.Stop was just called")
	}
	callInfo := struct {
	}{}
	mock.lockStop.Lock()
	mock.calls.Stop = append(mock.calls.Stop, callInfo)
	mock.lockStop.Unlock()
	mock.StopFunc()
}

// StopCalls gets all the calls that were made to Stop.
// Check the length with:
//
//	len(mockedAgent.StopCalls())
func (mock *AgentMock) StopCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockStop.RLock()
	calls = mock.calls.Stop
	mock.lockStop.RUnlock()
	return calls
}
