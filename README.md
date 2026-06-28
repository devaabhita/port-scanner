# Mini Port Scanner (Go)

## Deskripsi

Mini Port Scanner adalah program sederhana berbasis **Golang** yang digunakan untuk melakukan scanning port pada sebuah host (IP / domain). Program ini bekerja dengan mencoba melakukan koneksi ke berbagai port dan menentukan apakah port tersebut **OPEN atau CLOSED**, serta mencoba **mendeteksi service** yang berjalan di port tersebut.

Project ini merupakan versi sederhana yang terinspirasi dari tools seperti **Nmap**, dengan fokus pada:

* Pembelajaran networking dasar
* Pemahaman concurrency di Golang
* Eksplorasi komunikasi TCP

---

## Konsep Dasar

Dalam jaringan komputer:

* **Host (IP / Domain)** = alamat tujuan
* **Port** = pintu komunikasi
* **Service** = aplikasi yang berjalan di port

Contoh:

* Port 80 → HTTP (Website)
* Port 22 → SSH
* Port 3306 → MySQL

---

## Cara Kerja Program

Program akan:

1. Mengambil target (contoh: `localhost`)
2. Mengirim daftar port ke worker (goroutine)
3. Worker mencoba koneksi TCP ke tiap port
4. Jika berhasil:

   * Port dianggap **OPEN**
   * Program mencoba membaca **banner**
   * Service dideteksi
5. Hasil dikirim dan ditampilkan

---

## Alur Program

```
[ Daftar Port 1 - 9000 ]
            |
            v
      [ Channel Ports ]
            |
            v
   ---------------------
   |   Worker Pool     |
   |  (100 goroutine)  |
   ---------------------
     |    |     |    |
     v    v     v    v
   Scan Scan  Scan Scan
     |    |     |    |
     ---------------
            |
            v
     [ Channel Results ]
            |
            v
       Output Terminal
```

---

## Penjelasan Kode

### 1. Struct Result

```go
type Result struct {
	Port    int
	Service string
}
```

Digunakan untuk menyimpan hasil scanning.

---

### 2. detectService()

```go
func detectService(conn net.Conn, port int) string
```

Fungsi ini:

* Membaca data dari koneksi (banner)
* Menentukan service berdasarkan:

  * Port (hardcoded)
  * Isi banner

---

### 3. worker()

```go
func worker(wg *sync.WaitGroup, ports <-chan int, results chan<- Result, target string)
```

Worker bertugas:

* Mengambil port dari channel
* Mencoba koneksi TCP (`net.DialTimeout`)
* Jika berhasil → kirim hasil ke channel results

---

### 4. Worker Pool

```go
numWorkers := 100
```

Artinya:

* Maksimal hanya **100 koneksi berjalan bersamaan**
* Mencegah overload sistem

---

### 5. Channel

```go
ports := make(chan int, 100)
results := make(chan Result)
```

* `ports` → mengirim daftar port
* `results` → menerima hasil scan

---

### 6. Main Flow

```go
go func() {
	for port := startPort; port <= endPort; port++ {
		ports <- port
	}
	close(ports)
}()
```

Mengirim semua port ke worker.

```go
go func() {
	wg.Wait()
	close(results)
}()
```

Menutup channel setelah semua worker selesai.

---

## Konfigurasi yang Bisa Diubah

### Target

```go
target := "localhost"
```

Contoh:

* `"127.0.0.1"`
* `"scanme.nmap.org"`

---

### Range Port

```go
startPort := 1
endPort := 9000
```

---

### Jumlah Worker

```go
numWorkers := 100
```

* Kecil → lebih stabil
* Besar → lebih cepat tapi riskan

---

### Timeout

```go
net.DialTimeout("tcp", address, 1*time.Second)
```

---

## Cara Menjalankan

### 1. Pastikan Go sudah terinstall

```
go version
```

---

### 2. Jalankan program

```
go run main.go
```

---

### 3. Contoh output

```
PORT    STATUS  SERVICE
22      OPEN    SSH
80      OPEN    HTTP
8080    OPEN    HTTP
3306    OPEN    MySQL
```

---

## Testing (Disarankan)

Jalankan server lokal di terminal:

```
python3 -m http.server 8080
```

Lalu scan:

```
target := "127.0.0.1"
```
atau
```
target := "localhost"
```
Keduanya sama saja

## Tujuan Pembelajaran

Project ini membantu memahami:

* TCP networking dasar
* Cara kerja port & service
* Concurrency (goroutine & channel)
* Worker pool pattern

---

