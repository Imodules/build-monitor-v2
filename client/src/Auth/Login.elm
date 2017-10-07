module Auth.Login exposing (view)

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
        [ div [ class "column" ] [ loginForm model.auth ]
        ]


loginForm : Auth.Model -> Html Msg
loginForm model =
    let
        disableButton =
            not (isFormValid model)
    in
    form [ onSubmit (AuthMsg OnSubmitLogin) ]
        [ textField model.username "text" "username" "Username or Email" "fa-user" (ChangeUsername >> AuthMsg)
        , textField model.password "password" "password" "Password" "fa-lock" (ChangePassword >> AuthMsg)
        , div [ class "field" ]
            [ div [ class "control" ]
                [ button [ class "button is-primary", disabled disableButton ]
                    [ span [ class "icon" ] [ i [ class "fa fa-sign-in" ] [] ]
                    , span [] [ text "Login" ]
                    ]
                ]
            ]
        ]


isFormValid : Auth.Model -> Bool
isFormValid model =
    model.username.isValid
        && model.password.isValid
