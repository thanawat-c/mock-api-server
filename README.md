# 🚀 Mock API Generator Server

A modern, fast, and self-hosted **Full RESTful Mock API Server** built with **Go (Gin)** and **Nuxt 3 (Vue 3 + Element Plus)**. 

Allows Frontend & Mobile developers to build, test, and prototype UIs immediately without waiting for the Backend team to implement APIs.

> 🇹🇭 **สรุปโปรเจกต์ (Project Overview):**  
> แพลตฟอร์มจำลอง RESTful Mock API แบบ Self-hosted ที่พัฒนาด้วย Go (Gin) และ Nuxt 3 ช่วยให้นักพัฒนาฝั่ง Frontend และ Mobile สามารถสร้าง Mock Endpoints, กำหนด Schema, สุ่มข้อมูลจำลองด้วย Faker, และทดสอบยิง API ได้ทันทีแบบไม่ต้องรอ Backend พร้อมรองรับการจำลอง Network Latency, HTTP Status, Custom Headers ตลอดจน Import/Export ข้อมูลด้วย JSON

---

## ✨ Features

- **Full RESTful Resource Handling**:
  - Automatically provisions full CRUD endpoints for any resource path:
    - `GET /mock/:path` - List records with built-in **Pagination** (`?page=1&limit=10`) and **Search filter** (`?search=query`)
    - `GET /mock/:path/:id` - Retrieve individual record by ID
    - `POST /mock/:path` - Insert new records with schema validation
    - `PUT /mock/:path/:id` - Update existing records
    - `DELETE /mock/:path/:id` - Delete records
- **Visual Schema Builder & Realistic Data Generator**:
  - Define custom schemas visually with various Faker data types: Full Name, Email, Avatar URL, UUID, Phone, Date, Price, Address, Boolean, and more.
- **Network & Status Simulation**:
  - **Latency Simulation**: Custom delay slider (0 - 3000 ms) to test loading spinners.
  - **HTTP Status Simulation**: Simulate error codes (400, 401, 404, 500) via Query parameter `?mock_status=...` or Header `x-mock-status: ...`.
  - **Custom Response Headers**: Simulate response headers such as `X-Total-Count`, `Cache-Control`, and `X-API-Version`.
- **Interactive Postman-like Test Console**:
  - Send live test requests with custom methods, payloads, and simulated headers directly from the browser.
  - Inspect Status codes, Response Headers, and JSON output with single-click copy buttons.
- **Data Management & Reset**:
  - Visual table view to edit individual records or generate additional fake records on the fly.
  - One-click reset to restore default mock data.
- **JSON Import & Export**:
  - Export single or all endpoints into JSON format for sharing with teammates or backups.
  - Import JSON files to instantly replicate mock environments.

---

## 🛠️ Tech Stack

- **Backend**: Go 1.22+, Gin Web Framework, GORM, PostgreSQL
- **Frontend**: Nuxt 3, Vue 3, TypeScript, TailwindCSS, Element Plus, Lucide Icons, Pinia
- **Database**: PostgreSQL (with JSONB support for dynamic schemas)

---

## 🚀 Getting Started

### 1. Prerequisites
- [Go](https://golang.org/) (version 1.20+)
- [Node.js](https://nodejs.org/) (version 18+) & `pnpm` / `npm`
- [PostgreSQL](https://www.postgresql.org/) database running locally or remotely

---

### 2. Backend Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Create your `.env` configuration from the example:
   ```bash
   cp .env.example .env
   ```
3. Edit `.env` with your PostgreSQL credentials:
   ```env
   DB_URL=host=localhost user=your_user password=your_password dbname=mockapi_db port=5432 sslmode=disable
   ```
4. Run database migrations and start the backend API server:
   ```bash
   go run cmd/api/main.go
   ```
   *The backend will run on `http://localhost:8080`*

---

### 3. Frontend Setup

1. Open a new terminal and navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   pnpm install
   # or
   npm install
   ```
3. Start the Nuxt 3 development server:
   ```bash
   npm run dev
   ```
   *The frontend dashboard will be available at `http://localhost:3000`*

---

## 📖 How to Use

1. **Create Endpoint**:
   - Go to `http://localhost:3000` and click **"New Endpoint"**.
   - Enter your Resource Path (e.g. `/users`), select initial record count, delay, and optional custom response headers.
   - Define your schema fields and click **"Generate Mock API"**.
2. **Consume in Frontend / Apps**:
   - Call your endpoint directly in your frontend application:
     ```javascript
     const res = await fetch('http://localhost:8080/mock/users?page=1&limit=10')
     const data = await res.json()
     ```
3. **Test in Browser**:
   - Open the **Test Console** tab on any endpoint detail page to test `GET`, `POST`, `PUT`, `DELETE` operations interactively.

---

## 📄 License
This project is open-source and free to use for personal and team development.
