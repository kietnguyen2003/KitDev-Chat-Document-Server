from functools import lru_cache
from typing import Any

@lru_cache(maxsize=1)
def get_embedding_model() -> Any:
    from langchain_community.embeddings import HuggingFaceEmbeddings
    return HuggingFaceEmbeddings(model_name="intfloat/multilingual-e5-base2")
