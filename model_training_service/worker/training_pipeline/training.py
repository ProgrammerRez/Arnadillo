from typing import Literal, Any
from logger import training_logger as log
from .training_tools import (
                        PIPELINE_SELECTION, 
                        MODELS_DICT_REGRESSOR, 
                        MODELS_DICT_CLASSIFIER, 
                        MODIFIERS,
                        DEEP_MODELS_DICT,
                        LAYERS_DICT, 
                        CALLBACKS_DICT,
                        OPTIMIZERS_DICT
                        )
from models import DataParams, TrainingParameters
from tensorflow.keras.layers import InputLayer
from sklearn.pipeline import Pipeline


# Now creating the main training pipeline class

class TrainingPipeline:
    def __init__(self, data_params: DataParams, training_params: TrainingParameters, preprocessor = None):
        log.info("Initializing Pipeline")

        # Setting the class attributes
        self.data: DataParams = data_params
        self.training_params: TrainingParameters = training_params

        self.preprocessor = preprocessor
        
        # For Traditional Models
        self.pipeline = None

        # For Deep Models
        self.model = None
        self.callbacks = None
        self.training_config = None

    
    def load_model_and_params(self):

        try:
            """This function is suppossed to load the model pipeline using provided info and params"""

            log.info("Making and Compiling the Pipeline from Model and Params Info")
            
            # Doing dict empty check
            model_name = self.training_params.get("model_name", "")
            model_type = self.training_params.get("model_type", "")
            model_params = self.training_params.get("model_params", {})

            if not model_type:
                log.error("Model Type is Empty")
                raise ValueError("Model Type is Empty")
            
            if not model_name:
                log.error("Model Name Empty")
                raise ValueError("Model Type is Empty")
            
            if not model_params:
                log.error("Model Params Empty")
                raise ValueError("Model Params Empty")
            
            type_model = str(model_type).strip().lower()

            if  type_model == "regressor" or type_model == "classifier":
                self._make_traditional_pipeline(model_name, type_model, model_params)
            elif type_model == "deep":
                self._make_deep_pipeline(model_name, model_params)

            else:
                log.error("No Proper Model Type Provided")
                raise ValueError("No Proper Model Type Provided")
    
        except Exception as e:
            log.error(f"Unexpected Error Occurred during loading model and params: {e}")    
            raise ValueError(f"Unexpected Error Occurred during loading model and params: {e}")
        
        finally:
            log.info("Loaded Model and Params into a Pipeline")
        



    def _make_traditional_pipeline(self, 
                                    model_name: list[str], 
                                    model_type: Literal["regressor", "classifier"],
                                    model_params: dict[str, Any]):
        
        try:

            # Get all the requisite parameters
            pipeline_config: dict = model_params.get("pipeline_params", {})
            modifier_config: dict = model_params.get("modifier_params", {})
            all_params: dict = model_params.get("model_params", {})

            # Checking for empty parameters

            # No pipeline config is an error
            if not pipeline_config:
                log.error("No Pipeline Parameters Provided")
                return None
            
            if not modifier_config:
                log.warning("No Modifier Parameters Provided")

            if not all_params:
                log.warning("No Model Finetuning Parameters Provided")

            log.info("Loaded all the Provided Parameters")
            
            # Creating Pipeline Steps to append 
            pipeline_steps = []
            if self.preprocessor:
                pipeline_steps.append(("preprocessor", self.preprocessor))
                log.info("Added Preprocessor As well")
            
            # Then appending the modifiers
            for modifier, modifier_params in modifier_config.items():
                selected_modifier = MODIFIERS.get(modifier)
                if not selected_modifier:
                    log.error(f"Designated modifier doesn't exist: {modifier}")
                    continue 
                if isinstance(modifier_params, dict):
                    pipeline_steps.append((modifier, selected_modifier(**modifier_params)))
                    log.info(f"Found Modifier Params for modifier: {modifier}")
                else:
                    pipeline_steps.append((modifier, selected_modifier()))
                log.info(f"Added {modifier} to modifiers")


            # Now applying a placeholder model in the pipeline
            default_dict = MODELS_DICT_REGRESSOR if model_type == "regressor"  else MODELS_DICT_CLASSIFIER
            default_model = default_dict.get(model_name[0])
            if not default_model:
                log.error(f"No default model found for name: {model_name[0]}")
                return None
        
            pipeline_steps.append(("model", default_model))
        
            # Now Pipeline is Complete
            pipeline = Pipeline(steps=pipeline_steps)

            log.info("Pipeline Placeholders Set")

            # First Creating the parameters grid
            params_grid = []

            # Checking if we have parameters available
            
            for model in model_name:
                
                # Selecting the model using a one-liner
                selected_model = default_dict.get(model)

                # Checking if model exists
                if not selected_model:
                    log.error(f"Designated Model not found: {model}")
                    continue

                # Initialzing params_grid
                params_grid: dict[str, Any] = {
                    "model": [selected_model]
                }

                # Getting Model Parameters
                model_config = all_params.get(model, {})

                for param_name, params in model_config.items():
                    params_grid[f"model__{param_name}"] = params
                    

                # Then getting modifier parameters

                for modifier_name, modifier_params in modifier_config.items():
                    if isinstance(modifier_params, dict):
                        for mod_p_name, mod_p_vals in modifier_params.items():
                            params_grid[f"{modifier_name}__{mod_p_name}"] = mod_p_vals
            
            log.info("Param Grid Completed")

            # Finally creating the training pipeline instance

            if not isinstance(pipeline_config, dict):
                log.error("pipeline config is not a dict")
                return None

            try:
                pipeline_name, pipeline_kwargs = next(iter(pipeline_config.items()))
            except (StopIteration, AttributeError):
                log.error("Invalid pipeline_config structure")
                return None
            
            selected_pipeline = PIPELINE_SELECTION[pipeline_name]
            if not selected_pipeline:
                log.error(f"Pipeline withe name: {pipeline_name} is not available")
                return None

            self.pipeline = selected_pipeline(
                estimator=pipeline,
                param_grid = params_grid,
                **pipeline_kwargs
            )
        except Exception as e:
            log.error(f"Traditional Pipeline Error {e}")
            return None
        
        finally:
            log.info("Traditional Pipeline Compilation Completed")

    def _make_deep_pipeline(self, model_name: list[str], model_params: dict[str, Any]):
        try:
            # Getting the model parameters and checking for empty ones
            input_shape = DEEP_MODELS_DICT.get("input_shape", None)
            selected_model = DEEP_MODELS_DICT.get(model_name[0])
            layers_config = model_params.get("layers", [])
            
            compile_conifg = model_params.get("compile_params", {})
            train_config = model_params.get("train_params", {})
            
            optimizer_config = model_params.get("optimizer_params", {})
            callbacks_config = model_params.get("callback_params", {})

            # Check for None Type Model
            if not selected_model:
                log.error(f"{model_name[0]} not found / registered")
                return None

            # Checking for empty layers
            if not layers_config:
                log.error("No Layers Config Provided")
                return None
            
            # Check for empty compilation config
            if not compile_conifg:
                log.error("No Compilation Config Provided")
                return None
            
            # Check for training config
            if not train_config:
                log.error("No Training Config Provided")
                return None

            # Check for optimizer config
            if not optimizer_config:
                log.error("No Optimizer Config Provided")
                return None

            # Check for callbacks config
            if not callbacks_config:
                log.error("No Callbacks Config Provided")
                return None

            # First adding layers to the selected model
            final_layers = []

            if input_shape:
                final_layers.append(InputLayer(input_shape=input_shape))
                
            for idx, layer_tuple in enumerate(layers_config):
                neurons, layer_type, activation = layer_tuple

                layer_class = LAYERS_DICT.get(layer_type, None)
                if not layer_class:
                    log.error(f"Layer Class doesn't exist: {layer_class}")
                    continue
                
                if idx == 0 and not input_shape:
                    final_layers.append(layer_class(neurons, activation=activation, input_shape=input_shape))
                else:
                    final_layers.append(layer_class(neurons, activation=activation))


            # Build the model now
            model_builder = selected_model(final_layers)
            if not model_builder:
                log.error("Model Is None Type Here")
                raise ValueError("Model Is None Type Here")

            # Get Compilation Parameters
            loss_fn = compile_conifg.get("loss", None)

            try:
                optimizer, optimzier_params = next(iter(optimizer_config.items()))
            except (StopIteration, AttributeError):
                log.error("optimizer_config structure is bad")
                return None
            
            optimizer_instance = OPTIMIZERS_DICT.get(optimizer, None)
            if not optimizer_instance:
                log.error(f"{optimizer} not registered / found")
                return None
            
            # Finally Creating the optimizer
            final_optimizer = optimizer_instance(**optimzier_params)
            
            # Getting the metrics list
            metrics_list = compile_conifg.get("metrics_list", ["accuracy"])

            # Finally Compile the model
            model_builder.compile(
                optimizer=final_optimizer,
                loss=loss_fn,
                metrics=metrics_list
            )

            # Now Adding the callbacks
            callbacks_list = []
            for cb_name, cb_params in callbacks_config.items():
                cb = CALLBACKS_DICT.get(cb_name, None)
                if not cb:
                    log.error(f"Callback: {cb_name} not found / registered")
                    continue
                
                cb_instance = cb(**cb_params)
                callbacks_list.append(cb_instance)

                log.info(f"Loaded Callback with name: {cb_name}")
        
            # Storing all the metadata for later running of the pipeline
            self.model = model_builder
            self.callbacks = callbacks_list
            self.training_config = train_config
            

            log.info(f"Deep Model Pipeline for {model_name[0]} successfully built and compiled")
        
        except Exception as e:
            log.error(f"Deep Model Pipeline Error: {e}")
            raise ValueError(f"Deep Model Pipeline Error: {e}")
            
        finally:
            log.info("Completed the Deep Learning Pipeline Compilation")
            return None