from typing import Any
from pandas.core.frame import DataFrame
from sklearn.base import BaseEstimator
from sklearn.compose import ColumnTransformer
from sklearn.pipeline import Pipeline
from sklearn.preprocessing import StandardScaler, OneHotEncoder, LabelEncoder
from collections.abc import Callable
from logger import dp_logger as log
from models import DataParams
import pandas as pd

# Operation Dictionaries
NUM_OPS_REGISTRY: dict[str, Callable[pd.DataFrame, pd.DataFrame]] = {}
CAT_OPS_REGISTRY: dict[str, Callable[[pd.DataFrame, str], pd.DataFrame]] = {}

# Operation Decorator functions
def regitser_num_op(name: str):
    def decorator(func: Callable[pd.DataFrame, pd.DataFrame]):
        NUM_OPS_REGISTRY[name] = func
        return func
    return decorator

def register_cat_op(name: str) :
    def decorator(func: Callable[[pd.DataFrame, str], pd.DataFrame])
        CAT_OPS_REGISTRY[name] = func
        return func
    return decorator

# Creating the functions
@register_num_op("standard_scaler")
def standard_scaler(df: pd.Dataframe)-> pd.DataFrame:
    scaler: BaseEstimator = StandardScaler().set_output("pandas")
    return scaler.fit_transform(df)

@register_cat_op("label_encoding")
def label_encoding(df: pd.DataFrame, col: str) -> pd.DataFrame:
    df[col] = LabelEncoder().fit_transform(df[col])
    return df

@register_cat_op("one_hot_encoding")
def one_hot_encoding(df: pd.DataFrame, col: str) -> pd.DataFrame:
    one = OneHotEncoder(sparse_output=False, handle_unknown="ignore")
    encoded_data = one.fit_transform(df[col])
    encoded_df = pd.DataFrame(
        data = encoded_data,
        columns= one.get_feature_names_out(col)
    )
    final_df: DataFrame = pd.concat([df, encoded_df], axis=1).drop([col], axis=1)
    return final_df

@register_cat_op("dummy_encoding")
def dummy_encoding(df: pd.DataFrame, col: str) -> pd.DataFrame:
    dummies: DataFrame = pd.get_dummies(
        df,
        [col],
        drop_first=True
    )
    return dummies


# This object will be the main construct for the data preprocessing pipeline
class DataPreporcessorCosntructor:
    def __init__(self, params: DataParams) -> None:
        # Initializing the 
        self.params = params 
        self.df: pd.DataFrame | None = None
        self.pipeline: Pipeline | None = None
    
    def load_data(self) -> None:
        try:
            log.info("Loading the Dataframe")
            # Checking if Data is present
            if not self.params.ColumnList or not self.params.Data:
                log.error("No Column or Row Data found in the data object")
                raise ValueError("No Column or Row Data found in the data object")
            
            self.df = pd.DataFrame(data=self.params.Data, columns=self.params.ColumnList)
            log.info(f"Loaded dataframe with shape: {self.df.shape}")

            # Converting Numerical Columns
            for col in self.params.NumericalColumns:
                if col in self.df.columns:
                    self.df[col] = pd.to_numeric(self.df[col], errors="coerce")
                    log.info(f"Converted Numerical Column to Numerical Type: {col}")

            # Converting Categorical Columns
            for col in self.params.CategoricalColumns:
                if col in self.df.columns:
                    self.df[col] = self.df[col].astype(str)
                    log.info(f"Converted Categorical Column to Object Type: {col}")

            # Doing a final Check
            log.info(f"Final Data Type Check: {self.df.dtypes}") 

        except Exception as e:
            log.error(f"Unexpected Error Occured: {e}")
            raise ValueError(e)

    def build_pipeline(self) -> None:
        """This function will build the sklearn pipeline for the job"""
        
        try:

            # Makes sure data is loaded
            if self.df is None:
                self.load_data()

            log.info("Starting Pipeline Build")

            # Parsing Numerical Steps
            num_steps: list[[str, Callable[pd.DataFrame, pd.DataFrame]]] = []

            for ops in self.params.NumericalOperations:
                if ops in NUM_OPS_REGISTRY:
                    num_steps.append((ops, NUM_OPS_REGISTRY[ops]))
                    log.info(f"Added Step in num_steps: {ops}")
                else:
                    log.warning(f"Step doesn't exist in registry: {ops}")
            
            log.info("Completed Numerical Steps Building")

            # Parsing Categorical Steps
            cat_steps: list[[str, Callable[[pd.DataFrame, str], pd.DataFrame]]] = []
            for col, ops_list in self.params.CategoricalColumns.items():
                ops: list[Any] | Any = ops_list if isinstance(ops, list) else [ops_list]
                for op in ops:
                    if op in CAT_OPS_REGISTRY:
                        cat_steps.append((op, CAT_OPS_REGISTRY[op]))
                        log.info(f"Added Step in cat_steps: {op}")
                    else:
                        log.warning(f"Step doesn't exist in registry: {op}")
            
            log.info("Completed Categorical Steps Building")

            # Creating Numerical and Categorical Pipelines
            num_pipeline = Pipeline(steps=num_steps)
            cat_pipeline = Pipeline(steps=cat_steps)

            # Adding them under a single column transformer
            preprocessor = ColumnTransformer(
                transformers=[
                    ("num_ops", num_pipeline, self.params.NumericalColumns), 
                    ("cat_ops", cat_pipeline, self.params.CategoricalColumns),
                    ],
                    remainder="passthrough"
            )
            log.info("Completed the Column Transformer part")
            # Finally Compiling under a single Pipeline
            self.pipeline = Pipeline(
                steps=[
                    ("preprocessor", preprocessor)
                ]
            )
            log.info("Preprocessor Pipeline was Built Sucessfully")
        except Exception as e:
            log.error(e)
            raise Exception(e)

