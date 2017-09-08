module Pages.Profile exposing (..)

import Html exposing (Html, div, text)
import Html.Attributes exposing (class)
import Models exposing (Model)
import Msgs exposing (Msg)


view : Model -> Html Msg
view model =
    text "Pages.Profile.elm"
