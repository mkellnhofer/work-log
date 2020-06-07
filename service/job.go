package service

import (
	"context"
	"fmt"
	"time"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
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

type jobFunc func(context.Context) *e.Error

func scheduleJob(jobName string, f jobFunc, interval time.Duration) {
	go func() {
		for true {
			log.Infof("Starting %s ...", jobName)
			jErr := f(context.Background())
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
