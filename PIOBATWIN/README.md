# PIOBATWIN 1.0 (Universal)

Utility ringan berbasis PowerShell untuk inspeksi kesehatan baterai, pemantauan status daya, dan pengaturan charging threshold laptop pada sistem operasi Windows.

---

## Author & Credits

* **Developer**: Cyberly Dev
* **Publisher**: Cyberly Code Dev
* **Repository**: [Cyberly-Dev (branch: piobat)](https://github.com/cyberlycode/Cyberly-Dev/tree/piobat)
* **Suite**: PIO Automation & Optimization Suite

---

## Dukungan OS

Didesain dengan metode WMI universal agar kompatibel lintas generasi Windows:

* **Windows 10 & 11** (Full Support)
* **Windows 8 / 8.1**
* **Windows 7**
* **Windows Vista** *(memerlukan PowerShell 2.0 / WMF)*

---

## Cara Penggunaan

Tersedia dalam 2 mode penggunaan: **GUI** dan **CLI**.

### 1. Mode GUI (Klik 2x)
Langsung *double-click* file **`piobat.exe`**. Jendela konsol interaktif akan terbuka dan menampilkan seluruh informasi hardware serta kesehatan baterai secara otomatis tanpa perlu mengetik perintah apapun.

### 2. Mode CLI (Khusus `piobat.ps1` via PowerShell)
Buka PowerShell di folder lokasi file untuk menjalankan parameter spesifik:

```powershell
# Cek semua info & diagnostik utama
.\piobat.ps1 all

# Cek detail hardware & kapasitas baterai
.\piobat.ps1 info

# Cek persentase kesehatan baterai (Health mWh)
.\piobat.ps1 health

# Cek status daya (AC Charger / Discharging)
.\piobat.ps1 status

# Cek sisa persentase daya baterai saat ini
.\piobat.ps1 capacity

# Set batas pengisian baterai (misal: dibatasi 80%)
.\piobat.ps1 threshold 80


Catatan: Jika eksekusi skrip .ps1 terblokir oleh Windows Execution Policy, jalankan perintah ini sekali di PowerShell:

Set-ExecutionPolicy -ExecutionPolicy Unrestricted -Scope Process