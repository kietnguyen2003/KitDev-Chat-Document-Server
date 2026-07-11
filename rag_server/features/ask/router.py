from fastapi import APIRouter, Depends
from fastapi.responses import StreamingResponse

from typing import Literal
import json

from langchain_core.documents import Document
from langchain_chroma import Chroma
from langchain_ollama import ChatOllama

from .schema import AksDocument
from ...infracstruture.chroma_client import get_vector_store
from ...infracstruture.ollama_client import get_llm

router = APIRouter(tags=["Ask"])


@router.post("/api/ask/{language}")
async def ask_model_stream(
    language: Literal["vi", "en"],
    req: AksDocument,
    vector_store: Chroma = Depends(get_vector_store),
    llm: ChatOllama = Depends(get_llm),
):
    result = vector_store._collection.get(
        include=["documents", "metadatas"]
    )
    print(result["metadatas"])
    print(req)
    retrieved = vector_store.similarity_search_with_score(
        query=req.msg,
        k=8,
        filter={  # type: ignore
            "$and": [
                {"category": req.category},
                {"user_id": req.user_id},
            ]
        },
    )
    
    print(retrieved)

    docs = []

    for doc, score in retrieved:
        print(f"score={score:.4f} page={doc.metadata.get('page')}")

        # score càng nhỏ càng giống
        docs.append(doc)

    docs = dedupe_docs(docs)

    context, citations = build_context(docs)

    prompt = build_prompt(
        question=req.msg,
        context_text=context,
        language=language,
    )

    def generate():

        buffer = ""

        for chunk in llm.stream(prompt):

            if not isinstance(chunk.content, str):
                continue

            buffer += chunk.content

            if buffer.endswith((".", "!", "?", "\n")):

                yield (
                    "event: token\n"
                    f"data: {json.dumps({'text': buffer}, ensure_ascii=False)}\n\n"
                )

                buffer = ""

        if buffer:

            yield (
                "event: token\n"
                f"data: {json.dumps({'text': buffer}, ensure_ascii=False)}\n\n"
            )

        yield (
            "event: sources\n"
            f"data: {json.dumps({'citations': citations}, ensure_ascii=False)}\n\n"
        )

        yield "event: done\ndata:[DONE]\n\n"

    return StreamingResponse(
        generate(),
        media_type="text/event-stream",
    )
    
    
def build_context(
    docs: list[Document],
) -> tuple[str, list[dict]]:

    context_parts = []

    citations = []

    for idx, doc in enumerate(docs, start=1):

        metadata = doc.metadata or {}

        citation = {
            "citation_id": idx,
            "document_id": metadata.get("document_id"),
            "page": metadata.get("page"),
            "category": metadata.get("category"),
            "object": metadata.get("object"),
        }

        citations.append(citation)

        content = doc.page_content.strip()

        context_parts.append(
f"""
==============================
Chunk [{idx}]

Document:
{metadata.get("object")}

Page:
{metadata.get("page")}

Content:

{content}
"""
)

    return "\n".join(context_parts), citations


def dedupe_docs(
    docs: list[Document],
):

    seen = set()

    unique = []

    for doc in docs:

        key = (
            doc.metadata.get("object"),
            doc.metadata.get("page"),
            hash(doc.page_content),
        )

        if key in seen:
            continue

        seen.add(key)

        unique.append(doc)

    return unique

def build_prompt(
    question: str,
    context_text: str,
    language: str,
):

    return f"""
/no_think

You are a Retrieval-Augmented Generation assistant.

Rules:

- Only answer using the CONTEXT.
- Never use outside knowledge.
- Answer directly.
- Explain naturally.
- Do not say:
    - According to the context...
    - According to the paper...
    - Refer to...
    - See Figure...

When information comes from multiple chunks,
combine them into one coherent answer.

If the question asks for a summary,
summarize the whole document.

If the answer cannot be found, reply exactly:

I can not find any information in this document.

Always append citations:

[1]
[2]
[1][3]

Language:

{language}

==============================

CONTEXT

{context_text}

==============================

QUESTION

{question}

==============================

ANSWER
""".strip()