from functools import lru_cache
from typing import Any




@lru_cache(maxsize=1)
def get_caption_model() -> Any:
    from transformers import BlipForConditionalGeneration

    return BlipForConditionalGeneration.from_pretrained("Salesforce/blip-image-captioning-base")
