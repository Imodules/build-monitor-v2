module Msgs exposing (..)

import Dashboards.Models exposing (Dashboard, DashboardDetails)
import Http
import Models exposing (BuildType, Project, User)
import Navigation exposing (Location)
import RemoteData exposing (WebData)
import Routes exposing (Route)
import Time exposing (Time)
import Types exposing (..)


type Msg
    = DoNothing
    | ChangeLocation Route
    | OnLocationChange Location
    | GoBack
    | SetTokenStorage Token
    | GotTokenFromStorage Token
    | AuthMsg AuthMsg
    | DashboardMsg DashboardMsg
    | OnSignUp (Result Http.Error User)
    | OnLogin (Result Http.Error User)
    | OnReAuth (Result Http.Error User)
    | Logout
    | RefreshPageData Time
    | OnFetchProjects (WebData (List Project))
    | OnFetchBuildTypes (WebData (List BuildType))


type AuthMsg
    = ChangeUsername String
    | ChangeEmail String
    | ChangePassword String
    | ChangeConfirm String
    | OnSubmitSignUp
    | OnSubmitLogin


type DashboardMsg
    = ChangeDashboardName String
    | OnFetchDashboards (WebData (List Dashboard))
    | OnFetchDetails (WebData DashboardDetails)
    | ClickBuildType Id
    | StartCreateDashboard
    | CreateDashboard
    | StartEditDashboard Id
    | EditDashboard
    | OnCreateDashboard (Result Http.Error Dashboard)
    | ChangeBuildAbbreviation Id String
    | OnSelectTabClick
    | OnConfigureTabClick
