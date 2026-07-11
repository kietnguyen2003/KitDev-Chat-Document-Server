from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from .features.indexing.router import router as document_router
from .features.ask.router import router as ask_router

app = FastAPI(title="KitDev Rag")

app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "http://localhost:3000",
        "http://localhost:5173",
    ],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(document_router)
app.include_router(ask_router)
