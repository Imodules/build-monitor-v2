module Dashboards.Decoders exposing (..)

import Dashboards.Models as Dashboards exposing (Dashboard)
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
        |> required "ownerId" Decode.string
        |> required "buildTypeIds" (Decode.list string)


updateDashboardEncoder : Dashboards.Model -> Encode.Value
updateDashboardEncoder model =
    let
        attributes =
            [ ( "name", Encode.string model.dashboardForm.name.value )
            , ( "buildTypeIds", Encode.list <| List.map Encode.string <| model.dashboardForm.buildTypeIds )
            ]
    in
    Encode.object attributes
