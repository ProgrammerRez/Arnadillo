from ast import Pass
from typing import Any, Dict, Callable
from sklearn.base import BaseEstimator
from sklearn.linear_model import LinearRegression, LogisticRegression
from sklearn.ensemble import RandomForestRegressor, RandomForestClassifier
from sklearn.model_selection import GridSearchCV
from sklearn.model_selection import ShuffleSplit, StratifiedKFold
from tensorflow.keras.callbacks import ModelCheckpoint, EarlyStopping
from tensorflow.keras.optimizers import Adam, Nadam
from tensorflow.keras.layers import Dense, ReLU, Dropout, Flatten
from tensorflow.keras.models import Sequential

# First Creating the Dictionaries

# Traditional Models
PIPELINE_SELECTION: Dict[str, Any] = {}
MODELS_DICT_REGRESSOR: Dict[str, Callable[..., BaseEstimator]] = {}
MODELS_DICT_CLASSIFIER: Dict[str, Callable[..., BaseEstimator]] = {}
MODIFIERS: Dict[str, Any] = {}

# Deep Models
DEEP_MODELS_DICT: Dict[str, Any] = {}
LAYERS_DICT: Dict[str, Any] = {}
OPTIMIZERS_DICT: Dict[str, Any] = {}
CALLBACKS_DICT: Dict[str, Any] = {}

# Then the decorator methods
def register_pipeline(name: str):
    def decorator(pipeline: Any):
        PIPELINE_SELECTION[name] = pipeline
        return pipeline
    return decorator

def regsiter_modifiers(name: str):
    def decorator(modifier: Any):
        MODIFIERS[name] = modifier
        return modifier
    return decorator

def register_model_regressor(name: str):
    def decorator(model: Any):
        MODELS_DICT_REGRESSOR[name] = model
        return model
    return decorator

def register_model_classifier(name: str):
    def decorator(model: Any):
        MODELS_DICT_CLASSIFIER[name] = model
        return model
    return decorator

def register_deep_model(name: str):
    def decorator(model: Any):
        DEEP_MODELS_DICT[name] = model
        return model
    return decorator

def register_layer(name: str):
    def decorator(layer: Any):
        LAYERS_DICT[name] = layer
        return layer
    return decorator

def register_optimizer(name: str):
    def decorator(optimizer: Any):
        OPTIMIZERS_DICT[name] = optimizer
        return optimizer
    return decorator

def register_callbacks(name: str):
    def decorator(callback: Any):
        CALLBACKS_DICT[name] = callback
        return callback
    return decorator


# Registering the Pipelines
@register_pipeline("grid_search")
def grid_search(GridSearchCV):
    pass

# Register Regressor Models
@register_model_regressor("linear_regression")
def lr_regress(LinearRegression):
    pass

@register_model_regressor("rfr")
def rfr_regress(RandomForestRegressor):
    pass

# Register Classifer Models
@register_model_classifier("logistic_regression")
def log_classify(LogisticRegression):
    pass

@register_model_classifier("rfc")
def rfc_classify(RandomForestClassifier):
    pass

# Register Modifier
@regsiter_modifiers("startified_k_fold")
def stratify_k_fold(StratifiedKFold):
    pass

@regsiter_modifiers("shuffle_shift")
def shuffle_split(ShuffleSplit):
    pass

# Register Deep Models
@register_deep_model("sequential")
def sequential_model(Sequential):
    pass

# Register Layers
@register_layer("flatten")
def flatten_layer(Flatten):
    pass

@register_layer("dense")
def dense_layer(Dense):
    pass

@register_layer("dropout")
def dropout_layer(Dropout):
    pass

@register_layer("relu")
def relu_layer(ReLU):
    pass

# Register Optimizers
@register_optimizer("Adam")
def adam_optimizer(Adam):
    pass

@register_optimizer("Nadam")
def nadam_optimizer(Nadam):
    pass

# Register Callbacks
@register_callbacks("model_checkpoint")
def model_check(ModelCheckpoint):
    pass

@register_callbacks("early_stopping")
def early_stop(EarlyStopping):
    pass

