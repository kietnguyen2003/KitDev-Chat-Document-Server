from functools import lru_cache
import os
from langchain_ollama import ChatOllama

from ..core.config import get_settings

@lru_cache(maxsize=1)
def get_llm() -> ChatOllama:
    setting = get_settings()
    os.environ["OLLAMA_HOST"] = setting.ollama_base_url
    return ChatOllama(
        model=setting.ollama_model,
        num_predict=2048,
        temperature=0.2,
        keep_alive=setting.ollama_keep_alive
    )