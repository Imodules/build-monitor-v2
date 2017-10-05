module Dashboards.View exposing (..)

import Dashboards.Lib exposing (findVisibleBranch)
import Dashboards.Models exposing (Branch, Build, BuildStatus(Failure, Running, Success), ConfigDetail, DashboardDetails)
import Date exposing (Date)
import Date.Distance as DateDistance
import Date.Extra.Core as DateExtra
import Html exposing (Html, a, div, h2, h4, i, section, text)
import Html.Attributes exposing (class, href, id)
import List.Extra exposing (getAt)
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
    div [ class "columns is-multiline build-items" ] (List.map (\c -> configItem model details c) details.configs)


configItem : Model -> DashboardDetails -> ConfigDetail -> Html Msg
configItem model details cd =
    let
        branchIndex =
            let
                vb =
                    findVisibleBranch cd.id model.dashboards.visibleBranches
            in
            case vb of
                Just visibleBranch ->
                    visibleBranch.index

                _ ->
                    0

        branch =
            case getAt branchIndex cd.branches of
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

        itemSize =
            "is-" ++ toString details.columnCount
    in
    div [ class ("column " ++ itemSize ++ " is-paddingless buildItem") ]
        [ div [ class wrapperClass ]
            [ biTitle cd.abbreviation
            , biSubTitle (getSubtitleText cd branch)
            , buildRow branch.builds
            , bottomRow model branch.builds
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


bottomRow : Model -> List Build -> Html Msg
bottomRow model builds =
    let
        lastBuild =
            List.head builds
    in
    div [ class "level" ]
        [ div [ class "level-left leftStatus is-size-3" ] [ leftStatus model lastBuild ]
        ]


getDate : Model -> Date
getDate model =
    DateExtra.fromTime (round model.currentTime)


getAgoText : Model -> Build -> String
getAgoText model build =
    let
        dateAgo =
            DateDistance.inWords build.startDate (getDate model)
    in
    dateAgo ++ " ago"


leftStatus : Model -> Maybe Build -> Html Msg
leftStatus model maybeBuild =
    case maybeBuild of
        Just build ->
            case build.status of
                Running ->
                    div [ class "level-item" ] [ text build.statusText ]

                _ ->
                    div [ class "level-item" ] [ text (getAgoText model build) ]

        _ ->
            div [ class "level-item" ] [ text "no info" ]


buildItem : Build -> Html Msg
buildItem build =
    case build.status of
        Success ->
            successItem

        Running ->
            buildingItem

        _ ->
            failureItem


successItem : Html Msg
successItem =
    div [ class "column is-1 bhLabel bh-succ" ] [ i [ class "fa fa-trophy" ] [] ]


failureItem : Html Msg
failureItem =
    div [ class "column is-1 bhLabel bh-fail" ] [ i [ class "fa fa-trash-o" ] [] ]


buildingItem : Html Msg
buildingItem =
    div [ class "column is-1 bhLabel bh-succ" ] [ i [ class "fa fa-circle-o-notch faa-spin animated" ] [] ]


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
