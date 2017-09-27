module View exposing (..)

import Auth.Login as Login
import Auth.SignUp as SignUp
import Dashboards.AddEdit as DashboardAddEdit
import Dashboards.List as DashboardList
import Dashboards.View as DashboardView
import Html exposing (..)
import Html.Attributes exposing (class)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (smallIconLinkButton)
import Routes exposing (Route(..))
import Routing exposing (needsToLogin)


view : Model -> Html Msg
view model =
    if needsToLogin model model.route then
        noAccess model
    else
        case model.route of
            SignUpRoute ->
                SignUp.view model |> contentWrapper

            LoginRoute ->
                Login.view model |> contentWrapper

            DashboardRoute id ->
                DashboardView.view model

            EditDashboardRoute id ->
                DashboardAddEdit.edit model id |> contentWrapper

            NewDashboardRoute ->
                DashboardAddEdit.add model |> contentWrapper

            DashboardsRoute ->
                DashboardList.view model |> contentWrapper

            NotFoundRoute ->
                notFoundView model |> contentWrapper


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


contentWrapper : Html Msg -> Html Msg
contentWrapper content_ =
    div [ class "container is-fluid wrapper" ] [ content_ ]
