package main_test

import (
	"encoding/json"
	"fmt"
	"mts/jobs"
	log "mts/logging"
	"reflect"
	"testing"
	"mts/utils"
)


// This function is for testing the Add Job Handler
func TestGetEndpoint(t *testing.T){

	// Start the logger
	log.Init()
	log.SetLevel(log.InfoLevel)
	log.Info("Testing the Add Job Handler")


	// Generating Valid Tests First

	validDataMap := map[string]any{
		"data": []any{
			[]any{"row1_col1", "row1_col2"},
			[]any{"row2_col1", "row2_col2"},
		},
		"column_list": []any{"col1", "col2"},
	}

	validRegistryMap := map[string]any{
		"registry": []any{map[string]any{
			"data": validDataMap,
			"target_col": "target",
			},
		},
	}

	// Creating Working Struct
	working_struct := map[string]any{
		"name": "session1",
		"urgency": 1,
		"data_params": validRegistryMap,
		"train_params": map[string]any{
			"model_type": "traditional",
			"model_name": "LogisticRegression",
			"model_params": map[string]any{
				"limit": 1000,
			},
		},
		"eval_params": map[string]any{
			"problem_type": "regression",
			"tests_required": map[string]float64{
				"r2_score": 0.98,
				"r2_scor": 0.98,
				"r2_sco": 0.98,
			},
		},
		"do_over_params": map[string]any{		
			"max_retries": 3,
			"max_deviance": 0.5,
			"fallback_params": map[string]any{
				"jobs": -1.0,
			},
		}, 
	}

	api_struct := map[string]any{
		"name": "session1",
		"urgency": 1,
		"data_params": map[string]any{},
		"train_params": map[string]any{
			"model_type": "traditional",
			"model_name": "LogisticRegression",
			"model_params": map[string]any{
				"limit": 1000,
			},
		},
		"eval_params": map[string]any{
			"problem_type": "regression",
			"tests_required": map[string]float64{
				"r2_score": 0.98,
				"r2_scor": 0.98,
				"r2_sco": 0.98,
			},
		},
		"do_over_params": map[string]any{		
			"max_retries": 3,
			"max_deviance": 0.5,
			"fallback_params": map[string]any{
				"jobs": -1.0,
			},
		}, 
	}


	// Working Job Result
	working_result_slice := []jobs.Job{
		{
			Status: "queued",
			ID: 1,
			DataParams: []jobs.DataParameters{
				{
					Data: [][]string{
							{"row1_col1", "row1_col2"},
							{"row2_col1", "row2_col2"},
							},
					ColumnList: []string{"col1", "col2"},
					TargetLabel: "target",
				},	
			},
			TrainingParams: jobs.TrainingParameters{
				ModelType: "traditional",
				ModelName: "LogisticRegression",
				ModelParams: map[string]any{
					"limit":1000.0, 
				},
			},
			EvalParams: jobs.EvaluationParameters{
				ProblemType: "regression",
				TestsRequired: map[string]float64{
					"r2_score": 0.98,
					"r2_scor": 0.98,
					"r2_sco": 0.98,				
				},
			},
			DoOverParams: jobs.DoOverParameters{
				FallbackParams: map[string]any{
					"jobs": -1.0,
				},
				MaxRetries: 3,
				MaximumDeviance: 0.5,
			},
		},

		{
			Status: "queued",
			ID: 2,
			DataParams: nil,
			TrainingParams: jobs.TrainingParameters{
				ModelType: "traditional",
				ModelName: "LogisticRegression",
				ModelParams: map[string]any{
					"limit":1000.0, 
				},
			},
			EvalParams: jobs.EvaluationParameters{
				ProblemType: "regression",
				TestsRequired: map[string]float64{
					"r2_score": 0.98,
					"r2_scor": 0.98,
					"r2_sco": 0.98,				
				},
			},
			DoOverParams: jobs.DoOverParameters{
				FallbackParams: map[string]any{
					"jobs": -1.0,
				},
				MaxRetries: 3,
				MaximumDeviance: 0.5,
			},
		},
	}

	get_jobs_filter := map[string]any{
		"filter": 1,
	}

	// // Set Ops Tests
	set_ops_test := map[string]any{
		"id": 1,
		"file_id": 0,
		"filter": 1,
		"num_ops": []string{"standard_scaler"},
		"cat_ops": map[string][]string{
			"col1": []string{"label_encoding"},
			"col2": []string{"one_hot_encoding"},
		},
	}

	ops_result_struct := []jobs.Job{
		{
			Status: "queued",
			ID: 1,
			DataParams: []jobs.DataParameters{
				{
					Data: [][]string{
							{"row1_col1", "row1_col2"},
							{"row2_col1", "row2_col2"},
							},
					ColumnList: []string{"col1", "col2"},
					TargetLabel: "target",
					NumericalOperations: []string{"standard_scaler"},
					ObjectOperations: map[string][]string{
						"col1": []string{"label_encoding"},
						"col2": []string{"one_hot_encoding"},
					},
				},	
			},
			TrainingParams: jobs.TrainingParameters{
				ModelType: "traditional",
				ModelName: "LogisticRegression",
				ModelParams: map[string]any{
					"limit":1000.0, 
				},
			},
			EvalParams: jobs.EvaluationParameters{
				ProblemType: "regression",
				TestsRequired: map[string]float64{
					"r2_score": 0.98,
					"r2_scor": 0.98,
					"r2_sco": 0.98,				
				},
			},
			DoOverParams: jobs.DoOverParameters{
				FallbackParams: map[string]any{
					"jobs": -1.0,
				},
				MaxRetries: 3,
				MaximumDeviance: 0.5,
			},
		},

		{
			Status: "queued",
			ID: 2,
			DataParams: nil,
			TrainingParams: jobs.TrainingParameters{
				ModelType: "traditional",
				ModelName: "LogisticRegression",
				ModelParams: map[string]any{
					"limit":1000.0, 
				},
			},
			EvalParams: jobs.EvaluationParameters{
				ProblemType: "regression",
				TestsRequired: map[string]float64{
					"r2_score": 0.98,
					"r2_scor": 0.98,
					"r2_sco": 0.98,				
				},
			},
			DoOverParams: jobs.DoOverParameters{
				FallbackParams: map[string]any{
					"jobs": -1.0,
				},
				MaxRetries: 3,
				MaximumDeviance: 0.5,
			},
		},
	}


	// Creating Test Scenarios

	tests := []utils.TestingStuct{
		{
			Name: "Add Job Struct",
			Endpoint: "/data_params",
			RequestType: "POST",
			Input: working_struct,
			ExpectedOutput: "Worked :>",
			ExpectedError: false,
		},

		{
			Name: "API Switch Job Struct",
			Endpoint: "/data_params",
			RequestType: "POST",
			Input: api_struct,
			ExpectedOutput: "Worked :>",
			ExpectedError: false,
		},

		{
			Name: "Get Job Struct",
			Endpoint: "/get_jobs",
			RequestType: "POST",
			Input: get_jobs_filter,
			ExpectedOutput: working_result_slice,
			ExpectedError: false,
		},

		{
			Name: "Set Ops Struct",
			Endpoint: "/set_ops",
			RequestType: "POST",
			Input: set_ops_test,
			ExpectedOutput: "Worked :>",
			ExpectedError: false,
		},

		{
			Name: "Get Job Struct",
			Endpoint: "/get_jobs",
			RequestType: "POST",
			Input: get_jobs_filter,
			ExpectedOutput: ops_result_struct,
			ExpectedError: false,
		},
	}

	runTests(t, tests)

}


// This function runs all the provided tests
func runTests(t *testing.T, tests []utils.TestingStuct){
	for _, test :=  range(tests){
		t.Run(test.Name, func(t *testing.T) {
	
			// Run the command
			data, err := utils.GetEndpoint(test.Endpoint, test.RequestType, test.Input)

			if test.ExpectedError{
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

			if test.ExpectedOutput != nil{
				json_bytes, err := json.Marshal(data)
				if err != nil{
					log.Error(err.Error())
					t.Fatalf("%s", fmt.Sprintf("%v",err.Error()))
					return
				}
				target_type := reflect.TypeOf(test.ExpectedOutput)
				val := reflect.New(target_type)

				if json.Unmarshal(json_bytes, val.Interface()) == nil{
					data = val.Elem().Interface()
				}
			}

			if !reflect.DeepEqual(data, test.ExpectedOutput){
				log.Error("Result Mismatch")
				t.Fatalf("Result Mismatch: Wanted: %v, Got: %v", test.ExpectedOutput, data)
			}
		})
	}
}