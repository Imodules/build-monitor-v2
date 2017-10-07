module Dashboards.AddEdit exposing (add, edit)

import Dashboards.Components exposing (cancelLink, dashboardColumnCountField, dashboardNameField, dateFormatField, failedIconField, runningIconField, saveButton, successIconField)
import Dashboards.Configure as DashboardConfigure
import Dashboards.Lib exposing (getDate)
import Dashboards.Models as DashboardsModel exposing (DashboardForm, EditTab(Configure, Select))
import Dashboards.Select as DashboardSelect
import Html exposing (Html, a, div, hr, li, nav, p, span, text, ul)
import Html.Attributes exposing (class, href, id, target)
import Models exposing (BuildType, Model, Project)
import Msgs exposing (DashboardMsg(ChangeCenterDateFormat, ChangeLeftDateFormat, ChangeRightDateFormat, CreateDashboard, EditDashboard, OnConfigureTabClick, OnSelectTabClick), Msg(DashboardMsg))
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
    if dashForm.id == id then
        view model dashForm (DashboardMsg EditDashboard)
    else
        text "Loading..."


view : Model -> DashboardForm -> Msg -> Html Msg
view model dashForm saveMsg =
    let
        theDate =
            getDate model
    in
    div [ id "settings" ]
        [ nav [ class "navbar nav-fixed-top" ]
            [ div [ class "navbar-item" ] [ saveButton saveMsg (not (isFormValid model.dashboards)) ]
            , div [ class "navbar-item" ] [ cancelLink ]
            ]
        , div [ class "project-area wrapper" ]
            [ dashboardNameField dashForm
            , div [ class "columns is-marginless" ]
                [ div [ class "column" ] [ dashboardColumnCountField dashForm ]
                , div [ class "column" ] [ successIconField dashForm ]
                , div [ class "column" ] [ failedIconField dashForm ]
                , div [ class "column" ] [ runningIconField dashForm ]
                ]
            , div [ class "columns is-marginless" ]
                [ div [ class "column" ] [ dateFormatField dashForm.leftDateFormat "leftDateFormat" "Left Date Format" theDate (ChangeLeftDateFormat >> DashboardMsg) ]
                , div [ class "column" ] [ dateFormatField dashForm.centerDateFormat "centerDateFormat" "Center Date Format" theDate (ChangeCenterDateFormat >> DashboardMsg) ]
                , div [ class "column" ] [ dateFormatField dashForm.rightDateFormat "rightDateFormat" "Right Date Format (UTC)" theDate (ChangeRightDateFormat >> DashboardMsg) ]
                ]
            , div [ class "notification" ]
                [ div [ class "level" ]
                    [ div [ class "level-left" ]
                        [ div [ class "level-item" ]
                            [ text "Go"
                            , a [ class "button is-link", target "_blank", href "https://github.com/rluiten/elm-date-extra/blob/master/DocFormat.md" ] [ text "here" ]
                            , text "for formatting options"
                            ]
                        ]
                    ]
                ]
            , tabs dashForm
            , maybeProjects model
            ]
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
