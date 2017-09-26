module Dashboards.Decoders exposing (..)

import Dashboards.Models as Dashboards exposing (BuildConfig, Dashboard)
import Decoders exposing (ownerDecoder)
import Json.Decode as Decode exposing (Decoder, andThen, fail, string, succeed)
import Json.Decode.Pipeline exposing (decode, required)
import Json.Encode as Encode


dashboardsDecoder : Decode.Decoder (List Dashboard)
dashboardsDecoder =
    Decode.list dashboardDecoder


dashboardDecoder : Decode.Decoder Dashboard
dashboardDecoder =
    decode Dashboard
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "owner" ownerDecoder
        |> required "buildConfigs" (Decode.list buildConfigDecoder)


buildConfigDecoder : Decode.Decoder BuildConfig
buildConfigDecoder =
    decode BuildConfig
        |> required "id" Decode.string
        |> required "abbreviation" Decode.string


updateDashboardEncoder : Dashboards.Model -> Encode.Value
updateDashboardEncoder model =
    let
        attributes =
            [ ( "name", Encode.string model.dashboardForm.name.value )
            , ( "buildConfigs", Encode.list <| List.map buildConfigEncoder <| model.dashboardForm.buildConfigs )
            ]
    in
    Encode.object attributes


buildConfigEncoder : Dashboards.BuildConfigForm -> Encode.Value
buildConfigEncoder config =
    let
        attributes =
            [ ( "id", Encode.string config.id )
            , ( "abbriviation", Encode.string config.abbreviation.value )
            ]
    in
    Encode.object attributes
