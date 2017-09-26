module Dashboards.Configure exposing (..)

import Dashboards.Components exposing (cancelButton, dashboardNameField, saveButton)
import Dashboards.Models exposing (BuildConfig, BuildConfigForm)
import Html exposing (Html, div, h6, hr, text)
import Html.Attributes exposing (class, id)
import List.Extra exposing (find)
import Models exposing (BuildType, Model, Project, initialProject)
import Msgs exposing (DashboardMsg(ChangeDashboardName, EditDashboard), Msg(DashboardMsg))
import Pages.Components exposing (textBox, textField)
import RemoteData
import Types exposing (Id)


view : Model -> Id -> Html Msg
view model id_ =
    div [ id "configure" ]
        [ div [ class "button-area" ]
            [ saveButton (DashboardMsg EditDashboard) True
            , cancelButton
            ]
        , dashboardNameField model.dashboards.dashboardForm
        , hr [] []
        , h6 [ class "title is-6" ] [ text "Configure Builds" ]
        , maybeProjects model
        ]


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
            getBuildPath config projects buildTypes

        prefix =
            getDefaultPrefix buildPath
    in
    div [ class "box" ]
        [ text buildPath
        , div [ class "level" ]
            [ div [ class "level-left" ]
                [ textBox config.abbreviation (ChangeDashboardName >> Msgs.DashboardMsg) -- TODO: Need proper message
                ]
            ]
        ]


getBuildPath : BuildConfigForm -> List Project -> List BuildType -> String
getBuildPath config projects buildTypes =
    let
        maybeBuildType =
            find (\i -> i.id == config.id) buildTypes
    in
    case maybeBuildType of
        Just buildType ->
            getProjectPath buildType.projectId projects ++ " / " ++ buildType.name

        _ ->
            ""


getProjectPath : Id -> List Project -> String
getProjectPath id projects =
    let
        maybeParentProject =
            find (\i -> i.id == id) projects

        parentProject =
            case maybeParentProject of
                Just project ->
                    project

                _ ->
                    initialProject
    in
    if parentProject.parentProjectId /= "_Root" then
        parentProject.name ++ " / " ++ getProjectPath parentProject.parentProjectId projects
    else
        parentProject.name


getDefaultPrefix : String -> String
getDefaultPrefix path =
    let
        paths =
            String.split " / " path

        parts =
            List.map (\p -> getPathPart p ++ "-") paths
    in
    String.dropRight 1 (String.concat parts)


getPathPart : String -> String
getPathPart s =
    let
        words =
            String.words s

        letters =
            List.map getFirstLetter words
    in
    String.concat letters


getFirstLetter : String -> String
getFirstLetter s =
    String.slice 0 1 s
