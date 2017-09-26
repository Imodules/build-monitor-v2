module Types exposing (..)


type alias Id =
    String


type alias Username =
    String


type alias Email =
    String


type alias Token =
    String


type alias Icon =
    String


type alias TextField =
    { value : String
    , isValid : Bool
    , isDirty : Bool
    , error : String
    }


initTextField : TextField
initTextField =
    { value = ""
    , isValid = False
    , isDirty = False
    , error = ""
    }


initTextFieldValue : String -> TextField
initTextFieldValue value =
    { value = value
    , isValid = True
    , isDirty = False
    , error = ""
    }
