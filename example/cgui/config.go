package main

import (
    "gopkg.in/fzerorubigd/onion.v2"
    "gopkg.in/fzerorubigd/onion.v2/extraenv"
)

func readConfig() *onion.Onion {
    dl := onion.NewDefaultLayer()

    dl.SetDefault("APP_ID", "48841")
    dl.SetDefault("APP_TOKEN", "3151c01673d412c18c055f089128be50")
    
    cfg := onion.New()
    cfg.AddLayer(dl)
    cfg.AddLazyLayer(extraenv.NewExtraEnvLayer("TG"))

    return cfg
}
