# 📚 GOLANG-BOOKSERVİCE
* Bu proje kitap yönetimini kolaylaştırmak için tasarlanmış bir API sistemidir.
* Kullanıcı adı ve şifre ile kayıt olup o kullanıcı adı ve şifre ile giriş yaparak okuduğunuz kitapları kaydedebilirsiniz. Girilen bilgiler veri tabanına kaydedildiği için silinmez ve kaybolmaz.

## Özellikler
* Bu API CRUD işlemlerini başarıyla gerçekleşmektedir.
* Kitaplar üzerinde arama yapılabilmektedir.
* Json formatında veri döndürme işlemi yapılabilmektedir.
* Kullanıcı adı ve şifre girişlerinde şifreyi HASH algoritması ile şifreleyip veri tabanında HASH'li bir şekilde tutmaktadır.
* Bu projenin ön tarafı C# Form ile yazılmıştır. Bu uygulamayı arka planda çalıştırıp ayrıca [cSharp-BookServiceForm](https://github.com/grkmkly/cSharp-BookServiceForm) uygulamasını çalıştararak işlemlerinizi kolayca halledebilirsiniz.

## Kullanılan Teknolojiler
* Golang : API geliştirme
* MONGODB : Veri Tabanı
* MUX : Hızlı ve esnek router
* CRYPTO : Şifreleme sistemi

## 📦 Kurulum ve Çalıştırma  

### 1. Gerekli araçları yükleyin  

- Go ([https://go.dev/](https://go.dev/))  

### 2. Projeyi klonlayın  

```bash
git clone https://github.com/grkmkly/golang-bookservice.git
cd golang-bookservice
```

### 3. Bağımlılıkları yükleyin

```bash
go mod tidy
```

### 4. Yapılandırma
Uygulamayı çalıştırmadan önce bir `.env` dosyası oluşturun ve aşağıdaki değişkenleri doldurun.

```plaintext
mongoUri =<MongoDB bağlantı URI'si>
PORT =<Çalıştırılacak Port>
```

### 5. Uygulamayı Çalıştırın
```bash
go run src/main.go
```

* Uygulaamayı çalıştırdıktan sonra

```plaintext
http://localhost:<PORT>/<İstek atılacak handler> 
````
linkine girerek ya da o linke istek atarak ulaşabilirsiniz.

# 📫 İLETİŞİM

Herhangi bir sorunuz varsa veya katkıda bulunmak isterseniz

* E-posta: [kolaygorkem@outlook.com]
* GitHub: [grkmkly](github.com/grkmkly)

