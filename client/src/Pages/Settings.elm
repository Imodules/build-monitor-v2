module Pages.Settings exposing (..)

import Html exposing (Html, div, li, text, ul)
import Html.Attributes exposing (class)
import Models exposing (Model, Project)
import Msgs exposing (Msg)
import Pages.Components exposing (iconLinkButton)
import RemoteData
import Routes exposing (Route(DashboardRoute))


view : Model -> Html Msg
view model =
    div []
        [ div [] [ dashboardButton ]
        , div [] [ maybeProjects model model.projects ]
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
        -- TODO: Sort by name
        content =
            List.indexedMap projectRow (topProjects projects)
    in
    ul [] content


projectRow : Int -> Project -> Html Msg
projectRow id project =
    li [] [ text (toString id ++ " " ++ project.name) ]


topProjects : List Project -> List Project
topProjects projects =
    let
        hasRootParent n =
            n.parentObjectId == "_Root"
    in
    List.filter hasRootParent projects
