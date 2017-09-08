module View exposing (..)

import Auth.Login as Login
import Auth.SignUp as SignUp
import Html exposing (..)
import Html.Attributes exposing (class)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (smallIconLinkButton)
import Pages.Dashboard as Dashboard
import Pages.Profile as Profile
import Routes exposing (Route(..))
import Routing exposing (needsToLogin)


view : Model -> Html Msg
view model =
    div []
        [ section [ class "section" ]
            [ content model ]
        ]


content : Model -> Html Msg
content model =
    if needsToLogin model model.route then
        noAccess model
    else
        case model.route of
            SignUpRoute ->
                SignUp.view model

            LoginRoute ->
                Login.view model

            DashboardRoute ->
                Dashboard.view model

            SettingsRoute ->
                Profile.view model

            NotFoundRoute ->
                notFoundView model


notFoundView : Model -> Html Msg
notFoundView model =
    div [ class "notification is-info" ] [ text "I cannot find that page" ]


noAccess : Model -> Html Msg
noAccess model =
    div [ class "notification is-warning" ]
        [ text "You must "
        , smallIconLinkButton "is-dark is-outlined" LoginRoute "fa-sign-in" "Login"
        , text " or "
        , smallIconLinkButton "is-primary is-outlined" SignUpRoute "fa-check-square-o" "Sign Up"
        ]
