module Dashboards.List exposing (..)

import Dashboards.Lib exposing (isOwner)
import Dashboards.Models exposing (Dashboard)
import Html exposing (Html, div, h4, h5, i, li, text, ul)
import Html.Attributes exposing (class, disabled, id)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (icon, iconLinkButton, loginBanner, refreshProjectsButton)
import RemoteData
import Routes exposing (Route(DashboardRoute, EditDashboardRoute, NewDashboardRoute))
import Routing exposing (isLoggedIn)
import Types exposing (Id, Owner)


view : Model -> Html Msg
view model =
    div [ id "dashboards" ]
        [ loginBanner model
        , div [ class "button-area" ] [ newDashboardButton model, refreshProjectsButton model ]
        , div [ class "project-area" ] [ maybeDashboards model model.dashboards.dashboards ]
        ]


newDashboardButton : Model -> Html Msg
newDashboardButton model =
    let
        isDisabled =
            not (isLoggedIn model)
    in
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
    let
        content =
            List.map (\d -> dashboardListItem model d) dashboards
    in
    div [] content


dashboardListItem : Model -> Dashboard -> Html Msg
dashboardListItem model dashboard =
    div [ class "box" ]
        [ div [ class "level" ]
            [ div [ class "level-left" ]
                [ div [ class "level-item" ] [ h4 [ class "title is-4" ] [ icon "fa fa-tachometer fa-fw", text dashboard.name ] ]
                ]
            , div [ class "level-right" ]
                [ editButton model dashboard
                , viewButton dashboard.id
                ]
            ]
        ]


viewButton : Id -> Html Msg
viewButton id =
    div [ class "level-item" ] [ iconLinkButton "is-primary" (DashboardRoute id) "fa-eye" "View" ]


editButton : Model -> Dashboard -> Html Msg
editButton model dashboard =
    let
        disableEdit =
            not (isOwner model dashboard.owner)
    in
    if disableEdit then
        div [] []
    else
        div [ class "level-item" ] [ iconLinkButton "is-info" (EditDashboardRoute dashboard.id) "fa-edit" "Edit" ]
