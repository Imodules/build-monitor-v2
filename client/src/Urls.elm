module Urls exposing (..)

import Types exposing (Id)


signUp : String -> String
signUp baseApiUrl =
    baseApiUrl ++ "/signup"


login : String -> String
login baseApiUrl =
    baseApiUrl ++ "/login"


reAuthenticate : String -> String
reAuthenticate baseApiUrl =
    baseApiUrl ++ "/authenticate"


projects : String -> String
projects baseApiUrl =
    baseApiUrl ++ "/projects"


buildTypes : String -> String
buildTypes baseApiUrl =
    baseApiUrl ++ "/buildTypes"


dashboards : String -> String
dashboards baseApiUrl =
    baseApiUrl ++ "/dashboards"


dashboard : String -> Id -> String
dashboard baseApiUrl id =
    baseApiUrl ++ "/dashboards/" ++ id
