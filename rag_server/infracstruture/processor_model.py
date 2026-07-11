from functools import lru_cache
from typing import Any




@lru_cache(maxsize=1)
def get_processor_model() -> Any:
    from transformers import BlipProcessor

    return BlipProcessor.from_pretrained("Salesforce/blip-image-captioning-base")
