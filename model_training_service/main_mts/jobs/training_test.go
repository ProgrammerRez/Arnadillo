package jobs

import (
	log "mts/logging"
	"reflect"
	"testing"
)

func TestCreateTrainParamObject(t *testing.T){
	// Initialize Logging
	log.Init()
	log.SetLevel(log.InfoLevel)

	// Creating Working Struct
	working_struct := map[string]any{
		"model_type": "traditional",
		"model_name": []string{"LogisticRegression"},
		"model_params": map[string]any{
			"limit": 1000,
		},
	}
	// Creating Working Output
	working_output := &TrainingParameters{
		ModelType: "traditional",
		ModelName: []string{"LogisticRegression"},
		ModelParams: map[string]any{
			"limit": 1000,
		},
	}

	tests := []struct{
		name 			string
		input 			map[string]any
		expectedOutput	any
		expectedError	bool
	}{
		{
			name: "Success - Works",
			input: working_struct,
			expectedOutput: working_output,
			expectedError: false,
		},

		{
			name: "Error - Model Type Doesn't Exist",
			input: map[string]any{},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Model Type Invalid Type",
			input: map[string]any{
				"model_type": 1000,
			},
			expectedOutput: nil,
			expectedError: true,
		},

				{
			name: "Error - Model Type Empty String",
			input: map[string]any{
				"model_type": "",
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Model Name Doesn't Exist",
			input: map[string]any{
				"model_type": "model_type",
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Model Name Invalid Type",
			input: map[string]any{
				"model_type": "model_type",
				"model_name": 1000,
			},
			expectedOutput: nil,
			expectedError: true,
		},

				{
			name: "Error - Model Name Empty String",
			input: map[string]any{
				"model_type": "model_type",
				"model_name": "",
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Model Params Don't Exist",
			input: map[string]any{
				"model_type": "model_type",
				"model_name": []string{"model_name"},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Model Params Invalid Type",
			input: map[string]any{
				"model_type": "model_type",
				"model_name": []string{"model_name"},
				"model_params": []any{"stuff"},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Model Params Empty",
			input: map[string]any{
				"model_type": "model_type",
				"model_name": []string{"model_name"},
				"model_params": map[string]any{},
			},
			expectedOutput: nil,
			expectedError: true,
		},
	}

	// Now to execute the tests

	for _, test :=  range(tests){
		t.Run(test.name, func(t *testing.T) {

			// First run the function
			train_params, err := CreateTrainParamObject(test.input)

			if test.expectedError{
				if err == nil{
					log.Error("Error was expected but didn't occur")
					t.Fatalf("Error was expected but didn't occur")
				}
				return
			}

			if err != nil{
				log.Error(err.Error())
				t.Fatalf("Unexpected Error Occurred: %v", err)
			}

			if !reflect.DeepEqual(train_params, test.expectedOutput){
				log.Error("Result Mismatch")
				t.Fatalf("Result Mismatch: Wanted: %v, Got: %v", test.expectedOutput, train_params)
			}
		})
	}
}