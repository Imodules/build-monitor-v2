module Dashboards.List exposing (..)

import Dashboards.Components exposing (cancelButton, deleteButton)
import Dashboards.Lib exposing (isOwner)
import Dashboards.Models exposing (Dashboard)
import Html exposing (Html, div, h4, h5, i, li, nav, text, ul)
import Html.Attributes exposing (class, disabled, id)
import Models exposing (Model)
import Msgs exposing (DashboardMsg(CancelDeleteDashboard, ConfirmDeleteDashboard, DeleteDashboard), Msg)
import Pages.Components exposing (icon, iconLinkButton, loginBanner, refreshProjectsButton)
import RemoteData
import Routes exposing (Route(DashboardRoute, EditDashboardRoute, NewDashboardRoute))
import Routing exposing (isLoggedIn)
import Types exposing (Id, Owner)


view : Model -> Html Msg
view model =
    div [ id "dashboards" ]
        [ nav [ class "navbar nav-fixed-top" ]
            [ div [ class "navbar-item" ] [ newDashboardButton model ]
            , div [ class "navbar-item" ] [ refreshProjectsButton model ]
            ]
        , div [ class "project-area" ] [ loginBanner model, maybeDashboards model model.dashboards.dashboards ]
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
    let
        buttonAreaContent =
            if model.dashboards.deleteDashboardId == dashboard.id then
                confirmDashboardDelete
            else
                dashboardButtons
    in
    div [ class "box" ]
        [ div [ class "level" ]
            [ div [ class "level-left" ]
                [ div [ class "level-item" ] [ h4 [ class "title is-4" ] [ icon "fa fa-tachometer fa-fw", text dashboard.name ] ]
                ]
            , buttonAreaContent model dashboard
            ]
        ]


dashboardButtons : Model -> Dashboard -> Html Msg
dashboardButtons model dashboard =
    div [ class "level-right buttons" ]
        [ editButton model dashboard
        , viewButton dashboard.id
        , deleteButton_ model dashboard
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
        div [ class "level-item" ] [ iconLinkButton "is-default" (EditDashboardRoute dashboard.id) "fa-edit" "Edit" ]


deleteButton_ : Model -> Dashboard -> Html Msg
deleteButton_ model dashboard =
    let
        disableButton =
            not (isOwner model dashboard.owner)
    in
    div [ class "level-item" ] [ deleteButton (Msgs.DashboardMsg (DeleteDashboard dashboard.id)) disableButton ]


confirmDashboardDelete : Model -> Dashboard -> Html Msg
confirmDashboardDelete model dashboard =
    let
        disableButton =
            not (isOwner model dashboard.owner)
    in
    div [ class "level-right notification is-warning is-marginless" ]
        [ div [ class "level-item" ] [ text "Are you shour you want to delete this dashboard?" ]
        , div [ class "level-item buttons" ]
            [ deleteButton (Msgs.DashboardMsg (ConfirmDeleteDashboard dashboard.id)) disableButton
            , cancelButton (Msgs.DashboardMsg CancelDeleteDashboard) False
            ]
        ]
