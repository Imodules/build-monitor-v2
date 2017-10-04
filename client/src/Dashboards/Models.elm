module Dashboards.Models exposing (..)

import RemoteData exposing (WebData)
import Time.DateTime as DateTime exposing (DateTime)
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
    , columnCount = initTextField
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
    , columnCount: Int
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
    , columnCount: TextField
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
    , columnCount: Int
    , configs : List ConfigDetail
    }


type alias ConfigDetail =
    { id : Id
    , name : String
    , abbreviation : String
    , branches : List Branch
    }


type alias VisibleBranch =
    { id : Id
    , size : Int
    , index : Int
    }


type alias Branch =
    { name : String
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
    , startDate : DateTime
    , finishDate : DateTime
    }
