module Dashboards.Models exposing (..)

import RemoteData exposing (WebData)
import Types exposing (Id, Owner, TextField, initTextField)


type alias Model =
    { dashboards : WebData (List Dashboard)
    , dashboardForm : DashboardForm
    }


initialModel : Model
initialModel =
    { dashboards = RemoteData.NotAsked
    , dashboardForm = initialFormModel
    }


initialFormModel : DashboardForm
initialFormModel =
    { id = ""
    , name = initTextField
    , buildConfigs = []
    , isDirty = False
    }


initialBuildConfig : Id -> String -> BuildConfig
initialBuildConfig id abbr =
    { id = id
    , abbreviation = abbr
    }


type alias Dashboard =
    { id : Id
    , name : String
    , owner : Owner
    , buildConfigs : List BuildConfig
    }


type alias BuildConfig =
    { id : Id
    , abbreviation : String
    }


type alias DashboardForm =
    { id : Id
    , name : TextField
    , buildConfigs : List BuildConfig
    , isDirty : Bool
    }
