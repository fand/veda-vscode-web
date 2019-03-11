package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/getlantern/systray"
	"github.com/phayes/freeport"
)

func launchCodeServer() (*exec.Cmd, int) {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	// Run code-server
	cmd := exec.Command("code-server", "--no-auth", "-p", strconv.Itoa(port))
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return cmd, port
}

func launchWebServer(wrapperPort, port int) *http.Server {
	// Start HTTP server for IPC
	srv := &http.Server{Addr: ":" + strconv.Itoa(wrapperPort)}

	http.HandleFunc("/port", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, strconv.Itoa(port))
	})

	// Run server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	return srv
}

func initMenu() {
	// systray.SetIcon(getIcon("assets/clock.ico"))
	systray.SetTitle("VEDA")
	systray.SetTooltip("VEDA for VSCode Web Server")
	mQuit := systray.AddMenuItem("Quit", "Quit VSCode Web Server")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func main() {
	wrapperPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	println(wrapperPort)

	cmd, port := launchCodeServer()

	srv := launchWebServer(wrapperPort, port)

	systray.Run(initMenu, func() {
		// Kill HTTP server
		err = srv.Close()
		if err != nil {
			log.Fatal(err)
		}

		// Kill code-server before quit
		cmd.Process.Kill()
	})
}
