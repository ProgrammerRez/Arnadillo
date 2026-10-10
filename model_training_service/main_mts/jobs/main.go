package jobs

// package main

import (
	// "encoding/json"
	"errors"
	"fmt"
	log "mts/logging"
)

// 1. This file is suppossed to prepare and sanitize user requests
// 2. Make sure all available models fit the user requirements.
// 3. Translate it into a query.
// 4. Just Making sure user understands everything

// Let's start by initializing the Job Queue
func InitializeQueue() *JobQeueu{
	log.Info("Initializing Job Queue")
	return &JobQeueu{
		LowPriority: []*Job{},
		MehPriority: []*Job{},
		HighPrioroty: []*Job{},
	}
}

// Then by Adding Jobs to the Queue
func (queue *JobQeueu) CreateNewJob(
	data_params []DataParameters,
	training_params TrainingParameters,
	eval_params EvaluationParameters,
	do_over_params DoOverParameters,
	urgency int,
	) error {

		// 	Locking Mutex to avoid race conditions
		queue.Mu.Lock()
		defer queue.Mu.Unlock()

		log.Info("Creating New Job..... ")

		job := &Job{
			DataParams: data_params,
			TrainingParams: training_params,
			EvalParams: eval_params,
			DoOverParams: do_over_params,
			Status: JobQeued,
		}
		
		// Increment for next job addition
		queue.NextID++

		// Set the id for the job (should be unique)
		job.ID = queue.NextID

		// Urgency Switch Logic
		switch urgency {
		case 0:
			queue.LowPriority = append(queue.LowPriority, job)
			log.Info("Queueing Job at Low Priority")
		case 1:
			queue.MehPriority = append(queue.MehPriority, job)
			log.Info("Queueing Job at Meh Priority")
		case 2:
			queue.HighPrioroty = append(queue.HighPrioroty, job)
			log.Info("Queueing Job at High Priority")
		default:
			log.Error("Unable to Assign Severity Level. Falling back to Meh Priority")
			queue.MehPriority = append(queue.MehPriority, job)
			log.Warn("Falling Back to Queueing Job at Meh Priority. No proper level provided")
			return errors.New("Unable to Assign Severity Level. Falling back to Meh Priority")			
		}

		return nil
	}

// This function returns the job list based on set filter level
func (q *JobQeueu) GetJobs(filter int32) ([]*Job, error){
	q.Mu.Lock()
	defer q.Mu.Unlock()
	
	log.Info(fmt.Sprintf("Getting Jobs with filter: %d", filter))
	
	list := []*Job{}

	switch filter{
	case 0:
		log.Info(fmt.Sprintf("List %v", q.LowPriority))
		return q.LowPriority, nil
	case 1:
		log.Info(fmt.Sprintf("List %v", q.MehPriority))
		return q.MehPriority, nil
	case 2:
		log.Info(fmt.Sprintf("List %v", q.HighPrioroty))
		return q.HighPrioroty, nil
	default:
		log.Error("Please provide proper filters")
		return list, errors.New("Please provide proper filters")
	}
}