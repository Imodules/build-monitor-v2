module Dashboards.List exposing (..)

import Dashboards.Models exposing (Dashboard)
import Html exposing (Html, div, h4, h5, i, li, text, ul)
import Html.Attributes exposing (class, id)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (iconLinkButton)
import RemoteData
import Routes exposing (Route(NewDashboardRoute))


view : Model -> Html Msg
view model =
    div [ id "dashboards" ]
        [ div [ class "button-area" ] [ newDashboardButton ]
        , div [ class "project-area" ] [ maybeDashboards model model.dashboards.dashboards ]
        ]


newDashboardButton : Html Msg
newDashboardButton =
    iconLinkButton "is-success" NewDashboardRoute "fa-plus" "New Dashboard"


maybeDashboards : Model -> RemoteData.WebData (List Dashboard) -> Html Msg
maybeDashboards model response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success dashboards ->
            dashboardList model dashboards

        RemoteData.Failure error ->
            text (toString error)


dashboardList : Model -> List Dashboard -> Html Msg
dashboardList model dashboards =
    div [] [ text "Dashboards" ]
