============================================================
                   PIO BAT v1.0 (Windows CLI)                   
        Battery Detection & Management Tool for Windows         
============================================================

[ OVERVIEW ]
PIO BAT adalah utility CLI berbasis PowerShell & Batch Runner untuk 
melakukan inspeksi, diagnostik kesehatan, pemantauan status, dan 
manajemen charging threshold baterai laptop pada Windows 10 & 11.

[ AUTHOR & CREDITS ]
* Developer   : Cyberly Dev
* GitHub Repo : https://github.com/cyberlycode
* Project URL : https://github.com/cyberlycode/Cyberly-Dev/tree/piobat
* Suite       : PIO Automation & Optimization Suite

[ COMMANDS / CARA PENGGUNAAN ]
* piobat.bat all         : Menampilkan seluruh inspeksi & diagnostik baterai
* piobat.bat info        : Menampilkan rincian hardware & kapasitas baterai
* piobat.bat health      : Menampilkan persentase kesehatan sel baterai (mWh / Design)
* piobat.bat status      : Menampilkan status daya saat ini
* piobat.bat capacity    : Menampilkan persentase daya baterai saat ini
* piobat.bat threshold 80: Mengatur charging threshold pengisian (misal: 80%)

[ CARA PENGGUNAAN ]
1. Buka PowerShell / Terminal di folder ini.
2. Jika terblokir kebijakan eksekusi skrip, jalankan sekali:
   Set-ExecutionPolicy -ExecutionPolicy Unrestricted -Scope Process
3. Jalankan perintah runner:
   .\piobat.bat all
