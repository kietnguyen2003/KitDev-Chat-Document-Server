import httpx

async def update_document_status(
    base_url: str,
    doc_id: int,
    status: str,
):
    async with httpx.AsyncClient(timeout=10) as client:
        res = await client.put(
            f"{base_url}/api/documents/{doc_id}",
            json={
                "status": status
            },
            headers={"KIT-DEV-USERNAME": "ngkiet2611@gmail.com",}
        )

        res.raise_for_status()