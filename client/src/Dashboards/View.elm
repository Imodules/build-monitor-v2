module Dashboards.View exposing (..)

import Dashboards.Lib exposing (findVisibleBranch)
import Dashboards.Models exposing (Branch, Build, BuildStatus(Failure, Running, Success), ConfigDetail, DashboardDetails)
import Date exposing (Date)
import Date.Distance as DateDistance
import Date.Extra.Config.Config_en_us exposing (config)
import Date.Extra.Core as DateExtra
import Date.Extra.Duration as DateDuration
import Date.Extra.Format as DateFormat
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
    let
        theDate =
            getDate model
    in
    div [ class "level top-bar" ]
        [ div [ class "level-item" ] [ text (DateFormat.format config "%H:%M:%S%:z" theDate) ]
        , div [ class "level-item" ] [ text (DateFormat.format config "%a, %B %-@d %Y" theDate) ]
        , div [ class "level-item" ] [ text (DateFormat.formatUtc config "%H:%M:%S UTC" theDate) ]
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

        RemoteData.Success details ->
            detailsPage model details

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
                    { name = "NO BRANCHES", isRunning = False, builds = [] }

        branchCount =
            List.length cd.branches

        wrapperClass =
            "bi-wrapper "
                ++ (if isLastBuildError branch.builds then
                        "error"
                    else
                        "success"
                   )

        itemSize =
            "is-" ++ toString details.columnCount

        itemBaseClass =
            "column " ++ itemSize ++ " is-paddingless buildItem"

        itemClass =
            if branch.isRunning then
                itemBaseClass ++ " blink_me"
            else
                itemBaseClass
    in
    div [ class itemClass ]
        [ div [ class wrapperClass ]
            [ biTitle cd.abbreviation
            , biSubTitle (getSubtitleText cd branch (branchIndex + 1) branchCount)
            , buildRow branch.builds details.successIcon details.failedIcon details.runningIcon
            , bottomRow model branch.builds
            ]
        ]


biTitle : String -> Html Msg
biTitle t =
    div [ class "bi-title" ] [ text t ]


biSubTitle : String -> Html Msg
biSubTitle t =
    div [ class "bi-sub-title" ] [ text t ]


buildRow : List Build -> String -> String -> String -> Html Msg
buildRow builds sIcon fIcon rIcon =
    let
        biWithIcons build =
            buildItem build sIcon fIcon rIcon
    in
    div [ class "columns is-marginless" ] (List.map biWithIcons builds)


bottomRow : Model -> List Build -> Html Msg
bottomRow model builds =
    let
        maybeLastBuild =
            List.head builds
    in
    case maybeLastBuild of
        Just lastBuild ->
            div [ class "columns bottom-row" ]
                [ div [ class "column is-10 left-side is-size-3" ] [ leftStatus model lastBuild ]
                , div [ class "column is-2 right-side is-size-3" ] [ rightStatus model lastBuild ]
                ]

        _ ->
            div [ class "bottom-row" ] [ text "no info" ]


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


leftStatus : Model -> Build -> Html Msg
leftStatus model build =
    case build.status of
        Running ->
            text build.statusText

        _ ->
            text (getAgoText model build)


rightStatus : Model -> Build -> Html Msg
rightStatus model build =
    let
        duration =
            DateDuration.diff build.finishDate build.startDate

        durationMinText =
            zeroPad duration.minute ++ ":" ++ zeroPad duration.second

        durationText =
            if duration.hour > 0 then
                zeroPad duration.hour ++ ":" ++ durationMinText
            else
                durationMinText
    in
    case build.status of
        Running ->
            text (toString build.progress ++ " %")

        _ ->
            text durationText


zeroPad : Int -> String
zeroPad v =
    if v > 9 then
        toString v
    else
        "0" ++ toString v


buildItem : Build -> String -> String -> String -> Html Msg
buildItem build sIcon fIcon rIcon =
    case build.status of
        Success ->
            successItem sIcon

        Running ->
            runningItem rIcon

        _ ->
            failureItem fIcon


successItem : String -> Html Msg
successItem icon_ =
    div [ class "column is-1 bhLabel bh-succ" ] [ i [ class icon_ ] [] ]


failureItem : String -> Html Msg
failureItem icon_ =
    div [ class "column is-1 bhLabel bh-fail" ] [ i [ class icon_ ] [] ]


runningItem : String -> Html Msg
runningItem icon_ =
    div [ class "column is-1 bhLabel bh-succ" ] [ i [ class icon_ ] [] ]


isLastBuildError : List Build -> Bool
isLastBuildError builds =
    case List.head builds of
        Just build ->
            build.status == Failure

        _ ->
            False


getSubtitleText : ConfigDetail -> Branch -> Int -> Int -> String
getSubtitleText cd branch i branchCount =
    let
        branchNameString =
            if String.length branch.name > 0 then
                " (" ++ branch.name ++ ")"
            else
                ""

        branchText =
            if branchCount == 1 then
                branchNameString
            else
                branchNameString ++ " [" ++ toString i ++ "/" ++ toString branchCount ++ "]"
    in
    cd.name ++ branchText
