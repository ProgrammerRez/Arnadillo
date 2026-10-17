package jobs

import (
	log "mts/logging"
	"reflect"
	"testing"
)

func TestCreateDoOverParamObjects(t *testing.T) {

	// Starting the logging
	log.Init()
	log.SetLevel(log.InfoLevel)

	// Creating a working struct

	working_struct := map[string]any{
		"max_retries": 3.0,
		"max_deviance": 0.5,
		"fallback_params": map[string]any{
			"jobs": -1.0,
		},
	}

	working_result := &DoOverParameters{
		MaxRetries: float64(3.0),
		MaximumDeviance: float64(.5),
		FallbackParams: map[string]any{
			"jobs": -1.0,
		},
	}

	// Creating Tests Array
	tests := []struct{
		name 			string
		input 			map[string]any
		expectedOutput	any
		expectedError	bool
	}{
		{
			name: "Success - Works",
			input: working_struct,
			expectedOutput: working_result,
			expectedError: false,
		},

		{
			name: "Max Retries is Empty",
			input: map[string]any{
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Max Retries is Invalid Type",
			input: map[string]any{
				"max_retries": 1,
			},
			expectedOutput: nil,
			expectedError: true,
		},

				{
			name: "Max Deviance is Empty",
			input: map[string]any{
				"max_retries": 1.0,
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Max Deviance is Invalid Type",
			input: map[string]any{
				"max_retries": 1.0,
				"max_deviance": 0.5,
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Fallback params is Empty",
			input: map[string]any{
				"max_retries": 1.0,
				"max_deviance": 0.5,
				"fallback_params": map[string]any{},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Fallback params is Invalid Type",
			input: map[string]any{
				"max_retries": 1.0,
				"max_deviance": 0.5,
				"fallback_params": map[string]float32{},
			},
			expectedOutput: nil,
			expectedError: true,
		},
	}

	for _, test := range(tests){
		t.Run(test.name, func(t *testing.T) {

			// First Run the function
			do_over_params, err := CreateDoOverParamObjects(test.input)

			// Then check for all error cases
			if test.expectedError == true{
				if err == nil{
					log.Error("Error was expected but it didn't occur")
					t.Fatalf("Error was expected but it didn't occur")
				}
				return
			}

			if err != nil{
				log.Error(err.Error())
				t.Fatalf("Unexpected Error Occured, %v", err)
			}

			// Then verify output
			if !reflect.DeepEqual(test.expectedOutput, do_over_params){
				log.Error("Result Mismatch")
				t.Fatalf("Result Mismatch: Wanted: %v, Got: %v", test.expectedOutput, do_over_params)
			}
		})
	}
}