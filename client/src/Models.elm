module Models exposing (..)

import Auth.Models as Auth
import Dashboards.Models as Dashboards
import Date exposing (Date)
import RemoteData exposing (WebData)
import Routes exposing (Route)
import Time exposing (Time)
import Types exposing (..)


initialModel : Flags -> Route -> Model
initialModel flags route =
    { flags = flags
    , currentTime = 0
    , route = route
    , user = Nothing
    , auth = Auth.initialModel
    , dashboards = Dashboards.initialModel
    , projects = RemoteData.NotAsked
    , buildTypes = RemoteData.NotAsked
    }


type alias Flags =
    { apiUrl : String
    }


type alias Model =
    { flags : Flags
    , currentTime : Time
    , route : Route
    , user : Maybe User
    , auth : Auth.Model
    , dashboards : Dashboards.Model
    , projects : WebData (List Project)
    , buildTypes : WebData (List BuildType)
    }


type alias User =
    { id : Id
    , createdAt : Date
    , modifiedAt : Date
    , username : Username
    , email : Email
    , token : Token
    , lastLoginAt : Date
    }


type alias Project =
    { id : Id
    , name : String
    , description : String
    , parentProjectId : Id
    }


initialProject : Project
initialProject =
    { id = ""
    , name = "INIT"
    , description = ""
    , parentProjectId = "_Root"
    }


type alias BuildType =
    { id : Id
    , name : String
    , description : String
    , projectId : Id
    }
