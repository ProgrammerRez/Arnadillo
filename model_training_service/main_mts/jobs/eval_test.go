package jobs

import (
	log "mts/logging"
	"reflect"
	"testing"
)

func TestCreateEvalParamObjects(t *testing.T) {

	// Starting the logging
	log.Init()
	log.SetLevel(log.InfoLevel)

	// Creating a working struct

	working_struct := map[string]any{
		"problem_type": "regression",
		"tests_required": map[string]float64{
			"r2_score": 0.98,
		},
	}

	working_result := &EvaluationParameters{
		ProblemType: "regression",
		TestsRequired: map[string]float64{
			"r2_score": 0.98,
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
			name: "Problem_Type is Empty",
			input: map[string]any{
				"problem_type": "",
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Problem_Type is Invalid Type",
			input: map[string]any{
				"problem_type": 1.0,
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Tests Requred is Empty",
			input: map[string]any{
				"problem_type": "regression",
				"tests_required": map[string]float32{},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Tests Requred is Invalid Type",
			input: map[string]any{
				"problem_type": "regression",
				"tests_required": []float32{
					
				},
			},
			expectedOutput: nil,
			expectedError: true,
		},
	}

	for _, test := range(tests){
		t.Run(test.name, func(t *testing.T) {

			// First Run the function
			eval_params, err := CreateEvalParamObjects(test.input)

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
			if !reflect.DeepEqual(test.expectedOutput, eval_params){
				log.Error("Result Mismatch")
				t.Fatalf("Result Mismatch: Wanted: %v, Got: %v", test.expectedOutput, eval_params)
			}
		})
	}
}