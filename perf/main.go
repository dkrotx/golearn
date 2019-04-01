package main

import "strings"

type TaskInfo struct {
	Application string
	Zone string
	Host string
	Revision string
	Instance int
}

func (ti *TaskInfo) MatchAny(s string) bool {
	return strings.Index(ti.Application, s) != -1 ||
		strings.Index(ti.Zone, s) != -1 ||
		strings.Index(ti.Host, s) != -1 ||
		strings.Index(ti.Zone, s) != -1 ||
		strings.Index(ti.Revision, s) != -1
}

type ServiceName string

type ServiceRepository struct {
	TaskByServices map[ServiceName][]TaskInfo
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		TaskByServices: make(map[ServiceName][]TaskInfo),
	}
}

func (rep *ServiceRepository) AddTask(svc ServiceName, inf TaskInfo) {
	rep.TaskByServices[svc] = append(rep.TaskByServices[svc], inf)
}

func (rep *ServiceRepository) MatchAny(svc ServiceName, sub string) (res []*TaskInfo) {
	for _, tsk := range rep.TaskByServices[svc] {
		if tsk.MatchAny(sub) {
			res = append(res, &tsk)
		}
	}

	return
}