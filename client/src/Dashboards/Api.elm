module Dashboards.Api exposing (..)

import Api exposing (authGet, authPost, post)
import Dashboards.Decoders exposing (dashboardDecoder, dashboardsDecoder, updateDashboardEncoder)
import Dashboards.Models as Dashboards
import Http
import Msgs exposing (DashboardMsg(OnCreateDashboard, OnFetchDashboards), Msg)
import RemoteData
import Types exposing (Token)
import Urls


fetchDashboards : String -> String -> Cmd Msg
fetchDashboards baseApiUrl token =
    authGet (Urls.dashboards baseApiUrl) token dashboardsDecoder
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
