module Dashboards.Lib exposing (..)

import Dashboards.Models exposing (BuildConfig)
import Types exposing (Id)


configInList : Id -> List BuildConfig -> Bool
configInList id configs =
    List.any (\config -> config.id == id) configs
