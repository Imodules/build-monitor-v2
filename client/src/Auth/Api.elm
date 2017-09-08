module Auth.Api exposing (..)

import Api exposing (authGet, post)
import Auth.Decoders exposing (loginRequestEncoder, signUpRequestEncoder)
import Auth.Models as Auth
import Decoders exposing (profileDecoder)
import Http
import Msgs exposing (Msg(..))
import Types exposing (..)
import Urls


signUp : String -> Auth.Model -> Cmd Msg
signUp baseApiUrl model =
    let
        requestBody =
            signUpRequestEncoder model
                |> Http.jsonBody

        request =
            post (Urls.signUp baseApiUrl) requestBody profileDecoder
    in
    Http.send OnSignUp request


login : String -> Auth.Model -> Cmd Msg
login baseApiUrl model =
    let
        requestBody =
            loginRequestEncoder model
                |> Http.jsonBody

        request =
            post (Urls.login baseApiUrl) requestBody profileDecoder
    in
    Http.send OnLogin request


reAuthenticate : String -> Token -> Cmd Msg
reAuthenticate baseApiUrl token =
    let
        request =
            authGet (Urls.reAuthenticate baseApiUrl) token profileDecoder
    in
    Http.send OnReAuth request
