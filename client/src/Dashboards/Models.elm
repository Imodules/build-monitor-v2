module Dashboards.Models exposing (..)

import RemoteData exposing (WebData)
import Types exposing (Id, Owner, TextField, initTextField, initTextFieldValue)


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
    , tab = Select
    }


initialBuildConfigForm : Id -> String -> BuildConfigForm
initialBuildConfigForm id abbr =
    { id = id
    , abbreviation = initTextFieldValue abbr
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


type alias BuildConfigForm =
    { id : Id
    , abbreviation : TextField
    }


type alias DashboardForm =
    { id : Id
    , name : TextField
    , buildConfigs : List BuildConfigForm
    , isDirty : Bool
    , tab : EditTab
    }


type EditTab
    = Select
    | Configure


buildConfigToForm : BuildConfig -> BuildConfigForm
buildConfigToForm bc =
    { id = bc.id
    , abbreviation = initTextFieldValue bc.abbreviation
    }
