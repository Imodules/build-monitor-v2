module Pages.Components exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Msgs exposing (Msg)
import Routes exposing (Route(..))
import Routing exposing (onLinkClick, toPath)
import Types exposing (..)


link : String -> Route -> String -> Html Msg
link aClass route aText =
    a [ class aClass, href (toPath route), onLinkClick (Msgs.ChangeLocation route) ] [ text aText ]


iconLink : String -> Route -> Icon -> Html Msg
iconLink aClass route icon =
    a [ class aClass, href (toPath route), onLinkClick (Msgs.ChangeLocation route) ] [ i [ class icon ] [] ]


iconTextLink : String -> Route -> Icon -> String -> Html Msg
iconTextLink aClass route icon aText =
    a [ class aClass, href (toPath route), onLinkClick (Msgs.ChangeLocation route) ] [ i [ class icon ] [], text aText ]


iconLinkButton : String -> Route -> Icon -> String -> Html Msg
iconLinkButton aClasses route icon aText =
    a [ class ("button " ++ aClasses), href (toPath route), onLinkClick (Msgs.ChangeLocation route) ]
        [ span [ class "icon" ] [ i [ class ("fa " ++ icon) ] [] ]
        , span [] [ text aText ]
        ]


smallIconLinkButton : String -> Route -> Icon -> String -> Html Msg
smallIconLinkButton aClasses route icon aText =
    a [ class ("button is-small " ++ aClasses), href (toPath route), onLinkClick (Msgs.ChangeLocation route) ]
        [ span [ class "icon is-small" ] [ i [ class ("fa " ++ icon) ] [] ]
        , span [] [ text aText ]
        ]


icon : String -> Html Msg
icon v =
    span [ class "icon" ] [ i [ class v ] [] ]


textField : TextField -> String -> String -> String -> Icon -> (String -> Msg) -> Html Msg
textField field fieldType id_ labelText icon msg_ =
    let
        inputClass =
            if field.isDirty then
                if not field.isValid then
                    "input is-danger"
                else
                    "input is-success"
            else
                "input"
    in
    div [ class "field" ]
        [ label [ class "label", for id_ ] [ text labelText ]
        , div [ class "control has-icons-left" ]
            [ input
                [ class inputClass
                , id id_
                , type_ fieldType
                , value field.value
                , onInput msg_
                , required True
                ]
                []
            , span [ class "icon is-small is-left" ] [ i [ class ("fa " ++ icon) ] [] ]
            ]
        , p [ class "help is-danger" ] [ text field.error ]
        ]


textArea : TextField -> String -> String -> (String -> Msg) -> Html Msg
textArea field id_ labelText msg_ =
    let
        inputClass =
            if field.isDirty then
                if not field.isValid then
                    "textarea is-danger"
                else
                    "textarea is-success"
            else
                "textarea"
    in
    div [ class "field" ]
        [ label [ class "label", for id_ ] [ text labelText ]
        , div [ class "control" ]
            [ textarea
                [ class inputClass
                , id id_
                , value field.value
                , onInput msg_
                , required True
                ]
                []
            ]
        , p [ class "help is-danger" ] [ text field.error ]
        ]
