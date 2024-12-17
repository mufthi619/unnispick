# E-commerce API

A simple e-commerce API built with Go and Echo framework that helps you manage products and brands.

## Getting Started

First, make sure we have:
- Go installed
- PostgreSQL running
- Environment variables set up

### 1. Setup & Installation

Clone the repository and install dependencies:
```bash
go mod download
```

### 2. Config

Change & Replace ```config.yaml``` base on your env

### 2. Database Migration

Run the migrations to set up your database:
```bash
go run main.go -command migrate up
```

If you need to rollback:
```bash
go run main.go -command migrate down
```

## Working with Brands

### 1. Create Brand
We need to create a brand first before adding products:
```bash
curl --location 'http://localhost:4000/api/v1/brands' \
  --header 'Content-Type: application/json' \
  --data '{
      "brand_name": "SANGCLI"
  }'
```

### 2. List Brands
Get a list of all brands with pagination:
```bash
curl --location 'http://localhost:4000/api/v1/brands'
```

### 3. Get Brand by ID
Retrieve a specific brand using its ID:
```bash
curl --location 'http://localhost:4000/api/v1/brands/fe36e6e9-0b3f-4dac-943a-b60e093166f9'
```

## Working with Products

### 1. Create Product
After creating a brand, you can add products to it:
```bash
curl --location 'http://localhost:4000/api/v1/products' \
  --header 'Content-Type: application/json' \
  --data '{
    "product_name": "SANGCLI - Oat Barrier Sparkling Spa Bath Barm",
    "price": 460600,
    "quantity": 100,
    "brand_id": "fe36e6e9-0b3f-4dac-943a-b60e093166f9"
  }'
```

# ESSAY Answer
1. Mungkin saya akan menjelaskan terlebih dahulu project planning sesuai dengan pengalaman saya.
Project Planning biasanya akan diawali dengan permintaan user yang akan diwakili oleh Product Owner (PO), yang mana source Product Owner itu sendiri adalah orang bisnis dari perusahaan.
Setelah PO memiliki spek project, PO akan membuat PRD & MVP Project.btw sebelumnya saya ingin menjelaskan contoh kasus disini menggunakan SLDC Agile.
Nah once PRD selesai, kita akan melakukan Sprint Grooming untuk mendiskusikan PRD & MVP dari PO tersebut, setelah disepakati kita akan melakukan Sprint Planning yang dipimpin oleh Scrum Master / Project Manager.
After sprint planning maka akan keluar tiket tiket yang akan diselesaikan developer, after tiket selesai, QA akan melakukan SIT & UAT, beserta stres test dan test lainnya. After OK dari QA, kita akan naikan ke env Pre-Prod & Prod oleh Devops.
Setelah itu aplikasi release dan akan dilakukan internal review oleh orang bisnis dan internal testing oleh QA. Namun jika kita mempunyai cyber security, kita akan melakukan pentest juga.

Lalu untuk project review maksudnya code review kah ? Jika code review yang dimaksud, simple nya techlead akan melakukan Code Review pada pull request sebelum code naik ke env staging.

2. Menurut pemahaman saya load balancer adalah proses dimana kita melakukan distribusi request pada server yang bertujuan untuk memastikan server dapat berjalan baik sesuai bebannya masing-masing.
Contoh simple nya, kita akan menggunakan case kubernetes. Kubernetes sendiri bisa melakukan replika ketika resource dari satu nodes sudah mencapai maksimalnya, let say kubernetes melakukan replica sebanyak 3 nodes diawal,
load balancer akan membagi request masuk ke 3 nodes tersebut secara rata, proses pembagian secara rata tersebut adalah salah satu strategi Round Robin pada load balancer.
Lalu security group pada AWS EC2 sendiri sepengalaman saya adalah kumpulan konfigurasi VM, karena pada aplikasinya EC2 kebanyakan digunakan untuk membuat VM.
contoh security group pada EC2 itu sendiri adalah inbound dan outbound rules, yan mana contoh implementasinya adalah kita bisa membatasi port yang dapat diakses dari luar.
Lalu ada juga settingan VPC

3. Sebelum menjelaskan tentang step-step cara menangani issue memory leak di golang. Saya akan memberikan penjelasan singkat beserta contoh memory leak itu sendiri.
Simpelnya, memory leak adalah keadaan dimana garbage collector dari Go gagal mengembalikan memory yang telah dialokasikan kembali ke OS.
berbeda dengan Rust, Go menyediakan GC (Garbage Collector) yang dapat melakukan otomatisasi manajemen memory, nah GC ini sendiri bisa saja gagal melakukan release alokasi memori.
Berikut saya akan memberikan contoh beserta step-step cara menangani issue memory leak nya.

- Goroutine : 
function goroutine sendiri bisa saja terus berjalan di background tanpa henti dan tentu saja menyebabkan memory leak.
untuk mencegah memory leak pada goroutine kita bisa melakukan monitoring dengan pprof dan ```runtime.NumGoroutine```
caranya simple, kita hanya pelu memantau Goroutine yang aktif dengan pprof, atau secara manual melakukan print goroutine yang aktif dengan ```runtime.NumGoroutine```

- Long Live Reference
contoh simple nya, kita membuat looping ```for{}``` tanpa ada ttl dan tanpa ada parameter. Hal ini menyebabkan Go terus melakukan alokasi memory
dan GC tidak akan menghapus variable tersebur karena looping belum berhenti. Lalu contoh case lainnya, ketika kita melakukan goroutine dan terus melakukan
pengiriman data ke channel, yang mana channel tersebut sudah penuh dan tidak bisa menerima data lagi, namun goroutine terus mengirim data ke channel karena infinite loop tersebut.
Pencegahan pada kasus ini juga sangat simple, kita hanya perlu memberikan ttl ataupun parameter untuk berhenti disetiap melakukan perulangan

Pada akhirnya untuk kasus lain pada memory leak di Golang, kita bisa melakukan profiling dengan pprof dan tools lainnya.
Setelah mendapatkan masalahnya, kita bisa memberikan parameter berhenti, agar GC bisa melakukan alokasi memory.