# KitDev Chat Document Server

Hệ thống RAG cho upload tài liệu, lập chỉ mục nội dung và hỏi đáp theo tài liệu bằng mô hình cục bộ.

## Kiến trúc

Repo hiện có 3 service chính:

- `gateway` (Go): API gateway, kiểm tra JWT và proxy request vào các service phía sau.
- `server` (Go): xử lý auth, category, document metadata, quota lưu trữ, Postgres và MinIO.
- `rag_server` (Python/FastAPI): indexing tài liệu, vector store Chroma và hỏi đáp RAG qua Ollama.

Luồng cơ bản:

1. Client gọi vào `gateway`.
2. `gateway` xác thực JWT và forward request sang `server` hoặc `rag_server`.
3. `server` upload file lên MinIO, lưu metadata vào Postgres rồi gọi `rag_server` để index bất đồng bộ.
4. `rag_server` đọc file từ MinIO, tạo embeddings, lưu vào Chroma và callback cập nhật trạng thái document.

## Cấu trúc thư mục

```text
.
├── gateway/        # API gateway viết bằng Go
├── server/         # Backend nghiệp vụ viết bằng Go
├── rag_server/     # Dịch vụ RAG viết bằng Python/FastAPI
├── docker-compose.yml
├── api.json        # Mô tả API hiện tại
└── Diagram/        # Tài liệu và sơ đồ tham khảo
```

## Yêu cầu

- Go `1.26.1`
- Python `3.11+`
- Docker và Docker Compose
- Ollama đang chạy local hoặc trong mạng nội bộ mà `rag_server` truy cập được
- Tesseract OCR nếu bạn muốn index ảnh có chữ

## Cổng mặc định

- `gateway`: `http://localhost:8080`
- `server`: `http://localhost:8081`
- `rag_server`: `http://localhost:8000`
- `postgres`: `localhost:5432`
- `minio api`: `http://localhost:9000`
- `minio console`: `http://localhost:9001`

## Biến môi trường

### `gateway/.env`

```env
PORT=8080
JWT_SECRET=my_super_secret_key
SERVER_URL=http://localhost:8081
RAG_URL=http://localhost:8000
```

### `server/.env`

```env
PORT=8081
JWT_SECRET=my_super_secret_key
DB_URL=postgres://user:password@localhost:5432/rag_server?sslmode=disable
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=admin
MINIO_SECRET_KEY=password123
MINIO_BUCKET_NAME=documents
MINIO_USE_SSL=false
RAG_URL=http://127.0.0.1:8000
```

### `rag_server`

`rag_server` hiện chủ yếu đọc config từ biến môi trường hệ thống:

```env
HOST=0.0.0.0
PORT=8000
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_MODEL=qwen2.5:3b
OLLAMA_KEEP_ALIVE=10m
```

Ngoài ra service này còn dùng:

- MinIO ở `localhost:9000`
- Chroma persist tại `rag_server/chroma_db`
- OCR và image caption cho file ảnh

## Khởi động nhanh

### 1. Chạy Postgres và MinIO

```bash
docker compose up -d
```

### 2. Chạy `server`

```bash
cd server
go run ./cmd
```

### 3. Chạy `gateway`

```bash
cd gateway
go run ./cmd
```

### 4. Chạy `rag_server`

Nếu dùng virtualenv có sẵn trong repo:

```bash
cd rag_server
source venv/bin/activate
uvicorn rag_server.main:app --host 0.0.0.0 --port 8000 --reload
```

Nếu tự cài mới:

```bash
cd rag_server
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
uvicorn rag_server.main:app --host 0.0.0.0 --port 8000 --reload
```

## API client-facing qua gateway

Base URL:

```text
http://localhost:8080
```

### Public

- `POST /api/auth/register`
- `POST /api/auth/sign-in`

### Protected

- `POST /api/categories`
- `GET /api/categories`
- `POST /api/documents`
- `DELETE /api/documents`
- `GET /api/documents/:cateName`
- `PUT /api/documents/:id`
- `POST /api/ask/:language`

`language` hiện hỗ trợ:

- `vi`
- `en`

Chi tiết request/response được mô tả trong [api.json](./api.json).

## Các service nội bộ

### `server`

Vai trò:

- đăng ký và đăng nhập
- quản lý category
- upload file lên MinIO
- lưu metadata document
- gọi `rag_server` để index hoặc xóa vector

### `rag_server`

Vai trò:

- index PDF bằng loader và chunking
- index ảnh bằng OCR + image caption
- lưu embeddings vào Chroma
- truy vấn vector store và stream câu trả lời bằng SSE

Các route nội bộ chính:

- `POST /api/documents`
- `POST /api/documents/index`
- `DELETE /api/documents`
- `POST /api/ask/{language}`

## Ghi chú vận hành

- `gateway` hiện là lớp bảo vệ chính cho API client-facing.
- `server` tin vào các header nội bộ như `KIT-DEV-USERNAME`, vì vậy không nên expose trực tiếp `server` ra public.
- `rag_server` callback cập nhật trạng thái document sau khi index xong.
- Repo hiện đang chứa cả `rag_server/venv` và dữ liệu `rag_server/chroma_db`, phù hợp cho local development nhưng chưa tối ưu để đưa lên git lâu dài.

## Tài liệu tham khảo

- Mô tả API: [api.json](./api.json)
- Sơ đồ và ghi chú: thư mục [Diagram](./Diagram)
