module Dashboards.Api exposing (..)

import Api exposing (authDelete, authGet, authPost, authPut, post)
import Dashboards.Decoders exposing (dashboardDecoder, dashboardsDecoder, detailsDecoder, updateDashboardEncoder)
import Dashboards.Models as Dashboards
import Http
import Msgs exposing (DashboardMsg(OnCreateDashboard, OnDeleteDashboard, OnFetchDashboards, OnFetchDetails), Msg)
import RemoteData
import Types exposing (Id, Token)
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


dashboardDetails : String -> Id -> Cmd Msg
dashboardDetails baseApiUrl id =
    Http.get (Urls.dashboard baseApiUrl id) detailsDecoder
        |> RemoteData.sendRequest
        |> Cmd.map (OnFetchDetails >> Msgs.DashboardMsg)


deleteDashboard : String -> Token -> Id -> Cmd Msg
deleteDashboard baseApiUrl token id =
    let
        request =
            authDelete (Urls.dashboard baseApiUrl id) token
    in
    Http.send (OnDeleteDashboard >> Msgs.DashboardMsg) request
