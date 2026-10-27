from sklearn.feature_selection import f_classif
from sklearn.feature_selection import SelectKBest
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
CV_STEPS: Dict[str, Any] = {}

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

def register_cv_steps(name: str):
    def decorator(cv: Any):
        CV_STEPS[name] = cv
        return cv
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
def grid_search(**kwargs):
    return GridSearchCV(**kwargs)

# Register Regressor Models
@register_model_regressor("linear_regression")
def lr_regress():
    return LinearRegression()

@register_model_regressor("rfr")
def rfr_regress():
    return RandomForestRegressor()

# Register Classifer Models
@register_model_classifier("logistic_regression")
def log_classify():
    return LogisticRegression()

@register_model_classifier("rfc")
def rfc_classify():
    return RandomForestClassifier()

# Register Modifier
@regsiter_modifiers("select_k_best")
def select_k_best(**kwargs):
    return SelectKBest(score_func=f_classif, **kwargs)

# Register CV Steps
@register_cv_steps("startified_k_fold")
def stratify_k_fold(**kwargs):
    return StratifiedKFold(**kwargs)

@register_cv_steps("shuffle_split")
def shuffle_split(**kwargs):
    return ShuffleSplit(**kwargs)

# Register Deep Models
@register_deep_model("sequential")
def sequential_model(layers):
    return Sequential(layers)

# Register Layers
@register_layer("flatten")
def flatten_layer(*args, **kwargs):
    return Flatten(*args, **kwargs)

@register_layer("dense")
def dense_layer(*args, **kwargs):
    return Dense(*args, **kwargs)

@register_layer("dropout")
def dropout_layer(*args, **kwargs):
    return Dropout(*args, **kwargs)

@register_layer("relu")
def relu_layer(*args, **kwargs):
    return ReLU(*args, **kwargs)

# Register Optimizers
@register_optimizer("Adam")
def adam_optimizer(*args, **kwargs):
    return Adam(*args, **kwargs)

@register_optimizer("Nadam")
def nadam_optimizer(*args, **kwargs):
    return Nadam(*args, **kwargs)

# Register Callbacks
@register_callbacks("model_checkpoint")
def model_check(*args, **kwargs):
    return ModelCheckpoint(*args, **kwargs)

@register_callbacks("early_stopping")
def early_stop(*args, **kwargs):
    return EarlyStopping(*args, **kwargs)

