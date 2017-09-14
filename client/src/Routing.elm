module Routing exposing (..)

import Api
import Html exposing (Attribute)
import Html.Events exposing (onWithOptions)
import Json.Decode as Decode
import List
import Models exposing (Model)
import Msgs exposing (Msg)
import Navigation exposing (Location)
import Routes exposing (..)
import UrlParser exposing (..)


matchers : Parser (Route -> a) a
matchers =
    oneOf
        [ map DashboardRoute top
        , map SignUpRoute (s "signup")
        , map LoginRoute (s "login")
        , map SettingsRoute (s "settings")
        ]


toPath : Route -> String
toPath route =
    case route of
        SignUpRoute ->
            signUp

        LoginRoute ->
            login

        DashboardRoute ->
            dashboard

        SettingsRoute ->
            settings

        _ ->
            "not found"


getLocationCommand : Model -> Route -> Cmd Msg
getLocationCommand model route =
    let
        token =
            case model.user of
                Just user ->
                    user.token

                _ ->
                    ""
    in
    if needsToLogin model route then
        Cmd.none
    else
        case route of
            SettingsRoute ->
                Api.fetchProjects model.flags.apiUrl token

            _ ->
                Cmd.none


authRoutes : List Route
authRoutes =
    [ DashboardRoute
    , SettingsRoute
    ]


requiresAuth : Route -> Bool
requiresAuth route =
    List.member route authRoutes


needsToLogin : Model -> Route -> Bool
needsToLogin model route =
    requiresAuth route && (model.user == Nothing)


onLinkClick : msg -> Attribute msg
onLinkClick message =
    let
        options =
            { stopPropagation = False
            , preventDefault = True
            }
    in
    onWithOptions "click" options (Decode.succeed message)


parseLocation : Location -> Route
parseLocation location =
    case parsePath matchers location of
        Just route ->
            route

        Nothing ->
            NotFoundRoute
