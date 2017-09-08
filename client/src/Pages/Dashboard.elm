module Pages.Dashboard exposing (..)

import Html exposing (Html, div, section, text)
import Html.Attributes exposing (class, id)
import Models exposing (Model)
import Msgs exposing (Msg)


view : Model -> Html Msg
view model =
    text "Pages.Dashboard.elm"
