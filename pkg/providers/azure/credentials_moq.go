// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package azure

import (
	"context"
	v1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	"sync"
)

// Ensure, that CredentialManagerMock does implement CredentialManager.
// If this is not the case, regenerate this file with moq.
var _ CredentialManager = &CredentialManagerMock{}

// CredentialManagerMock is a mock implementation of CredentialManager.
//
// 	func TestSomethingThatUsesCredentialManager(t *testing.T) {
//
// 		// make and configure a mocked CredentialManager
// 		mockedCredentialManager := &CredentialManagerMock{
// 			ReconcileBucketOwnerCredentialsFunc: func(ctx context.Context, name string, ns string) (*Credentials, *v1.CredentialsRequest, error) {
// 				panic("mock out the ReconcileBucketOwnerCredentials method")
// 			},
// 			ReconcileCredentialsFunc: func(ctx context.Context, name string, ns string) (*v1.CredentialsRequest, *Credentials, error) {
// 				panic("mock out the ReconcileCredentials method")
// 			},
// 			ReconcileProviderCredentialsFunc: func(ctx context.Context, ns string) (*Credentials, error) {
// 				panic("mock out the ReconcileProviderCredentials method")
// 			},
// 			ReconcileSESCredentialsFunc: func(ctx context.Context, name string, ns string) (*Credentials, error) {
// 				panic("mock out the ReconcileSESCredentials method")
// 			},
// 		}
//
// 		// use mockedCredentialManager in code that requires CredentialManager
// 		// and then make assertions.
//
// 	}
type CredentialManagerMock struct {
	// ReconcileBucketOwnerCredentialsFunc mocks the ReconcileBucketOwnerCredentials method.
	ReconcileBucketOwnerCredentialsFunc func(ctx context.Context, name string, ns string) (*Credentials, *v1.CredentialsRequest, error)

	// ReconcileCredentialsFunc mocks the ReconcileCredentials method.
	ReconcileCredentialsFunc func(ctx context.Context, name string, ns string) (*v1.CredentialsRequest, *Credentials, error)

	// ReconcileProviderCredentialsFunc mocks the ReconcileProviderCredentials method.
	ReconcileProviderCredentialsFunc func(ctx context.Context, ns string) (*Credentials, error)

	// ReconcileSESCredentialsFunc mocks the ReconcileSESCredentials method.
	ReconcileSESCredentialsFunc func(ctx context.Context, name string, ns string) (*Credentials, error)

	// calls tracks calls to the methods.
	calls struct {
		// ReconcileBucketOwnerCredentials holds details about calls to the ReconcileBucketOwnerCredentials method.
		ReconcileBucketOwnerCredentials []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Ns is the ns argument value.
			Ns string
		}
		// ReconcileCredentials holds details about calls to the ReconcileCredentials method.
		ReconcileCredentials []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Ns is the ns argument value.
			Ns string
		}
		// ReconcileProviderCredentials holds details about calls to the ReconcileProviderCredentials method.
		ReconcileProviderCredentials []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Ns is the ns argument value.
			Ns string
		}
		// ReconcileSESCredentials holds details about calls to the ReconcileSESCredentials method.
		ReconcileSESCredentials []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Ns is the ns argument value.
			Ns string
		}
	}
	lockReconcileBucketOwnerCredentials sync.RWMutex
	lockReconcileCredentials            sync.RWMutex
	lockReconcileProviderCredentials    sync.RWMutex
	lockReconcileSESCredentials         sync.RWMutex
}

// ReconcileBucketOwnerCredentials calls ReconcileBucketOwnerCredentialsFunc.
func (mock *CredentialManagerMock) ReconcileBucketOwnerCredentials(ctx context.Context, name string, ns string) (*Credentials, *v1.CredentialsRequest, error) {
	if mock.ReconcileBucketOwnerCredentialsFunc == nil {
		panic("CredentialManagerMock.ReconcileBucketOwnerCredentialsFunc: method is nil but CredentialManager.ReconcileBucketOwnerCredentials was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
		Ns   string
	}{
		Ctx:  ctx,
		Name: name,
		Ns:   ns,
	}
	mock.lockReconcileBucketOwnerCredentials.Lock()
	mock.calls.ReconcileBucketOwnerCredentials = append(mock.calls.ReconcileBucketOwnerCredentials, callInfo)
	mock.lockReconcileBucketOwnerCredentials.Unlock()
	return mock.ReconcileBucketOwnerCredentialsFunc(ctx, name, ns)
}

// ReconcileBucketOwnerCredentialsCalls gets all the calls that were made to ReconcileBucketOwnerCredentials.
// Check the length with:
//     len(mockedCredentialManager.ReconcileBucketOwnerCredentialsCalls())
func (mock *CredentialManagerMock) ReconcileBucketOwnerCredentialsCalls() []struct {
	Ctx  context.Context
	Name string
	Ns   string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
		Ns   string
	}
	mock.lockReconcileBucketOwnerCredentials.RLock()
	calls = mock.calls.ReconcileBucketOwnerCredentials
	mock.lockReconcileBucketOwnerCredentials.RUnlock()
	return calls
}

// ReconcileCredentials calls ReconcileCredentialsFunc.
func (mock *CredentialManagerMock) ReconcileCredentials(ctx context.Context, name string, ns string) (*v1.CredentialsRequest, *Credentials, error) {
	if mock.ReconcileCredentialsFunc == nil {
		panic("CredentialManagerMock.ReconcileCredentialsFunc: method is nil but CredentialManager.ReconcileCredentials was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
		Ns   string
	}{
		Ctx:  ctx,
		Name: name,
		Ns:   ns,
	}
	mock.lockReconcileCredentials.Lock()
	mock.calls.ReconcileCredentials = append(mock.calls.ReconcileCredentials, callInfo)
	mock.lockReconcileCredentials.Unlock()
	return mock.ReconcileCredentialsFunc(ctx, name, ns)
}

// ReconcileCredentialsCalls gets all the calls that were made to ReconcileCredentials.
// Check the length with:
//     len(mockedCredentialManager.ReconcileCredentialsCalls())
func (mock *CredentialManagerMock) ReconcileCredentialsCalls() []struct {
	Ctx  context.Context
	Name string
	Ns   string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
		Ns   string
	}
	mock.lockReconcileCredentials.RLock()
	calls = mock.calls.ReconcileCredentials
	mock.lockReconcileCredentials.RUnlock()
	return calls
}

// ReconcileProviderCredentials calls ReconcileProviderCredentialsFunc.
func (mock *CredentialManagerMock) ReconcileProviderCredentials(ctx context.Context, ns string) (*Credentials, error) {
	if mock.ReconcileProviderCredentialsFunc == nil {
		panic("CredentialManagerMock.ReconcileProviderCredentialsFunc: method is nil but CredentialManager.ReconcileProviderCredentials was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Ns  string
	}{
		Ctx: ctx,
		Ns:  ns,
	}
	mock.lockReconcileProviderCredentials.Lock()
	mock.calls.ReconcileProviderCredentials = append(mock.calls.ReconcileProviderCredentials, callInfo)
	mock.lockReconcileProviderCredentials.Unlock()
	return mock.ReconcileProviderCredentialsFunc(ctx, ns)
}

// ReconcileProviderCredentialsCalls gets all the calls that were made to ReconcileProviderCredentials.
// Check the length with:
//     len(mockedCredentialManager.ReconcileProviderCredentialsCalls())
func (mock *CredentialManagerMock) ReconcileProviderCredentialsCalls() []struct {
	Ctx context.Context
	Ns  string
} {
	var calls []struct {
		Ctx context.Context
		Ns  string
	}
	mock.lockReconcileProviderCredentials.RLock()
	calls = mock.calls.ReconcileProviderCredentials
	mock.lockReconcileProviderCredentials.RUnlock()
	return calls
}

// ReconcileSESCredentials calls ReconcileSESCredentialsFunc.
func (mock *CredentialManagerMock) ReconcileSESCredentials(ctx context.Context, name string, ns string) (*Credentials, error) {
	if mock.ReconcileSESCredentialsFunc == nil {
		panic("CredentialManagerMock.ReconcileSESCredentialsFunc: method is nil but CredentialManager.ReconcileSESCredentials was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
		Ns   string
	}{
		Ctx:  ctx,
		Name: name,
		Ns:   ns,
	}
	mock.lockReconcileSESCredentials.Lock()
	mock.calls.ReconcileSESCredentials = append(mock.calls.ReconcileSESCredentials, callInfo)
	mock.lockReconcileSESCredentials.Unlock()
	return mock.ReconcileSESCredentialsFunc(ctx, name, ns)
}

// ReconcileSESCredentialsCalls gets all the calls that were made to ReconcileSESCredentials.
// Check the length with:
//     len(mockedCredentialManager.ReconcileSESCredentialsCalls())
func (mock *CredentialManagerMock) ReconcileSESCredentialsCalls() []struct {
	Ctx  context.Context
	Name string
	Ns   string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
		Ns   string
	}
	mock.lockReconcileSESCredentials.RLock()
	calls = mock.calls.ReconcileSESCredentials
	mock.lockReconcileSESCredentials.RUnlock()
	return calls
}
