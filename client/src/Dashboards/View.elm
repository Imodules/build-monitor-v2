module Dashboards.View exposing (..)

import Html exposing (Html, a, div, i, section, text)
import Html.Attributes exposing (class, href, id)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (iconLink)
import Routes exposing (Route(DashboardsRoute))


view : Model -> Html Msg
view model =
    div [ id "dashboard" ] [ configLink ]


configLink : Html Msg
configLink =
    div [ id "configLink" ] [ iconLink "button is-link" DashboardsRoute "fa fa-cogs" ]
