module Lib exposing (createCommand)

import Msgs exposing (Msg)
import Task


createCommand : Msg -> Cmd Msg
createCommand msg =
    Task.succeed msg
        |> Task.perform identity
