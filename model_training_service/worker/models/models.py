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

