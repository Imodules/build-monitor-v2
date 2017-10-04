module Dashboards.AddEdit exposing (add, edit)

import Dashboards.Components exposing (cancelButton, dashboardNameField, dashboardColumnCountField, saveButton)
import Dashboards.Configure as DashboardConfigure
import Dashboards.Models as DashboardsModel exposing (DashboardForm, EditTab(Configure, Select))
import Dashboards.Select as DashboardSelect
import Html exposing (Html, a, div, hr, li, span, text, ul)
import Html.Attributes exposing (class, id)
import Models exposing (BuildType, Model, Project)
import Msgs exposing (DashboardMsg(CreateDashboard, EditDashboard, OnConfigureTabClick, OnSelectTabClick), Msg(DashboardMsg))
import Pages.Components exposing (icon)
import RemoteData
import Routing exposing (onLinkClick)
import Types exposing (Id)


add : Model -> Html Msg
add model =
    let
        dashForm =
            model.dashboards.dashboardForm
    in
    view model dashForm (DashboardMsg CreateDashboard)


edit : Model -> Id -> Html Msg
edit model id =
    let
        dashForm =
            model.dashboards.dashboardForm
    in
    view model dashForm (DashboardMsg EditDashboard)


view : Model -> DashboardForm -> Msg -> Html Msg
view model dashForm saveMsg =
    div [ id "settings" ]
        [ div [ class "button-area" ] [ saveButton saveMsg (not (isFormValid model.dashboards)), cancelButton ]
        , dashboardNameField dashForm
        , dashboardColumnCountField dashForm
        , tabs dashForm
        , div [ class "project-area" ] [ maybeProjects model ]
        ]


isFormValid : DashboardsModel.Model -> Bool
isFormValid model =
    model.dashboardForm.isDirty && model.dashboardForm.name.isValid


maybeProjects : Model -> Html Msg
maybeProjects model =
    case model.projects of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success projects ->
            maybeBuildTypes model projects

        RemoteData.Failure error ->
            text (toString error)


maybeBuildTypes : Model -> List Project -> Html Msg
maybeBuildTypes model projects =
    case model.buildTypes of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success buildTypes ->
            projectArea model projects buildTypes

        RemoteData.Failure error ->
            text (toString error)


tabs : DashboardForm -> Html Msg
tabs dashForm =
    div [ class "tabs is-fullwidth is-medium" ]
        [ ul []
            [ tab "Select" "fa fa-check-square-o fa-fw" (dashForm.tab == Select) OnSelectTabClick
            , tab "Configure" "fa fa-cog fa-fw" (dashForm.tab == Configure) OnConfigureTabClick
            ]
        ]


tab : String -> String -> Bool -> DashboardMsg -> Html Msg
tab tabText img isActive msg =
    let
        tabClass =
            if isActive then
                "is-active"
            else
                ""
    in
    li [ class tabClass ] [ a [ onLinkClick (DashboardMsg msg) ] [ icon img, span [] [ text tabText ] ] ]


projectArea : Model -> List Project -> List BuildType -> Html Msg
projectArea model projects buildTypes =
    if model.dashboards.dashboardForm.tab == Select then
        DashboardSelect.view model projects buildTypes
    else
        DashboardConfigure.view model projects buildTypes
