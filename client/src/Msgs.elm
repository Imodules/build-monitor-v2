module Msgs exposing (..)

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
