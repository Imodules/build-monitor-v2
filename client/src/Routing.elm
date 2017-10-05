module Routing exposing (..)

import Api
import Dashboards.Api as DashboardsApi
import Html exposing (Attribute)
import Html.Events exposing (onWithOptions)
import Json.Decode as Decode
import Lib exposing (createCommand)
import List
import Models exposing (Model)
import Msgs exposing (DashboardMsg(StartCreateDashboard, StartEditDashboard), Msg(DashboardMsg))
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
        , map EditDashboardRoute (s "dashboards" </> string </> s "edit")
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

        EditDashboardRoute id ->
            editDashboard id

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
        routeCommand =
            getLocationCommands model route

        refreshCommand =
            getLocationRefreshCommand model route

        cmdList =
            List.append routeCommand refreshCommand
    in
        case cmdList of
            [] ->
                Cmd.none

            _ ->
                Cmd.batch cmdList


getLocationCommands : Model -> Route -> List (Cmd Msg)
getLocationCommands model route =
    let
        token =
            getToken model
    in
        case route of
            NewDashboardRoute ->
                cmdsIfLoggedIn model [ createCommand (DashboardMsg StartCreateDashboard) ]

            EditDashboardRoute id ->
                cmdsIfLoggedIn model [ createCommand (DashboardMsg (StartEditDashboard id)) ]

            DashboardsRoute ->
                [ DashboardsApi.fetchDashboards model.flags.apiUrl ]

            _ ->
                []


getLocationRefreshCommand : Model -> Route -> List (Cmd Msg)
getLocationRefreshCommand model route =
    let
        token =
            getToken model
    in
        case route of
            NewDashboardRoute ->
                cmdsIfLoggedIn model
                    [ Api.fetchProjects model.flags.apiUrl
                    , Api.fetchBuildTypes model.flags.apiUrl
                    ]

            EditDashboardRoute id ->
                cmdsIfLoggedIn model
                    [ DashboardsApi.fetchDashboards model.flags.apiUrl
                    , Api.fetchProjects model.flags.apiUrl
                    , Api.fetchBuildTypes model.flags.apiUrl
                    ]

            DashboardsRoute ->
                [ DashboardsApi.fetchDashboards model.flags.apiUrl ]

            DashboardRoute id ->
                [ DashboardsApi.dashboardDetails model.flags.apiUrl id ]

            _ ->
                []


cmdsIfLoggedIn : Model -> List (Cmd Msg) -> List (Cmd Msg)
cmdsIfLoggedIn model cmds =
    if isLoggedIn model then
        cmds
    else
        []


isLoggedIn : Model -> Bool
isLoggedIn model =
    model.user /= Nothing


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
