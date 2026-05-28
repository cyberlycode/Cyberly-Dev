# ASUS X441UV Linux Pop!_OS 24.04 Optimization & Troubleshooting Guide

Panduan komprehensif ini berisi kumpulan solusi teknis (*tweak*) untuk mengatasi berbagai masalah hardware dan visual pada laptop **ASUS X441UV (Intel Core + NVIDIA GeForce 920MX)** yang menjalankan **Pop!_OS 24.04 (KDE Plasma/SDDM Setup)**. 

Panduan ini mendokumentasikan perbaikan nyata untuk masalah tombol hantu (*ghosting*), optimasi performa storage, transisi *booting* yang bersih, serta stabilisasi layar dari *flicker*.

---

## 🖥️ Spesifikasi Target Perangkat
* **Model Perangkat:** ASUS X441UV Series
* **Sistem Operasi:** Pop!_OS 24.04 LTS (Arsitektur kernel-stub / systemd-boot)
* **Lingkungan Desktop:** KDE Plasma / SDDM
* **Spesifikasi Diuji:** RAM 12GB / Dual-GPU Intel + NVIDIA / SSD Storage

---

## 🛠️ 1. Memperbaiki Tombol Hantu (Ghosting Keyboard)
### Gejala Masalah:
Saat mesin laptop dalam kondisi dingin atau pasca-boot, keyboard internal mendadak mengirimkan sinyal interupsi palsu (mengetik sendiri secara acak tombol `C`, `Ctrl`, atau `F`), dipicu oleh kebocoran kode pemicu (*scancode* kotor `-1`) dari *firmware* BIOS/ACPI bawaan ASUS.

### Solusi Kernel:
Memaksa kernel Linux berkomunikasi langsung dengan pengontrol interupsi fisik hardware (`i8042`) tanpa melalui emulasi atau *handshake* firmware BIOS.

Jalankan perintah berikut untuk menyuntikkan parameter ke konfigurasi *bootloader*:
```bash
sudo kernelstub -o "quiet splash loglevel=0 systemd.show_status=false rd.systemd.show_status=false rd.udev.log_level=0 nvidia-drm.modeset=1 libata.force=6.0Gbps,ncq i8042.direct"

2. Akselerasi SSD SATA III & Aktivasi NCQ
Masalah:
Kecepatan transfer data pada SSD internal sering kali tidak berjalan di batas maksimal akibat negosiasi otomatis antarjalur bus data yang kurang optimal pada pengontrol bawaan laptop lama.

Solusi:
Mengunci kecepatan bus pengontrol ATA (libata) secara paksa di standar SATA 3.0 (6.0 Gbps) sekaligus mengaktifkan fitur antrean perintah (Native Command Queuing / NCQ) untuk mendongkrak performa kecepatan baca/tulis (I/O throughput).

Catatan: Parameter pengunci ini sudah disatukan ke dalam perintah kernelstub pada Langkah 1 di atas lewat argumen libata.force=6.0Gbps,ncq.

3. Konfigurasi "Silent Boot" Sempurna
Masalah:
Meskipun parameter visual quiet splash sudah aktif, layar logo booting Plymouth bawaan Pop!_OS sering kali kecolongan menampilkan baris teks indikator status [ OK ] warna hijau (terutama dari servis colord, cups, dan status penutupan Plymouth) sesaat sebelum menu login SDDM muncul.

Solusi:
Mengubah perilaku pengelola layanan (Systemd) agar menyembunyikan status informasi sukses dan menunda tirai transisi TTY Plymouth hingga pemuatan komponen printer dan profil warna selesai sepenuhnya di latar belakang.

Langkah A: Override Konfigurasi Systemd Manager
Jalankan perintah ini untuk membuat instruksi pembungkaman log sukses secara permanen:

Bash
sudo mkdir -p /etc/systemd/system.conf.d/
echo -e "[Manager]\nShowStatus=no\nLogLevel=notice" | sudo tee /etc/systemd/system.conf.d/99-silent-boot.conf
Langkah B: Mengatur Jeda Penutupan Tirai TTY Plymouth
Memaksa sistem untuk menahan proses pembersihan layar (plymouth-hide-tty.service) agar menutupi proses pemuatan layanan subsistem yang sering bocor ke layar utama:

Bash
sudo mkdir -p /etc/systemd/system/plymouth-hide-tty.service.d/
echo -e "[Unit]\nAfter=colord.service cups.service cups-browsed.service" | sudo tee /etc/systemd/system/plymouth-hide-tty.service.d/99-delay-hide.conf
Langkah C: Sinkronisasi Daemon Systemd
Muat ulang manajer servis agar seluruh arsitektur penundaan baru langsung diterapkan oleh sistem:

Bash
sudo systemctl daemon-reload
🔏 4. Tahap Akhir: Kompilasi Citra Ramdisk
Agar seluruh parameter kernel dan aturan baru yang telah dibuat di atas dapat dimuat sejak detik pertama laptop dinyalakan, Anda wajib membangun ulang citra ramdisk utama sistem (initrd).

Jalankan perintah regenerasi ini:

Bash
sudo update-initramfs -u
Setelah proses penyalinan data kernelstub ke dalam partisi EFI (ESP) selesai tanpa memicu pesan error, lakukan proses menyalakan ulang laptop:

Bash
sudo reboot
Sistem laptop ASUS X441UV Anda sekarang akan melakukan proses booting dengan sangat tenang, bersih, fungsional, dan bebas dari gangguan tombol hantu.
