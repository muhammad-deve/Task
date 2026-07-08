# Driver Registry Service

Taxi driver management service built with Go, PostgreSQL, and Echo framework.

## Ishga tushirish (Quick Start)

### Requirements
- Docker & Docker Compose (recommended)
- OR Go 1.22+ (for local development)

### Docker orqali ishga tushirish (Recommended)

```bash
# Copy environment file
cp app/.env.example app/.env

# Start all services (app + postgres + redis + minio)
docker compose up -d

# Loglarni ko'rish
docker compose logs -f app
```

Service starts on `http://localhost:8080`

### Migratsiya (Migrations)

Migratsiyalar dastur yonganda avtomatik ishga tushadi (`main.go` ichida). Lekin qo'lda ishga tushirish uchun:
```bash
cd app
make migrate-up
```

## API misollari (API Examples)

Ilovada barcha API misollarini o'z ichiga olgan `requests.http` fayli mavjud. Uni VS Code'dagi REST Client yordamida ishlatishingiz yoki Swagger UI orqali tekshirishingiz mumkin:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Health Check**: `http://localhost:8080/healthz`

*CURL misollari uchun loyihada mavjud `requests.http` fayliga qarang.*

## Qarorlar va trade-off'lar

Loyiha "Clean Architecture" tamoyillari asosida Handler, Service va Repository qatlamlariga ajratilgan holda qurildi. Asosiy freymvork sifatida tezkor va sodda bo'lgan Echo tanlandi, ma'lumotlar bazasi bilan ishlashda esa ORM o'rniga xavfsiz va tezkor `sqlc` ishlatildi. Xavfsizlikni ta'minlash maqsadida API JWT token orqali himoyalandi va parollar heshlab saqlandi. Unikal maydonlar (telefon, haydovchilik guvohnomasi) tranzaksiyadan oldin EXISTS orqali tekshirilib, bazada ham unikal indekslar bilan mustahkamlandi. Tizimda rate limiter (xotiraga asoslangan) va "soft delete" (ma'lumotlarni o'chirmasdan faqat belgilab qo'yish) logikasi qo'llanildi. Qidiruv qismi sodda bo'lishi uchun PostgreSQL'ning `ILIKE` funksiyasi yordamida amalga oshirildi, katta masshtabda buni Full-Text Search'ga o'zgartirish kerak bo'ladi. Shuningdek, loyihada xatoliklarni qaytarish (Error Handling) barcha API'lar uchun yagona standart (HTTP Code bilan) ko'rinishga keltirildi.

## Mening g'oyam

Challenge qismi uchun **Driver Status Statistics** (Haydovchilarning holati bo'yicha statistika) API'sini qo'shdim.
Bu tizimda qancha haydovchi faol (`active`), nofaol (`inactive`) yoki bloklangan (`blocked`) ekanligini hisoblab beradi.

**Nima uchun bu muhim?**
Taksi kompaniyalari operatorlari kunlik qancha haydovchi liniyada ekanligini yoki nechta haydovchi qoidabuzarlik uchun bloklanganini doimiy monitoring qilib borishlari kerak. Barcha haydovchilar ro'yxatini to'liq yuklab olmasdan, faqatgina holati bo'yicha sonini olish server va baza resurslarini tejaydi hamda dashboard'lar uchun juda qulay hisoblanadi.

**API:**
`GET /api/v1/drivers/stats/active?status=blocked`

## Sarflangan taxminiy vaqt

- Boshlang'ich struktura va Docker muhit: 30 daqiqa
- CRUD API va JWT Auth: 1.5 soat
- Validation va Error handling (kodlarni standartlashtirish): 1 soat
- "Mening g'oyam" (Statistika): 45 daqiqa
- Swagger doc, Testing va Debugging: 1 soat
**Umumiy sarflangan vaqt: ~4.5 soat**
