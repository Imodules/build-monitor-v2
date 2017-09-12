module Pages.Settings exposing (..)

import Html exposing (Html, div, text)
import Html.Attributes exposing (class)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (iconLinkButton)
import Routes exposing (Route(DashboardRoute))


view : Model -> Html Msg
view model =
    dashboardButton


dashboardButton : Html Msg
dashboardButton =
    iconLinkButton "is-primary" DashboardRoute "fa-tachometer" "Dashboard"
