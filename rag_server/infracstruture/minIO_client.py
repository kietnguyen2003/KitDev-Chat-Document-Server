from minio import Minio
from functools import lru_cache

@lru_cache(maxsize=1)
def get_storage():
    return Minio(
    "localhost:9000",
    access_key="admin",
    secret_key="password123",
    secure=False,
)