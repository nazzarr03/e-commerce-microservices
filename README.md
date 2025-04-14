# 💲 E-Commerce Microservices Architecture

Bu proje, mikroservis mimarisi ile oluşturulmuş bir e-ticaret sistemidir. Her servis bağımsız çalışır ve Docker Compose ile orkestrasyonu sağlanır. Event-driven yapı için RabbitMQ, loglama ve gözlemlenebilirlik için Elasticsearch, Kibana, Grafana ve özel logger servisi kullanılmaktadır.

---

## 🧱 Proje Yapısı

```
.
├── api-gateway         # Tüm isteklerin yönlendirildiği giriş noktası
├── user-service        # Kullanıcı yönetimi (auth, register vs.)
├── product-service     # Ürün listeleme, ekleme vs.
├── order-service       # Sipariş işlemleri
├── logger              # Logları Elasticsearch'e gönderen servis
└── docker-compose.yml  # Tüm sistemin orkestrasyon dosyası
```

---

## 🚀 Başlatma

Projeyi başlatmak için aşağıdaki adımları takip edin:

```bash
docker-compose up
```

---

## 🔌 Kullanılan Servisler

| Servis | Açıklama | Port |
| ------ | -------- | ---- |
| **API Gateway**     | Mikroservisleri yöneten giriş noktası | `8080`                            |
| **User Service**    | Kullanıcı işlemleri servisi           | `8081`                            |
| **Product Service** | Ürün yönetimi servisi                 | `8082`                            |
| **Order Service**   | Sipariş yönetimi servisi              | `8083`                            |
| **Logger**          | Logları Elasticsearch'e gönderir      | `8084`                            |
| **PostgreSQL**      | Veritabanı servisi                    | `5432`                            |
| **RabbitMQ**        | Event-driven iletişim                 | `5672` / `15672` (Yönetim Paneli) |
| **Elasticsearch**   | Log veri deposu                       | `9200`                            |
| **Kibana**          | Log görüntüleme arayüzü               | `5601`                            |
| **Grafana**         | Monitoring paneli                     | `3000`                            |

---

## 🧪 Servisler Arası İletişim

- Servisler arasında doğrudan HTTP veya RabbitMQ ile haberleşme yapılır.
- `logger` servisi tüm logları alır ve `elasticsearch`'e yollar.
- `api-gateway`, tüm dış istekleri ilgili mikroservislere yönlendirir.
- RabbitMQ, `product-service` ve `order-service` arasında event iletişimini sağlar.

---

## 💠 Geliştirme Notları

- Her servis kendi içinde bağımsız çalışabilir ve ölçeklenebilir.
- Ortak ağ olarak `postgres_network` kullanılmıştır.
- Logger servisi JSON formatında log bekler. (Örnek: `{"service": "user-service", "level": "info", "message": "User created"}`)

---

## 📈 Monitoring

- **Grafana Paneli:** [http://localhost:3000](http://localhost:3000)
- **Kibana Arayüzü:** [http://localhost:5601](http://localhost:5601)
- **RabbitMQ Yönetim Paneli:** [http://localhost:15672](http://localhost:15672)\
  Kullanıcı: `guest`   Şifre: `guest`

---

## 📄 Lisans

Bu proje bir öğrenme projesidir ve açık kaynak olarak paylaşılmıştır.

---

## 👩‍💻 Hazırlayan

**Adı Soyadı:** Nazar Arık  
**İletişim:** [nazararik585@gmail.com](nazararik585@gmail.com)  
**GitHub:** [github.com/nazzarr03](https://github.com/nazzarr03)

