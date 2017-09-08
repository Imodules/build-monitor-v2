module Auth.Models exposing (Model, initialModel)

import Types exposing (TextField, initTextField)


type alias Model =
    { username : TextField
    , email : TextField
    , password : TextField
    , confirm : TextField
    }


initialModel : Model
initialModel =
    { username = initTextField
    , email = initTextField
    , password = initTextField
    , confirm = initTextField
    }
