// Binary piobat is a battery detection and management utility for Linux laptops.
package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	rtdebug "runtime/debug"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/sys/unix"
)

type Service struct {
	Event, Path, Shell string
	Threshold          int
}

type Target struct {
	Unit string `json:"unit"`
}

const appVersion = "1.0"

var globalThresholdPaths = []string{
	"/sys/class/power_supply/BAT0/charge_control_end_threshold",
	"/sys/class/power_supply/BAT0/charge_stop_threshold",
	"/sys/class/power_supply/BAT0/charge_end_threshold",
	"/sys/class/power_supply/BAT1/charge_control_end_threshold",
	"/sys/class/power_supply/BAT1/charge_stop_threshold",
	"/sys/devices/platform/asus-wmi/charge_control_end_threshold",
	"/sys/devices/platform/asus-nb-wmi/charge_control_end_threshold",
}

var (
	events = map[string]struct{}{
		"hibernate":              {},
		"hybrid-sleep":           {},
		"multi-user":             {},
		"suspend":                {},
		"suspend-then-hibernate": {},
	}

	services = filepath.Join("/", "etc", "systemd", "system")

	//go:embed bat.service
	unit string

	//go:embed help.txt
	usage string
)

type battery struct {
	root string
}

func (b *battery) has(variable string) (bool, error) {
	_, err := os.Stat(b.path(variable))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (b *battery) path(variable string) string {
	return filepath.Join(b.root, variable)
}

func (b *battery) read(variable string) (string, error) {
	contents, err := os.ReadFile(b.path(variable))
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(contents)), nil
}

func (b *battery) readOrDefault(variable, defaultValue string) string {
	val, err := b.read(variable)
	if err != nil {
		return defaultValue
	}
	return val
}

// Mengambil rincian kapasitas fisik (Max Design vs Max Current dalam mAh/mWh)
func (b *battery) getCapacityDetails() (designCap string, currentFullCap string, healthPct string) {
	var (
		err  error
		v, w string
		unitStr = "mAh"
	)
	s, t := "charge_full", "charge_full_design"
	v, err = b.read(s)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			s, t = "energy_full", "energy_full_design"
			unitStr = "mWh"
			v, err = b.read(s)
			if err != nil {
				return "N/A", "N/A", "N/A"
			}
			w, err = b.read(t)
			if err != nil {
				return "N/A", "N/A", "N/A"
			}
		} else {
			return "N/A", "N/A", "N/A"
		}
	} else {
		w, err = b.read(t)
		if err != nil {
			return "N/A", "N/A", "N/A"
		}
	}

	x, err1 := strconv.Atoi(v)
	y, err2 := strconv.Atoi(w)
	if err1 != nil || err2 != nil || y == 0 {
		return "N/A", "N/A", "N/A"
	}

	// Nilai sysfs biasanya menggunakan micro-ampere/watt hour (µAh/µWh), dikonversi ke mAh/mWh
	designCap = fmt.Sprintf("%d %s", y/1000, unitStr)
	currentFullCap = fmt.Sprintf("%d %s", x/1000, unitStr)
	healthPct = fmt.Sprintf("%d%%", x*100/y)

	return designCap, currentFullCap, healthPct
}

func (b *battery) write(variable string, contents []byte) error {
	return os.WriteFile(b.path(variable), contents, 0o644)
}

func getThresholdPath(batRoot string) (string, bool) {
	for _, cand := range []string{"charge_control_end_threshold", "charge_stop_threshold", "charge_end_threshold"} {
		p := filepath.Join(batRoot, cand)
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	for _, p := range globalThresholdPaths {
		if _, err := os.Stat(p); err == nil {
			return p, true
		}
	}
	return "", false
}

func main() {
	const ignore = ""
	var (
		d, debug   = flag.Bool("d", false, ignore), flag.Bool("debug", false, ignore)
		h, help    = flag.Bool("h", false, ignore), flag.Bool("help", false, ignore)
		v, version = flag.Bool("v", false, ignore), flag.Bool("version", false, ignore)
	)
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	if *h || *help {
		flag.Usage()
		return
	}

	if *v || *version {
		fmt.Printf("PIO BAT v%s\nBattery Detection & Management Tool for Linux\nPIO Automation & Optimization Suite\n", appVersion)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			var message string
			if *d || *debug {
				message = fmt.Sprintf("%s\n\n%s", err, string(rtdebug.Stack()))
			} else {
				message = "A fatal error occurred. Please rerun the command with the `--debug` flag enabled."
			}
			fmt.Fprintln(os.Stderr, message)
		}
	}()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	batteries, err := filepath.Glob(filepath.Join("/", "sys", "class", "power_supply", "BAT?"))
	if err != nil {
		panic(err)
	}
	if len(batteries) == 0 {
		fmt.Fprintln(
			os.Stderr,
			"This program is not compatible with your system (No battery device detected).",
		)
		os.Exit(1)
	}

	bat := &battery{root: batteries[0]}

	switch subcommand := flag.Arg(0); subcommand {
	case "all", "check":
		vendor := bat.readOrDefault("manufacturer", "N/A")
		model := bat.readOrDefault("model_name", "N/A")
		tech := bat.readOrDefault("technology", "N/A")
		status := bat.readOrDefault("status", "N/A")
		capacity := bat.readOrDefault("capacity", "N/A")
		cycles := bat.readOrDefault("cycle_count", "N/A")
		designCap, maxCap, health := bat.getCapacityDetails()

		fmt.Println("==========================================")
		fmt.Printf("      PIO BAT v%s - FEATURE DETECTION     \n", appVersion)
		fmt.Println("==========================================")
		fmt.Printf("[+] Device Vendor      : %s\n", vendor)
		fmt.Printf("[+] Model Name         : %s\n", model)
		fmt.Printf("[+] Battery Technology : %s\n", tech)
		fmt.Printf("[+] Current Level      : %s%%\n", capacity)
		fmt.Printf("[+] Charging Status    : %s\n", status)
		fmt.Printf("[+] Max Design Cap     : %s (Original Factory)\n", designCap)
		fmt.Printf("[+] Max Current Cap    : %s (Full Charge Available)\n", maxCap)
		fmt.Printf("[+] Battery Health     : %s\n", health)
		fmt.Printf("[+] Cycle Count        : %s\n", cycles)

		if path, ok := getThresholdPath(bat.root); ok {
			contents, err := os.ReadFile(path)
			if err == nil {
				fmt.Printf("[+] Threshold Control  : SUPPORTED (%s%%)\n", strings.TrimSpace(string(contents)))
				fmt.Printf("    └─ Node Path       : %s\n", path)
			}
		} else {
			fmt.Println("[-] Threshold Control  : NOT SUPPORTED (Hardware/Kernel limit)")
		}
		fmt.Println("==========================================")

	case "capacity", "status":
		v, err := bat.read(subcommand)
		if err != nil {
			panic(err)
		}
		fmt.Println(v)

	case "health":
		_, _, health := bat.getCapacityDetails()
		fmt.Println(health)

	case "info":
		vendor := bat.readOrDefault("manufacturer", "N/A")
		model := bat.readOrDefault("model_name", "N/A")
		tech := bat.readOrDefault("technology", "N/A")
		cycles := bat.readOrDefault("cycle_count", "N/A")
		status := bat.readOrDefault("status", "N/A")
		capacity := bat.readOrDefault("capacity", "N/A")
		designCap, maxCap, health := bat.getCapacityDetails()

		fmt.Printf("=== PIO BAT v%s Info ===\n", appVersion)
		fmt.Printf("Vendor          : %s\n", vendor)
		fmt.Printf("Model           : %s\n", model)
		fmt.Printf("Technology      : %s\n", tech)
		fmt.Printf("Status          : %s (%s%%)\n", status, capacity)
		fmt.Printf("Max Factory Cap : %s\n", designCap)
		fmt.Printf("Max Current Cap : %s\n", maxCap)
		fmt.Printf("Health          : %s\n", health)
		fmt.Printf("Cycle Count     : %s\n", cycles)

		if path, ok := getThresholdPath(bat.root); ok {
			contents, err := os.ReadFile(path)
			if err == nil {
				fmt.Printf("Threshold       : %s%% (node: %s)\n", strings.TrimSpace(string(contents)), path)
			}
		} else {
			fmt.Println("Threshold       : Not supported directly via sysfs")
		}

	case "profile":
		if flag.NArg() < 2 {
			fmt.Fprintln(os.Stderr, "Usage: piobat profile <desktop|balanced|travel>")
			os.Exit(1)
		}

		thresholdPath, ok := getThresholdPath(bat.root)
		if !ok {
			fmt.Fprintln(os.Stderr, "Charging threshold setting not supported on this device.")
			os.Exit(1)
		}

		var targetVal string
		profName := strings.ToLower(flag.Arg(1))

		switch profName {
		case "desktop":
			targetVal = "60"
		case "balanced":
			targetVal = "80"
		case "travel":
			targetVal = "100"
		default:
			fmt.Fprintln(os.Stderr, "Invalid profile. Available: desktop (60%), balanced (80%), travel (100%).")
			os.Exit(1)
		}

		if err := os.WriteFile(thresholdPath, []byte(targetVal), 0o644); err != nil {
			if errors.Is(err, unix.EACCES) {
				fmt.Fprintln(os.Stderr, "Permission denied. Try running this command with `sudo`.")
				os.Exit(1)
			}
			panic(err)
		}
		fmt.Printf("Profile '%s' applied. Charging threshold set to %s%%.\n", profName, targetVal)

	case "persist":
		thresholdPath, ok := getThresholdPath(bat.root)
		if !ok {
			fmt.Fprintln(os.Stderr, "Charging threshold setting not found.")
			os.Exit(1)
		}

		output, err := exec.Command("systemctl", "--version").CombinedOutput()
		if err != nil {
			panic(err)
		}
		var revision int
		_, err = fmt.Sscanf(string(output), "systemd %d", &revision)
		if err != nil {
			panic(err)
		}
		if revision < 244 {
			fmt.Fprintln(os.Stderr, "Requires systemd version 243-rc1 or later.")
			os.Exit(1)
		}

		shell, err := exec.LookPath("sh")
		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				fmt.Fprintln(os.Stderr, "Could not find `sh` in your `$PATH`.")
				os.Exit(1)
			}
			panic(err)
		}

		contents, err := os.ReadFile(thresholdPath)
		if err != nil {
			panic(err)
		}
		current, err := strconv.Atoi(strings.TrimSpace(string(contents)))
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("systemctl", "list-units", "--type", "target", "--all", "--plain", "--output", "json")
		output, err = cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		targets := make([]Target, 0)
		if err = json.Unmarshal(output, &targets); err != nil {
			panic(err)
		}
		available := make([]string, 0)
		for _, target := range targets {
			event := strings.TrimSuffix(target.Unit, ".target")
			_, ok := events[event]
			if ok {
				available = append(available, event)
			}
		}
		tmpl := template.Must(template.New("unit").Parse(unit))
		for _, event := range available {
			service := "piobat-" + event + ".service"
			f, err := os.Create(filepath.Join(services, service))
			if err != nil {
				if errors.Is(err, unix.EACCES) {
					fmt.Fprintln(os.Stderr, "Permission denied. Try running this command with `sudo`.")
					os.Exit(1)
				}
				panic(err)
			}
			s := Service{
				Event:     event,
				Path:      thresholdPath,
				Shell:     shell,
				Threshold: current,
			}
			if err = tmpl.Execute(f, s); err != nil {
				panic(err)
			}
			if err := exec.Command("systemctl", "enable", service).Run(); err != nil {
				panic(err)
			}
			f.Close()
		}
		fmt.Println("Persistence of the current charging threshold enabled for PIO BAT.")

	case "threshold":
		thresholdPath, ok := getThresholdPath(bat.root)
		if !ok {
			fmt.Fprintln(os.Stderr, "Charging threshold setting not found.")
			os.Exit(1)
		}
		switch flag.NArg() {
		case 1:
			contents, err := os.ReadFile(thresholdPath)
			if err != nil {
				panic(err)
			}
			fmt.Println(strings.TrimSpace(string(contents)))
		case 2:
			var utsname unix.Utsname
			if err = unix.Uname(&utsname); err != nil {
				panic(err)
			}
			var maj, min int
			_, err = fmt.Sscanf(string(utsname.Release[:]), "%d.%d", &maj, &min)
			if err != nil {
				panic(err)
			}
			if maj <= 5 && (maj != 5 || min < 4) {
				fmt.Fprintln(os.Stderr, "Requires Linux kernel version 5.4 or later.")
				os.Exit(1)
			}

			setting := flag.Arg(1)
			i, err := strconv.Atoi(setting)
			if err != nil {
				if errors.Is(err, strconv.ErrSyntax) {
					fmt.Fprintln(os.Stderr, "Argument should be an integer.")
					os.Exit(1)
				}
				panic(err)
			}
			if i < 1 || i > 100 {
				fmt.Fprintln(os.Stderr, "Threshold value should be between 1 and 100.")
				os.Exit(1)
			}
			if err := os.WriteFile(thresholdPath, []byte(setting), 0o644); err != nil {
				if errors.Is(err, unix.EACCES) {
					fmt.Fprintln(os.Stderr, "Permission denied. Try running this command with `sudo`.")
					os.Exit(1)
				}
				panic(err)
			}
			fmt.Println("Charging threshold set.\n" +
				"Run `sudo piobat persist` to persist the setting between restarts.")
		default:
			fmt.Fprintln(os.Stderr, "Invalid number of arguments.")
			flag.Usage()
			os.Exit(1)
		}

	case "reset":
		for event := range events {
			service := "piobat-" + event + ".service"
			output, err := exec.Command("systemctl", "disable", service).CombinedOutput()
			if err != nil {
				switch {
				case bytes.Contains(output, []byte("authentication required")):
					fmt.Fprintln(os.Stderr, "Permission denied. Try running this command with `sudo`.")
					os.Exit(1)
				case bytes.Contains(output, []byte("service does not exist")):
					continue
				default:
					panic(string(output))
				}
			}
			err = os.Remove(filepath.Join(services, service))
			if err != nil && !errors.Is(err, unix.ENOENT) {
				if errors.Is(err, unix.EACCES) {
					fmt.Fprintln(os.Stderr, "Permission denied. Try running this command with `sudo`.")
					os.Exit(1)
				}
				panic(err)
			}
		}
		fmt.Println("Charging threshold persistence reset.")

	default:
		fmt.Fprintf(
			os.Stderr,
			"There is no `%s` command. Run `piobat --help` to see a list of available commands.\n",
			subcommand,
		)
		os.Exit(1)
	}
}
