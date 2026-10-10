package jobs

import (
	"errors"
	log "mts/logging"
)

// This function will create the parameters definition of Do Over Params Objects
func CreateDoOverParamObjects(params map[string]any) (*DoOverParameters, error){

	log.SetLevel(log.InfoLevel)
	log.Info("Creating the Do Over Parametres")


	// First Getting Problem Type
	raw_type, err := getValue(params, "max_retries")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	max_retries, ok := raw_type.(int)
	if !ok{
		log.Error("Max Retries is Empty or Invalid Type")
		return nil, errors.New("Max Retries is Empty or Invalid Type")
	}

	// Then getting the max_deviance
	raw_type, err = getValue(params, "max_deviance")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	max_deviance, ok := raw_type.(float64)
	if !ok{
		log.Error("Max Deviance is Empty or Invalid Type")
		return nil, errors.New("Max Deviance is Empty or Invalid Type")
	}


	// Then extracting the required tests map
	raw_type, err = getValue(params, "fallback_params")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	fallaback_params, ok := raw_type.(map[string]any)
	if !ok || len(fallaback_params) == 0{
		log.Error("Fallback Params is Empty or Invalid Type")
		return nil, errors.New("Fallback Params Tests is Empty or Invalid Type")
	}

	return &DoOverParameters{
		MaxRetries: max_retries,
		MaximumDeviance: max_deviance,
		FallbackParams: fallaback_params,
	}, nil
}