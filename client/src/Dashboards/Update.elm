module Dashboards.Update exposing (..)

import Dashboards.Models as Dashboards
import Models exposing (Model)
import Msgs exposing (DashboardMsg(..), Msg(DashboardMsg))
import Types exposing (TextField)


update : DashboardMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        ( newDashboards, cmd ) =
            update_ model.flags.apiUrl msg model.dashboards
    in
    ( { model | dashboards = newDashboards }, cmd )


update_ : String -> DashboardMsg -> Dashboards.Model -> ( Dashboards.Model, Cmd Msg )
update_ baseUrl msg model =
    case msg of
        OnFetchDashboards response ->
            ( { model | dashboards = response }, Cmd.none )

        ChangeDashboardName value ->
            let
                newDashboardForm old =
                    { old | name = updateDashboardName value }
            in
            ( { model | dashboardForm = newDashboardForm model.dashboardForm }, Cmd.none )


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
