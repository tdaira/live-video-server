package main
import (
	"github.com/tdaira/live-video-server/internal/config"
	"github.com/tdaira/live-video-server/internal/domain"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type LiveManifestHandler struct {
	cfg config.Config
}

func (l LiveManifestHandler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mediaPL, err := domain.CurrentManifest(filepath.Join(l.cfg.Source.VideoPath, "/video.m3u8"))
	if err != nil {
		log.Fatalf("Unable to get current manifest: %#v", err)
	}
	w.Header().Set("Content-Type", "application/x-mpegurl")
	_, err = w.Write(mediaPL.Encode().Bytes())
	if err != nil {
		log.Fatalf("Responce write was failed: %#v", err)
	}
}

func main() {
	// 設定の初期化
	cfg, err := config.SetupConfig()
	if err != nil {
		log.Fatalf("Unable to get config: %#v", err)
	}

	fileServer := http.StripPrefix("/video/", http.FileServer(http.Dir(cfg.Source.VideoPath)))
	err = http.ListenAndServe(cfg.Server.Port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/video/live.m3u8") {
			LiveManifestHandler{cfg}.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/video/") {
			fileServer.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	log.Fatalf("Http server crashed: %#v", err)
}
