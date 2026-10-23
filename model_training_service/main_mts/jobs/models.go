package jobs

import "sync"

// package main

// Creating Job Status Constants
const(
	JobQeued = "queued"
	JobRunning = "running"
	JobFailed = "failed"
	JobDoOver = "do-over"
	JobFinished	= "fininshed"
)

// Creating JobQueue Object
type JobQeueu struct{
	HighPrioroty	[]*Job			`json:"high_priority"`
	MehPriority 	[]*Job			`json:"meh_priority"`
	LowPriority		[]*Job			`json:"low_priority"`
	NextID			int				`json:"next_id"`
	Mu				sync.RWMutex	`json:"mutex: Dont touch"`
}

// Creating Main Job Request
type Job struct{
	DataParams			[]DataParameters		`json:"data_params"`
	Status 				string					`json:"job_status"`
	TrainingParams 		TrainingParameters		`json:"training_params"`
	EvalParams			EvaluationParameters	`json:"eval_params"`
	DoOverParams		DoOverParameters		`json:"do_over_params"`
	ID					int						`json:"id"`
}


// Data Prameters Struct (Would be a carry over from the data management service)
type DataParameters struct{
	NumericalColumns 	[]string			`json:"num_cols"`		// Storing for Preprocessing Pipeline
	CategoricalColumns	[]string			`json:"cat_cols"`		// Storing for Preprocessing Pipeline
	// Actual Data to be Trained on
	ColumnList			[]string			`json:"column_list"`
	Data 				[][]string			`json:"data"`
	NumericalOperations	[]string			`json:"num_ops"`		//The operations for preprocessing numerical data in the preprocessing pipeline
	TargetLabel			string				`json:"target_col"`		// Column that's the target variable
	ObjectOperations	map[string][]string	`json:"cat_ops"`		//The operations for preprocessing categorical data in the preprocessing pipeline
	ID					int					`json:"id"`				// `csv_id` of the object from `/session` endpoint in the endpoint
	TestSplit			float64				`json:"test_split"`		// The ratio of test to train data
}

// Training Parameters Struct for the Model Training
type TrainingParameters struct{
	// Specific Model Name like `Logistic Regression`, `CNN`, etc.
	ModelName 		[]string				`json:"model_name"`
	// These parameters would decide between different classes and models of neural networks and traditional alogirithms.
	// Would be `traditional`, `deep`, `XGBoost/Catboost`
	ModelType		string				`json:"model_type"`
	// Since every model has specific parameters I'll just use a dictionary of possible parameters
	ModelParams		map[string]any		`json:"model_params"`
}


// Creating Evaluation Parameters for the Model Evaluation Part
type EvaluationParameters struct{
// This contains mostly tests and data to test and custom parameters for judging model performance
	ProblemType		string				`json:"problem_type"`	// Aiming for a binary answer of `classification` or `regression`
	TestsRequired	map[string]float64 	`json:"tests_required"`// Aiming for a test: minimum_proposed_value relationship like f1_score: 0.9 or 90%
}

// These are the Do Over Parameters for Restarting a failed Job.
type DoOverParameters struct{
	FallbackParams	map[string]any	`json:"fallback_params"`	 // In case of not working; will move to fallback parameters provided
	// This contains mostly parameters judge model performace during and after training to automate training reruns
	MaxRetries		float64  		`json:"max_retries"`	// Maximum Retires When Failing
	MaximumDeviance	float64 		`json:"max_deviance"`// Maximum Deviance from minimum_proposed_value before switching parametes or retrying training
}