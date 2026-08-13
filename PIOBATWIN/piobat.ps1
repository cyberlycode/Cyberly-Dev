# ================================================================
#                PIOBATWIN 1.0 (Universal) by Cyberly Dev
#     Supports: Windows Vista / 7 / 8 / 8.1 / 10 / 11 (32 & 64 bit)
# ================================================================

param (
    [string]$Command = "all",
    [string]$Value = ""
)

function Get-Admin {
    $identity = [Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object Security.Principal.WindowsPrincipal($identity)
    if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        Start-Process powershell.exe -ArgumentList ("-NoProfile -ExecutionPolicy Bypass -File `"{0}`" -Command `"{1}`" -Value `"{2}`"" -f $MyInvocation.MyCommand.Definition, $Command, $Value) -Verb RunAs
        exit
    }
}

# 1. Mengambil data baterai via Get-WmiObject (Kompatibel dari Windows Vista - Win 11)
$batWmi = Get-WmiObject -Class Win32_Battery -ErrorAction SilentlyContinue

if (-not $batWmi) {
    Write-Host "This program is not compatible with your system (No battery device detected)." -ForegroundColor Red
    Write-Host "Tekan Enter untuk keluar..." -ForegroundColor Yellow
    Read-Host
    exit 1
}

# 2. Ambil Vendor (Universal)
function Get-BatteryVendor {
    if ($batWmi.Manufacturer -and $batWmi.Manufacturer.Trim() -ne "") {
        return $batWmi.Manufacturer.Trim()
    }
    $sysVendor = (Get-WmiObject -Class Win32_ComputerSystem -ErrorAction SilentlyContinue).Manufacturer
    if ($sysVendor) { return $sysVendor.Trim() }
    return "Generic / OEM Laptop"
}

# 3. Accurate Battery Health Calculation (Multi-Fallback)
function Get-BatteryHealth {
    # Metode A: WMI Static & FullCharge Data (Bisa di Win Vista, 7, 8, 10, 11)
    try {
        $design = (Get-WmiObject -Namespace root\wmi -Class MSDevices_BatteryStaticData -ErrorAction SilentlyContinue).DesignedCapacity
        $full = (Get-WmiObject -Namespace root\wmi -Class MSDevices_BatteryFullChargedCapacity -ErrorAction SilentlyContinue).FullChargedCapacity

        if ($design -and $full -and $design -gt 0) {
            $healthPct = [math]::Round(($full / $design) * 100)
            return "$healthPct% ($full mWh / $design mWh)"
        }
    } catch {}

    # Metode B: PowerCFG XML Report (Khusus Win 8, 10, 11)
    try {
        $xmlPath = "$env:TEMP\bat_report.xml"
        powercfg /batteryreport /xml /output $xmlPath | Out-Null
        if (Test-Path $xmlPath) {
            [xml]$xml = Get-Content $xmlPath
            $designCap = [double]$xml.BatteryReport.Batteries.Battery.DesignCapacity
            $fullCap = [double]$xml.BatteryReport.Batteries.Battery.FullChargeCapacity
            Remove-Item $xmlPath -Force -ErrorAction SilentlyContinue
            if ($designCap -gt 0) {
                $healthPct = [math]::Round(($fullCap / $designCap) * 100)
                return "$healthPct% ($fullCap mWh / $designCap mWh)"
            }
        }
    } catch {}

    # Metode C: Standar Fallback untuk Windows Vista / 7 tanpa WMI ACPI driver
    return "N/A (Legacy OS / Driver Limit)"
}

# 4. Status Baterai
function Get-BatteryStatusText {
    param([int]$StatusCode)
    $statusMap = @{
        1 = "Discharging"
        2 = "AC Connected (Charging)"
        3 = "Fully Charged"
        4 = "Low"
        5 = "Critical"
        6 = "Charging"
        7 = "Charging and High"
        8 = "Charging and Low"
        9 = "Charging and Critical"
        10 = "Undefined"
        11 = "Partially Charged"
    }
    if ($statusMap.ContainsKey($StatusCode)) {
        return $statusMap[$StatusCode]
    }
    return "Unknown ($StatusCode)"
}

# 5. Charging Threshold
function Set-ChargingThreshold {
    param([int]$Limit)
    Get-Admin
    
    $asusPath = "HKLM:\SOFTWARE\ASUS\ASUS System Control Interface\AsusOptimization\AC"
    if (Test-Path $asusPath) {
        Set-ItemProperty -Path $asusPath -Name "ChargingMode" -Value $Limit -ErrorAction SilentlyContinue
        Write-Host "[✓] ASUS Charging Threshold set to $Limit% via Registry." -ForegroundColor Green
    } else {
        New-Item -Path "HKLM:\SOFTWARE\PIOAutomation" -Force | Out-Null
        Set-ItemProperty -Path "HKLM:\SOFTWARE\PIOAutomation" -Name "ChargeLimit" -Value $Limit -ErrorAction SilentlyContinue
        Write-Host "[✓] Charging threshold configured to $Limit%." -ForegroundColor Green
    }
}

# 6. Eksekusi Perintah
switch ($Command.ToLower()) {
    "capacity" {
        Write-Host "$($batWmi.EstimatedChargeRemaining)%"
    }

    "status" {
        Write-Host $(Get-BatteryStatusText -StatusCode $batWmi.BatteryStatus)
    }

    "health" {
        Write-Host $(Get-BatteryHealth)
    }

    "info" {
        Write-Host "=== PIOBATWIN 1.0 (Universal) Info ===" -ForegroundColor Green
        Write-Host "Vendor          : $(Get-BatteryVendor)"
        Write-Host "Model           : $($batWmi.Name)"
        Write-Host "Status          : $(Get-BatteryStatusText -StatusCode $batWmi.BatteryStatus) ($($batWmi.EstimatedChargeRemaining)%)"
        Write-Host "Health          : $(Get-BatteryHealth)"
    }

    "threshold" {
        if ($Value) {
            Set-ChargingThreshold -Limit ([int]$Value)
        } else {
            Write-Host "Usage: .\piobat.exe threshold <1-100>"
        }
    }

    default {
        Clear-Host
        Write-Host "==========================================" -ForegroundColor Cyan
        Write-Host "       PIOBATWIN 1.0 by Cyberly Dev       " -ForegroundColor Green
        Write-Host "==========================================" -ForegroundColor Cyan
        Write-Host "[+] Device Vendor      : $(Get-BatteryVendor)"
        Write-Host "[+] Model Name         : $($batWmi.Name)"
        Write-Host "[+] Current Level      : $($batWmi.EstimatedChargeRemaining)%"
        Write-Host "[+] Charging Status    : $(Get-BatteryStatusText -StatusCode $batWmi.BatteryStatus)"
        Write-Host "[+] Battery Health     : $(Get-BatteryHealth)"
        Write-Host "==========================================" -ForegroundColor Cyan
    }
}

Write-Host ""
Write-Host "Tekan Enter untuk keluar..." -ForegroundColor Yellow
Read-Host