module Routes exposing (..)


type Route
    = NotFoundRoute
    | SignUpRoute
    | LoginRoute
    | DashboardRoute
    | SettingsRoute


dashboard : String
dashboard =
    "/"


signUp : String
signUp =
    "/signup"


login : String
login =
    "/login"


profile : String
profile =
    "/user"


settings : String
settings =
    "/user/settings"
