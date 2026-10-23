from datetime import datetime
from pydantic import BaseModel
from typing import Tuple, Dict, Any, Optional

# This model is used to describe a Job
class Job(BaseModel):
    ID: int
    StartTime: Optional[datetime] = None
    EstimatedResources: Optional[Dict[str, Any]] = None
    DataParameters: Optional[Dict[str, Any]] = None
    TrainingParameters: Optional[Dict[str, Any]] = None
    EvalParameters: Optional[Dict[str, Any]] = None
    DoOverParameters: Optional[Dict[str, Any]] = None

# This Model is for the Data Parameters
class DataParams(BaseModel):
    NumericalColumns:    Optional[list[str]] = None
    CategoricalColumns: Optional[list[str]] = None
    TargetLabel: Optional[str] = None
    ColumnList: Optional[list[str]] = None
    Data: Optional[list[list[Any]]] = None
    NumericalOperations: Optional[list[str]] = None
    CategoricalOperations: Optional[Dict[str, list[str]]] = None
    ID: int
    TestSplit: float


# This model is for the Training Parameters
class TrainingParameters(BaseModel):
    ModelType: Optional[str] = None
    # If deep it's going to be just one, while many for traditional
    ModelName: Optional[list[str]] = None
    ModelParams: Optional[Dict[str, Any]] = None