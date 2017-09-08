module Update exposing (..)

import Auth.Api exposing (reAuthenticate)
import Auth.Models
import Auth.Update as Auth
import Debug exposing (log)
import Http
import Lib
import Models exposing (..)
import Msgs exposing (Msg(..))
import Navigation exposing (back, newUrl)
import Ports exposing (logout, setTokenStorage)
import Routes exposing (Route(DashboardRoute))
import Routing exposing (getLocationCommand, parseLocation, toPath)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        DoNothing ->
            ( model, Cmd.none )

        ChangeLocation route ->
            ( model, newUrl (toPath route) )

        GoBack ->
            ( model, back 1 )

        OnLocationChange location ->
            let
                newRoute =
                    parseLocation location

                newCommand =
                    getLocationCommand model newRoute
            in
            ( { model | route = newRoute }, newCommand )

        Poll _ ->
            ( model, getLocationCommand model model.route )

        SetTokenStorage token ->
            ( model, setTokenStorage token )

        GotTokenFromStorage token ->
            ( model, reAuthenticate model.flags.apiUrl token )

        Logout ->
            ( { model | user = Nothing }, logout "" )

        AuthMsg msg_ ->
            Auth.update msg_ model

        OnSignUp result ->
            handleAuth model result

        OnLogin result ->
            handleAuth model result

        OnReAuth result ->
            handleReAuth model result


handleAuth : Model -> Result Http.Error User -> ( Model, Cmd Msg )
handleAuth model result =
    case result of
        Ok user_ ->
            ( { model | user = Just user_, auth = Auth.Models.initialModel }
            , Cmd.batch
                [ setTokenStorage user_.token
                , Lib.createCommand (ChangeLocation DashboardRoute)
                ]
            )

        Err r ->
            let
                x =
                    Debug.log "error on authentication" r
            in
            ( model, Cmd.none )


handleReAuth : Model -> Result Http.Error User -> ( Model, Cmd Msg )
handleReAuth model result =
    case result of
        Ok user_ ->
            ( { model | user = Just user_ }, setTokenStorage user_.token )

        Err r ->
            let
                x =
                    Debug.log "error on authentication" r
            in
            ( model, Cmd.none )
