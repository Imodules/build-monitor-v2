module Decoders exposing (..)

import Date exposing (Date)
import Date.Extra.Utils as DateExtraUtils
import Json.Decode as Decode exposing (Decoder, andThen, fail, string, succeed)
import Json.Decode.Pipeline exposing (decode, hardcoded, optional, required)
import Models exposing (BuildType, Project, User, Model)
import Types exposing (Owner)
import Json.Encode as Encode


dateTimeDecoder : Decoder Date
dateTimeDecoder =
    let
        convert : String -> Decoder Date
        convert raw =
            succeed (DateExtraUtils.unsafeFromString raw)

        --        convert raw =
        --            case DateTime.fromISO8601 raw of
        --                Ok date ->
        --                    succeed date
        --
        --                Err error ->
        --                    fail error
    in
        string |> andThen convert


profileDecoder : Decode.Decoder User
profileDecoder =
    decode User
        |> required "id" Decode.string
        |> required "createdAt" dateTimeDecoder
        |> required "modifiedAt" dateTimeDecoder
        |> required "username" Decode.string
        |> required "email" Decode.string
        |> required "token" Decode.string
        |> required "lastLoginAt" dateTimeDecoder


ownerDecoder : Decode.Decoder Owner
ownerDecoder =
    decode Owner
        |> required "id" Decode.string
        |> required "username" Decode.string


projectsDecoder : Decode.Decoder (List Project)
projectsDecoder =
    Decode.list projectDecoder


projectDecoder : Decode.Decoder Project
projectDecoder =
    decode Project
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "description" Decode.string
        |> required "parentProjectId" Decode.string


buildTypesDecoder : Decode.Decoder (List BuildType)
buildTypesDecoder =
    Decode.list buildTypeDecoder


buildTypeDecoder : Decode.Decoder BuildType
buildTypeDecoder =
    decode BuildType
        |> required "id" Decode.string
        |> required "name" Decode.string
        |> required "description" Decode.string
        |> required "projectId" Decode.string


refreshRequestEncoder : Model -> Encode.Value
refreshRequestEncoder model =
    let
        attributes =
            [ ( "t", Encode.float model.currentTime )
            ]
    in
        Encode.object attributes
