package main

import (
    "github.com/aaronkvanmeerten/junebug"
    "github.com/caarlos0/env"
    "github.com/rs/zerolog"
    "os"
    "time"
)

var log = zerolog.New(os.Stdout).With().
    Timestamp().
    Logger()

type appCfg struct {
    // server config
    Issuer   string `env:"CONFLUENCE_ISSUER" envDefault:"jitsi"`
    Audience string `env:"CONFLUENCE_AUDIENCE" envDefault:"confluence"`
    Secret   string `env:"CONFLUENCE_SECRET" envDefault:"s3krit"`
    Qsh      string `env:"CONFLUENCE_QUERY_HASH" envDefault:"79ccdc28e5f25b5ee15d7dfcfcc7977848375ea060057ed4e17b7b4aef756694"`
}

func main() {
    app := appCfg{}
    err := env.Parse(&app)
    if err != nil {
        log.Fatal().Err(err).Msg("service is misconfigured")
    }
    tg := junebug.TokenGenerator{
        Lifetime: time.Hour * 24,
        Secret:   app.Secret,
        Issuer:   app.Issuer,
        Audience: app.Audience,
    }
    token, err := tg.CreateJWT(app.Qsh)
    if err != nil {
        log.Fatal().Err(err).Msg("token failed to be created")
    }
    log.Info().Msgf("testing token: %s", token)
    cc := junebug.ConfluenceClient{URL: "https://pi-dev-sandbox.atlassian.net/wiki"}
    cp, err := cc.CreatePage(token, "some more testing", "MEET")
    if err != nil {
        log.Fatal().Err(err).Msg("confluence page failed to be created")
    }
    log.Info().Msgf("page results: %s %s %s", cp.ID, cp.GetWebURL(), cp.GetEditURL())

}
