package jobs

import (
	log "mts/logging"
	"reflect"
	"testing"
)

func TestCreateDataParamObject(t *testing.T) {
	
	// Start the logging first
	log.Init()
	log.SetLevel(log.InfoLevel)

	// Creating Working Data Structures

	validDataMap := map[string]any{
		"data": []any{
			[]any{"row1_col1", "row1_col2"},
			[]any{"row2_col1", "row2_col2"},
		},
		"column_list": []any{"col1", "col2"},
	}

	validRegistryMap := map[string]any{
		"data": validDataMap,
		"target_col": "target",
	}

	validExpectedResult := []DataParameters{
			{
			Data: [][]string{
					{"row1_col1", "row1_col2"},
					{"row2_col1", "row2_col2"},
					},
			ColumnList: []string{"col1", "col2"},
			TargetLabel: "target",
			},	
	}

	// Now creating test scenarios
	tests := []struct{
		name 			string
		input			map[string]any
		expectedOutput	[]DataParameters
		expectedError 	bool
	}{
		{
			name : "Success - NO Error - Single Valid Register Entry",
			input: map[string]any{
				"registry": []any{validRegistryMap},
			},
			expectedOutput: validExpectedResult,
			expectedError: false,
		},

		{
			name: "Success - NO Error - Multiple Valid Registry Entry",
			input: map[string]any{
				"registry": []any{validRegistryMap, validRegistryMap},
			},
			expectedOutput: append(validExpectedResult, validExpectedResult[0]),
			expectedError: false,
		},

		{
			name: "Error - NO Registry key",
			input: map[string]any{},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Invalid Registry Type",
			input: map[string]any{
				"registry": "invlaid-string-type",
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Empty Slice",
			input: map[string]any{
				"registry": []any{},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Not a Map",
			input: map[string]any{
				"registry": []any{"not-a-map"},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Missing `data` key in `regsitry`",
			input: map[string]any{
				"registry": map[string]any{
					"target_col": "target",
				},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Malinformed Row Data",
			input: map[string]any{
				"registry": map[string]any{
					"data": map[string]any{
						"data": []any{"1","2"},
						"column_list": []any{"col_1", "col_2"},
					},
					"target_col": "target",
				},
			},
			expectedOutput: nil,
			expectedError: true,
		},

		{
			name: "Error - Malinformed Column List",
			input: map[string]any{
				"registry": map[string]any{
					"data": map[string]any{
						"data": [][]any{{"row_1"}, {"row_2"}},
						"column_list": []any{100,200},	
					},
					"target_col": "target",
				},
			},
			expectedOutput: nil,
			expectedError: true,
		},
	}

	// Now to execute testing
	for _, test := range tests{
		t.Run(test.name, func(t *testing.T) {

			result, err := CreateDataParamObject(test.input)

			if test.expectedError{
				if err == nil{
					t.Fatalf("Expected an error and got nil")
					}
				return //This mean test was passed
				}

			if err != nil{
				t.Fatalf("Unexpected Error Occurred: %v", err)
				}

			if !reflect.DeepEqual(result, test.expectedOutput){
				t.Fatalf("Output Mistmatched. Wanted: %v, Got: %v", test.expectedOutput, result)
			}
		})
	}
}