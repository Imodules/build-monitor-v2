module Dashboards.Models exposing (..)

import Date exposing (Date)
import RemoteData exposing (WebData)
import Types exposing (Id, Owner, TextField, initTextField, initTextFieldValue)


type alias Model =
    { dashboards : WebData (List Dashboard)
    , dashboardForm : DashboardForm
    , details : WebData DashboardDetails
    , visibleBranches : List VisibleBranch
    }


initialModel : Model
initialModel =
    { dashboards = RemoteData.NotAsked
    , dashboardForm = initialFormModel
    , details = RemoteData.NotAsked
    , visibleBranches = []
    }


initialFormModel : DashboardForm
initialFormModel =
    { id = ""
    , name = initTextField
    , columnCount = initTextFieldValue "6"
    , successIcon = initTextFieldValue "fa fa-check"
    , failedIcon = initTextFieldValue "fa fa-exclamation"
    , runningIcon = initTextFieldValue "fa fa-circle-o-notch faa-spin animated"
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
    , columnCount : Int
    , successIcon : String
    , failedIcon : String
    , runningIcon : String
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
    , columnCount : TextField
    , successIcon : TextField
    , failedIcon : TextField
    , runningIcon : TextField
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


type alias DashboardDetails =
    { id : Id
    , name : String
    , columnCount : Int
    , successIcon : String
    , failedIcon : String
    , runningIcon : String
    , configs : List ConfigDetail
    }


type alias ConfigDetail =
    { id : Id
    , name : String
    , abbreviation : String
    , isRunning : Bool
    , branches : List Branch
    }


type alias VisibleBranch =
    { id : Id
    , size : Int
    , index : Int
    }


type alias Branch =
    { name : String
    , isRunning : Bool
    , builds : List Build
    }


type BuildStatus
    = Unknown
    | Success
    | Running
    | Failure


type alias Build =
    { id : Int
    , number : String
    , status : BuildStatus
    , statusText : String
    , progress : Int
    , startDate : Date
    , finishDate : Date
    }
