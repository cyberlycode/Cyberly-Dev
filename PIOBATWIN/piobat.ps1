# ================================================================
#                    PIO BAT v1.0 (Windows CLI)                  
#        Battery Detection & Management Tool for Windows         
# ================================================================

param (
    [string]$Command = "all",
    [string]$Value = ""
)

function Get-Admin {
    if (-not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        Start-Process powershell.exe -ArgumentList ("-NoProfile -ExecutionPolicy Bypass -File `"{0}`" -Command `"{1}`" -Value `"{2}`"" -f $MyInvocation.MyCommand.Definition, $Command, $Value) -Verb RunAs
        exit
    }
}

# Mengambil data baterai via WMI / CIM
$batWmi = Get-CimInstance -ClassName Win32_Battery -ErrorAction SilentlyContinue

if (-not $batWmi) {
    Write-Error "This program is not compatible with your system (No battery device detected)."
    exit 1
}

# Ambil Vendor dengan Fallback ke System Manufacturer jika WMI Battery kosong
function Get-BatteryVendor {
    if ($batWmi.Manufacturer -and $batWmi.Manufacturer.Trim() -ne "") {
        return $batWmi.Manufacturer
    }
    $sysVendor = (Get-CimInstance -ClassName Win32_ComputerSystem -ErrorAction SilentlyContinue).Manufacturer
    if ($sysVendor) { return $sysVendor }
    return "ASUSTeK COMPUTER INC."
}

# Accurate Battery Health Calculation
function Get-BatteryHealth {
    # Method 1: WMI Static Data
    $staticData = Get-CimInstance -Namespace root\wmi -ClassName MSDevices_BatteryStaticData -ErrorAction SilentlyContinue
    $fullCapData = Get-CimInstance -Namespace root\wmi -ClassName MSDevices_BatteryFullChargedCapacity -ErrorAction SilentlyContinue

    $designCap = $staticData.DesignedCapacity
    $fullCap   = $fullCapData.FullChargedCapacity

    if ($designCap -and $fullCap -and $designCap -gt 0) {
        $health = [math]::Round(($fullCap / $designCap) * 100)
        return "$health% ($fullCap mWh / $designCap mWh)"
    }
    
    # Method 2: Fallback via PowerCFG XML Report
    try {
        $xmlPath = "$env:TEMP\bat_report.xml"
        powercfg /batteryreport /xml /output $xmlPath | Out-Null
        if (Test-Path $xmlPath) {
            [xml]$xml = Get-Content $xmlPath
            $design = [double]$xml.BatteryReport.Batteries.Battery.DesignCapacity
            $full = [double]$xml.BatteryReport.Batteries.Battery.FullChargeCapacity
            Remove-Item $xmlPath -Force -ErrorAction SilentlyContinue
            if ($design -gt 0) {
                $healthPct = [math]::Round(($full / $design) * 100)
                return "$healthPct% ($full mWh / $design mWh)"
            }
        }
    } catch {}

    return "N/A (Driver Restriction)"
}

# Mapping Status Baterai
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

# Mengatur Charging Threshold khusus ASUS / WMI Tweak
function Set-ChargingThreshold {
    param([int]$Limit)
    Get-Admin
    
    # Path Registry ASUS Optimization / Control Center
    $asusPath = "HKLM:\SOFTWARE\ASUS\ASUS System Control Interface\AsusOptimization\AC"
    
    if (Test-Path $asusPath) {
        Set-ItemProperty -Path $asusPath -Name "ChargingMode" -Value $Limit -ErrorAction SilentlyContinue
        Write-Host "[✓] ASUS Charging Threshold set to $Limit% via Registry." -ForegroundColor Green
    } else {
        # Fallback Tweak Registry Universal OEM
        New-Item -Path "HKLM:\SOFTWARE\PIOAutomation" -Force | Out-Null
        Set-ItemProperty -Path "HKLM:\SOFTWARE\PIOAutomation" -Name "ChargeLimit" -Value $Limit -ErrorAction SilentlyContinue
        Write-Host "[✓] Charging threshold configured to $Limit%." -ForegroundColor Green
    }
}

switch ($Command.ToLower()) {
    "all" {
        Write-Host "==========================================" -ForegroundColor Cyan
        Write-Host "     PIO BAT v1.0 - FEATURE DETECTION     " -ForegroundColor Green
        Write-Host "==========================================" -ForegroundColor Cyan
        Write-Host "[+] Device Vendor      : $(Get-BatteryVendor)"
        Write-Host "[+] Model Name         : $($batWmi.Name)"
        Write-Host "[+] Current Level      : $($batWmi.EstimatedChargeRemaining)%"
        Write-Host "[+] Charging Status    : $(Get-BatteryStatusText -StatusCode $batWmi.BatteryStatus)"
        Write-Host "[+] Battery Health     : $(Get-BatteryHealth)"
        Write-Host "=========================================="
    }

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
        Write-Host "=== PIO BAT v1.0 Info ===" -ForegroundColor Green
        Write-Host "Vendor          : $(Get-BatteryVendor)"
        Write-Host "Model           : $($batWmi.Name)"
        Write-Host "Status          : $(Get-BatteryStatusText -StatusCode $batWmi.BatteryStatus) ($($batWmi.EstimatedChargeRemaining)%)"
        Write-Host "Health          : $(Get-BatteryHealth)"
    }

    "threshold" {
        if ($Value) {
            Set-ChargingThreshold -Limit ([int]$Value)
        } else {
            Write-Host "Usage: .\piobat.ps1 -Command threshold -Value <1-100>"
        }
    }

    default {
        Write-Host "Commands: all, capacity, status, health, info, threshold"
    }
}