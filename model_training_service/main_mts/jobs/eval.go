package jobs

import (
	"errors"
	"fmt"
	log "mts/logging"
	// "reflect"
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

	rawMap, ok := raw_type.(map[string]any)
	if !ok || len(rawMap) == 0{
		log.Error(fmt.Sprintf("Required Tests is Empty or Invalid Type. Length: %d", len(rawMap)))
		return nil, errors.New("Required Tests is Empty or Invalid Type")
	}

	// Works Just needed to add another assertion step for tests_required map

	tests_required := make(map[string]float64)
	for k, v := range(rawMap){
		floatval, valid := v.(float64)
		if !valid{
			log.Error(fmt.Sprintf("Test Value for key %s is not a valid float", k))
			return nil, errors.New("Invalid Type for Test Metrics")
		}
		tests_required[k] = floatval
	}

	return &EvaluationParameters{
		ProblemType: problem_type,
		TestsRequired: tests_required,
	}, nil
}