module Dashboards.View exposing (..)

import Dashboards.Models exposing (Branch, Build, BuildStatus(Failure, Success), ConfigDetail, DashboardDetails)
import Html exposing (Html, a, div, h2, h4, i, section, text)
import Html.Attributes exposing (class, href, id)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (iconLink)
import RemoteData
import Routes exposing (Route(DashboardsRoute))


view : Model -> Html Msg
view model =
    div [ id "dashboard" ]
        [ topBar model
        , maybeDetails model
        ]


topBar : Model -> Html Msg
topBar model =
    div [ class "level top-bar" ]
        [ div [ class "level-left" ] []
        , div [ class "level-right" ] [ div [ class "level-item" ] [ configLink ] ]
        ]


configLink : Html Msg
configLink =
    div [ id "configLink" ] [ iconLink "button is-link" DashboardsRoute "fa fa-cogs" ]


maybeDetails : Model -> Html Msg
maybeDetails model =
    case model.dashboards.details of
        RemoteData.NotAsked ->
            text ""

        RemoteData.Loading ->
            text "Loading..."

        RemoteData.Success projects ->
            detailsPage model projects

        RemoteData.Failure error ->
            text (toString error)


detailsPage : Model -> DashboardDetails -> Html Msg
detailsPage model details =
    div [ class "columns is-multiline build-items" ] (List.map configItem details.configs)


configItem : ConfigDetail -> Html Msg
configItem cd =
    let
        branch =
            case List.head cd.branches of
                Just branch ->
                    branch

                _ ->
                    { name = "NO BRANCHES", builds = [] }

        wrapperClass =
            "bi-wrapper "
                ++ (if isLastBuildError branch.builds then
                        "error"
                    else
                        "success"
                   )
    in
    div [ class "column is-6 is-paddingless buildItem" ]
        [ div [ class wrapperClass ]
            [ biTitle cd.abbreviation
            , biSubTitle (getSubtitleText cd branch)
            , buildRow branch.builds
            ]
        ]


biTitle : String -> Html Msg
biTitle t =
    div [ class "bi-title" ] [ text t ]


biSubTitle : String -> Html Msg
biSubTitle t =
    div [ class "bi-sub-title" ] [ text t ]


buildRow : List Build -> Html Msg
buildRow builds =
    div [ class "columns is-marginless" ] (List.map buildItem builds)


buildItem : Build -> Html Msg
buildItem build =
    if build.status == Success then
        successItem
    else
        failureItem


successItem : Html Msg
successItem =
    div [ class "column is-1 bhLabel bh-succ" ] [ i [ class "fa fa-trophy" ] [] ]


failureItem : Html Msg
failureItem =
    div [ class "column is-1 bhLabel bh-fail" ] [ i [ class "fa fa-trash-o" ] [] ]


isLastBuildError : List Build -> Bool
isLastBuildError builds =
    case List.head builds of
        Just build ->
            build.status == Failure

        _ ->
            False


getSubtitleText : ConfigDetail -> Branch -> String
getSubtitleText cd branch =
    let
        branchString =
            if String.length branch.name > 0 then
                " (" ++ branch.name ++ ")"
            else
                ""
    in
    cd.name ++ branchString
