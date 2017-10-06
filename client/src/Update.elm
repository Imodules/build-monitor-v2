module Update exposing (..)

import Api exposing (refreshServerProjects)
import Auth.Api exposing (reAuthenticate)
import Auth.Models
import Auth.Update as Auth
import Dashboards.Update as Dashboards
import Debug exposing (log)
import Http
import Lib
import Models exposing (..)
import Msgs exposing (Msg(..))
import Navigation exposing (back, newUrl)
import Ports exposing (logout, setTokenStorage)
import Routes exposing (Route(DashboardRoute, DashboardsRoute))
import Routing exposing (getLocationCommand, getLocationRefreshCommand, getToken, parseLocation, toPath)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        DoNothing ->
            ( model, Cmd.none )

        GotTime time ->
            ( { model | currentTime = time }, Cmd.none )

        ChangeLocation route ->
            ( model, newUrl (toPath route) )

        GoBack ->
            ( model, back 1 )

        OnLocationChange location ->
            let
                newRoute =
                    parseLocation location
            in
            ( { model | route = newRoute }, getLocationCommand model newRoute )

        SetTokenStorage token ->
            ( model, setTokenStorage token )

        GotTokenFromStorage token ->
            ( model, reAuthenticate model.flags.apiUrl token )

        Logout ->
            ( { model | user = Nothing }, logout "" )

        AuthMsg msg_ ->
            Auth.update msg_ model

        DashboardMsg msg_ ->
            Dashboards.update msg_ model

        OnSignUp result ->
            handleAuth model result

        OnLogin result ->
            handleAuth model result

        OnReAuth result ->
            handleReAuth model result

        RefreshServerProjects ->
            ( model, refreshServerProjects model.flags.apiUrl (getToken model) )

        OnRefreshServerProjects _ ->
            ( model, Cmd.none )

        RefreshPageData _ ->
            let
                cmdList =
                    getLocationRefreshCommand model model.route

                cmd =
                    case cmdList of
                        [] ->
                            Cmd.none

                        _ ->
                            Cmd.batch cmdList
            in
            ( model, cmd )

        OnFetchProjects response ->
            ( { model | projects = response }, Cmd.none )

        OnFetchBuildTypes response ->
            ( { model | buildTypes = response }, Cmd.none )


handleAuth : Model -> Result Http.Error User -> ( Model, Cmd Msg )
handleAuth model result =
    case result of
        Ok user_ ->
            ( { model | user = Just user_, auth = Auth.Models.initialModel }
            , Cmd.batch
                [ setTokenStorage user_.token
                , Lib.createCommand (ChangeLocation DashboardsRoute)
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
            let
                newModel =
                    { model | user = Just user_ }
            in
            ( newModel
            , Cmd.batch
                [ setTokenStorage user_.token
                , Lib.createCommand (RefreshPageData 0)
                ]
            )

        Err r ->
            let
                x =
                    Debug.log "error on authentication" r
            in
            ( model, Cmd.none )
