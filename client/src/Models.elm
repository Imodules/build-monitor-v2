module Models exposing (..)

import Auth.Models as Auth
import Dashboards.Models as Dashboards
import RemoteData exposing (WebData)
import Routes exposing (Route)
import Time.DateTime as DateTime exposing (DateTime)
import Types exposing (..)


initialModel : Flags -> Route -> Model
initialModel flags route =
    { flags = flags
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
    , route : Route
    , user : Maybe User
    , auth : Auth.Model
    , dashboards : Dashboards.Model
    , projects : WebData (List Project)
    , buildTypes : WebData (List BuildType)
    }


type alias User =
    { id : Id
    , createdAt : DateTime
    , modifiedAt : DateTime
    , username : Username
    , email : Email
    , token : Token
    , lastLoginAt : DateTime
    }


type alias Project =
    { id : Id
    , name : String
    , description : String
    , parentProjectId : Id
    }


type alias BuildType =
    { id : Id
    , name : String
    , description : String
    , projectId : Id
    }
