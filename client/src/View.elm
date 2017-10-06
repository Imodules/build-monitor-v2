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
import Pages.Components exposing (loginBanner, smallIconLinkButton)
import Routes exposing (Route(..))
import Routing exposing (isLoggedIn)


view : Model -> Html Msg
view model =
    case model.route of
        SignUpRoute ->
            SignUp.view model |> contentWrapper

        LoginRoute ->
            Login.view model |> contentWrapper

        DashboardRoute id ->
            DashboardView.view model

        EditDashboardRoute id ->
            if isLoggedIn model then
                DashboardAddEdit.edit model id |> contentWrapper
            else
                loginBanner model |> contentWrapper

        NewDashboardRoute ->
            if isLoggedIn model then
                DashboardAddEdit.add model |> contentWrapper
            else
                loginBanner model |> contentWrapper

        DashboardsRoute ->
            DashboardList.view model |> contentWrapper

        NotFoundRoute ->
            notFoundView model |> contentWrapper


notFoundView : Model -> Html Msg
notFoundView model =
    div [ class "notification is-info" ] [ text "I cannot find that page" ]


contentWrapper : Html Msg -> Html Msg
contentWrapper content_ =
    div [ class "container is-fluid wrapper" ] [ content_ ]
