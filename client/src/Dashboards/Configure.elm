module Dashboards.Configure exposing (..)

import Dashboards.Components exposing (cancelButton, dashboardNameField, saveButton)
import Dashboards.Lib exposing (getBuildPath)
import Dashboards.Models as DashboardsModel exposing (BuildConfig, BuildConfigForm)
import Html exposing (Html, div, h6, hr, text)
import Html.Attributes exposing (class, id)
import Models exposing (BuildType, Model, Project, initialProject)
import Msgs exposing (DashboardMsg(ChangeBuildAbbreviation, ChangeDashboardName, EditDashboard), Msg(DashboardMsg))
import Pages.Components exposing (textBox, textField)
import RemoteData
import Types exposing (Id)


-- TODO: Change this so it is not a route, but a view that can be toggled with a tab


view : Model -> Id -> Html Msg
view model id_ =
    div [ id "configure" ]
        [ div [ class "button-area" ]
            [ saveButton (DashboardMsg EditDashboard) (not (isFormValid model.dashboards))
            , cancelButton
            ]
        , dashboardNameField model.dashboards.dashboardForm
        , hr [] []
        , h6 [ class "title is-6" ] [ text "Configure Builds" ]
        , maybeProjects model
        ]


isFormValid : DashboardsModel.Model -> Bool
isFormValid model =
    let
        isValid =
            not (List.any (\c -> not c.abbreviation.isValid) model.dashboardForm.buildConfigs)
    in
    model.dashboardForm.name.isValid && isValid


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
            configsArea model projects buildTypes

        RemoteData.Failure error ->
            text (toString error)


configsArea : Model -> List Project -> List BuildType -> Html Msg
configsArea model projects buildTypes =
    let
        content =
            List.map (\i -> configRow i projects buildTypes) model.dashboards.dashboardForm.buildConfigs
    in
    div [] content


configRow : BuildConfigForm -> List Project -> List BuildType -> Html Msg
configRow config projects buildTypes =
    let
        buildPath =
            getBuildPath config.id projects buildTypes
    in
    div [ class "box" ]
        [ text buildPath
        , div [ class "level" ]
            [ div [ class "level-left" ]
                [ textBox config.abbreviation (ChangeBuildAbbreviation config.id >> Msgs.DashboardMsg)
                ]
            ]
        ]
