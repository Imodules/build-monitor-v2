module Routes exposing (..)

import Types exposing (Id)


type Route
    = NotFoundRoute
    | SignUpRoute
    | LoginRoute
    | DashboardRoute Id
    | NewDashboardRoute
    | EditDashboardRoute Id
    | ConfigureDashboardRoute Id
    | DashboardsRoute


dashboards : String
dashboards =
    "/"


dashboard : Id -> String
dashboard id =
    "/dashboards/" ++ id


editDashboard : Id -> String
editDashboard id =
    "/dashboards/" ++ id ++ "/edit"


configureDashboard : Id -> String
configureDashboard id =
    "/dashboards/" ++ id ++ "/configure"


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
