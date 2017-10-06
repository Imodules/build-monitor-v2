module Dashboards.Decoders exposing (..)

import Dashboards.Models as Dashboards exposing (Branch, Build, BuildConfig, BuildStatus(Failure, Running, Success, Unknown), ConfigDetail, Dashboard, DashboardDetails)
import Decoders exposing (dateTimeDecoder, ownerDecoder)
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline exposing (decode, optional, required)
import Json.Encode as Encode


dashboardsDecoder : Decoder (List Dashboard)
dashboardsDecoder =
    Decode.list dashboardDecoder


dashboardDecoder : Decoder Dashboard
dashboardDecoder =
    decode Dashboard
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "columnCount" Decode.int
        |> required "successIcon" Decode.string
        |> required "failedIcon" Decode.string
        |> required "runningIcon" Decode.string
        |> required "owner" ownerDecoder
        |> required "buildConfigs" (Decode.list buildConfigDecoder)


buildConfigDecoder : Decoder BuildConfig
buildConfigDecoder =
    decode BuildConfig
        |> required "id" Decode.string
        |> required "abbreviation" Decode.string


updateDashboardEncoder : Dashboards.Model -> Encode.Value
updateDashboardEncoder model =
    let
        attributes =
            [ ( "name", Encode.string model.dashboardForm.name.value )
            , ( "successIcon", Encode.string model.dashboardForm.successIcon.value )
            , ( "failedIcon", Encode.string model.dashboardForm.failedIcon.value )
            , ( "runningIcon", Encode.string model.dashboardForm.runningIcon.value )
            , ( "columnCount", Encode.int (Result.withDefault 0 (String.toInt model.dashboardForm.columnCount.value)) )
            , ( "buildConfigs", Encode.list <| List.map buildConfigEncoder <| model.dashboardForm.buildConfigs )
            ]
    in
    Encode.object attributes


buildConfigEncoder : Dashboards.BuildConfigForm -> Encode.Value
buildConfigEncoder config =
    let
        attributes =
            [ ( "id", Encode.string config.id )
            , ( "abbreviation", Encode.string config.abbreviation.value )
            ]
    in
    Encode.object attributes


detailsDecoder : Decoder DashboardDetails
detailsDecoder =
    decode DashboardDetails
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "columnCount" Decode.int
        |> required "successIcon" Decode.string
        |> required "failedIcon" Decode.string
        |> required "runningIcon" Decode.string
        |> optional "details" (Decode.list configDetailDecoder) []


configDetailDecoder : Decoder ConfigDetail
configDetailDecoder =
    decode ConfigDetail
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "abbreviation" Decode.string
        |> required "isRunning" Decode.bool
        |> optional "branches" (Decode.list branchDecoder) []


branchDecoder : Decoder Branch
branchDecoder =
    decode Branch
        |> optional "name" Decode.string ""
        |> required "isRunning" Decode.bool
        |> optional "builds" (Decode.list buildDecoder) []


buildDecoder : Decoder Build
buildDecoder =
    decode Build
        |> required "id" Decode.int
        |> required "number" Decode.string
        |> required "status" buildStatusDecoder
        |> optional "statusText" Decode.string ""
        |> optional "progress" Decode.int 0
        |> required "startDate" dateTimeDecoder
        |> required "finishDate" dateTimeDecoder


buildStatusDecoder : Decoder BuildStatus
buildStatusDecoder =
    Decode.int
        |> Decode.andThen
            (\statusInt ->
                case statusInt of
                    1 ->
                        Decode.succeed Success

                    2 ->
                        Decode.succeed Running

                    3 ->
                        Decode.succeed Failure

                    somethingElse ->
                        Decode.fail <| "Unknown build status: " ++ toString somethingElse
            )
