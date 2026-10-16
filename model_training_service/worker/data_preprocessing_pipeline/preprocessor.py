from models import DataParams
from logger import dp_logger as log
import pandas as pd

# This object will be the main construct for the data preprocessing pipeline
class DataPreporcessorCosntructor:
    def __init__(self, params: DataParams) -> None:
        # Initializing the 
        self.params = params 
        self.df: pd.DataFrame | None = None
    
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
