package main

import (
    "github.com/caarlos0/env"
    "github.com/jmacelroy/junebug"
    "github.com/rs/zerolog"
    "os"
    "time"
)

var log = zerolog.New(os.Stdout).With().
    Timestamp().
    Logger()

type appCfg struct {
    // server config
    ConfluenceIssuer   string `env:"CONFLUENCE_ISSUER" envDefault:"junebug-confluence-writer"`
    ConfluenceAudience string `env:"CONFLUENCE_AUDIENCE,required"`
    ConfluenceSecret   string `env:"CONFLUENCE_SECRET,required"`
    ConfluenceQsh      string `env:"CONFLUENCE_QUERY_HASH" envDefault:"79ccdc28e5f25b5ee15d7dfcfcc7977848375ea060057ed4e17b7b4aef756694"`
    ConfluenceSpace    string `env:"CONFLUENCE_SPACE" envDefault:"MEET"`
    ConfluenceURL      string `env:"CONFLUENCE_URL,required"`
}

func main() {
    app := appCfg{}
    err := env.Parse(&app)
    if err != nil {
        log.Fatal().Err(err).Msg("service is misconfigured")
    }
    tg := &junebug.ConfluenceTokenGenerator{Lifetime: time.Hour * 24, Issuer: app.ConfluenceIssuer, Audience: app.ConfluenceAudience, Secret: app.ConfluenceSecret}
    cc := junebug.NewConfluenceClient(app.ConfluenceURL, app.ConfluenceSpace, app.ConfluenceQsh, tg)

    token, err := tg.CreateJWT(app.ConfluenceQsh)
    if err != nil {
        log.Fatal().Err(err).Msg("token failed to be created")
    }
    log.Info().Msgf("testing token: %s", token)

    cp, err := cc.CreatePage("some more testing")
    if err != nil {
        log.Fatal().Err(err).Msg("confluence page failed to be created")
    }
    log.Info().Msgf("page results: %s %s %s", cp.ID, cp.GetWebURL(), cp.GetEditURL())

}
