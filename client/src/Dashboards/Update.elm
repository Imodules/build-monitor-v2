module Dashboards.Update exposing (..)

import Dashboards.Api as Api
import Dashboards.Models as Dashboards
import List.Extra exposing (find)
import Models exposing (Model)
import Msgs exposing (DashboardMsg(..), Msg(DashboardMsg))
import RemoteData
import Routing exposing (getToken)
import Types exposing (TextField, Token, initTextFieldValue)


update : DashboardMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        ( newDashboards, cmd ) =
            update_ model.flags.apiUrl (getToken model) msg model.dashboards
    in
    ( { model | dashboards = newDashboards }, cmd )


update_ : String -> Token -> DashboardMsg -> Dashboards.Model -> ( Dashboards.Model, Cmd Msg )
update_ baseUrl token msg model =
    case msg of
        OnFetchDashboards response ->
            ( { model | dashboards = response }, Cmd.none )

        ChangeDashboardName value ->
            let
                newDashboardForm old =
                    { old | name = updateDashboardName value, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        ClickBuildType id ->
            let
                updatedList old =
                    if List.member id old then
                        List.filter (\i -> i /= id) old
                    else
                        id :: old

                newDashboardForm old =
                    { old | buildTypeIds = updatedList old.buildTypeIds, isDirty = True }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        CreateDashboard ->
            ( model, Api.createDashboard baseUrl token model )

        EditDashboard ->
            ( model, Api.editDashboard baseUrl token model )

        StartEditDashboard id ->
            let
                dashboards =
                    case model.dashboards of
                        RemoteData.Success dashboards ->
                            dashboards

                        _ ->
                            []

                maybeDashboardToEdit =
                    find (\i -> i.id == id) dashboards

                newDashboardForm old =
                    case maybeDashboardToEdit of
                        Just dashboard ->
                            if String.isEmpty old.id || old.id /= id then
                                { id = dashboard.id
                                , name = initTextFieldValue dashboard.name
                                , buildTypeIds = dashboard.buildTypeIds
                                , isDirty = False
                                }
                            else
                                old

                        _ ->
                            Dashboards.initialFormModel
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )

        OnCreateDashboard result ->
            ( model, Cmd.none )


updateDashboardName : String -> TextField
updateDashboardName value =
    let
        ( isValid, error ) =
            if String.length value < 5 then
                ( False, "Name must be at least 5 characters" )
            else
                ( True, "" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }
