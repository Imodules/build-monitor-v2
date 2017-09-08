module Auth.Update exposing (update)

import Auth.Api as Api
import Auth.Models as Auth
import Models exposing (Model)
import Msgs exposing (AuthMsg(..), Msg)
import Regex
import String
import Types exposing (TextField)


update : AuthMsg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        ( newAuth, cmd ) =
            update_ model.flags.apiUrl msg model.auth
    in
    ( { model | auth = newAuth }, cmd )


update_ : String -> AuthMsg -> Auth.Model -> ( Auth.Model, Cmd Msg )
update_ baseUrl msg model =
    case msg of
        ChangeUsername value ->
            ( { model | username = updateUsername value }, Cmd.none )

        ChangeEmail value ->
            ( { model | email = updateEmail value }, Cmd.none )

        ChangePassword value ->
            ( { model | password = updatePassword value }, Cmd.none )

        ChangeConfirm value ->
            ( { model | confirm = updateConfirm model.password.value value }, Cmd.none )

        OnSubmitSignUp ->
            ( model, Api.signUp baseUrl model )

        OnSubmitLogin ->
            ( model, Api.login baseUrl model )


updateUsername : String -> TextField
updateUsername value =
    let
        ( isValid, error ) =
            if String.length value < 5 then
                ( False, "Username must be at least 5 characters" )
            else
                ( True, "" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }


updateEmail : String -> TextField
updateEmail value =
    let
        ( isValid, error ) =
            if isValidEmail value then
                ( True, "" )
            else
                ( False, "You must enter a valid email address" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }


isValidEmail : String -> Bool
isValidEmail =
    let
        validEmail =
            Regex.regex "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
                |> Regex.caseInsensitive
    in
    Regex.contains validEmail


updatePassword : String -> TextField
updatePassword value =
    let
        ( isValid, error ) =
            if String.length value < 8 then
                ( False, "Password must be at least 8 characters" )
            else
                ( True, "" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }


updateConfirm : String -> String -> TextField
updateConfirm password value =
    let
        ( isValid, error ) =
            if password == value then
                ( True, "" )
            else
                ( False, "Confirm does not match the password" )
    in
    { value = value
    , isValid = isValid
    , isDirty = True
    , error = error
    }
