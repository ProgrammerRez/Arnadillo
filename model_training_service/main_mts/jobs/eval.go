package jobs

import (
	"fmt"
	"errors"
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

	tests_required, ok := raw_type.(map[string]float64)
	if !ok || len(tests_required) == 0{
		log.Error(fmt.Sprintf("Required Tests is Empty or Invalid Type. Length: %d", len(tests_required)))
		return nil, errors.New("Required Tests is Empty or Invalid Type")
	}

	return &EvaluationParameters{
		ProblemType: problem_type,
		TestsRequired: tests_required,
	}, nil
}