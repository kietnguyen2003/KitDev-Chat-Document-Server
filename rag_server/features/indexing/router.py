from fastapi import APIRouter, Depends

from .schema import DocumentResponse, IndexDocumentRequest, DeleteDocumentRequeset
from .services import index_pdf_document, delete_document, index_image_document

from ...infracstruture.chroma_client import get_vector_store
from ...infracstruture.minIO_client import get_storage
from ...infracstruture.caption_model import get_caption_model
from ...infracstruture.processor_model import get_processor_model

router = APIRouter(tags=["indexing"])

@router.post("/api/documents", response_model=DocumentResponse)
@router.post("/api/documents/index")
async def indexing_document(
    req: IndexDocumentRequest,
    processor = Depends(get_processor_model),
    caption_model = Depends(get_caption_model),
    vector_store = Depends(get_vector_store),
    storage = Depends(get_storage)
) -> DocumentResponse:
    print("Type of file: ", req.type)
    if req.type == "application/pdf":
        return await index_pdf_document(
            req=req,
            vector_store=vector_store,
            storage=storage
        )
    else:
        return await index_image_document(
            req=req,
            processor=processor,
            caption_model=caption_model,
            vector_store=vector_store,
            storage=storage
        )
        


@router.delete("/api/documents", response_model=DocumentResponse)
async def deleting_document(
    req: DeleteDocumentRequeset,
    vector_store = Depends(get_vector_store)
) -> DocumentResponse:
    return await delete_document(
        req=req,
        vector_store=vector_store,
    )

