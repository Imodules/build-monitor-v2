module Dashboards.Components exposing (..)

import Dashboards.Models as DashboardsModel exposing (DashboardForm)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Msgs exposing (DashboardMsg(ChangeDashboardColumnCount, ChangeDashboardName, ChangeFailedIcon, ChangeRunningIcon, ChangeSuccessIcon), Msg)
import Pages.Components exposing (icon, iconLinkButton, textField)
import Routes exposing (Route(DashboardsRoute))
import Types exposing (TextField)
import Date exposing (Date)
import Date.Extra.Config.Config_en_us exposing (config)
import Date.Extra.Format as DateFormat


saveButton : Msg -> Bool -> Html Msg
saveButton saveMsg disabled_ =
    actionButton "is-success" "fa fa-check-square-o" "Save" saveMsg disabled_


deleteButton : Msg -> Bool -> Html Msg
deleteButton deleteMsg disabled_ =
    actionButton "is-danger" "fa fa-remove" "Delete" deleteMsg disabled_


cancelButton : Msg -> Bool -> Html Msg
cancelButton msg disabled_ =
    actionButton "is-default" "fa fa-times-circle-o" "Cancel" msg disabled_


actionButton : String -> String -> String -> Msg -> Bool -> Html Msg
actionButton class_ icon_ text_ msg disabled_ =
    button [ class ("button " ++ class_), disabled disabled_, onClick msg ]
        [ icon icon_
        , span [] [ text text_ ]
        ]


cancelLink : Html Msg
cancelLink =
    iconLinkButton "" DashboardsRoute "fa-times-circle-o" "Cancel"


dashboardNameField : DashboardForm -> Html Msg
dashboardNameField dashForm =
    div [] [ textField dashForm.name "text" "dashboardName" "Dashboard Name" "fa-tachometer" (ChangeDashboardName >> Msgs.DashboardMsg) ]


dashboardColumnCountField : DashboardForm -> Html Msg
dashboardColumnCountField dashForm =
    div [] [ textField dashForm.columnCount "text" "columnCount" "# of Columns for each Build (out of 12)" "fa-tachometer" (ChangeDashboardColumnCount >> Msgs.DashboardMsg) ]


successIconField : DashboardForm -> Html Msg
successIconField dashForm =
    div [] [ textField dashForm.successIcon "text" "successIcon" "Success Icon" dashForm.successIcon.value (ChangeSuccessIcon >> Msgs.DashboardMsg) ]


failedIconField : DashboardForm -> Html Msg
failedIconField dashForm =
    div [] [ textField dashForm.failedIcon "text" "successIcon" "Failed Icon" dashForm.failedIcon.value (ChangeFailedIcon >> Msgs.DashboardMsg) ]


runningIconField : DashboardForm -> Html Msg
runningIconField dashForm =
    div [] [ textField dashForm.runningIcon "text" "successIcon" "Running Icon" dashForm.runningIcon.value (ChangeRunningIcon >> Msgs.DashboardMsg) ]


dateFormatField : TextField -> String -> String -> Date -> (String -> Msg) -> Html Msg
dateFormatField field id_ labelText date_ msg_ =
    div [ class "field" ]
        [ label [ class "label", for id_ ] [ text labelText ]
        , div [ class "control has-icons-left" ]
            [ input
                [ class "input"
                , id id_
                , type_ "text"
                , value field.value
                , onInput msg_
                , required True
                ]
                []
            , span [ class "icon is-small is-left" ] [ i [ class "fa fa-calendar" ] [] ]
            ]
        , div [ class "has-text-centered has-text-grey" ] [ text (DateFormat.format config field.value date_) ]
        ]
