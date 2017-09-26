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
import Routes exposing (Route(DashboardRoute, DashboardsRoute))
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

        RefreshPageData _ ->
            ( model, getLocationCommand model model.route )

        OnFetchProjects response ->
            ( { model | projects = response }, Cmd.none )

        OnFetchBuildTypes response ->
            ( { model | buildTypes = response }, Cmd.none )

        OnFetchDashboards response ->
            ( { model | dashboards = response }, Cmd.none )

        ChangeDashboardName name ->
            ( model, Cmd.none )


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
