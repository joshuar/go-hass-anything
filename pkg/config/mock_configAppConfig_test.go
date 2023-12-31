// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package config

import (
	"sync"
)

// Ensure, that AppConfigMock does implement AppConfig.
// If this is not the case, regenerate this file with moq.
var _ AppConfig = &AppConfigMock{}

// AppConfigMock is a mock implementation of AppConfig.
//
//	func TestSomethingThatUsesAppConfig(t *testing.T) {
//
//		// make and configure a mocked AppConfig
//		mockedAppConfig := &AppConfigMock{
//			DeleteFunc: func(s string) error {
//				panic("mock out the Delete method")
//			},
//			GetFunc: func(s string, ifaceVal interface{}) error {
//				panic("mock out the Get method")
//			},
//			IsRegisteredFunc: func() bool {
//				panic("mock out the IsRegistered method")
//			},
//			RegisterFunc: func() error {
//				panic("mock out the Register method")
//			},
//			SetFunc: func(s string, ifaceVal interface{}) error {
//				panic("mock out the Set method")
//			},
//			UnRegisterFunc: func() error {
//				panic("mock out the UnRegister method")
//			},
//		}
//
//		// use mockedAppConfig in code that requires AppConfig
//		// and then make assertions.
//
//	}
type AppConfigMock struct {
	// DeleteFunc mocks the Delete method.
	DeleteFunc func(s string) error

	// GetFunc mocks the Get method.
	GetFunc func(s string, ifaceVal interface{}) error

	// IsRegisteredFunc mocks the IsRegistered method.
	IsRegisteredFunc func() bool

	// RegisterFunc mocks the Register method.
	RegisterFunc func() error

	// SetFunc mocks the Set method.
	SetFunc func(s string, ifaceVal interface{}) error

	// UnRegisterFunc mocks the UnRegister method.
	UnRegisterFunc func() error

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
		// IsRegistered holds details about calls to the IsRegistered method.
		IsRegistered []struct {
		}
		// Register holds details about calls to the Register method.
		Register []struct {
		}
		// Set holds details about calls to the Set method.
		Set []struct {
			// S is the s argument value.
			S string
			// IfaceVal is the ifaceVal argument value.
			IfaceVal interface{}
		}
		// UnRegister holds details about calls to the UnRegister method.
		UnRegister []struct {
		}
	}
	lockDelete       sync.RWMutex
	lockGet          sync.RWMutex
	lockIsRegistered sync.RWMutex
	lockRegister     sync.RWMutex
	lockSet          sync.RWMutex
	lockUnRegister   sync.RWMutex
}

// Delete calls DeleteFunc.
func (mock *AppConfigMock) Delete(s string) error {
	if mock.DeleteFunc == nil {
		panic("AppConfigMock.DeleteFunc: method is nil but AppConfig.Delete was just called")
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
//	len(mockedAppConfig.DeleteCalls())
func (mock *AppConfigMock) DeleteCalls() []struct {
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
func (mock *AppConfigMock) Get(s string, ifaceVal interface{}) error {
	if mock.GetFunc == nil {
		panic("AppConfigMock.GetFunc: method is nil but AppConfig.Get was just called")
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
//	len(mockedAppConfig.GetCalls())
func (mock *AppConfigMock) GetCalls() []struct {
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

// IsRegistered calls IsRegisteredFunc.
func (mock *AppConfigMock) IsRegistered() bool {
	if mock.IsRegisteredFunc == nil {
		panic("AppConfigMock.IsRegisteredFunc: method is nil but AppConfig.IsRegistered was just called")
	}
	callInfo := struct {
	}{}
	mock.lockIsRegistered.Lock()
	mock.calls.IsRegistered = append(mock.calls.IsRegistered, callInfo)
	mock.lockIsRegistered.Unlock()
	return mock.IsRegisteredFunc()
}

// IsRegisteredCalls gets all the calls that were made to IsRegistered.
// Check the length with:
//
//	len(mockedAppConfig.IsRegisteredCalls())
func (mock *AppConfigMock) IsRegisteredCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockIsRegistered.RLock()
	calls = mock.calls.IsRegistered
	mock.lockIsRegistered.RUnlock()
	return calls
}

// Register calls RegisterFunc.
func (mock *AppConfigMock) Register() error {
	if mock.RegisterFunc == nil {
		panic("AppConfigMock.RegisterFunc: method is nil but AppConfig.Register was just called")
	}
	callInfo := struct {
	}{}
	mock.lockRegister.Lock()
	mock.calls.Register = append(mock.calls.Register, callInfo)
	mock.lockRegister.Unlock()
	return mock.RegisterFunc()
}

// RegisterCalls gets all the calls that were made to Register.
// Check the length with:
//
//	len(mockedAppConfig.RegisterCalls())
func (mock *AppConfigMock) RegisterCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockRegister.RLock()
	calls = mock.calls.Register
	mock.lockRegister.RUnlock()
	return calls
}

// Set calls SetFunc.
func (mock *AppConfigMock) Set(s string, ifaceVal interface{}) error {
	if mock.SetFunc == nil {
		panic("AppConfigMock.SetFunc: method is nil but AppConfig.Set was just called")
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
//	len(mockedAppConfig.SetCalls())
func (mock *AppConfigMock) SetCalls() []struct {
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

// UnRegister calls UnRegisterFunc.
func (mock *AppConfigMock) UnRegister() error {
	if mock.UnRegisterFunc == nil {
		panic("AppConfigMock.UnRegisterFunc: method is nil but AppConfig.UnRegister was just called")
	}
	callInfo := struct {
	}{}
	mock.lockUnRegister.Lock()
	mock.calls.UnRegister = append(mock.calls.UnRegister, callInfo)
	mock.lockUnRegister.Unlock()
	return mock.UnRegisterFunc()
}

// UnRegisterCalls gets all the calls that were made to UnRegister.
// Check the length with:
//
//	len(mockedAppConfig.UnRegisterCalls())
func (mock *AppConfigMock) UnRegisterCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockUnRegister.RLock()
	calls = mock.calls.UnRegister
	mock.lockUnRegister.RUnlock()
	return calls
}
