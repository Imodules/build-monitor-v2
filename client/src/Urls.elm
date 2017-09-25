module Urls exposing (..)


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
