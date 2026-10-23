package jobs

import (
	"errors"
	"fmt"
	log "mts/logging"
	"strings"
)

// This function will create the parameters definition of Eval Params Objects
func CreateEvalParamObjects(params map[string]any) (*EvaluationParameters, error){

	log.SetLevel(log.InfoLevel)
	log.Info("Creating the Evaluation Parametres")


	// First Getting Problem Type
	raw_type, err := getValue(params, "problem_type")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	problem_type, ok := raw_type.(string)
	if !ok || strings.TrimSpace(problem_type) == ""{
		log.Error("Problem Type is Empty or Invalid Type")
		return nil, errors.New("Problem Type is Empty or Invalid Type")
	}

	// Then extracting the required tests map
	raw_type, err = getValue(params, "tests_required")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	test_required := make(map[string]float64)

	switch v := raw_type.(type){
	case map[string]float64:
		test_required = v
	case map[string]any:
		for k, val := range(v){
			floatVal, valid := val.(float64)
			if !valid{
				log.Error(fmt.Sprintf("Test Value for key %s is not a valid float", k))
				return nil, errors.New("Invalid Type for Test Metrics")
			}
			test_required[k] = floatVal
		}
	default:
		log.Error("Required Tests is Invalid Type")
		return nil, errors.New("Required Tests is Invalid Type")
	}


	if len(test_required) == 0{
		log.Error(fmt.Sprintf("Required Tests is Empty. Length: %d", len(test_required)))
		return nil, errors.New("Required Tests is Empty")
	}


	return &EvaluationParameters{
		ProblemType: problem_type,
		TestsRequired: test_required,
	}, nil
}