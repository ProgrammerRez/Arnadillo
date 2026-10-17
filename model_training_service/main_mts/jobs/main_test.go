package jobs

import (
	// "fmt"
	log "mts/logging"
	"testing"
)

// Testing the Get Jobs Function due to some bugs
func TestGetJobs(t *testing.T) {

	// Initialize Logging so we can debug later
	log.Init()
	log.SetLevel(log.InfoLevel)

	// Create Objects and initialize them for each priority level 
	low_job := &Job{}
	meh_job := &Job{}
	high_job :=	&Job{}
	
	// Initialize the Qeue
	q := InitializeQueue()
	q.LowPriority = append(q.LowPriority, low_job)
	q.MehPriority = append(q.MehPriority, meh_job)
	q.HighPrioroty = append(q.HighPrioroty, high_job)

	// Creating the test struct
	test_struct := []struct{
		name 			string
		filter 			int
		expectedJobs	[]*Job
		expectedError	bool
	}{
		{
			name: "low_priority",
			filter: 0,
			expectedJobs: q.LowPriority,
			expectedError: false,
		},

		{
			name: "meh_priority",
			filter: 1,
			expectedJobs: q.MehPriority,
			expectedError: false,
		},

		{
			name: "high_priority",
			filter: 2,
			expectedJobs: q.HighPrioroty,
			expectedError: false,
		},

		{
			name: "invalid priority",
			filter: 99,
			expectedJobs: nil,
			expectedError: true,
		},
	}

	// Now to run the tests

	for _, test := range test_struct{
		t.Run(test.name, func(t *testing.T) {

			// Get the data from the function
			actualData, err := q.GetJobs(test.filter)
			if test.expectedError && err != nil{
				log.Info("Planted Error Occured")
			}

			if !test.expectedError && err != nil{
				log.Error("An unexpected error occurred")
				t.Fatalf("An unexpected error occurred: %v", err)
			}

			if len(actualData) != len(test.expectedJobs){
				log.Error("Data Length Mismatch")
				t.Fatalf("Data Length Mismatch. Expected: %d, Resulted: %d", 
						len(test.expectedJobs), 
						len(actualData))

			for i := range test.expectedJobs{

				if actualData[i] != test.expectedJobs[i]{
					log.Error("Data Placement Mismatch")
					t.Fatalf("Data Placmeent Mismatch. Expected: %v at %d, Resulted: %v at %d",
							test.expectedJobs[i],
							i,
							actualData[i],
							i)		
					}
				}
			}
		})
	}
}