# KitDev Chat Document Server

Hệ thống RAG cho upload tài liệu, lập chỉ mục nội dung và hỏi đáp theo tài liệu bằng mô hình cục bộ.

## Kiến trúc tổng thể

Repo hiện có 3 service chính:

- `gateway` (Go): API gateway, kiểm tra JWT, rate limit, load balancing và proxy request vào các service phía sau.
- `server` (Go): xử lý auth, category, document metadata, quota lưu trữ, Postgres và MinIO.
- `rag_server` (Python/FastAPI): indexing tài liệu, vector store Chroma và hỏi đáp RAG qua Ollama.

Luồng cơ bản:

1. Client gọi vào `gateway`.
2. `gateway` xác thực JWT và forward request sang `server` hoặc `rag_server`.
3. `server` upload file lên MinIO, lưu metadata vào Postgres rồi gọi `rag_server` để index bất đồng bộ.
4. `rag_server` đọc file từ MinIO, tạo embeddings, lưu vào Chroma và callback cập nhật trạng thái document.

### Sơ đồ luồng request

```text
Client
  |
  v
Gateway (JWT + Rate Limit + Load Balancer)
  |-------------------------> Server (Go)
  |                            |
  |                            +--> Postgres
  |                            +--> MinIO
  |                            +--> Rag Server (index/delete vector)
  |
  +-------------------------> Rag Server (ask)
                               |
                               +--> Chroma
                               +--> Ollama
```

### Vai trò từng service

#### `gateway`

- entrypoint public duy nhất cho client
- xác thực JWT cho các route protected
- áp dụng rate limit theo `Client IP`
- phân phối request theo round-robin khi có nhiều upstream `server` hoặc `rag_server`
- gắn các header nội bộ như `KIT-DEV-USER-ID`, `KIT-DEV-USERNAME`, `KIT-DEV-ROLE-ID`

#### `server`

- đăng ký và đăng nhập
- quản lý category
- quản lý metadata document
- quản lý quota lưu trữ
- upload/xóa file trong MinIO
- gọi `rag_server` để index hoặc xóa vector

#### `rag_server`

- index PDF bằng loader + chunking
- index ảnh bằng OCR + caption model
- lưu embeddings vào Chroma
- truy vấn tài liệu và stream câu trả lời qua SSE

### Gateway middleware hiện tại

Thứ tự xử lý trong gateway:

1. CORS middleware
2. Rate limit middleware
3. JWT auth middleware cho protected routes
4. Reverse proxy tới upstream đã chọn bởi load balancer

### Load balancing

- `SERVER_URL` và `RAG_URL` có thể chứa 1 hoặc nhiều upstream
- nhiều upstream được phân tách bằng dấu phẩy
- chiến lược hiện tại là `round-robin`
- nếu upstream lỗi, gateway trả về `502 Bad Gateway`

Ví dụ:

```env
SERVER_URL=http://localhost:8081,http://localhost:8082
RAG_URL=http://localhost:8000,http://localhost:8001
```

### Rate limiting

- áp dụng ở gateway cho toàn bộ request client-facing
- dùng bộ đếm in-memory theo `Client IP`
- mặc định `120 request / 60 giây`
- response có các header:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `X-RateLimit-Reset`
  - `Retry-After` khi vượt ngưỡng

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
SERVER_URL=http://localhost:8081,http://localhost:8082
RAG_URL=http://localhost:8000,http://localhost:8001
RATE_LIMIT_REQUESTS=120
RATE_LIMIT_WINDOW_SECONDS=60
RATE_LIMIT_CLEANUP_SECONDS=300
```

Ghi chú:

- `SERVER_URL` và `RAG_URL` hỗ trợ nhiều upstream, phân tách bằng dấu phẩy.
- Gateway hiện dùng round-robin để cân bằng tải giữa các upstream trong danh sách.
- Rate limit hiện được áp dụng theo `Client IP` cho toàn bộ request đi qua gateway.
- Nếu không truyền nhiều upstream, gateway sẽ hoạt động như reverse proxy một đích bình thường.

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

## Luồng nghiệp vụ chính

### Upload và indexing tài liệu

1. Client gọi `POST /api/documents` qua gateway.
2. Gateway xác thực JWT, áp rate limit, rồi forward sang `server`.
3. `server` kiểm tra user/category, trừ quota, upload file lên MinIO và lưu metadata vào Postgres.
4. `server` gọi bất đồng bộ sang `rag_server` để index.
5. `rag_server` tải file từ MinIO, chunk/embedding và lưu vào Chroma.
6. `rag_server` callback về `server` để cập nhật trạng thái `indexed` hoặc `failed`.

### Hỏi đáp tài liệu

1. Client gọi `POST /api/ask/:language` qua gateway.
2. Gateway xác thực JWT, áp rate limit, chọn một `rag_server` upstream và forward request.
3. `rag_server` truy vấn Chroma theo `category` và `user_id`.
4. `rag_server` gọi Ollama để sinh câu trả lời và stream kết quả về client bằng SSE.

## Ghi chú vận hành

- `gateway` hiện là lớp bảo vệ chính cho API client-facing.
- `server` tin vào các header nội bộ như `KIT-DEV-USERNAME`, vì vậy không nên expose trực tiếp `server` ra public.
- `rag_server` callback cập nhật trạng thái document sau khi index xong.
- Repo hiện đang chứa cả `rag_server/venv` và dữ liệu `rag_server/chroma_db`, phù hợp cho local development nhưng chưa tối ưu để đưa lên git lâu dài.

## Tài liệu tham khảo

- Mô tả API: [api.json](./api.json)
- Sơ đồ và ghi chú: thư mục [Diagram](./Diagram)
