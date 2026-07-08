# Take-Home Topshiriq — Middle Go Backend Developer

**Loyiha:** Haydovchilar ro'yxati servisi (Driver Registry Service)
**Taxminiy hajm:** 4–6 soatlik ish
**Muddat:** topshiriqni ko'rib chiqib, **qachon topshirishingizni o'zingiz ayting** — va o'zingiz belgilagan muddatda tayyorlab bering. Muddatga rioya qilish ham baholanadigan jihatlardan biri.
**Topshirish:** GitHub/GitLab repozitoriy havolasi (public yoki bizga access berilgan private)

---

## 1. Vazifa tavsifi

Taksi kompaniyasi uchun haydovchilarni ro'yxatga olish va boshqarish bo'yicha kichik backend servis yozing. Servis REST API orqali haydovchilar ustida CRUD amallarini bajaradi va ma'lumotlarni PostgreSQL'da saqlaydi.

Bu real production servisining kichraytirilgan modeli. Bizni funksionalning o'zi emas, **kodni qanday tashkil qilganingiz va muhandislik qarorlaringiz** ko'proq qiziqtiradi.

---

## 2. Funksional talablar

### 2.1. Ma'lumotlar modeli

`Driver` obyekti quyidagi maydonlarga ega:

| Maydon | Turi | Izoh |
|---|---|---|
| `id` | UUID | Server tomonida generatsiya qilinadi |
| `full_name` | string | Majburiy, 3–100 belgi |
| `phone` | string | Majburiy, unikal, `+998XXXXXXXXX` formatida |
| `license_number` | string | Majburiy, unikal |
| `car_model` | string | Majburiy |
| `car_plate` | string | Majburiy, unikal |
| `status` | enum | `active`, `inactive`, `blocked`. Default: `active` |
| `created_at` | timestamp | Server tomonida |
| `updated_at` | timestamp | Server tomonida |

### 2.2. API endpointlar

| Metod | Yo'l | Tavsif |
|---|---|---|
| `POST` | `/api/v1/drivers` | Yangi haydovchi yaratish |
| `GET` | `/api/v1/drivers/{id}` | Bitta haydovchini olish |
| `GET` | `/api/v1/drivers` | Ro'yxat (pagination + filtr) |
| `PATCH` | `/api/v1/drivers/{id}` | Qisman yangilash |
| `DELETE` | `/api/v1/drivers/{id}` | O'chirish (soft delete afzal, lekin majburiy emas) |
| `PATCH` | `/api/v1/drivers/{id}/status` | Statusni o'zgartirish |
| `GET` | `/healthz` | Health check (DB ulanishini ham tekshirsin) |

### 2.3. Ro'yxat endpointi talablari

`GET /api/v1/drivers` quyidagi query parametrlarni qo'llab-quvvatlashi kerak:

- `page` va `limit` (default: `page=1`, `limit=20`, maksimal `limit=100`)
- `status` bo'yicha filtr (masalan `?status=active`)
- `search` — `full_name` yoki `phone` bo'yicha qidiruv (qisman moslik yetarli)

Javob formati:

```json
{
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 137
  }
}
```

### 2.4. Xatolar formati

Barcha xatolar yagona JSON formatda qaytishi kerak, masalan:

```json
{
  "error": {
    "code": "DRIVER_NOT_FOUND",
    "message": "driver with id ... not found"
  }
}
```

To'g'ri HTTP status kodlar ishlatilishi kerak: `400` (validatsiya), `404` (topilmadi), `409` (unikal maydon konflikti — masalan, telefon raqami band), `500` (ichki xato). Ichki xatolarda DB xabari yoki stack trace mijozga chiqmasligi kerak.

---

## 3. Texnik talablar

### Majburiy

1. **Go 1.22+**, standart `net/http` yoki istalgan router (`chi`, `gin`, `echo` — tanlov sizniki, README'da nima uchun tanlaganingizni bir jumlada yozing).
2. **PostgreSQL** — `database/sql` + `pgx`, yoki `sqlc`/`squirrel`. ORM (GORM) ishlatish mumkin, lekin xom SQL bilan ishlay olishingizni ko'rsatsangiz — plyus.
3. **Migratsiyalar** — `golang-migrate`, `goose` yoki shunga o'xshash vosita bilan, alohida fayllar ko'rinishida.
4. **Konfiguratsiya** — environment o'zgaruvchilar orqali (`.env.example` fayli bo'lsin). Kodda hardcode qilingan parol/port bo'lmasin.
5. **Docker Compose** — `docker compose up` bitta buyruq bilan servis + DB + migratsiyalarni ko'tarsin.
6. **Strukturali log** — `log/slog`, `zap` yoki `zerolog`. Har bir so'rov uchun metod, yo'l, status va davomiylik loglansin.
7. **Graceful shutdown** — `SIGTERM`/`SIGINT` da joriy so'rovlar tugatilib, DB pool yopilsin.
8. **Context** — so'rov konteksti handler'dan DB qatlamigacha uzatilsin.
9. **Testlar** — kamida quyidagilar:
   - Validatsiya logikasiga unit testlar
   - Kamida bitta handler testi (`httptest` bilan)
   - Table-driven uslub qo'llanilsin
10. **README.md** — ishga tushirish yo'riqnomasi, API misollari (curl yoki .http fayl), qabul qilingan qarorlar va trade-off'lar haqida qisqa bo'lim ("Nimani ataylab soddalashtirdim va nima uchun").

### Ixtiyoriy (bonus, majburiy emas)

Vaqtingiz qolsa, quyidagilardan **bir-ikkitasini** tanlang (hammasini qilish shart emas):

- Integratsion testlar `testcontainers-go` bilan (real PostgreSQL'da)
- `Makefile` (run, test, lint, migrate buyruqlari)
- `golangci-lint` konfiguratsiyasi va toza lint natijasi
- Rate limiting middleware (masalan, IP bo'yicha)
- OpenAPI/Swagger spetsifikatsiyasi
- Prometheus metrikalar endpointi (`/metrics`)
- Soft delete + o'chirilganlarni filtrdan chiqarish

> **Eslatma:** bonus qismlarsiz ham topshiriq to'liq hisoblanadi. Yaxshi bajarilgan majburiy qism > chala qilingan 7 ta bonus.

### 3.1. Challenge — o'z g'oyangiz (majburiy)

Yuqoridagi talablardan tashqari, **o'zingiz o'ylab topgan kamida bitta funksiyani** qo'shing. Bu spetsifikatsiyada yozilmagan, lekin sizningcha real taksi kompaniyasi servisiga foydali bo'ladigan narsa bo'lsin.

Bu qism sizning mahsulotga qiziqishingiz va yangilikka ishtiyoqingizni ko'rsatadi — biz shunchaki "topshiriqni bajaradigan" emas, "loyihaga g'oya qo'shadigan" odamni qidiryapmiz.

**Qoidalar:**

- G'oya o'zingizniki bo'lsin — bonus ro'yxatidagi punktlarni takrorlash hisoblanmaydi
- README'da alohida **"Mening g'oyam"** bo'limi oching: nima qo'shdingiz, **nima uchun** bu foydali deb o'ylaysiz, va uni kelajakda qanday rivojlantirish mumkin
- Ko'lami katta bo'lishi shart emas — kichik, lekin o'ylangan funksiya katta va chala qilinganidan yaxshiroq

**Yo'nalish uchun misollar** (bularni takrorlamang, o'zingiznikini o'ylab toping): haydovchi faoliyati tarixi (audit log), statuslar bo'yicha statistika endpointi, guvohnoma muddati tugashi haqida ogohlantirish mexanizmi va h.k.

---

## 4. Cheklovlar va qoidalar

- Kod **Go idiomatik uslubida** bo'lsin (`gofmt` majburiy).
- Tayyor boilerplate/template repozitoriylardan nusxa ko'chirmang — loyiha strukturasi o'zingizning qaroringiz bo'lsin.
- Commit tarixi mazmunli bo'lsin — bitta "final commit" emas, ish jarayoni ko'rinsin.
- Savollar tug'ilsa — yozing. To'g'ri savol berish ham baholanadi.

---

## 5. Himoya suhbati

Topshiriq qabul qilingandan keyin 20–30 daqiqalik qisqa suhbat o'tkaziladi: kod bo'yicha savollar beriladi, qabul qilgan qarorlaringizni asoslashingiz va kodning ayrim qismlarini joyida o'zgartirishingiz so'ralishi mumkin.

---

## 6. Topshirish tartibi

1. Topshiriqni olgach, **qancha vaqtda topshirishingizni bizga yozib yuboring** (masalan: "3 kun ichida"). Muddatni real baholang — o'zingiz aytgan muddatga rioya qilish muhim.
2. Tayyor bo'lgach, repozitoriy havolasini yuboring.
3. README'da quyidagilar bo'lsin:
   - Ishga tushirish: `docker compose up` + migratsiya buyrug'i (agar alohida bo'lsa)
   - API misollari (curl buyruqlari yoki `requests.http` fayl)
   - "Qarorlar va trade-off'lar" bo'limi (5–10 jumla yetarli)
   - "Mening g'oyam" bo'limi (Challenge qismi uchun)
   - Sarflangan taxminiy vaqt (halol yozing — bu baholashga salbiy ta'sir qilmaydi)

Omad! Savollar bo'lsa, bemalol murojaat qiling.
