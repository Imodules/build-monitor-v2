module Dashboards.Api exposing (..)

import Api exposing (authGet, authPost, authPut, post)
import Dashboards.Decoders exposing (dashboardDecoder, dashboardsDecoder, updateDashboardEncoder)
import Dashboards.Models as Dashboards
import Http
import Msgs exposing (DashboardMsg(OnCreateDashboard, OnFetchDashboards), Msg)
import RemoteData
import Types exposing (Token)
import Urls


fetchDashboards : String -> Cmd Msg
fetchDashboards baseApiUrl =
    Http.get (Urls.dashboards baseApiUrl) dashboardsDecoder
        |> RemoteData.sendRequest
        |> Cmd.map (OnFetchDashboards >> Msgs.DashboardMsg)


createDashboard : String -> Token -> Dashboards.Model -> Cmd Msg
createDashboard baseApiUrl token model =
    let
        requestBody =
            updateDashboardEncoder model
                |> Http.jsonBody

        request =
            authPost (Urls.dashboards baseApiUrl) token requestBody dashboardDecoder
    in
    Http.send (OnCreateDashboard >> Msgs.DashboardMsg) request


editDashboard : String -> Token -> Dashboards.Model -> Cmd Msg
editDashboard baseApiUrl token model =
    let
        requestBody =
            updateDashboardEncoder model
                |> Http.jsonBody

        request =
            authPut (Urls.dashboard baseApiUrl model.dashboardForm.id) token requestBody dashboardDecoder
    in
    Http.send (OnCreateDashboard >> Msgs.DashboardMsg) request
