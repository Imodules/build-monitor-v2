module Decoders exposing (..)

import Json.Decode as Decode exposing (Decoder, andThen, fail, string, succeed)
import Json.Decode.Pipeline exposing (decode, hardcoded, optional, required)
import Models exposing (User)
import Time.DateTime as DateTime exposing (DateTime)


dateTimeDecoder : Decoder DateTime
dateTimeDecoder =
    let
        convert : String -> Decoder DateTime
        convert raw =
            case DateTime.fromISO8601 raw of
                Ok date ->
                    succeed date

                Err error ->
                    fail error
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
