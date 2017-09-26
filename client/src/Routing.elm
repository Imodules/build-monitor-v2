module Routing exposing (..)

import Api
import Dashboards.Api as DashboardsApi
import Html exposing (Attribute)
import Html.Events exposing (onWithOptions)
import Json.Decode as Decode
import List
import Models exposing (Model)
import Msgs exposing (Msg)
import Navigation exposing (Location)
import Routes exposing (..)
import Types exposing (Token)
import UrlParser exposing (..)


matchers : Parser (Route -> a) a
matchers =
    oneOf
        [ map DashboardsRoute top
        , map DashboardsRoute (s "dashboards")
        , map NewDashboardRoute (s "dashboards" </> s "new")
        , map DashboardRoute (s "dashboards" </> string)
        , map SignUpRoute (s "signup")
        , map LoginRoute (s "login")
        ]


toPath : Route -> String
toPath route =
    case route of
        SignUpRoute ->
            signUp

        LoginRoute ->
            login

        DashboardRoute id ->
            dashboard id

        DashboardsRoute ->
            dashboards

        NewDashboardRoute ->
            newDashboard

        _ ->
            "not found"


getToken : Model -> Token
getToken model =
    case model.user of
        Just user ->
            user.token

        _ ->
            ""


getLocationCommand : Model -> Route -> Cmd Msg
getLocationCommand model route =
    let
        token =
            getToken model
    in
    if needsToLogin model route then
        Cmd.none
    else
        case route of
            NewDashboardRoute ->
                Cmd.batch
                    [ Api.fetchProjects model.flags.apiUrl token
                    , Api.fetchBuildTypes model.flags.apiUrl token
                    ]

            DashboardsRoute ->
                DashboardsApi.fetchDashboards model.flags.apiUrl token

            _ ->
                Cmd.none


authRoutes : List Route
authRoutes =
    [ DashboardsRoute
    , NewDashboardRoute
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
