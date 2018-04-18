module View exposing (..)

import Auth.Login as Login
import Auth.SignUp as SignUp
import Dashboards.AddEdit as DashboardAddEdit
import Dashboards.List as DashboardList
import Dashboards.View as DashboardView
import Html exposing (..)
import Html.Attributes exposing (class, id)
import Models exposing (Model)
import Msgs exposing (Msg)
import Pages.Components exposing (loginBanner, smallIconLinkButton)
import Routes exposing (Route(..))
import Routing exposing (isLoggedIn)


view : Model -> Html Msg
view model =
    case model.route of
        SignUpRoute ->
            SignUp.view model |> fluidContainer

        LoginRoute ->
            Login.view model |> fluidContainer

        DashboardRoute id ->
            DashboardView.view model

        EditDashboardRoute id ->
            if isLoggedIn model then
                DashboardAddEdit.edit model id |> fluidContainer
            else
                loginBanner model |> fluidContainer

        NewDashboardRoute ->
            if isLoggedIn model then
                DashboardAddEdit.add model |> fluidContainer
            else
                loginBanner model |> fluidContainer

        DashboardsRoute ->
            DashboardList.view model |> fluidContainer

        NotFoundRoute ->
            notFoundView model |> fluidContainer


notFoundView : Model -> Html Msg
notFoundView model =
    div [ class "notification is-info" ] [ text "I cannot find that page" ]


fluidContainer : Html Msg -> Html Msg
fluidContainer content_ =
    div [ id "fluidWrapper", class "container is-fluid" ] [ content_ ]
