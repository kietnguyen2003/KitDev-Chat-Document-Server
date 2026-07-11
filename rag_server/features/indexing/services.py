
from .schema import IndexDocumentRequest, DocumentResponse, DeleteDocumentRequeset
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_core.documents import Document
from minio import Minio
from transformers import BlipProcessor, BlipForConditionalGeneration

from typing import cast
from PIL import Image
import torch
import pytesseract



from ...shared.helper import update_document_status


from langchain_chroma import Chroma
from langchain_community.document_loaders import PyPDFLoader, UnstructuredPDFLoader
import os
import tempfile


async def index_pdf_document(
    req: IndexDocumentRequest,
    vector_store: Chroma,
    storage: Minio
) -> DocumentResponse:
    print(req)

    suffix = f".pdf"

    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        tmp_path = tmp.name

    try:
        storage.fget_object(req.bucket, req.object, tmp_path)

        loader = UnstructuredPDFLoader(tmp_path)
        docs = loader.load()

        for doc in docs:
            doc.metadata["document_id"] = req.document_id
            doc.metadata["category"] = req.category
            doc.metadata["bucket"] = req.bucket
            doc.metadata["user_id"] = req.user_id
            doc.metadata["object"] = req.object

        vector_store._collection.delete(
            where={"document_id": req.document_id}
        )

        chunks = chunk_document(documents=docs)

        ids = [
            f"{req.document_id}:{i}"
            for i in range(len(chunks))
        ]

        vector_store.add_documents(chunks, ids=ids)

        await update_document_status(
            base_url="http://localhost:8081",
            doc_id=req.document_id,
            status="indexed",
        )

        return DocumentResponse(status=True)

    except Exception as e:
        print("index failed:", e)

        try:
            await update_document_status(
                base_url="http://localhost:8081",
                doc_id=req.document_id,
                status="failed",
            )
        except Exception as update_err:
            print("update status failed:", update_err)

        return DocumentResponse(status=False)

    finally:
        if os.path.exists(tmp_path):
            os.remove(tmp_path)   
                               

def chunk_document(documents: list[Document]):
    splitter=RecursiveCharacterTextSplitter(
        separators=[
            "\n\n",
            "\n",
            ". ",
            "? ",
            "! ",
            " "
        ],
        chunk_size=1000,
        chunk_overlap=100
    )
    
    return splitter.split_documents(documents) 


async def delete_document(
    req: DeleteDocumentRequeset,
    vector_store: Chroma,
) -> DocumentResponse:
    try:
        vector_store._collection.delete(
            where={
                "document_id": req.document_id
            }
        )
        
        return DocumentResponse(status=True)
    except:
        return DocumentResponse(status=False)


async def index_image_document(
    req: IndexDocumentRequest,
    processor:BlipProcessor,
    caption_model: BlipForConditionalGeneration,
    vector_store: Chroma,
    storage: Minio
) -> DocumentResponse:
    print(req)

    suffix = f".jpg"

    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        tmp_path = tmp.name

    try:
        # Lấy ảnh từ minio
        storage.fget_object(req.bucket, req.object, tmp_path)
        # xử lý ảnh
        docs = image_to_document(tmp_path, req.user_id, req.category, req.document_id, req.object,processor, caption_model)
        print(docs.page_content)
        vector_store.add_documents(
            [docs],
            ids=[f"image:{req.document_id}:0"],
        )
        await update_document_status(
            base_url="http://localhost:8081",
            doc_id=req.document_id,
            status="indexed",
        )
        
        return DocumentResponse(status=True)
 
        
    except Exception as e:
        print("index failed:", e)

        try:
            await update_document_status(
                base_url="http://localhost:8081",
                doc_id=req.document_id,
                status="failed",
            )
        except Exception as update_err:
            print("update status failed:", update_err)

        return DocumentResponse(status=False)
    
    
def image_to_document(
    image_path: str,
    user_id: int,
    category: str,
    image_id: int,
    object_key: str,
    processor:BlipProcessor,
    caption_model: BlipForConditionalGeneration,
) -> Document:
    # chuyển thành orc text => xử lý chữ trong ảnh
    ocr_text = ocr_image(image_path)
    # chuyển thành caption text => xử lý ảnh trong ảnh
    caption = caption_image(image_path,processor,caption_model)
    # ghép 2 cái lại
    page_content = f"""
OCR Text:
{ocr_text if ocr_text else "No text detected."}

Image Description:
{caption}
""".strip()

    return Document(
        page_content=page_content,
        metadata={
            "type": "image",
            "user_id": user_id,
            "category": category,
            "document_id": image_id,
            "object": object_key,
        },
    )
    
def ocr_image(image_path: str) -> str:
    image = Image.open(image_path).convert("RGB")
    text = pytesseract.image_to_string(image, lang="vie+eng")
    return text.strip()

def caption_image(
    image_path: str,
    processor: BlipProcessor,
    caption_model: BlipForConditionalGeneration,
) -> str:
    image = Image.open(image_path).convert("RGB")

    inputs = cast(
        dict[str, torch.Tensor],
        processor(images=image, return_tensors="pt")
    )

    pixel_values = inputs["pixel_values"]

    with torch.no_grad():
        output = caption_model.generate(
            pixel_values=pixel_values,  # type: ignore[arg-type]
            max_new_tokens=80,
        )

    return processor.decode(
        output[0],
        skip_special_tokens=True,
    ).strip()  