// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"

	"github.com/shurcooL/githubv4"
	"github.com/skarlso/caretaker/pkg/client"
)

type FakeClient struct {
	AddLabelStub        func(context.Context, string, githubv4.String) error
	addLabelMutex       sync.RWMutex
	addLabelArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 githubv4.String
	}
	addLabelReturns struct {
		result1 error
	}
	addLabelReturnsOnCall map[int]struct {
		result1 error
	}
	AssignIssueToProjectStub        func(context.Context, int, int) error
	assignIssueToProjectMutex       sync.RWMutex
	assignIssueToProjectArgsForCall []struct {
		arg1 context.Context
		arg2 int
		arg3 int
	}
	assignIssueToProjectReturns struct {
		result1 error
	}
	assignIssueToProjectReturnsOnCall map[int]struct {
		result1 error
	}
	LeaveCommentStub        func(context.Context, githubv4.String, string) error
	leaveCommentMutex       sync.RWMutex
	leaveCommentArgsForCall []struct {
		arg1 context.Context
		arg2 githubv4.String
		arg3 string
	}
	leaveCommentReturns struct {
		result1 error
	}
	leaveCommentReturnsOnCall map[int]struct {
		result1 error
	}
	PullRequestStub        func(context.Context, int) (client.PullRequest, error)
	pullRequestMutex       sync.RWMutex
	pullRequestArgsForCall []struct {
		arg1 context.Context
		arg2 int
	}
	pullRequestReturns struct {
		result1 client.PullRequest
		result2 error
	}
	pullRequestReturnsOnCall map[int]struct {
		result1 client.PullRequest
		result2 error
	}
	PullRequestsStub        func(context.Context) ([]client.PullRequest, error)
	pullRequestsMutex       sync.RWMutex
	pullRequestsArgsForCall []struct {
		arg1 context.Context
	}
	pullRequestsReturns struct {
		result1 []client.PullRequest
		result2 error
	}
	pullRequestsReturnsOnCall map[int]struct {
		result1 []client.PullRequest
		result2 error
	}
	RemoveLabelStub        func(context.Context, string, githubv4.String) error
	removeLabelMutex       sync.RWMutex
	removeLabelArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 githubv4.String
	}
	removeLabelReturns struct {
		result1 error
	}
	removeLabelReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateIssueStatusStub        func(context.Context, client.Issue) error
	updateIssueStatusMutex       sync.RWMutex
	updateIssueStatusArgsForCall []struct {
		arg1 context.Context
		arg2 client.Issue
	}
	updateIssueStatusReturns struct {
		result1 error
	}
	updateIssueStatusReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) AddLabel(arg1 context.Context, arg2 string, arg3 githubv4.String) error {
	fake.addLabelMutex.Lock()
	ret, specificReturn := fake.addLabelReturnsOnCall[len(fake.addLabelArgsForCall)]
	fake.addLabelArgsForCall = append(fake.addLabelArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 githubv4.String
	}{arg1, arg2, arg3})
	stub := fake.AddLabelStub
	fakeReturns := fake.addLabelReturns
	fake.recordInvocation("AddLabel", []interface{}{arg1, arg2, arg3})
	fake.addLabelMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) AddLabelCallCount() int {
	fake.addLabelMutex.RLock()
	defer fake.addLabelMutex.RUnlock()
	return len(fake.addLabelArgsForCall)
}

func (fake *FakeClient) AddLabelCalls(stub func(context.Context, string, githubv4.String) error) {
	fake.addLabelMutex.Lock()
	defer fake.addLabelMutex.Unlock()
	fake.AddLabelStub = stub
}

func (fake *FakeClient) AddLabelArgsForCall(i int) (context.Context, string, githubv4.String) {
	fake.addLabelMutex.RLock()
	defer fake.addLabelMutex.RUnlock()
	argsForCall := fake.addLabelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) AddLabelReturns(result1 error) {
	fake.addLabelMutex.Lock()
	defer fake.addLabelMutex.Unlock()
	fake.AddLabelStub = nil
	fake.addLabelReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) AddLabelReturnsOnCall(i int, result1 error) {
	fake.addLabelMutex.Lock()
	defer fake.addLabelMutex.Unlock()
	fake.AddLabelStub = nil
	if fake.addLabelReturnsOnCall == nil {
		fake.addLabelReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addLabelReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) AssignIssueToProject(arg1 context.Context, arg2 int, arg3 int) error {
	fake.assignIssueToProjectMutex.Lock()
	ret, specificReturn := fake.assignIssueToProjectReturnsOnCall[len(fake.assignIssueToProjectArgsForCall)]
	fake.assignIssueToProjectArgsForCall = append(fake.assignIssueToProjectArgsForCall, struct {
		arg1 context.Context
		arg2 int
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.AssignIssueToProjectStub
	fakeReturns := fake.assignIssueToProjectReturns
	fake.recordInvocation("AssignIssueToProject", []interface{}{arg1, arg2, arg3})
	fake.assignIssueToProjectMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) AssignIssueToProjectCallCount() int {
	fake.assignIssueToProjectMutex.RLock()
	defer fake.assignIssueToProjectMutex.RUnlock()
	return len(fake.assignIssueToProjectArgsForCall)
}

func (fake *FakeClient) AssignIssueToProjectCalls(stub func(context.Context, int, int) error) {
	fake.assignIssueToProjectMutex.Lock()
	defer fake.assignIssueToProjectMutex.Unlock()
	fake.AssignIssueToProjectStub = stub
}

func (fake *FakeClient) AssignIssueToProjectArgsForCall(i int) (context.Context, int, int) {
	fake.assignIssueToProjectMutex.RLock()
	defer fake.assignIssueToProjectMutex.RUnlock()
	argsForCall := fake.assignIssueToProjectArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) AssignIssueToProjectReturns(result1 error) {
	fake.assignIssueToProjectMutex.Lock()
	defer fake.assignIssueToProjectMutex.Unlock()
	fake.AssignIssueToProjectStub = nil
	fake.assignIssueToProjectReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) AssignIssueToProjectReturnsOnCall(i int, result1 error) {
	fake.assignIssueToProjectMutex.Lock()
	defer fake.assignIssueToProjectMutex.Unlock()
	fake.AssignIssueToProjectStub = nil
	if fake.assignIssueToProjectReturnsOnCall == nil {
		fake.assignIssueToProjectReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.assignIssueToProjectReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) LeaveComment(arg1 context.Context, arg2 githubv4.String, arg3 string) error {
	fake.leaveCommentMutex.Lock()
	ret, specificReturn := fake.leaveCommentReturnsOnCall[len(fake.leaveCommentArgsForCall)]
	fake.leaveCommentArgsForCall = append(fake.leaveCommentArgsForCall, struct {
		arg1 context.Context
		arg2 githubv4.String
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.LeaveCommentStub
	fakeReturns := fake.leaveCommentReturns
	fake.recordInvocation("LeaveComment", []interface{}{arg1, arg2, arg3})
	fake.leaveCommentMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) LeaveCommentCallCount() int {
	fake.leaveCommentMutex.RLock()
	defer fake.leaveCommentMutex.RUnlock()
	return len(fake.leaveCommentArgsForCall)
}

func (fake *FakeClient) LeaveCommentCalls(stub func(context.Context, githubv4.String, string) error) {
	fake.leaveCommentMutex.Lock()
	defer fake.leaveCommentMutex.Unlock()
	fake.LeaveCommentStub = stub
}

func (fake *FakeClient) LeaveCommentArgsForCall(i int) (context.Context, githubv4.String, string) {
	fake.leaveCommentMutex.RLock()
	defer fake.leaveCommentMutex.RUnlock()
	argsForCall := fake.leaveCommentArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) LeaveCommentReturns(result1 error) {
	fake.leaveCommentMutex.Lock()
	defer fake.leaveCommentMutex.Unlock()
	fake.LeaveCommentStub = nil
	fake.leaveCommentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) LeaveCommentReturnsOnCall(i int, result1 error) {
	fake.leaveCommentMutex.Lock()
	defer fake.leaveCommentMutex.Unlock()
	fake.LeaveCommentStub = nil
	if fake.leaveCommentReturnsOnCall == nil {
		fake.leaveCommentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.leaveCommentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) PullRequest(arg1 context.Context, arg2 int) (client.PullRequest, error) {
	fake.pullRequestMutex.Lock()
	ret, specificReturn := fake.pullRequestReturnsOnCall[len(fake.pullRequestArgsForCall)]
	fake.pullRequestArgsForCall = append(fake.pullRequestArgsForCall, struct {
		arg1 context.Context
		arg2 int
	}{arg1, arg2})
	stub := fake.PullRequestStub
	fakeReturns := fake.pullRequestReturns
	fake.recordInvocation("PullRequest", []interface{}{arg1, arg2})
	fake.pullRequestMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) PullRequestCallCount() int {
	fake.pullRequestMutex.RLock()
	defer fake.pullRequestMutex.RUnlock()
	return len(fake.pullRequestArgsForCall)
}

func (fake *FakeClient) PullRequestCalls(stub func(context.Context, int) (client.PullRequest, error)) {
	fake.pullRequestMutex.Lock()
	defer fake.pullRequestMutex.Unlock()
	fake.PullRequestStub = stub
}

func (fake *FakeClient) PullRequestArgsForCall(i int) (context.Context, int) {
	fake.pullRequestMutex.RLock()
	defer fake.pullRequestMutex.RUnlock()
	argsForCall := fake.pullRequestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) PullRequestReturns(result1 client.PullRequest, result2 error) {
	fake.pullRequestMutex.Lock()
	defer fake.pullRequestMutex.Unlock()
	fake.PullRequestStub = nil
	fake.pullRequestReturns = struct {
		result1 client.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PullRequestReturnsOnCall(i int, result1 client.PullRequest, result2 error) {
	fake.pullRequestMutex.Lock()
	defer fake.pullRequestMutex.Unlock()
	fake.PullRequestStub = nil
	if fake.pullRequestReturnsOnCall == nil {
		fake.pullRequestReturnsOnCall = make(map[int]struct {
			result1 client.PullRequest
			result2 error
		})
	}
	fake.pullRequestReturnsOnCall[i] = struct {
		result1 client.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PullRequests(arg1 context.Context) ([]client.PullRequest, error) {
	fake.pullRequestsMutex.Lock()
	ret, specificReturn := fake.pullRequestsReturnsOnCall[len(fake.pullRequestsArgsForCall)]
	fake.pullRequestsArgsForCall = append(fake.pullRequestsArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.PullRequestsStub
	fakeReturns := fake.pullRequestsReturns
	fake.recordInvocation("PullRequests", []interface{}{arg1})
	fake.pullRequestsMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) PullRequestsCallCount() int {
	fake.pullRequestsMutex.RLock()
	defer fake.pullRequestsMutex.RUnlock()
	return len(fake.pullRequestsArgsForCall)
}

func (fake *FakeClient) PullRequestsCalls(stub func(context.Context) ([]client.PullRequest, error)) {
	fake.pullRequestsMutex.Lock()
	defer fake.pullRequestsMutex.Unlock()
	fake.PullRequestsStub = stub
}

func (fake *FakeClient) PullRequestsArgsForCall(i int) context.Context {
	fake.pullRequestsMutex.RLock()
	defer fake.pullRequestsMutex.RUnlock()
	argsForCall := fake.pullRequestsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) PullRequestsReturns(result1 []client.PullRequest, result2 error) {
	fake.pullRequestsMutex.Lock()
	defer fake.pullRequestsMutex.Unlock()
	fake.PullRequestsStub = nil
	fake.pullRequestsReturns = struct {
		result1 []client.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) PullRequestsReturnsOnCall(i int, result1 []client.PullRequest, result2 error) {
	fake.pullRequestsMutex.Lock()
	defer fake.pullRequestsMutex.Unlock()
	fake.PullRequestsStub = nil
	if fake.pullRequestsReturnsOnCall == nil {
		fake.pullRequestsReturnsOnCall = make(map[int]struct {
			result1 []client.PullRequest
			result2 error
		})
	}
	fake.pullRequestsReturnsOnCall[i] = struct {
		result1 []client.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) RemoveLabel(arg1 context.Context, arg2 string, arg3 githubv4.String) error {
	fake.removeLabelMutex.Lock()
	ret, specificReturn := fake.removeLabelReturnsOnCall[len(fake.removeLabelArgsForCall)]
	fake.removeLabelArgsForCall = append(fake.removeLabelArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 githubv4.String
	}{arg1, arg2, arg3})
	stub := fake.RemoveLabelStub
	fakeReturns := fake.removeLabelReturns
	fake.recordInvocation("RemoveLabel", []interface{}{arg1, arg2, arg3})
	fake.removeLabelMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) RemoveLabelCallCount() int {
	fake.removeLabelMutex.RLock()
	defer fake.removeLabelMutex.RUnlock()
	return len(fake.removeLabelArgsForCall)
}

func (fake *FakeClient) RemoveLabelCalls(stub func(context.Context, string, githubv4.String) error) {
	fake.removeLabelMutex.Lock()
	defer fake.removeLabelMutex.Unlock()
	fake.RemoveLabelStub = stub
}

func (fake *FakeClient) RemoveLabelArgsForCall(i int) (context.Context, string, githubv4.String) {
	fake.removeLabelMutex.RLock()
	defer fake.removeLabelMutex.RUnlock()
	argsForCall := fake.removeLabelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) RemoveLabelReturns(result1 error) {
	fake.removeLabelMutex.Lock()
	defer fake.removeLabelMutex.Unlock()
	fake.RemoveLabelStub = nil
	fake.removeLabelReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) RemoveLabelReturnsOnCall(i int, result1 error) {
	fake.removeLabelMutex.Lock()
	defer fake.removeLabelMutex.Unlock()
	fake.RemoveLabelStub = nil
	if fake.removeLabelReturnsOnCall == nil {
		fake.removeLabelReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.removeLabelReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) UpdateIssueStatus(arg1 context.Context, arg2 client.Issue) error {
	fake.updateIssueStatusMutex.Lock()
	ret, specificReturn := fake.updateIssueStatusReturnsOnCall[len(fake.updateIssueStatusArgsForCall)]
	fake.updateIssueStatusArgsForCall = append(fake.updateIssueStatusArgsForCall, struct {
		arg1 context.Context
		arg2 client.Issue
	}{arg1, arg2})
	stub := fake.UpdateIssueStatusStub
	fakeReturns := fake.updateIssueStatusReturns
	fake.recordInvocation("UpdateIssueStatus", []interface{}{arg1, arg2})
	fake.updateIssueStatusMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) UpdateIssueStatusCallCount() int {
	fake.updateIssueStatusMutex.RLock()
	defer fake.updateIssueStatusMutex.RUnlock()
	return len(fake.updateIssueStatusArgsForCall)
}

func (fake *FakeClient) UpdateIssueStatusCalls(stub func(context.Context, client.Issue) error) {
	fake.updateIssueStatusMutex.Lock()
	defer fake.updateIssueStatusMutex.Unlock()
	fake.UpdateIssueStatusStub = stub
}

func (fake *FakeClient) UpdateIssueStatusArgsForCall(i int) (context.Context, client.Issue) {
	fake.updateIssueStatusMutex.RLock()
	defer fake.updateIssueStatusMutex.RUnlock()
	argsForCall := fake.updateIssueStatusArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) UpdateIssueStatusReturns(result1 error) {
	fake.updateIssueStatusMutex.Lock()
	defer fake.updateIssueStatusMutex.Unlock()
	fake.UpdateIssueStatusStub = nil
	fake.updateIssueStatusReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) UpdateIssueStatusReturnsOnCall(i int, result1 error) {
	fake.updateIssueStatusMutex.Lock()
	defer fake.updateIssueStatusMutex.Unlock()
	fake.UpdateIssueStatusStub = nil
	if fake.updateIssueStatusReturnsOnCall == nil {
		fake.updateIssueStatusReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateIssueStatusReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addLabelMutex.RLock()
	defer fake.addLabelMutex.RUnlock()
	fake.assignIssueToProjectMutex.RLock()
	defer fake.assignIssueToProjectMutex.RUnlock()
	fake.leaveCommentMutex.RLock()
	defer fake.leaveCommentMutex.RUnlock()
	fake.pullRequestMutex.RLock()
	defer fake.pullRequestMutex.RUnlock()
	fake.pullRequestsMutex.RLock()
	defer fake.pullRequestsMutex.RUnlock()
	fake.removeLabelMutex.RLock()
	defer fake.removeLabelMutex.RUnlock()
	fake.updateIssueStatusMutex.RLock()
	defer fake.updateIssueStatusMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
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

var _ client.Client = new(FakeClient)
