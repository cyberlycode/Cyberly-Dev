================================================================
                    PIO BAT v1.0 (Linux CLI)                     
         Battery Detection & Management Tool for Linux          
================================================================

[ OVERVIEW ]
PIO BAT adalah utility CLI berbasis Go untuk melakukan inspeksi,
diagnostik kesehatan, pemantauan status, dan manajemen charging 
threshold baterai laptop pada sistem operasi Linux (Pop!_OS, Ubuntu,
Debian, Arch Linux, Fedora, Linux Mint, dll).

[ AUTHOR & CREDITS ]
* Developer   : Cyberly Dev
* GitHub Repo : https://github.com/cyberlycode
* Project URL : https://github.com/cyberlycode/Cyberly-Dev/tree/piobat
* Suite       : PIO Automation & Optimization Suite

[ FITUR UTAMA ]
1. Auto-Detection Lengkap:
   - Identifikasi Vendor, Model, & Teknologi Sel Baterai.
   - Perhitungan Akurat Persentase Kesehatan Baterai (Battery Health).
   - Pembandingan Kapasitas Pabrikan (Max Design Cap) vs Kapasitas Maksimal Saat Ini (Max Current Cap).
   - Pemantauan Siklus Pengisian (Cycle Count) & Status Daya.
2. Pengaturan Threshold (Khusus Perangkat Terdukung):
   - Deteksi otomatis node sysfs & driver platform (ASUS/Lenovo/ThinkPad/WMI).
   - Preset mode profil pengisian (Desktop 60%, Balanced 80%, Travel 100%).
   - Persistensi pengaturan otomatis via Systemd service.

[ COMMANDS / CARA PENGGUNAAN ]
* piobat all      : Menampilkan seluruh inspeksi & diagnostik fitur baterai (Full Check)
* piobat info     : Menampilkan rincian informasi hardware & kapasitas baterai
* piobat health   : Menampilkan persentase kesehatan sel baterai (Health/Wear Level)
* piobat status   : Menampilkan status daya saat ini (Charging / Discharging)
* piobat capacity : Menampilkan persentase daya baterai saat ini
* piobat profile  : Memilih preset threshold (piobat profile desktop|balanced|travel)
* piobat threshold: Mengatur angka threshold secara manual (misal: piobat threshold 80)
* piobat persist  : Menyimpan konfigurasi threshold agar aktif otomatis saat boot
* piobat reset    : Menghapus konfigurasi persistensi threshold

================================================================
           PANDUAN CLONING, COMPILE, & INSTALLASI
================================================================

PRASYARAT (PREREQUISITES):
- OS Linux (Pop!_OS, Ubuntu, Debian, Arch, Fedora, dll)
- Golang / Go Compiler terpasang (Cek dengan: go version)
- Git terpasang (Cek dengan: git --version)

1. CLONE REPOSITORI (Branch piobat):
   Buka terminal dan clone branch 'piobat' dari repositori Cyberly-Dev:

   git clone -b piobat https://github.com/cyberlycode/Cyberly-Dev.git piobat-tool
   cd piobat-tool

2. COMPILE / BUILD BINER:
   Jalankan kompilasi menggunakan Go:

   go build -o piobat main_linux.go

3. PEMASANGAN KE SISTEM GLOBAL (SIAP PAKAI):
   Pindahkan biner hasil kompilasi ke direktori biner sistem agar bisa dipanggil dari direktori mana saja:

   sudo cp piobat /usr/local/bin/

4. PENGUJIAN & EKSEKUSI:
   Jalankan perintah berikut di terminal mana saja tanpa perlu memasukkan path folder:

   piobat all

----------------------------------------------------------------
CARA CEPAT 1-LINE INSTALL (QUICK ONE-LINER):
Gunakan perintah satu baris ini di terminal untuk langsung clone, build, dan install:

git clone -b piobat https://github.com/cyberlycode/Cyberly-Dev.git piobat-tool && cd piobat-tool && go build -o piobat main_linux.go && sudo cp piobat /usr/local/bin/ && piobat all

================================================================
                   Created with <3 by Cyberly Dev
================================================================
