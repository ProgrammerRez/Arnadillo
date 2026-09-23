package jobs

// Creating Job Status Constants
const(
	JobQeued = "queued"
	JobRunning = "running"
	JobFailed = "failed"
	JobDoOver = "do-over"
	JobFinished	= "fininshed"
)


// Creating Main Job Request
type Job struct{
	Status 				string				`json:"job_status"`
	DataParams			DataParameters		`json:"data_params"`
	TrainingParams 		TrainingParameters	`json:"training_params"`
	EvalParams			map[string]string	`json:"eval_params"`
	DoOverParams		map[string]string	`json:"do_over_params"`
	ID					int					`json:"id"`
}


// Data Prameters Struct (Would be a carry over from the data management service)
type DataParameters struct{
	NumericalColumns 	[]string	`json:"num_cols"`		// Storing for Preprocessing Pipeline
	CategoricalColumns	[]string	`json:"cat_cols"`		// Storing for Preprocessing Pipeline
	// Actual Data to be Trained on
	ColumnList			[]string	`json:"column_list"`
	Data 				[][]string	`json:"row_data"`
	TargetLabel			string		`json:"target_label"`	// Column that's the target variable
	ID					int			`json:"id"`				// csv_id of the object from `/session` endpoint in the endpoint
	TestSplit			float32		`json:"test_split"`		// The ratio of test to train data
}

// Training Parameters Struct for the Model Training
type TrainingParameters struct{

	// These parameters would decide between different classes and models of neural networks and traditional alogirithms.
	// Would be `traditional`, `deep`, `XGBoost/Catboost`
	ModelType		string				`json:"model_type"`
	// Specific Model Name like `Logistic Regression`, `CNN`, etc.
	ModelName 		string				`json:"model_name"`
	// Since every model has specific parameters I'll just use a dictionary of possible parameters
	ModelParams		map[string]string	`json:"model_params"`
}


// Creating Evaluation Parameters for the Model Evaluation Part
type EvaluationParameters struct{
// This contains mostly tests and data to test and custom parameters for judging model performance
}

// These are the Do Over Parameters for Restarting a failed Job.
type DoOverParameters struct{
	// This contains mostly parameters judge model performace during and after training to automate training reruns
}