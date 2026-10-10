package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	log "mts/logging"
	"strings"
)

// This file is meant to streamline the pathway for training parameters

func CreateTrainParamObject(params map[string]any) (*TrainingParameters, error){

	log.SetLevel(log.InfoLevel)
	log.Info("Extracting Training Parameters")

	// Extracting model_type fist
	raw_type, err := getValue(params, "model_type")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	model_type, ok := raw_type.(string)
	if !ok || strings.TrimSpace(model_type) == ""{
		log.Error("Model Name is Empty or Invalid Type")
		return nil, errors.New("Model Name is Empty or Invalid Type")
	}


	// Then extracting Model Name
	raw_name, err := getValue(params, "model_name")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	model_name, ok := raw_name.(string)
	if !ok || strings.TrimSpace(model_name) == ""{
		log.Error("Model Name is Empty or Invalid Type")
		return nil, errors.New("Model Name is Empty or Invalid Type")
	}


	// Finally extracting and verifying if model_params is okay
	raw_params, err := getValue(params, "model_params")
	if err != nil{
		log.Error(err.Error())
		return nil, err
	}

	var model_params map[string]any

	// Switching Based on Input i.e, JSON, simple map[string]any
	switch v := raw_params.(type){
	case map[string]any:
		model_params = v
		model_params, ok = raw_params.(map[string]any);
		if !ok || len(model_params) == 0{
			log.Error("Model Params is Empty or of Invalid Type")
			return nil, errors.New("Model Params is Empty or of Invalid Type")
		}
	case string:
		if err := json.Unmarshal([]byte(v), &model_params); err != nil{
			log.Error(fmt.Sprintf("Failed to Parse Model Params String into  map %v", err))
			return nil, errors.New("Model Params is Empty or of Invalid Type")		
		}
	default:
		log.Error(fmt.Sprintf("Unexpected type for model_params %T", raw_params))
		return nil, errors.New("Model Params is Empty or of Invalid Type")
	}
	
	if len(model_params) == 0{
		log.Error("Model Params is Empty")
		return nil, errors.New("Model Params is Empty")
	}


	return &TrainingParameters{
		ModelType: model_type,
		ModelName: model_name,
		ModelParams: model_params,
	}, nil
}