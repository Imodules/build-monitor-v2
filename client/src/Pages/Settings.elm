module Pages.Settings exposing (..)

import Html exposing (Html, div, li, text, ul)
import Html.Attributes exposing (class, id)
import Models exposing (Model, Project)
import Msgs exposing (Msg)
import Pages.Components exposing (iconLinkButton)
import RemoteData
import Routes exposing (Route(DashboardRoute))


view : Model -> Html Msg
view model =
    div [ id "settings" ]
        [ div [ class "button-area" ] [ dashboardButton ]
        , div [ class "project-area" ] [ maybeProjects model model.projects ]
        ]


dashboardButton : Html Msg
dashboardButton =
    iconLinkButton "is-primary" DashboardRoute "fa-tachometer" "Dashboard"


maybeProjects : Model -> RemoteData.WebData (List Project) -> Html Msg
maybeProjects model response =
    case response of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success projects ->
            projectTree model projects

        RemoteData.Failure error ->
            text (toString error)


projectTree : Model -> List Project -> Html Msg
projectTree model projects =
    let
        content =
            List.map (\p -> topProjectRow p projects) (topProjects projects)
    in
    div [] content


topProjectRow : Project -> List Project -> Html Msg
topProjectRow project projects =
    let
        content =
            List.map (\p -> topProjectRow p projects) (projectsByParent project.id projects)
    in
    div [ class "box" ]
        [ text project.name
        , div [ class "box" ] content
        ]


topProjects : List Project -> List Project
topProjects =
    projectsByParent "_Root"


projectsByParent : String -> List Project -> List Project
projectsByParent parent projects =
    List.filter (\n -> n.parentObjectId == parent) projects
