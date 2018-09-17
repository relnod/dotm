// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/relnod/dotm/pkg/dotfiles (interfaces: Action)

package mock

import (
	"reflect"

	pegomock "github.com/petergtz/pegomock"
)

type MockAction struct {
	fail func(message string, callerSkip ...int)
}

func NewMockAction() *MockAction {
	return &MockAction{fail: pegomock.GlobalFailHandler}
}

func (mock *MockAction) Run(_param0 string, _param1 string, _param2 string) error {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockAction().")
	}
	params := []pegomock.Param{_param0, _param1, _param2}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Run", params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(error)
		}
	}
	return ret0
}

func (mock *MockAction) VerifyWasCalledOnce() *VerifierAction {
	return &VerifierAction{mock, pegomock.Times(1), nil}
}

func (mock *MockAction) VerifyWasCalled(invocationCountMatcher pegomock.Matcher) *VerifierAction {
	return &VerifierAction{mock, invocationCountMatcher, nil}
}

func (mock *MockAction) VerifyWasCalledInOrder(invocationCountMatcher pegomock.Matcher, inOrderContext *pegomock.InOrderContext) *VerifierAction {
	return &VerifierAction{mock, invocationCountMatcher, inOrderContext}
}

type VerifierAction struct {
	mock                   *MockAction
	invocationCountMatcher pegomock.Matcher
	inOrderContext         *pegomock.InOrderContext
}

func (verifier *VerifierAction) Run(_param0 string, _param1 string, _param2 string) *Action_Run_OngoingVerification {
	params := []pegomock.Param{_param0, _param1, _param2}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Run", params)
	return &Action_Run_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type Action_Run_OngoingVerification struct {
	mock              *MockAction
	methodInvocations []pegomock.MethodInvocation
}

func (c *Action_Run_OngoingVerification) GetCapturedArguments() (string, string, string) {
	_param0, _param1, _param2 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1], _param1[len(_param1)-1], _param2[len(_param2)-1]
}

func (c *Action_Run_OngoingVerification) GetAllCapturedArguments() (_param0 []string, _param1 []string, _param2 []string) {
	params := pegomock.GetGenericMockFrom(c.mock).GetInvocationParams(c.methodInvocations)
	if len(params) > 0 {
		_param0 = make([]string, len(params[0]))
		for u, param := range params[0] {
			_param0[u] = param.(string)
		}
		_param1 = make([]string, len(params[1]))
		for u, param := range params[1] {
			_param1[u] = param.(string)
		}
		_param2 = make([]string, len(params[2]))
		for u, param := range params[2] {
			_param2[u] = param.(string)
		}
	}
	return
}
