module Auth.Decoders exposing (..)

import Auth.Models as Auth
import Json.Encode as Encode


signUpRequestEncoder : Auth.Model -> Encode.Value
signUpRequestEncoder model =
    let
        attributes =
            [ ( "username", Encode.string model.username.value )
            , ( "email", Encode.string model.email.value )
            , ( "password", Encode.string model.password.value )
            ]
    in
    Encode.object attributes


loginRequestEncoder : Auth.Model -> Encode.Value
loginRequestEncoder model =
    let
        attributes =
            [ ( "username", Encode.string model.username.value )
            , ( "password", Encode.string model.password.value )
            ]
    in
    Encode.object attributes
