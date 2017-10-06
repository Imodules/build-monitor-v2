module Dashboards.Components exposing (..)

import Dashboards.Models as DashboardsModel exposing (DashboardForm)
import Html exposing (Html, button, div, span, text)
import Html.Attributes exposing (class, disabled)
import Html.Events exposing (onClick)
import Msgs exposing (DashboardMsg(ChangeDashboardColumnCount, ChangeDashboardName, ChangeFailedIcon, ChangeRunningIcon, ChangeSuccessIcon), Msg)
import Pages.Components exposing (icon, iconLinkButton, textField)
import Routes exposing (Route(DashboardsRoute))


saveButton : Msg -> Bool -> Html Msg
saveButton saveMsg disabled_ =
    button [ class "button is-success", disabled disabled_, onClick saveMsg ]
        [ icon "fa fa-check-square-o"
        , span [] [ text "Save" ]
        ]


cancelButton : Html Msg
cancelButton =
    iconLinkButton "" DashboardsRoute "fa-times-circle-o" "Cancel"


dashboardNameField : DashboardForm -> Html Msg
dashboardNameField dashForm =
    div [] [ textField dashForm.name "text" "dashboardName" "Dashboard Name" "fa-tachometer" (ChangeDashboardName >> Msgs.DashboardMsg) ]


dashboardColumnCountField : DashboardForm -> Html Msg
dashboardColumnCountField dashForm =
    div [] [ textField dashForm.columnCount "text" "columnCount" "Column Count" "fa-tachometer" (ChangeDashboardColumnCount >> Msgs.DashboardMsg) ]


successIconField : DashboardForm -> Html Msg
successIconField dashForm =
    div [] [ textField dashForm.successIcon "text" "successIcon" "Success Icon" dashForm.successIcon.value (ChangeSuccessIcon >> Msgs.DashboardMsg) ]


failedIconField : DashboardForm -> Html Msg
failedIconField dashForm =
    div [] [ textField dashForm.failedIcon "text" "successIcon" "Failed Icon" dashForm.failedIcon.value (ChangeFailedIcon >> Msgs.DashboardMsg) ]


runningIconField : DashboardForm -> Html Msg
runningIconField dashForm =
    div [] [ textField dashForm.runningIcon "text" "successIcon" "Running Icon" dashForm.runningIcon.value (ChangeRunningIcon >> Msgs.DashboardMsg) ]
