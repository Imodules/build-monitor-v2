module Dashboards.Models exposing (..)

import RemoteData exposing (WebData)
import Types exposing (Id, TextField, initTextField)


type alias Model =
    { dashboards : WebData (List Dashboard)
    , dashboardForm : DashboardForm
    }


initialModel : Model
initialModel =
    { dashboards = RemoteData.NotAsked
    , dashboardForm =
        { id = ""
        , name = initTextField
        , buildTypeIds = []
        }
    }


type alias Dashboard =
    { id : Id
    , name : String
    , ownerId : Id
    , buildTypeIds : List Id
    }


type alias DashboardForm =
    { id : Id
    , name : TextField
    , buildTypeIds : List Id
    }
