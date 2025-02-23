// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/navidrome/navidrome/core"
	"github.com/navidrome/navidrome/core/transcoder"
	"github.com/navidrome/navidrome/persistence"
	"github.com/navidrome/navidrome/scanner"
	"github.com/navidrome/navidrome/server"
	"github.com/navidrome/navidrome/server/app"
	"github.com/navidrome/navidrome/server/events"
	"github.com/navidrome/navidrome/server/subsonic"
	"sync"
)

// Injectors from wire_injectors.go:

func CreateServer(musicFolder string) *server.Server {
	dataStore := persistence.New()
	serverServer := server.New(dataStore)
	return serverServer
}

func CreateAppRouter() *app.Router {
	dataStore := persistence.New()
	broker := GetBroker()
	router := app.New(dataStore, broker)
	return router
}

func CreateSubsonicAPIRouter() *subsonic.Router {
	dataStore := persistence.New()
	artworkCache := core.GetImageCache()
	artwork := core.NewArtwork(dataStore, artworkCache)
	transcoderTranscoder := transcoder.New()
	transcodingCache := core.GetTranscodingCache()
	mediaStreamer := core.NewMediaStreamer(dataStore, transcoderTranscoder, transcodingCache)
	archiver := core.NewArchiver(dataStore)
	players := core.NewPlayers(dataStore)
	externalMetadata := core.NewExternalMetadata(dataStore)
	scanner := GetScanner()
	router := subsonic.New(dataStore, artwork, mediaStreamer, archiver, players, externalMetadata, scanner)
	return router
}

func createScanner() scanner.Scanner {
	dataStore := persistence.New()
	artworkCache := core.GetImageCache()
	artwork := core.NewArtwork(dataStore, artworkCache)
	cacheWarmer := core.NewCacheWarmer(artwork, artworkCache)
	broker := GetBroker()
	scannerScanner := scanner.New(dataStore, cacheWarmer, broker)
	return scannerScanner
}

func createBroker() events.Broker {
	broker := events.NewBroker()
	return broker
}

// wire_injectors.go:

var allProviders = wire.NewSet(core.Set, subsonic.New, app.New, persistence.New)

// Scanner must be a Singleton
var (
	onceScanner     sync.Once
	scannerInstance scanner.Scanner
)

func GetScanner() scanner.Scanner {
	onceScanner.Do(func() {
		scannerInstance = createScanner()
	})
	return scannerInstance
}

// Broker must be a Singleton
var (
	onceBroker     sync.Once
	brokerInstance events.Broker
)

func GetBroker() events.Broker {
	onceBroker.Do(func() {
		brokerInstance = createBroker()
	})
	return brokerInstance
}
