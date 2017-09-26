module Dashboards.Configure exposing (..)

import Html exposing (Html, div, text)
import Models exposing (Model)
import Msgs exposing (Msg)
import Types exposing (Id)


view : Model -> Id -> Html Msg
view model id =
    div [] [ text ("Dashboards.Configure.elm: " ++ id) ]
