// Code generated by pegomock. DO NOT EDIT.
// Source: github.com/relnod/dotm/pkg/fileutil (interfaces: Visitor)

package mock

import (
	"reflect"
	"time"

	pegomock "github.com/petergtz/pegomock"
)

type MockVisitor struct {
	fail func(message string, callerSkip ...int)
}

func NewMockVisitor() *MockVisitor {
	return &MockVisitor{fail: pegomock.GlobalFailHandler}
}

func (mock *MockVisitor) Visit(_param0 string, _param1 string) error {
	if mock == nil {
		panic("mock must not be nil. Use myMock := NewMockVisitor().")
	}
	params := []pegomock.Param{_param0, _param1}
	result := pegomock.GetGenericMockFrom(mock).Invoke("Visit", params, []reflect.Type{reflect.TypeOf((*error)(nil)).Elem()})
	var ret0 error
	if len(result) != 0 {
		if result[0] != nil {
			ret0 = result[0].(error)
		}
	}
	return ret0
}

func (mock *MockVisitor) VerifyWasCalledOnce() *VerifierVisitor {
	return &VerifierVisitor{
		mock:                   mock,
		invocationCountMatcher: pegomock.Times(1),
	}
}

func (mock *MockVisitor) VerifyWasCalled(invocationCountMatcher pegomock.Matcher) *VerifierVisitor {
	return &VerifierVisitor{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
	}
}

func (mock *MockVisitor) VerifyWasCalledInOrder(invocationCountMatcher pegomock.Matcher, inOrderContext *pegomock.InOrderContext) *VerifierVisitor {
	return &VerifierVisitor{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		inOrderContext:         inOrderContext,
	}
}

func (mock *MockVisitor) VerifyWasCalledEventually(invocationCountMatcher pegomock.Matcher, timeout time.Duration) *VerifierVisitor {
	return &VerifierVisitor{
		mock:                   mock,
		invocationCountMatcher: invocationCountMatcher,
		timeout:                timeout,
	}
}

type VerifierVisitor struct {
	mock                   *MockVisitor
	invocationCountMatcher pegomock.Matcher
	inOrderContext         *pegomock.InOrderContext
	timeout                time.Duration
}

func (verifier *VerifierVisitor) Visit(_param0 string, _param1 string) *Visitor_Visit_OngoingVerification {
	params := []pegomock.Param{_param0, _param1}
	methodInvocations := pegomock.GetGenericMockFrom(verifier.mock).Verify(verifier.inOrderContext, verifier.invocationCountMatcher, "Visit", params, verifier.timeout)
	return &Visitor_Visit_OngoingVerification{mock: verifier.mock, methodInvocations: methodInvocations}
}

type Visitor_Visit_OngoingVerification struct {
	mock              *MockVisitor
	methodInvocations []pegomock.MethodInvocation
}

func (c *Visitor_Visit_OngoingVerification) GetCapturedArguments() (string, string) {
	_param0, _param1 := c.GetAllCapturedArguments()
	return _param0[len(_param0)-1], _param1[len(_param1)-1]
}

func (c *Visitor_Visit_OngoingVerification) GetAllCapturedArguments() (_param0 []string, _param1 []string) {
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
	}
	return
}
