module Routes exposing (..)

import Types exposing (Id)


type Route
    = NotFoundRoute
    | SignUpRoute
    | LoginRoute
    | DashboardRoute Id
    | NewDashboardRoute
    | DashboardsRoute
    | SettingsRoute


dashboards : String
dashboards =
    "/"


dashboard : Id -> String
dashboard id =
    "/dashboards/" ++ id


newDashboard : String
newDashboard =
    "/dashboards/new"


signUp : String
signUp =
    "/signup"


login : String
login =
    "/login"


profile : String
profile =
    "/user"
