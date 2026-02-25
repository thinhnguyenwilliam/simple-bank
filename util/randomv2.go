package util

import (
	"math/rand"
	"time"

	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomFloat64(min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}

func randomLaptopBrand() string {
	brands := []string{
		"Lenovo",
		"Dell",
		"Apple",
		"HP",
		"ASUS",
	}
	return brands[rng.Intn(len(brands))]
}

func randomLaptopName() string {
	names := []string{
		"ThinkPad X1 Carbon",
		"MacBook Pro",
		"XPS 15",
		"Spectre x360",
		"ROG Zephyrus",
	}
	return names[rng.Intn(len(names))]
}

type resolution struct {
	width  uint32
	height uint32
}

func randomScreenSize() float32 {
	sizes := []float32{13.3, 14.0, 15.6, 16.0}
	return sizes[rng.Intn(len(sizes))]
}

func randomResolution(size float32) resolution {
	switch size {
	case 13.3:
		return resolution{1920, 1080}
	case 14.0:
		return resolution{1920, 1200}
	case 15.6:
		return resolution{2560, 1440}
	case 16.0:
		return resolution{3840, 2400}
	default:
		return resolution{1920, 1080}
	}
}

func randomPanel() pb.Screen_Panel {
	panels := []pb.Screen_Panel{
		pb.Screen_PANEL_IPS,
		pb.Screen_PANEL_OLED,
	}
	return panels[rng.Intn(len(panels))]
}

func randomHDDSize() uint64 {
	sizes := []uint64{1024, 2048, 4096}
	return sizes[rng.Intn(len(sizes))]
}

func randomSSDSize() uint64 {
	sizes := []uint64{256, 512, 1024, 2048}
	return sizes[rng.Intn(len(sizes))]
}

func randomRAMSize() uint64 {
	sizes := []uint64{8, 16, 32, 64}
	return sizes[rng.Intn(len(sizes))]
}

func randomGPUBrand() string {
	brands := []string{"NVIDIA", "AMD", "Intel"}
	return brands[rng.Intn(len(brands))]
}

func randomGPUInfo(brand string) (name string, minGhz, maxGhz float64, vram uint64) {
	switch brand {
	case "NVIDIA":
		return "RTX 4070", 1.9, 2.5, 12
	case "AMD":
		return "Radeon RX 7800 XT", 1.8, 2.4, 16
	case "Intel":
		return "Arc A770", 2.1, 2.4, 16
	default:
		return "Unknown GPU", 1.5, 2.0, 8
	}
}

func randomCPUBrand() string {
	brands := []string{"Intel", "AMD", "Apple"}
	return brands[rng.Intn(len(brands))]
}

func randomCPUInfo(brand string) (name string, cores, threads uint32, minGhz, maxGhz float64) {
	switch brand {
	case "Intel":
		return "Core i7-13700H", 14, 20, 2.4, 5.0
	case "AMD":
		return "Ryzen 7 7840HS", 8, 16, 3.8, 5.1
	case "Apple":
		return "M2 Pro", 12, 12, 3.2, 3.6
	default:
		return "Unknown", 4, 8, 2.0, 3.0
	}
}

func randomBool() bool {
	return rng.Intn(2) == 1
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	layouts := []pb.Keyboard_Layout{
		pb.Keyboard_LAYOUT_QWERTY,
		pb.Keyboard_LAYOUT_QWERTZ,
		pb.Keyboard_LAYOUT_AZERTY,
	}
	return layouts[rng.Intn(len(layouts))]
}
