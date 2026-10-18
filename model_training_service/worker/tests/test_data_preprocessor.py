from models import DataParams
from data_pipeline import DataPreporcessorCosntructor
from logger import test_logger as log
from sklearn.pipeline import Pipeline
from unittest import TestCase
import pandas as pd
import numpy as np

# The main class for testing the preprocessor
class TestDataPreprocessorConstructor(TestCase):
    
    @classmethod
    def setUpClass(cls):
        """Setting Up Class Test Cases"""

        # Local Var for Data cuz it will not be used outside this function
        Data: DataParams = DataParams(**{
            "ID": 1,
            "TestSplit": 0.2,

            "Data":[["1", "a", "c", "e"], ["2", "b", "d", "f"]],
            "ColumnList": ["col1", "col2", "col3", "col4"],
            "NumericalColumns": ["col1"],
            "CategoricalColumns": ["col2", "col3", "col4"],
            "NumericalOperations": ["standard_scaler"],
            "CategoricalOperations": {
                "col2": ["one_hot_encoding"],
                "col3": ["label_encoding"],
                "col4": ["dummy_encoding"]
            }
        })

        # This one is a class global variable
        cls.DPC = DataPreporcessorCosntructor(
            params = Data
        )

        # This one is the expected result for load data test
        cls.output1 = pd.DataFrame(
            data=[[1, "a", "c", "e"], [2, "b", "d", "f"]], 
            columns=["col1", "col2", "col3", "col4"]
            )

    def test_load_data(self):
        """This function will test the load data class method"""
        log.info("Testing DPC.load_data().....")
        try:    
            self.DPC.load_data()
            self.assertIsInstance(self.DPC.df, pd.DataFrame)
            pd.testing.assert_frame_equal(self.DPC.df, self.output1)
        except Exception as e:
            log.error(f"DPC Test - Load Data Error: {e}")
            raise Exception(e)
        
        finally:
            log.info("Completed the Load Data Test")
    

    def test_build_pipeline(self):
        """This function will test the build pipeline class method"""
        try:
            log.info("Testing DPC.build_pipeline().....")
            self.DPC.build_pipeline()
            self.assertIsInstance(self.DPC.pipeline, Pipeline)

            self.DPC.pipeline.set_output(transform="pandas")
            result = self.DPC.pipeline.fit_transform(self.output1)
            self.assertIsInstance(result, pd.DataFrame)

        except Exception as e:
            log.error(f"DPC Error - Build Pipeline Error: {e}")
            raise Exception(e)
        finally:
            log.info(f"Build Pipleine Test Completed: {result}")

