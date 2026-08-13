# PIOBATWIN 1.0 (Universal WIndows Kompatibel)

Utility ringan untuk inspeksi kesehatan baterai, pemantauan status daya, dan pengaturan charging threshold laptop di Windows.

---

## Author & Credits

* **Developer**: Cyberly Dev
* **Publisher**: Cyberly Code Dev
* **Repository**: [Cyberly-Dev (branch: piobat)](https://github.com/cyberlycode/Cyberly-Dev/tree/piobat)
* **Suite**: PIO Automation & Optimization Suite

---

## Dukungan OS

* **Windows 10 & 11** (Full Support)
* **Windows 8 / 8.1**
* **Windows 7**
* **Windows Vista** *(perlu PowerShell 2.0)*

---

## Cara Penggunaan

### 1. Mode GUI (Klik 2x)
Langsung *double-click* file **`piobat.exe`**. Tampilan informasi hardware dan kesehatan baterai akan langsung muncul di jendela program tanpa perlu ketik perintah apapun.

### 2. Mode CLI (Khusus `piobat.ps1` via PowerShell)
Jika ingin menggunakan parameter/command spesifik, buka PowerShell di folder lokasi file lalu jalankan skrip `.ps1`:

```powershell
# Cek semua info & diagnostik
.\piobat.ps1 all

# Cek detail hardware & kapasitas
.\piobat.ps1 info

# Cek persentase kesehatan baterai (Health)
.\piobat.ps1 health

# Cek status daya (AC / Battery)
.\piobat.ps1 status

# Cek sisa persentase baterai
.\piobat.ps1 capacity

# Set batas pengisian baterai (misal: 80%)
.\piobat.ps1 threshold 80