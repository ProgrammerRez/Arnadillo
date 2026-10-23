from models import TrainingParameters, DataParams
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

        cls.traditional_training_parameters = {}
        cls.deep_training_parameters = {}
    



if __name__=='__main__':
    data = load_iris()

    df = pd.DataFrame(data=data.data, columns=data.feature_names)
    df["target"] = data.target
    df["target_names"] = df["target"].map({
        0: "setosa",
        1: "versicolor",
        2: "virginica"
    })

    data_params = DataParams(**{
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

    dpc = DataPreporcessorCosntructor(data_params)
    dpc.load_data()

    dpc.build_pipeline()
    print(dpc.pipeline)
    print(dpc.pipeline.fit_transform(df)[120])