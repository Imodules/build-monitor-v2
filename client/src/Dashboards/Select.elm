module Dashboards.Select exposing (view)

import Dashboards.Lib exposing (configInList)
import Html exposing (Html, div, h4, h5, text)
import Html.Attributes exposing (class)
import Html.Events exposing (onClick)
import Models exposing (BuildType, Model, Project)
import Msgs exposing (DashboardMsg(ClickBuildType), Msg(DashboardMsg))
import Pages.Components exposing (icon)
import Types exposing (Id)


view : Model -> List Project -> List BuildType -> Html Msg
view model projects buildTypes =
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
            configInList buildType.id model.dashboards.dashboardForm.buildConfigs

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
