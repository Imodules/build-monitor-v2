module Msgs exposing (..)

import Http
import Models exposing (Project, User)
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
    | Poll Time
    | SetTokenStorage Token
    | GotTokenFromStorage Token
    | AuthMsg AuthMsg
    | OnSignUp (Result Http.Error User)
    | OnLogin (Result Http.Error User)
    | OnReAuth (Result Http.Error User)
    | Logout
    | OnFetchProjects (WebData (List Project))


type AuthMsg
    = ChangeUsername String
    | ChangeEmail String
    | ChangePassword String
    | ChangeConfirm String
    | OnSubmitSignUp
    | OnSubmitLogin
