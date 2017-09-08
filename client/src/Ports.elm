port module Ports exposing (..)

import Types exposing (..)


port setTokenStorage : Token -> Cmd msg


port getTokenFromStorage : String -> Cmd msg


port gotTokenFromStorage : (Token -> msg) -> Sub msg


port logout : String -> Cmd msg
