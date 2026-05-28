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
