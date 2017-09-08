module Models exposing (..)

import Auth.Models as Auth
import Routes exposing (Route)
import Time.DateTime as DateTime exposing (DateTime)
import Types exposing (..)


initialModel : Flags -> Route -> Model
initialModel flags route =
    { flags = flags
    , route = route
    , user = Nothing
    , auth = Auth.initialModel
    }


type alias Flags =
    { apiUrl : String
    }


type alias Model =
    { flags : Flags
    , route : Route
    , user : Maybe User
    , auth : Auth.Model
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
