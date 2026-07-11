from functools import lru_cache
from pathlib import Path
import os


class Settings:
    def __init__(self) -> None:
        base_dir = Path(__file__).resolve().parent.parent
        self.host = os.getenv("HOST", "0.0.0.0")
        self.port = int(os.getenv("PORT", "8002"))
                
        self.document_root = Path(os.getenv("DOCUMENT_ROOT", str(base_dir / "documents")))
        self.internal_token = os.getenv("INTERNAL_TOKEN", "dev_internal_secret")
        self.folder_name=os.getenv("FOLDER_NAME_DOWNLOAD", "documents")

        
        self.ollama_base_url = os.getenv("OLLAMA_BASE_URL", "http://ollama:11434")
        self.ollama_model = os.getenv("OLLAMA_MODEL", "qwen2.5:3b")
        self.ollama_credit= (os.getenv("OLLAMA_CREDIT", "q384"))
        self.ollama_temperature= (os.getenv("OLLAMA_TEMPERATURE", "temperature"))
        self.ollama_keep_alive= os.getenv("OLLAMA_KEEP_ALIVE", "10m")

@lru_cache(maxsize=1)
def get_settings() -> Settings:
    return Settings()
