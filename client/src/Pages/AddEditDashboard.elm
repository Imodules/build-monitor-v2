module Pages.AddEditDashboard exposing (..)

import Html exposing (Html, div, h4, h5, h6, hr, i, li, text, ul)
import Html.Attributes exposing (class, id)
import Models exposing (BuildType, Model, Project)
import Msgs exposing (Msg(ChangeDashboardName))
import Pages.Components exposing (iconLinkButton, textField)
import RemoteData
import Routes exposing (Route(DashboardRoute, DashboardsRoute))
import Types exposing (Id)


view : Model -> Html Msg
view model =
    div [ id "settings" ]
        [ div [ class "button-area" ] [ saveButton ]
        , div [] [ textField model.dashboardAddEdit.name "text" "dashboardName" "Dashboard Name" "fa-tachometer" ChangeDashboardName ]
        , hr [] []
        , h6 [ class "title is-6" ] [ text "Choose Builds" ]
        , div [ class "project-area" ] [ maybeProjects model model.projects ]
        ]


saveButton : Html Msg
saveButton =
    iconLinkButton "is-success" DashboardsRoute "fa-save" "Save"


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
                List.map (\bt -> buildTypeRow bt) (buildTypesByProject project.id buildTypes)
    in
    div [ class "box" ]
        [ h4 [ class "title is-4" ] [ i [ class "fa fa-cubes fa-fw" ] [], text project.name ]
        , div [ class "box" ] content
        ]


buildTypeRow : BuildType -> Html Msg
buildTypeRow buildType =
    div [ class "box" ] [ h5 [ class "title is-5" ] [ i [ class "fa fa-cube fa-fw" ] [], text buildType.name ] ]


topProjects : List Project -> List Project
topProjects =
    projectsByParent "_Root"


projectsByParent : String -> List Project -> List Project
projectsByParent parent projects =
    List.filter (\n -> n.parentProjectId == parent) projects


buildTypesByProject : Id -> List BuildType -> List BuildType
buildTypesByProject projectId buildTypes =
    List.filter (\n -> n.projectId == projectId) buildTypes
