response có dạng
{
    code
    messeage
    data
}

auth

POST SIGN-IN    api/auth/sign-in
{
    username
    password
}
data
{
    access_token
    refresh_token
    expires_in
    user{
        fullname
        role
    }
}

POST REGISTER    api/auth/register
{
    username
    password
    fullname
}
data
{
    access_token
    request_token
    expires_in
    user{
        fullname
        role
    }
}

GET CATEGORY                api/categories
data
{
    [
       {
            category_id,
            category_name
        },
        {
            category_id,
            category_name
        },
        {
            category_id,
            category_name
        },
        {
            category_id,
            category_name
        },
    ]
}

POST CREATE CATEGORY        api/categories
req
{
    category_name
}
data
{
    category_id,
    category_name
},

PUT RENAME CATEGORY         api/categories/{category_id}
req
{
    category_new_name
}

DELETE CATEGORY             api/categories/{category_id}

GET DOCUMENT_BY_CATEGORY    api/document/{category_id}
data
{
   [
        {
            "document_id": "doc_001",
            "document_name": "RAG Notes.pdf",
            "document_type": "application/pdf",
            "status": "indexed",
            "size_bytes": 2458120,
            "page_count": 24,
            "chunk_count": 86,
            "created_at": "2026-06-06T10:00:00Z"
        },
        ...
   ]
}

GET DOCUMENT_DETAIL         api/document/{document_id}
{
    "document_id": "doc_001",
    "category_id": "cat_001",
    "document_name": "RAG Notes.pdf",
    "document_type": "application/pdf",
    "status": "indexed",
    "size_bytes": 2458120,
    "page_count": 24,
    "chunk_count": 86,
    "indexed_at": "2026-06-06T10:05:00Z",
    "created_at": "2026-06-06T10:00:00Z"
}

POST UPLOAD_DOCUMENT        api/document/{category_id}/documents
Content-Type: multipart/form-data
req
{
    file=
}
data
{
    "document_id": "doc_001",
    "category_id": "cat_001",
    "document_name": "RAG Notes.pdf",
    "document_type": "application/pdf",
    "status": "processing",
    "size_bytes": 2458120,
    "page_count": null,
    "chunk_count": null,
    "created_at": "2026-06-06T10:00:00Z"
}


POST RE-INDEX DOCUMENT       api/document/{category_id}/documents

POST ASK                     api/documents/{document_id}/ask
req
{
    messeage
}
res
streamdata{}

DELETE DOCUMENT             api/document/{document_id}