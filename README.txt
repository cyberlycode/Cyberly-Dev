================================================================
                PIO BAT v1.0 (Linux CLI)                     
     Battery Detection & Management Tool for Linux          
================================================================

[ OVERVIEW ]
PIO BAT adalah utility CLI berbasis Go untuk melakukan inspeksi,
diagnostik kesehatan, pemantauan status, dan manajemen charging 
threshold baterai laptop pada sistem operasi Linux (Pop!_OS, Ubuntu,
Debian, Arch, dll).

[ AUTHOR & CREDITS ]
* Developer   : Cyberly Dev
* GitHub Repo : https://github.com/cyberlycode
* Suite       : PIO Automation & Optimization Suite

[ COMMANDS ]
* piobat all      : Menampilkan seluruh inspeksi fitur baterai (Full Check)
* piobat info     : Menampilkan informasi detail hardware baterai
* piobat health   : Menampilkan persentase kesehatan baterai (Battery Health)
* piobat status   : Menampilkan status daya (Charging / Discharging)
* piobat capacity : Menampilkan persentase daya baterai saat ini
* piobat profile  : Memilih preset threshold (desktop|balanced|travel)
* piobat persist  : Menyimpan konfigurasi threshold via systemd service

[ INSTALLATION ]
1. Compile biner:
go build -o piobat main_linux.go

2. Pasang ke sistem:
sudo mv piobat /usr/local/bin/

================================================================
               Created with <3 by Cyberly Dev
