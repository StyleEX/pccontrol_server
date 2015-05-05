package main

import (
	"fmt"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"log"
	"net/http"
	"os/exec"
)

// TODO: Сервисы ержать в отдельном пакете.
// VolumeService - управление звуком.
type VolumeLevelArgs struct {
	Level int
}

type VolumeService struct{}

func (h *VolumeService) SetLevel(r *http.Request, args *VolumeLevelArgs, res *json2.EmptyResponse) error {
	log.Printf("Set volume level %d%%", args.Level)
	exec.Command("amixer", "sset", "'Master'", fmt.Sprintf("%d%%", args.Level)).Run()
	return nil
}

type ShutdownArgs struct {
	Minutes int
}

// SystemService - управление системой
type SystemService struct{}

func (h *SystemService) Shutdown(r *http.Request, args *ShutdownArgs, res *json2.EmptyResponse) error {
	if args.Minutes > 0 {
		log.Printf("Set shutdown after %d minutes", args.Minutes)
		output, _ := exec.Command("shutdown", "-P", "-h", fmt.Sprintf("+%d", args.Minutes)).CombinedOutput()
		log.Println(string(output[:]))
	} else {
		log.Println("Disable shutdown")
		exec.Command("shutdown", "-c").Run()
	}
	return nil
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")

	// TODO: Регистрировать платформенно зависимые ханлеры.
	s.RegisterService(new(VolumeService), "")
	s.RegisterService(new(SystemService), "")

	http.Handle("/rpc", s)

	// TODO: Использовать не HTTP, а держать TCP соединение.
	http.ListenAndServe(":10000", nil)
}
