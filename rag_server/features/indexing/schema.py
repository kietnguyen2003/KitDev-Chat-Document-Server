
from pydantic import BaseModel, ConfigDict, Field


class APIModel:
    model_config = ConfigDict(populate_by_name=True)

class IndexDocumentRequest(BaseModel):
    document_id: int = Field(...)
    bucket: str = Field(...)
    object: str = Field(...)
    user_id: int = Field(...)
    category: str = Field(...)
    type: str = Field(...)
    
    
class DeleteDocumentRequeset(BaseModel):
    document_id: int = Field(...)

    
class DocumentResponse(BaseModel):
    status: bool