from sklearn.model_selection import GridSearchCV
from sklearn.pipeline import Pipeline
from models import DataParams
from logger import test_logger as log
from data_pipeline import DataPreporcessorCosntructor
from training_pipeline import TrainingPipeline
from sklearn.datasets import load_iris
from unittest import TestCase
import pandas as pd


# This class will test the training pipeline code

class TestTrainingPipeline(TestCase):

    @classmethod
    def setUpClass(cls):
        
        data = load_iris()

        df = pd.DataFrame(data=data.data, columns=data.feature_names)
        df["target"] = data.target
        df["target_names"] = df["target"].map({
            "setosa": 0,
            "virginica": 1,
            "versicolor": 2
        })
        

        # First Creating the data requirements
        cls.data_params = DataParams(**{
                "Data": df.to_numpy(),
                "NumericalColumns": df.select_dtypes("number").columns,
                "CategoricalColumns": df.select_dtypes("object").columns,
                "TargetLabel": "target",
                "ColumnList": df.columns,
                "NumericalOperations": ["standard_scaler"],
                "CategoricalOperations":{
                    "target_names": ["one_hot_encoding"]
                },
                "ID": 1, 
                "TestSplit": 0.2
            }
        )

        dpc = DataPreporcessorCosntructor(cls.data_params)
        cls.preprocessor: Pipeline | None = dpc.pipeline

        # Then creating the training params mock
        cls.traditional_training_parameters_regressor = {
            "model_type": "regressor",
            "model_name": ["linear_regression", "rfr"],
            "model_params":{
                "model_params":{
                    "linear_regression":{},
                    "rfr":{
                        "n_estimators": [100, 50, 150, 200],
                        "criterion": ["squared_error", "absolute_error"],
                        "max_depth": [None, 1, 2]
                    }
                },
                "pipeline_params":{
                    "grid_search":{
                        "scoring": "r2",
                        "cv": {
                            "shuffle_split":{
                                "n_splits": 10,
                                "test_size": 0.2,
                                "random_state": 42
                            }
                        }, 
                        "refit": True
                    }
                },
                "modifier_params":{
                    "select_k_best": {
                        "k": [10, 20, 5]
                    }
                }
            }
        }

        cls.traditional_training_parameters_classifier = {
            "model_type": "classifier",
            "model_name": ["logistic_regression", "rfc"],
            "model_params":{
                "model_params":{
                    "logistic_regression":{},
                    "rfc":{
                        "n_estimators": [100, 50, 150, 200],
                        "criterion": ["log_loss", "gini"],
                        "max_depth": [None, 1, 2]
                    }
                },
                "pipeline_params":{
                    "grid_search":{
                        "scoring": "accuracy",
                        "cv": 5, 
                        "refit": True
                    }
                },
                "modifier_params":{
                    "select_k_best": {
                        "k": [10, 20, 5]
                    }
                }
            }
        }

        cls.deep_training_parameters = {
        "input_shape": (df.shape[1]-1,1),
        "model_type": "deep",
        "model_name": ["sequential"],
        "model_params":{
            "layers": [
                (64, "dense", "relu"),
                (128, "dense", "relu"),
                (256, "dense", "relu"),
                (32, "dense", "relu"),
                (1, "dense", "sigmoid")
            ],
            "compile_params":{
                "loss": "sparse_categorical_crossentropy",
                "metrics": ["accuracy", "loss"],
            },
            "optimizer_params":{
                "Adam": {
                    "learning_rate": 0.001
                }
            },
            "train_params":{
                "epochs": 50,
                "batch_size": 32,
                "validation_split": 0.1
            },

            "callback_params":{
                "early_stopping": {
                    "monitor": "val_loss",
                    "restore_best_weights": True,
                    "start_from_epoch": 10
                },
                "model_checkpoint": {
                    "filepath": "./best_model.keras",
                    "monitor": "val_accuracy",
                    "save_best_only": True
                }
            }
        }
    }

    # This function tests the traditional regressor pipeline creation
    def test_traditional_regressor_pipeline(self):    
        try:
            training_pipeline_instance = TrainingPipeline(self.data_params, self.traditional_training_parameters_regressor, self.preprocessor)
            training_pipeline_instance.load_model_and_params()
            self.assertIsInstance(training_pipeline_instance.pipeline, GridSearchCV)
            log.info(f"Instance Checked for Traditional Regressor: {type(training_pipeline_instance.pipeline)}")
        except Exception as e:
            log.error(f"Traditional Regressor Constructor Failed: {e}")
            raise e
    

    # This function tests the traditional classifier pipeline creation
    def test_traditional_classifier_pipeline(self):    
        try:
            training_pipeline_instance = TrainingPipeline(self.data_params, self.traditional_training_parameters_classifier, self.preprocessor)
            training_pipeline_instance.load_model_and_params()
            self.assertIsInstance(training_pipeline_instance.pipeline, GridSearchCV)
            log.info(f"Instance Checked for Traditional Classifier: {type(training_pipeline_instance.pipeline)}")
        except Exception as e:
            log.error(f"Traditional Classifier Constructor Failed: {e}")
            raise e
    

    # This function tests the deep model creation pipeline
    def test_deep_model_pipeline(self):
        try:
            training_pipeline_instance = TrainingPipeline(self.data_params, self.deep_training_parameters, self.preprocessor)
            training_pipeline_instance.load_model_and_params()

            self.assertIsNotNone(training_pipeline_instance.model)
            self.assertIsNotNone(training_pipeline_instance.callbacks)
            self.assertIsNotNone(training_pipeline_instance.training_config)

            log.info("All instances Checked for Deep Model")
    
        except Exception as e:
            log.error(f"Deep Model Constructor Failed: {e}")
            raise ValueError(e)
        
