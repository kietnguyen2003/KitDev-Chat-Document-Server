
from pydantic import BaseModel, ConfigDict, Field


class APIModel:
    model_config = ConfigDict(populate_by_name=True)

class AksDocument(BaseModel):
    msg: str = Field(...)
    document_id: int = Field(...)
    user_id: int = Field(...)
    category: str = Field(...)