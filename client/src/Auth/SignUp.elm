module Auth.SignUp exposing (view)

import Auth.Models as Auth
import Html exposing (Html, button, div, form, i, section, span, text)
import Html.Attributes exposing (class, disabled)
import Html.Events exposing (onSubmit)
import Models exposing (Model)
import Msgs exposing (AuthMsg(..), Msg(AuthMsg))
import Pages.Components exposing (textField)


view : Model -> Html Msg
view model =
    div [ class "columns wrapper" ]
        [ div [ class "column" ] [ signUpForm model.auth ]
        ]


signUpForm : Auth.Model -> Html Msg
signUpForm model =
    let
        disableButton =
            not (isFormValid model)
    in
    form [ onSubmit (AuthMsg OnSubmitSignUp) ]
        [ textField model.username "text" "username" "Username" "fa-user" (ChangeUsername >> AuthMsg)
        , textField model.email "email" "email" "Email" "fa-envelope" (ChangeEmail >> AuthMsg)
        , textField model.password "password" "password" "Password" "fa-lock" (ChangePassword >> AuthMsg)
        , textField model.confirm "password" "confirm" "Confirm" "fa-lock" (ChangeConfirm >> AuthMsg)
        , div [ class "field" ]
            [ div [ class "control" ]
                [ button [ class "button is-success", disabled disableButton ]
                    [ span [ class "icon" ] [ i [ class "fa fa-check-square-o" ] [] ]
                    , span [] [ text "Sign Up" ]
                    ]
                ]
            ]
        ]


isFormValid : Auth.Model -> Bool
isFormValid model =
    model.username.isValid
        && model.email.isValid
        && model.password.isValid
        && model.confirm.isValid
