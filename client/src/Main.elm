module Main exposing (..)

import Models exposing (Flags, Model, initialModel)
import Msgs exposing (Msg)
import Navigation exposing (Location)
import Ports exposing (getTokenFromStorage, gotTokenFromStorage)
import Routes exposing (Route)
import Routing exposing (getLocationCommand, parseLocation)
import Time exposing (Time, second)
import Update exposing (update)
import View exposing (view)


initialCommands : Model -> Route -> Cmd Msg
initialCommands model currentRoute =
    getLocationCommand model currentRoute


init : Flags -> Location -> ( Model, Cmd Msg )
init flags location =
    let
        currentRoute =
            parseLocation location

        model =
            initialModel flags currentRoute

        cmds =
            initialCommands model currentRoute
    in
    ( model, Cmd.batch [ getTokenFromStorage "", cmds ] )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.batch
        [ Time.every (20 * second) Msgs.RefreshPageData
        , Time.every (5 * second) (Msgs.ChangeBranches >> Msgs.DashboardMsg)
        , gotTokenFromStorage Msgs.GotTokenFromStorage
        ]



-- MAIN


main : Program Flags Model Msg
main =
    Navigation.programWithFlags Msgs.OnLocationChange
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }
