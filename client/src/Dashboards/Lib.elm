module Dashboards.Lib exposing (..)

import Dashboards.Models exposing (BuildConfigForm)
import Types exposing (Id)


configInList : Id -> List BuildConfigForm -> Bool
configInList id configs =
    List.any (\config -> config.id == id) configs
