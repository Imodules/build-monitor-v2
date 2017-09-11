module Pages.Settings exposing (..)

import Html exposing (Html, div, text)
import Html.Attributes exposing (class)
import Models exposing (Model)
import Msgs exposing (Msg)


view : Model -> Html Msg
view model =
    text "Pages.Settings.elm"
