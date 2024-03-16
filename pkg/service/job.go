package service

import (
	"context"
	"fmt"
	"time"

	"kellnhofer.com/work-log/pkg/constant"
	"kellnhofer.com/work-log/pkg/db/tx"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
)

const sessionsCleanUpInterval = 15 * time.Minute

// JobService contains job related logic.
type JobService struct {
	sServ *SessionService
}

// NewJobService create a new job service.
func NewJobService(ss *SessionService) *JobService {
	return &JobService{ss}
}

// --- Job functions ---

// ScheduleJobs schedules jobs.
func (s *JobService) ScheduleJobs() {
	s.scheduleSessionsCleanUpJob()
}

// ScheduleJobs schedules jobs.
func (s *JobService) scheduleSessionsCleanUpJob() {
	scheduleJob("sessions clean up job", s.sServ.DeleteExpiredSessions, sessionsCleanUpInterval)
}

type jobFunc func(context.Context) error

func scheduleJob(jobName string, f jobFunc, interval time.Duration) {
	go func() {
		for {
			log.Infof("Starting %s ...", jobName)
			jErr := f(createJobContext())
			if jErr == nil {
				log.Infof("Finished %s.", jobName)
			} else {
				err := e.WrapError(e.SysJobFailed, fmt.Sprintf("Job '%s' failed.", jobName), jErr)
				log.Error(err.StackTrace())
			}
			time.Sleep(interval)
		}
	}()
}

func createJobContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.ContextKeyTransactionHolder, &tx.TransactionHolder{})
	return context.WithValue(ctx, constant.ContextKeySecurityContext,
		model.GetSystemUserSecurityContext())
}
