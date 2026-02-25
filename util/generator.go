package util

import (
	"time"

	"github.com/google/uuid"
	pb "github.com/thinhcompany/simple-bank/pb/pb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewLaptop() *pb.Laptop {
	laptop := &pb.Laptop{
		Id:    uuid.NewString(),
		Brand: randomLaptopBrand(),
		Name:  randomLaptopName(),

		Cpu:  NewCPU(),
		Ram:  NewRAM(),
		Gpus: []*pb.GPU{NewGPU()},
		Storages: []*pb.Storage{
			NewSSD(),
			NewHDD(),
		},
		Screen:   NewScreen(),
		Keyboard: NewKeyboard(),

		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.2, 3.0),
		},

		PriceUsd:    randomFloat64(800, 4000),
		ReleaseYear: uint32(time.Now().Year()),
		UpdatedAt:   timestamppb.Now(),
	}

	return laptop
}

func NewScreen() *pb.Screen {
	size := randomScreenSize()
	res := randomResolution(size)

	return &pb.Screen{
		SizeInch: size,
		Resolution: &pb.Screen_Resolution{
			Width:  res.width,
			Height: res.height,
		},
		Panel:      randomPanel(),
		Multitouch: randomBool(),
	}
}

func NewHDD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_DRIVER_HDD,
		Memory: &pb.Memory{
			Value: randomHDDSize(),
			Unit:  pb.Memory_UNIT_GIGABYTE,
		},
	}
}

func NewSSD() *pb.Storage {
	return &pb.Storage{
		Driver: pb.Storage_DRIVER_SSD,
		Memory: &pb.Memory{
			Value: randomSSDSize(),
			Unit:  pb.Memory_UNIT_GIGABYTE,
		},
	}
}

func NewRAM() *pb.Memory {
	return &pb.Memory{
		Value: randomRAMSize(),
		Unit:  pb.Memory_UNIT_GIGABYTE,
	}
}

func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name, minGhz, maxGhz, vram := randomGPUInfo(brand)

	return &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: &pb.Memory{
			Value: vram,
			Unit:  pb.Memory_UNIT_GIGABYTE,
		},
	}
}

func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name, cores, threads, minGhz, maxGhz := randomCPUInfo(brand)

	return &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   cores,
		NumberThreads: threads,
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}
}

func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}
