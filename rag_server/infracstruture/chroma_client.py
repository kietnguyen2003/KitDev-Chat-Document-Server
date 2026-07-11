from functools import lru_cache
from typing import Any

from .embedding_model import get_embedding_model

@lru_cache(maxsize=1)
def get_vector_store() -> Any:
    from langchain_chroma import Chroma
    
    return Chroma(
        collection_name="rag_documents",
        embedding_function=get_embedding_model(),
        persist_directory="chroma_db",
    )