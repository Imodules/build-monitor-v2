module Dashboards.AddEdit exposing (..)

import Dashboards.Models as DashboardsModel
import Html exposing (Html, button, div, h4, h5, h6, hr, i, li, span, text, ul)
import Html.Attributes exposing (class, disabled, id)
import Html.Events exposing (onClick)
import Models exposing (BuildType, Model, Project)
import Msgs exposing (DashboardMsg(..), Msg(DashboardMsg))
import Pages.Components exposing (icon, iconLinkButton, textField)
import RemoteData
import Routes exposing (Route(DashboardRoute, DashboardsRoute))
import Types exposing (Id)


view : Model -> Html Msg
view model =
    div [ id "settings" ]
        [ div [ class "button-area" ] [ saveButton model.dashboards ]
        , div [] [ textField model.dashboards.dashboardForm.name "text" "dashboardName" "Dashboard Name" "fa-tachometer" (ChangeDashboardName >> DashboardMsg) ]
        , hr [] []
        , h6 [ class "title is-6" ] [ text "Choose Builds" ]
        , div [ class "project-area" ] [ maybeProjects model model.projects ]
        ]


saveButton : DashboardsModel.Model -> Html Msg
saveButton model =
    let
        disableButton =
            not (isFormValid model)
    in
    button [ class "button is-success", disabled disableButton, onClick (DashboardMsg CreateDashboard) ]
        [ icon "fa fa-check-square-o"
        , span [] [ text "Save" ]
        ]


isFormValid : DashboardsModel.Model -> Bool
isFormValid model =
    model.dashboardForm.name.isValid


maybeProjects : Model -> RemoteData.WebData (List Project) -> Html Msg
maybeProjects model response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success projects ->
            maybeBuildTypes model projects model.buildTypes

        RemoteData.Failure error ->
            text (toString error)


maybeBuildTypes : Model -> List Project -> RemoteData.WebData (List BuildType) -> Html Msg
maybeBuildTypes model projects response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success buildTypes ->
            projectTree model projects buildTypes

        RemoteData.Failure error ->
            text (toString error)


projectTree : Model -> List Project -> List BuildType -> Html Msg
projectTree model projects buildTypes =
    let
        content =
            List.map (\p -> topProjectRow model p projects buildTypes) (topProjects projects)
    in
    div [] content


topProjectRow : Model -> Project -> List Project -> List BuildType -> Html Msg
topProjectRow model project projects buildTypes =
    let
        childProjects =
            projectsByParent project.id projects

        content =
            if List.length childProjects > 0 then
                List.map (\p -> topProjectRow model p projects buildTypes) childProjects
            else
                List.map (\bt -> buildTypeRow model bt) (buildTypesByProject project.id buildTypes)
    in
    div [ class "box" ]
        [ h4 [ class "title is-4" ] [ icon "fa fa-cubes fa-fw", text project.name ]
        , div [ class "box" ] content
        ]


buildTypeRow : Model -> BuildType -> Html Msg
buildTypeRow model buildType =
    let
        handleClick =
            buildType.id |> (ClickBuildType >> DashboardMsg)

        isSelected =
            List.member buildType.id model.dashboards.dashboardForm.buildTypeIds

        rowClass =
            if isSelected then
                "box buildType selected"
            else
                "box buildType"
    in
    div [ class rowClass, onClick handleClick ] [ h5 [ class "title is-5" ] [ icon "fa fa-cube fa-fw", text buildType.name ] ]


topProjects : List Project -> List Project
topProjects =
    projectsByParent "_Root"


projectsByParent : String -> List Project -> List Project
projectsByParent parent projects =
    List.filter (\n -> n.parentProjectId == parent) projects


buildTypesByProject : Id -> List BuildType -> List BuildType
buildTypesByProject projectId buildTypes =
    List.filter (\n -> n.projectId == projectId) buildTypes
